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

package db

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/aheadaviation/Users/users"
)

type Database interface {
	Init() error
	GetUserByName(string) (users.User, error)
	GetUser(string) (users.User, error)
	GetUsers() ([]users.User, error)
	CreateUser(*users.User) error
	GetUserAttributes(*users.User) error
	GetAddress(string) (users.Address, error)
	GetAddresses() ([]users.Address, error)
	CreateAddress(*users.Address, string) error
	GetCard(string) (users.Card, error)
	GetCards() ([]users.Card, error)
	CreateCard(*users.Card, string) error
	Delete(string, string) error
	Ping() error
}

var (
	database              string
	DefaultDb             Database
	DBTypes               = map[string]Database{}
	ErrNoDatabaseFound    = "No database with name %v registered"
	ErrNoDatabaseSelected = errors.New("No DB selected")
)

func init() {
	flag.StringVar(&database, "database", os.Getenv("USERS_DATABASE"), "Database to use for Users")
}

func Init() error {
	if database == "" {
		return ErrNoDatabaseSelected
	}
	err := Set()
	if err != nil {
		return err
	}
	return DefaultDb.Init()
}

func Set() error {
	if v, ok := DBTypes[database]; ok {
		DefaultDb = v
		return nil
	}
	return fmt.Errorf(ErrNoDatabaseFound, database)
}

func Register(name string, db Database) {
	DBTypes[name] = db
}

func CreateUser(u *users.User) error {
	return DefaultDb.CreateUser(u)
}

func GetUserByName(n string) (users.User, error) {
	u, err := DefaultDb.GetUserByName(n)
	if err == nil {
		u.AddLinks()
	}
	return u, err
}

func GetUser(n string) (users.User, error) {
	u, err := DefaultDb.GetUser(n)
	if err == nil {
		u.AddLinks()
	}
	return u, err
}

func GetUsers() ([]users.User, error) {
	us, err := DefaultDb.GetUsers()
	for k, _ := range us {
		us[k].AddLinks()
	}
	return us, err
}

func GetUserAttributes(u *users.User) error {
	err := DefaultDb.GetUserAttributes(u)
	if err != nil {
		return err
	}
	for k, _ := range u.Addresses {
		u.Addresses[k].AddLinks()
	}
	for k, _ := range u.Cards {
		u.Cards[k].AddLinks()
	}
	return nil
}

func CreateAddress(a *users.Address, userid string) error {
	return DefaultDb.CreateAddress(a, userid)
}

func GetAddress(n string) (users.Address, error) {
	a, err := DefaultDb.GetAddress(n)
	if err == nil {
		a.AddLinks()
	}
	return a, err
}

func GetAddresses() ([]users.Address, error) {
	as, err := DefaultDb.GetAddresses()
	for k, _ := range as {
		as[k].AddLinks()
	}
	return as, err
}

func CreateCard(c *users.Card, userid string) error {
	return DefaultDb.CreateCard(c, userid)
}

func GetCard(n string) (users.Card, error) {
	return DefaultDb.GetCard(n)
}

func GetCards() ([]users.Card, error) {
	cs, err := DefaultDb.GetCards()
	for k, _ := range cs {
		cs[k].AddLinks()
	}
	return cs, err
}

func Delete(entity, id string) error {
	return DefaultDb.Delete(entity, id)
}

func Ping() error {
	return DefaultDb.Ping()
}
