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
	"os"
	"testing"

	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/dbtest"

	"github.com/aheadaviation/Users/users"
)

var (
	TestMongo  = Mongo{}
	TestServer = dbtest.DBServer{}
	TestUser   = users.User{
		FirstName: "firstname",
		LastName:  "lastname",
		Username:  "username",
		Password:  "mypword",
	}
)

func init() {
	TestServer.SetPath("/tmp")
}

func TestMain(m *testing.M) {
	TestMongo.Session = TestServer.Session()
	TestMongo.EnsureIndexes()
	TestMongo.Session.Close()
	exitTest(m.Run())
}

func exitTest(i int) {
	TestServer.Wipe()
	TestServer.Stop()
	os.Exit(i)
}

func TestInit(t *testing.T) {
	err := TestMongo.Init()
	if err.Error() != "no reachable servers" {
		t.Error("expected no reachable servers error")
	}
}

// func TestNew(t *testing.T) {
//   m := New()
// }

func TestAddUserIDs(t *testing.T) {
	m := New()
	uid := bson.NewObjectId()
	m.ID = uid
	m.AddUserIds()
	if m.UserID != uid.Hex() {
		t.Error("Expected matching User Hex")
	}
}

func TestCreate(t *testing.T) {
	TestMongo.Session = TestServer.Session()
	defer TestMongo.Session.Close()
	err := TestMongo.CreateUser(&TestUser)
	if err != nil {
		t.Error(err)
	}
	err = TestMongo.CreateUser(&TestUser)
	if err == nil {
		t.Error("Expected duplicate key error")
	}
}

func TestGetUserByName(t *testing.T) {
	TestMongo.Session = TestServer.Session()
	defer TestMongo.Session.Close()
	u, err := TestMongo.GetUserByName(TestUser.Username)
	if err != nil {
		t.Error(err)
	}
	if u.Username != TestUser.Username {
		t.Error("Expected equal usernames")
	}
	_, err = TestMongo.GetUserByName("bogususers")
	if err == nil {
		t.Error("expected not found error")
	}
}

func TestGetUser(t *testing.T) {
	TestMongo.Session = TestServer.Session()
	defer TestMongo.Session.Close()
	_, err := TestMongo.GetUser(TestUser.UserID)
	if err != nil {
		t.Error(err)
	}
}

func TestGetURL(t *testing.T) {
	name = "test"
	password = "password"
	host = "shouldnotexist:3038"
	u := getURL()
	if u.String() != "mongodb://test:password@shouldnotexist:3038/users" {
		t.Error("expected URL mismatch")
	}
}

func TestPing(t *testing.T) {
	TestMongo.Session = TestServer.Session()
	defer TestMongo.Session.Close()
	err := TestMongo.Ping()
	if err != nil {
		t.Error(err)
	}
}
