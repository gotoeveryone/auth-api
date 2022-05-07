package entity

import (
	"testing"
)

func TestActivate(t *testing.T) {
	a := Activate{
		Authenticate: Authenticate{
			Account:  "testtest",
			Password: "testtest",
		},
	}

	if a.NewPassword != "" {
		t.Errorf("NewPassword is not default value")
	}
}

func TestAuthenticate(t *testing.T) {
	a := Authenticate{}

	if a.Account != "" {
		t.Errorf("Account is not default value")
	}
	if a.Password != "" {
		t.Errorf("Password is not default value")
	}
}
