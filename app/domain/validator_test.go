package domain

import (
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
)

func TestActivateValidate(t *testing.T) {
	a := entity.Activate{
		Authenticate: entity.Authenticate{
			Account:  "testtest",
			Password: "testtest",
		},
	}

	// required
	if err := binding.Validator.ValidateStruct(a); err != nil {
		errs := err.(validator.ValidationErrors)
		messages := ValidationErrors(errs, &a)
		if !strings.Contains(messages["newPassword"], "required") {
			t.Errorf("NewPassword is required not specified")
		}
	}

	// min length
	a.NewPassword = "test"
	if err := binding.Validator.ValidateStruct(a); err != nil {
		errs := err.(validator.ValidationErrors)
		messages := ValidationErrors(errs, &a)
		if !strings.Contains(messages["newPassword"], "min") {
			t.Errorf("NewPassword is min length not specified")
		}
	}
}

func TestAuthenticateValidate(t *testing.T) {
	a := entity.Authenticate{}

	// required
	if err := binding.Validator.ValidateStruct(a); err != nil {
		errs := err.(validator.ValidationErrors)
		messages := ValidationErrors(errs, &a)
		if !strings.Contains(messages["account"], "required") {
			t.Errorf("Account is required not specified")
		}
		if !strings.Contains(messages["password"], "required") {
			t.Errorf("Password is required not specified")
		}
	}

	// min length
	a.Account = "testt"
	a.Password = "testtes"
	if err := binding.Validator.ValidateStruct(a); err != nil {
		errs := err.(validator.ValidationErrors)
		messages := ValidationErrors(errs, &a)
		if !strings.Contains(messages["account"], "min") {
			t.Errorf("Account is min length not specified")
		}
		if !strings.Contains(messages["password"], "min") {
			t.Errorf("Password is min length not specified")
		}
	}

	// max length
	a.Account = "testtesttes"
	if err := binding.Validator.ValidateStruct(a); err != nil {
		errs := err.(validator.ValidationErrors)
		messages := ValidationErrors(errs, &a)
		if !strings.Contains(messages["account"], "max") {
			t.Errorf("Account is max length not specified")
		}
	}
}
