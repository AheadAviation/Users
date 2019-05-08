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
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/aheadaviation/Users/db"
	"github.com/aheadaviation/Users/users"
)

var (
	ErrUnauthorized = errors.New("Unauthorized")
)

type Service interface {
	Login(username, password string) (users.User, error)
	Register(username, password, email, first, last string) (string, error)
	GetUsers(id string) ([]users.User, error)
	PostUser(u users.User) (string, error)
	GetAddresses(id string) ([]users.Address, error)
	PostAddress(a users.Address, userid string) (string, error)
	GetCards(id string) ([]users.Card, error)
	PostCard(c users.Card, userid string) (string, error)
	Delete(entity, id string) error
	Health() []Health
}

func NewFixedService() Service {
	return &fixedService{}
}

type fixedService struct{}

type Health struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Time    string `json:"time"`
}

func (s *fixedService) Login(username, password string) (users.User, error) {
	u, err := db.GetUserByName(username)
	if err != nil {
		return users.New(), err
	}
	if u.Password != calculatePassHash(password, u.Salt) {
		return users.New(), ErrUnauthorized
	}
	return u, nil
}

func (s *fixedService) Register(username, password, email, first, last string) (string, error) {
	u := users.New()
	u.Username = username
	u.Password = calculatePassHash(password, u.Salt)
	u.Email = email
	u.FirstName = first
	u.LastName = last
	err := db.CreateUser(&u)
	return u.UserID, err
}

func (s *fixedService) GetUsers(id string) ([]users.User, error) {
	if id == "" {
		us, err := db.GetUsers()
		for k, u := range us {
			us[k] = u
		}
		return us, err
	}
	u, err := db.GetUser(id)
	return []users.User{u}, err
}

func (s *fixedService) PostUser(u users.User) (string, error) {
	u.NewSalt()
	u.Password = calculatePassHash(u.Password, u.Salt)
	err := db.CreateUser(&u)
	return u.UserID, err
}

func (s *fixedService) GetAddresses(id string) ([]users.Address, error) {
	if id == "" {
		as, err := db.GetAddresses()
		for k, a := range as {
			a.AddLinks()
			as[k] = a
		}
		return as, err
	}
	a, err := db.GetAddress(id)
	a.AddLinks()
	return []users.Address{a}, err
}

func (s *fixedService) PostAddress(a users.Address, userid string) (string, error) {
	err := db.CreateAddress(&a, userid)
	return a.ID, err
}

func (s *fixedService) GetCards(id string) ([]users.Card, error) {
	if id == "" {
		cs, err := db.GetCards()
		for k, c := range cs {
			c.AddLinks()
			cs[k] = c
		}
		return cs, err
	}
	c, err := db.GetCard(id)
	c.AddLinks()
	return []users.Card{c}, err
}

func (s *fixedService) PostCard(c users.Card, userid string) (string, error) {
	err := db.CreateCard(&c, userid)
	return c.ID, err
}

func (s *fixedService) Delete(entity, id string) error {
	return db.Delete(entity, id)
}

func (s *fixedService) Health() []Health {
	var health []Health
	dbstatus := "OK"

	err := db.Ping()
	if err != nil {
		dbstatus = "err"
	}

	app := Health{"user", "OK", time.Now().String()}
	db := Health{"user-db", dbstatus, time.Now().String()}

	health = append(health, app)
	health = append(health, db)

	return health
}

func calculatePassHash(pass, salt string) string {
	h := sha1.New()
	io.WriteString(h, salt)
	io.WriteString(h, pass)
	return fmt.Sprintf("%x", h.Sum(nil))
}
