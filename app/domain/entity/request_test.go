package entity

import (
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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
		for _, fe := range err.(validator.ValidationErrors) {
			if fe.Field() == "Authenticate.Account" && !strings.Contains(fe.Tag(), "required") {
				t.Errorf("Account is required not specified")
			}
			if fe.Field() == "Authenticate.Password" && !strings.Contains(fe.Tag(), "required") {
				t.Errorf("Password is required not specified")
			}
		}
	}

	// min length
	a.Account = "testt"
	a.Password = "testtes"
	if err := binding.Validator.ValidateStruct(a); err != nil {
		for _, fe := range err.(validator.ValidationErrors) {
			if fe.Field() == "Authenticate.Account" && !strings.Contains(fe.Tag(), "min") {
				t.Errorf("Account is min length not specified")
			}
			if fe.Field() == "Authenticate.Password" && !strings.Contains(fe.Tag(), "min") {
				t.Errorf("Password is min length not specified")
			}
		}
	}

	// max length
	a.Account = "testtesttes"
	if err := binding.Validator.ValidateStruct(a); err != nil {
		for _, fe := range err.(validator.ValidationErrors) {
			if fe.Field() == "Authenticate.Account" && !strings.Contains(fe.Tag(), "max") {
				t.Errorf("Account is max length not specified")
			}
		}
	}

	a.Account = "testtest"
	a.Password = "testtest"
	ac := Activate{
		Authenticate: a,
	}
	// required
	if err := binding.Validator.ValidateStruct(ac); err != nil {
		for _, fe := range err.(validator.ValidationErrors) {
			if fe.Field() == "Authenticate.NewPassword" && !strings.Contains(fe.Tag(), "min") {
				t.Errorf("NewPassword is required not specified")
			}
		}
	}

	// min length
	ac.NewPassword = "test"
	if err := binding.Validator.ValidateStruct(ac); err != nil {
		for _, fe := range err.(validator.ValidationErrors) {
			if fe.Field() == "Authenticate.NewPassword" && !strings.Contains(fe.Tag(), "min") {
				t.Errorf("NewPassword is min length not specified")
			}
		}
	}
}
