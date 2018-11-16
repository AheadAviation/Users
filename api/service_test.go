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
	"testing"

	"github.com/aheadaviation/Users/users"
)

var (
	TestService  Service
	TestCustomer = users.User{
		Username: "testuser",
		Password: "",
	}
)

func init() {
	TestService = NewFixedService()
}

func TestCalculatePassHash(t *testing.T) {
	hash1 := calculatePassHash("eve", "c748112bc027878aa62812ba1ae00e40ad46d497")
	if hash1 != "c748112bc027878aa62812ba1ae00e40ad46d497" {
		t.Error("Eve's password failed hash test")
	}
}
