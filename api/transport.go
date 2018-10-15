// Copyright © 2018 Tim Curless <tim.curless@thinkahead.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os/user"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ErrInvalidRequest = errors.New("Invalid request")
)

func MakeHTTPHandler(e Endpoints, logger log.Logger, tracer stdopentracing.Tracer) *mux.Router {
	r := mux.NewRouter().StrictSlash(false)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("GET").Path("/login").Handler(httptransport.NewServer(
		e.LoginEndpoint,
		decodeLoginRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "GET /login", logger)))...,
	))
	r.Methods("POST").Path("/register").Handler(httptransport.NewServer(
		e.RegisterEndpoint,
		decodeRegisterRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "POST /register", logger)))...,
	))
	r.Methods("GET").Path("/customers").Handler(httptransport.NewServer(
		e.UserGetEndpoint,
		decodeGetRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "GET /customers", logger)))...,
	))
	r.Methods("POST").Path("/customers").Handler(httptransport.NewServer(
		e.UserPostEndpoint,
		decodeUserRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "POST /customers", logger)))...,
	))
	r.Methods("DELETE").Path("/").Handler(httptransport.NewServer(
		e.DeleteEndpoint,
		decodeDeleteRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "DELETE /", logger)))...,
	))
	r.Methods("GET").Path("/health").Handler(httptransport.NewServer(
		e.HealthEndpoint,
		decodeHealthRequest,
		encodeHealthResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "GET /health", logger)))...,
	))
	r.Handle("/metrics", promhttp.Handler())
	return r
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	code := http.StatusInternalServerError
	switch err {
	case ErrUnauthorized:
		code = http.StatusUnauthorized
	}
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/hal+json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":       err.Error(),
		"status_code": code,
		"status_text": http.StatusText(code),
	})
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	u, p, ok := r.BasicAuth()
	if !ok {
		return loginRequest{}, ErrUnauthorized
	}

	return loginRequest{
		Username: u,
		Password: p,
	}, nil
}

func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	reg := registerRequest{}
	err := json.NewDecoder(r.Body).Decode(&reg)
	if err != nil {
		return nil, err
	}
	return reg, nil
}

func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	d := deleteRequest{}
	u := strings.Split(r.URL.Path, "/")
	if len(u) == 3 {
		d.Entity = u[1]
		d.ID = u[2]
		return d, nil
	}
	return d, ErrInvalidRequest
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	g := GetRequest{}
	u := strings.Split(r.URL.Path, "/")
	if len(u) > 2 {
		g.ID = u[2]
		if len(u) > 3 {
			g.Attr = u[3]
		}
	}
	return g, nil
}

func decodeUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	u := user.User{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func decodeHealthRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return struct{}{}, nil
}

func encodeHealthResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return encodeResponse(ctx, w, response.(healthResponse))
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/hal+json")
	return json.NewEncoder(w).Encode(response)
}
