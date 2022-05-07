package server

import (
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
)

func TestRegistrationUserValidate(t *testing.T) {
	a := entity.RegistrationUser{
		Birthday: "20060102",
	}

	// date
	if err := binding.Validator.ValidateStruct(a); err != nil {
		errs := err.(validator.ValidationErrors)
		messages := ValidationErrors(errs, &a)
		if !strings.Contains(messages["birthday"], "invalid") {
			t.Errorf("Birthday is valid format")
		}
	}
}

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

	// password
	a.NewPassword = "hogefug"
	if err := binding.Validator.ValidateStruct(a); err != nil {
		errs := err.(validator.ValidationErrors)
		messages := ValidationErrors(errs, &a)
		if !strings.Contains(messages["newPassword"], "invalid") {
			t.Errorf("NewPassword is valid format")
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
	if err := binding.Validator.ValidateStruct(a); err != nil {
		errs := err.(validator.ValidationErrors)
		messages := ValidationErrors(errs, &a)
		if !strings.Contains(messages["account"], "min") {
			t.Errorf("Account is min length not specified")
		}
	}

	// max length
	a.Account = "testtesttesttesttestt"
	if err := binding.Validator.ValidateStruct(a); err != nil {
		errs := err.(validator.ValidationErrors)
		messages := ValidationErrors(errs, &a)
		if !strings.Contains(messages["account"], "max") {
			t.Errorf("Account is max length not specified")
		}
	}

	// password
	a.Password = "hogefug"
	if err := binding.Validator.ValidateStruct(a); err != nil {
		errs := err.(validator.ValidationErrors)
		messages := ValidationErrors(errs, &a)
		if !strings.Contains(messages["password"], "invalid") {
			t.Errorf("Password is valid format")
		}
	}
}
