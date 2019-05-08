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
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"

	"github.com/aheadaviation/Users/users"
)

type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) Login(username, password string) (user users.User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Login",
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Login(username, password)
}

func (mw loggingMiddleware) Register(username, password, email, first, last string) (string, error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Register",
			"username", username,
			"email", email,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Register(username, password, email, first, last)
}

func (mw loggingMiddleware) PostUser(user users.User) (id string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "PostUser",
			"username", user.Username,
			"email", user.Email,
			"result", id,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.PostUser(user)
}

func (mw loggingMiddleware) GetUsers(id string) (u []users.User, err error) {
	defer func(begin time.Time) {
		who := id
		if who == "" {
			who = "all"
		}
		mw.logger.Log(
			"method", "GetUsers",
			"id", who,
			"result", len(u),
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.GetUsers(id)
}

func (mw loggingMiddleware) PostAddress(a users.Address, id string) (string, error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "PostAddress",
			"street", a.Street,
			"number", a.Number,
			"user", id,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.PostAddress(a, id)
}

func (mw loggingMiddleware) GetAddresses(id string) (a []users.Address, err error) {
	defer func(begin time.Time) {
		who := id
		if who == "" {
			who = "all"
		}
		mw.logger.Log(
			"method", "GetAddresses",
			"id", who,
			"result", len(a),
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.GetAddresses(id)
}

func (mw loggingMiddleware) PostCard(c users.Card, id string) (string, error) {
	defer func(begin time.Time) {
		cc := c
		cc.MaskCC()
		mw.logger.Log(
			"method", "PostCard",
			"card", cc.LongNum,
			"user", id,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.PostCard(c, id)
}

func (mw loggingMiddleware) GetCards(id string) (c []users.Card, err error) {
	defer func(begin time.Time) {
		who := id
		if who == "" {
			who = "all"
		}
		mw.logger.Log(
			"method", "GetCards",
			"id", who,
			"result", len(c),
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.GetCards(id)
}

func (mw loggingMiddleware) Delete(entity, id string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Delete",
			"entity", entity,
			"id", id,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Delete(entity, id)
}

func (mw loggingMiddleware) Health() (health []Health) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Health",
			"result", len(health),
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Health()
}

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func NewInstrumentingService(requestCount metrics.Counter, requestLatency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   requestCount,
		requestLatency: requestLatency,
		Service:        s,
	}
}

func (s *instrumentingService) Login(username, password string) (users.User, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "login").Add(1)
		s.requestLatency.With("method", "login").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Login(username, password)
}

func (s *instrumentingService) Register(username, password, email, first, last string) (string, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "register").Add(1)
		s.requestLatency.With("method", "register").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Register(username, password, email, first, last)
}

func (s *instrumentingService) PostUser(user users.User) (string, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "postUser").Add(1)
		s.requestLatency.With("method", "postUser").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.PostUser(user)
}

func (s *instrumentingService) GetUsers(id string) (u []users.User, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "getUsers").Add(1)
		s.requestLatency.With("method", "getUsers").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetUsers(id)
}

func (s *instrumentingService) PostAddress(a users.Address, id string) (string, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "postAddress").Add(1)
		s.requestLatency.With("method", "postAddress").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.PostAddress(a, id)
}

func (s *instrumentingService) GetAddresses(id string) ([]users.Address, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "getAddresses").Add(1)
		s.requestLatency.With("method", "getAddresses").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetAddresses(id)
}

func (s *instrumentingService) PostCard(c users.Card, id string) (string, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "postCard").Add(1)
		s.requestLatency.With("method", "postCard").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.PostCard(c, id)
}

func (s *instrumentingService) GetCards(id string) ([]users.Card, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "getCards").Add(1)
		s.requestLatency.With("method", "getCards").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetCards(id)
}

func (s *instrumentingService) Delete(entity, id string) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "delete").Add(1)
		s.requestLatency.With("method", "delete").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Delete(entity, id)
}

func (s *instrumentingService) Health() []Health {
	defer func(begin time.Time) {
		s.requestCount.With("method", "health").Add(1)
		s.requestLatency.With("method", "health").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Health()
}
