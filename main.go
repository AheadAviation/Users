package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	corelog "log"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/sd/consul"
	stdconsul "github.com/hashicorp/consul/api"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	//commonMiddleware "github.com/weaveworks/common/middleware"

	"github.com/aheadaviation/Users/api"
	"github.com/aheadaviation/Users/db"
	"github.com/aheadaviation/Users/db/mongodb"
)

var (
	port        string
	zipkinV2URL string
	consulAddr  string
)

var (
	HTTPLatency = stdprometheus.NewHistogramVec(stdprometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Time (in seconds) spent serving HTTP requests.",
		Buckets: stdprometheus.DefBuckets,
	}, []string{"method", "path", "status_code", "isWS"})
)

const (
	ServiceName = "users"
)

func init() {
	stdprometheus.MustRegister(HTTPLatency)
	flag.StringVar(&zipkinV2URL, "zipkin", os.Getenv("ZIPKIN_V2_URL"), "zipkin v2 address")
	flag.StringVar(&port, "port", "8084", "Port on which to run")
	flag.StringVar(&consulAddr, "consul_addr", os.Getenv("CONSUL_ADDR"), "Address of consul agent")
	db.Register("mongodb", &mongodb.Mongo{})
}

func main() {

	flag.Parse()
	errc := make(chan error)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	if consulAddr == "" {
		logger.Log("error", "no consul address set")
		os.Exit(1)
	}

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	host := strings.Split(localAddr.String(), ":")[0]
	defer conn.Close()

	var tracer stdopentracing.Tracer
	{
		if zipkinV2URL == "" {
			tracer = stdopentracing.NoopTracer{}
		} else {
			logger := log.With(logger, "tracer", "Zipkin")
			logger.Log("addr", zipkinV2URL)
			collector, err := zipkin.NewHTTPCollector(
				zipkinV2URL,
				zipkin.HTTPLogger(logger),
			)
			if err != nil {
				logger.Log("err", err)
				os.Exit(1)
			}
			tracer, err = zipkin.NewTracer(
				zipkin.NewRecorder(collector, false, fmt.Sprintf("%v:%v", host, port), ServiceName),
			)
			if err != nil {
				logger.Log("err", err)
				os.Exit(1)
			}
		}
		stdopentracing.InitGlobalTracer(tracer)
	}

	dbconn := false
	for !dbconn {
		err := db.Init()
		if err != nil {
			if err == db.ErrNoDatabaseSelected {
				corelog.Fatal(err)
			}
			corelog.Print(err)
		} else {
			dbconn = true
		}
	}

	fieldKeys := []string{"method"}

	var service api.Service
	{
		service = api.NewFixedService()
		service = api.LoggingMiddleware(logger)(service)
		service = api.NewInstrumentingService(
			kitprometheus.NewCounterFrom(
				stdprometheus.CounterOpts{
					Namespace: "microservices_demo",
					Subsystem: "users",
					Name:      "request_count",
					Help:      "Number of requests received",
				},
				fieldKeys),
			kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
				Namespace: "microservices_demo",
				Subsystem: "user",
				Name:      "request_latency_microseconds",
				Help:      "Total duration of requests in microseconds.",
			}, fieldKeys),
			service,
		)
	}

	endpoints := api.MakeEndpoints(service, tracer)

	router := api.MakeHTTPHandler(endpoints, logger, tracer)

	// httpMiddleware := []commonMiddleware.Interface{
	// 	commonMiddleware.Instrument{
	// 		Duration:     HTTPLatency,
	// 		RouteMatcher: router,
	// 	},
	// }

	//handler := commonMiddleware.Merge(httpMiddleware...).Wrap(router)

	stdClient, err := stdconsul.NewClient(&stdconsul.Config{
		Address: consulAddr,
	})
	if err != nil {
		logger.Log("error", "couldn't connect to consul agent")
		os.Exit(1)
	}
	sdclient := consul.NewClient(stdClient)

	intPort, _ := strconv.Atoi(port)
	reg := &stdconsul.AgentServiceRegistration{
		ID:                "USR",
		Name:              "users",
		Tags:              []string{"app=bagshop"},
		Port:              intPort,
		Address:           localAddr.String(),
		EnableTagOverride: false,
	}

	registrar := consul.NewRegistrar(sdclient, reg, log.With(logger, "component", "registrar"))
	registrar.Register()
	defer registrar.Deregister()

	go func() {
		logger.Log("transport", "http", "port", port)
		errc <- http.ListenAndServe(fmt.Sprintf(":%v", port), router)
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("exit", <-errc)
}
