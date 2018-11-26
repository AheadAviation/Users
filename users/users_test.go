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

package users

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewUser(t *testing.T) {

	Convey("Given a new user", t, func() {
		u := New()

		Convey("When created", func() {

			c := len(u.Cards)
			Convey("Then cards should be initialized and empty", func() {
				So(c, ShouldEqual, 0)
			})

			a := len(u.Addresses)
			Convey("Then addresses should be initialized and empty", func() {
				So(a, ShouldEqual, 0)
			})
		})
	})

}

func TestNewUserNotValid(t *testing.T) {

	Convey("Given a new user", t, func() {
		u1 := User{
			FirstName: "",
			LastName:  "Test",
			Username:  "testuser",
			Password:  "testpass",
		}

		u2 := User{
			FirstName: "Test",
			LastName:  "User",
			Username:  "testuser",
			Password:  "testpass",
		}

		Convey("When validated", func() {
			err := u1.Validate()
			Convey("Then result should be Missing First Name error", func() {
				So(err, ShouldBeError, fmt.Errorf(ErrMissingField, "FirstName"))
			})

			err = u2.Validate()
			Convey("Then u2 result should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestCardsMasked(t *testing.T) {

	Convey("Given a new user", t, func() {
		u := New()

		Convey("When masking new cards", func() {
			u.Cards = append(u.Cards, Card{LongNum: "abcdefg"})
			u.MaskCCs()
			Convey("Then LongNum should be ***defg", func() {
				So(u.Cards[0].LongNum, ShouldEqual, "***defg")
			})

			u.Cards = append(u.Cards, Card{LongNum: "hijklmnopqrs"})
			u.MaskCCs()
			Convey("Then LongNum should be ********pqrs", func() {
				So(u.Cards[1].LongNum, ShouldEqual, "********pqrs")
			})

		})
	})
}
