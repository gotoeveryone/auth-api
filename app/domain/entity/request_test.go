package entity

import (
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
	validator "gopkg.in/go-playground/validator.v8"
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

func TestValidateStruct(t *testing.T) {
	a := Authenticate{}

	// required
	if err := binding.Validator.ValidateStruct(a); err != nil {
		errs := err.(validator.ValidationErrors)
		if !strings.Contains(errs["Authenticate.Account"].Tag, "required") {
			t.Errorf("Account is required not specified")
		}
		if !strings.Contains(errs["Authenticate.Password"].Tag, "required") {
			t.Errorf("Password is required not specified")
		}
	}

	// min length
	a.Account = "testt"
	a.Password = "testtes"
	if err := binding.Validator.ValidateStruct(a); err != nil {
		errs := err.(validator.ValidationErrors)
		if !strings.Contains(errs["Authenticate.Account"].Tag, "min") {
			t.Errorf("Account is min length not specified")
		}
		if !strings.Contains(errs["Authenticate.Password"].Tag, "min") {
			t.Errorf("Password is min length not specified")
		}
	}

	// max length
	a.Account = "testtesttes"
	if err := binding.Validator.ValidateStruct(a); err != nil {
		errs := err.(validator.ValidationErrors)
		if !strings.Contains(errs["Authenticate.Account"].Tag, "max") {
			t.Errorf("Account is max length not specified")
		}
	}

	a.Account = "testtest"
	a.Password = "testtest"
	ac := Activate{
		Authenticate: a,
	}
	// required
	if err := binding.Validator.ValidateStruct(ac); err != nil {
		errs := err.(validator.ValidationErrors)
		if !strings.Contains(errs["Activate.NewPassword"].Tag, "required") {
			t.Errorf("NewPassword is required not specified")
		}
	}

	// min length
	ac.NewPassword = "test"
	if err := binding.Validator.ValidateStruct(ac); err != nil {
		errs := err.(validator.ValidationErrors)
		if !strings.Contains(errs["Activate.NewPassword"].Tag, "min") {
			t.Errorf("NewPassword is min length not specified")
		}
	}
}
