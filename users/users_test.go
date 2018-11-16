package users

import (
	"fmt"
	"testing"
)

// func TestNew(t *testing.T) {
//   u := New()
// }

func TestValidate(t *testing.T) {
	u := New()
	err := u.Validate()
	if err.Error() != fmt.Sprintf(ErrMissingField, "FirstName") {
		t.Error("Expected missing first name error")
	}
	u.FirstName = "test"
	err = u.Validate()
	if err.Error() != fmt.Sprintf(ErrMissingField, "LastName") {
		t.Error("Expected missing last name error")
	}
	u.LastName = "test"
	err = u.Validate()
	if err.Error() != fmt.Sprintf(ErrMissingField, "Username") {
		t.Error("Expected missing username error")
	}
	u.Username = "test"
	err = u.Validate()
	if err.Error() != fmt.Sprintf(ErrMissingField, "Password") {
		t.Error("Expected missing password error")
	}
	u.Password = "test"
	err = u.Validate()
	if err != nil {
		t.Error(err)
	}
}
