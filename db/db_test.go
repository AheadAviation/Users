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

//
// import (
// 	"errors"
// 	"testing"
//
// 	"github.com/aheadaviation/Users/users"
// )
//
// var (
// 	TestDB       = fake{}
// 	ErrFakeError = errors.New("Fake Error")
// )
//
// func TestInit(t *testing.T) {
// 	err := Init()
// 	if err == nil {
// 		t.Error("Expected no registered db error")
// 	}
// 	Register("test", TestDB)
// 	database = "test"
// 	err = Init()
// 	if err != ErrFakeError {
// 		t.Error("Expected fake db error from init")
// 	}
// }
//
// func TestSet(t *testing.T) {
// 	database = "nodb"
// 	err := Set()
// 	if err == nil {
// 		t.Error("Expecting error for no database found")
// 	}
// 	Register("nodb2", TestDB)
// 	database = "nodb2"
// 	err = Set()
// 	if err != nil {
// 		t.Error(err)
// 	}
// }
//
// func TestRegister(t *testing.T) {
// 	l := len(DBTypes)
// 	Register("test2", TestDB)
// 	if len(DBTypes) != l+1 {
// 		t.Errorf("Expecting %v DB types, received %v", l+1, len(DBTypes))
// 	}
// 	l = len(DBTypes)
// 	Register("test2", TestDB)
// 	if len(DBTypes) != l {
// 		t.Errorf("Expecting %v DB types, received %v duplicate names", l, len(DBTypes))
// 	}
// }
//
// func TestCreateUser(t *testing.T) {
// 	err := CreateUser(&users.User{})
// 	if err != ErrFakeError {
// 		t.Error("expected fake db error from create")
// 	}
// }
//
// func TestGetUser(t *testing.T) {
// 	_, err := GetUser("test")
// 	if err != ErrFakeError {
// 		t.Error("expected fake db error from get")
// 	}
// }
//
// func TestGetUserByName(t *testing.T) {
// 	_, err := GetUserByName("test")
// 	if err != ErrFakeError {
// 		t.Error("expected fake db error from get")
// 	}
// }
//
// func TestDelete(t *testing.T) {
// 	err := Delete("test", "test_id")
// 	if err != ErrFakeError {
// 		t.Error("expected fake db error from delete")
// 	}
// }
//
// func TestPing(t *testing.T) {
// 	err := Ping()
// 	if err != ErrFakeError {
// 		t.Error("expected fake db error from ping")
// 	}
// }
//
// type fake struct{}
//
// func (f fake) Init() error {
// 	return ErrFakeError
// }
//
// func (f fake) GetUserByName(name string) (users.User, error) {
// 	return users.User{}, ErrFakeError
// }
//
// func (f fake) GetUser(id string) (users.User, error) {
// 	return users.User{}, ErrFakeError
// }
//
// func (f fake) GetUsers() ([]users.User, error) {
// 	return make([]users.User, 0), ErrFakeError
// }
//
// func (f fake) CreateUser(*users.User) error {
// 	return ErrFakeError
// }
//
// func (f fake) Delete(entity, id string) error {
// 	return ErrFakeError
// }
//
// func (f fake) Ping() error {
// 	return ErrFakeError
// }
