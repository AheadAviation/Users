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
	//GetUserAttributes(*users.User) error
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
	return u, err
}

func GetUser(n string) (users.User, error) {
	u, err := DefaultDb.GetUser(n)
	return u, err
}

func GetUsers() ([]users.User, error) {
	us, err := DefaultDb.GetUsers()
	return us, err
}

// func GetUserAttributes(u *users.User) error {
// 	err := DefaultDb.GetUserAttributes(u)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func Delete(entity, id string) error {
	return DefaultDb.Delete(entity, id)
}

func Ping() error {
	return DefaultDb.Ping()
}
