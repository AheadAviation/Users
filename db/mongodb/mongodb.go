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

package mongodb

import (
	"errors"
	"flag"
	"net/url"
	"os"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/aheadaviation/Users/users"
)

var (
	name            string
	password        string
	host            string
	db              = "users"
	ErrInvalidHexID = errors.New("Invalid Id Hex")
)

func init() {
	flag.StringVar(&name, "mongo-user", os.Getenv("MONGO_USER"), "Mongo Username")
	flag.StringVar(&password, "mongo-password", os.Getenv("MONGO_PASS"), "Mongo Password")
	flag.StringVar(&host, "mongo-host", os.Getenv("MONGO_HOST"), "Mongo Host")
}

type Mongo struct {
	Session *mgo.Session
}

func (m *Mongo) Init() error {
	u := getURL()
	var err error
	m.Session, err = mgo.DialWithTimeout(u.String(), time.Duration(5)*time.Second)
	if err != nil {
		return nil
	}
	return m.EnsureIndexes()
}

type MongoUser struct {
	users.User `bson:",inline"`
	ID         bson.ObjectId `bson:"_id"`
}

func New() MongoUser {
	u := users.New()
	return MongoUser{
		User: u,
	}
}

func (mu *MongoUser) AddUserIds() {
	mu.User.UserID = mu.ID.Hex()
}

func (m *Mongo) CreateUser(u *users.User) error {
	s := m.Session.Copy()
	defer s.Close()
	id := bson.NewObjectId()
	mu := New()
	mu.User = *u
	mu.ID = id
	c := s.DB("").C("customers")
	_, err := c.UpsertId(mu.ID, mu)
	if err != nil {
		//m.cleanAttributes(mu)
		return err
	}
	mu.User.UserID = mu.ID.Hex()
	*u = mu.User
	return nil
}

func (m *Mongo) GetUserByName(name string) (users.User, error) {
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C("customers")
	mu := New()
	err := c.Find(bson.M{"username": name}).One(&mu)
	mu.AddUserIds()
	return mu.User, err
}

func (m *Mongo) GetUser(id string) (users.User, error) {
	s := m.Session.Copy()
	defer s.Close()
	if !bson.IsObjectIdHex(id) {
		return users.New(), errors.New("Invalid id hex")
	}
	c := s.DB("").C("customers")
	mu := New()
	err := c.FindId(bson.ObjectIdHex(id)).One(&mu)
	mu.AddUserIds()
	return mu.User, err
}

func (m *Mongo) GetUsers() ([]users.User, error) {
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C("customers")
	var mus []MongoUser
	err := c.Find(nil).All(&mus)
	us := make([]users.User, 0)
	for _, mu := range mus {
		mu.AddUserIds()
		us = append(us, mu.User)
	}
	return us, err
}

func (m *Mongo) Delete(entity, id string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("invalid id hex")
	}
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C(entity)
	return c.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
}

func (m *Mongo) EnsureIndexes() error {
	s := m.Session.Copy()
	defer s.Close()
	i := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     false,
	}
	c := s.DB("").C("customers")
	return c.EnsureIndex(i)
}

func (m *Mongo) Ping() error {
	s := m.Session.Copy()
	defer s.Close()
	return s.Ping()
}

// func (m *Mongo) cleanAttributes(mu MongoUser) error {
// 	s := m.Session.Copy()
// 	defer s.Close()
// 	return err
// }

func getURL() url.URL {
	ur := url.URL{
		Scheme: "mongodb",
		Host:   host,
		Path:   db,
	}
	if name != "" {
		u := url.UserPassword(name, password)
		ur.User = u
	}
	return ur
}
