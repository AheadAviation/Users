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

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/zipkin"
	stdzipkin "github.com/openzipkin/zipkin-go"

	"github.com/aheadaviation/Users/users"
)

type Endpoints struct {
	LoginEndpoint    endpoint.Endpoint
	RegisterEndpoint endpoint.Endpoint
	UserGetEndpoint  endpoint.Endpoint
	UserPostEndpoint endpoint.Endpoint
	DeleteEndpoint   endpoint.Endpoint
	HealthEndpoint   endpoint.Endpoint
}

func MakeEndpoints(s Service, tracer *stdzipkin.Tracer) Endpoints {
	return Endpoints{
		LoginEndpoint:    zipkin.TraceEndpoint(tracer, "GET /login")(MakeLoginEndpoint(s)),
		RegisterEndpoint: zipkin.TraceEndpoint(tracer, "POST /register")(MakeRegisterEndpoint(s)),
		HealthEndpoint:   zipkin.TraceEndpoint(tracer, "GET /health")(MakeHealthEndpoint(s)),
		UserGetEndpoint:  zipkin.TraceEndpoint(tracer, "GET /customers")(MakeUserGetEndpoint(s)),
		UserPostEndpoint: zipkin.TraceEndpoint(tracer, "POST /customers")(MakeUserPostEndpoint(s)),
		DeleteEndpoint:   zipkin.TraceEndpoint(tracer, "DELETE /")(MakeDeleteEndpoint(s)),
	}
}

func MakeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(loginRequest)
		u, err := s.Login(req.Username, req.Password)
		return userResponse{User: u}, err
	}
}

func MakeRegisterEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(registerRequest)
		id, err := s.Register(req.Username, req.Password, req.Email, req.FirstName, req.LastName)
		return postResponse{ID: id}, err
	}
}

func MakeUserGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetRequest)

		usrs, err := s.GetUsers(req.ID)
		if req.ID == "" {
			return EmbedStruct{usersResponse{Users: usrs}}, err
		}
		if len(usrs) == 0 {
			return users.User{}, err
		}
		user := usrs[0]
		return user, err
	}
}

func MakeUserPostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(users.User)
		id, err := s.PostUser(req)
		return postResponse{ID: id}, err
	}
}

func MakeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteRequest)
		err = s.Delete(req.Entity, req.ID)
		if err == nil {
			return statusResponse{Status: true}, err
		}
		return statusResponse{Status: false}, err
	}
}

func MakeHealthEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		health := s.Health()
		return healthResponse{Health: health}, nil
	}
}

type GetRequest struct {
	ID   string
	Attr string
}

type loginRequest struct {
	Username string
	Password string
}

type userResponse struct {
	User users.User `json:"user"`
}

type usersResponse struct {
	Users []users.User `json:"customer"`
}

type registerRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type statusResponse struct {
	Status bool `json:"status"`
}

type postResponse struct {
	ID string `json:"id"`
}

type deleteRequest struct {
	Entity string
	ID     string
}

type healthRequest struct {
	//
}

type healthResponse struct {
	Health []Health `json:"health"`
}

type EmbedStruct struct {
	Embed interface{} `json:"_embedded"`
}
