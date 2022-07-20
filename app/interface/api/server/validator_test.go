package server

import (
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestRegistrationUserValidate(t *testing.T) {
	a := entity.RegistrationUser{
		Birthday: "20060102",
	}

	// date
	if err := binding.Validator.ValidateStruct(a); err != nil {
		errs := err.(validator.ValidationErrors)
		messages := ValidationErrors(errs, &a)
		assert.Contains(t, messages["birthday"], "invalid")
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
		assert.Contains(t, messages["newPassword"], "required")
	}

	// password
	a.NewPassword = "hogefug"
	if err := binding.Validator.ValidateStruct(a); err != nil {
		errs := err.(validator.ValidationErrors)
		messages := ValidationErrors(errs, &a)
		assert.Contains(t, messages["newPassword"], "invalid")
	}
}

func TestAuthenticateValidate(t *testing.T) {
	a := entity.Authenticate{}

	// required
	if err := binding.Validator.ValidateStruct(a); err != nil {
		errs := err.(validator.ValidationErrors)
		messages := ValidationErrors(errs, &a)
		assert.Contains(t, messages["account"], "required")
		assert.Contains(t, messages["password"], "required")
	}

	// min length
	a.Account = "testt"
	if err := binding.Validator.ValidateStruct(a); err != nil {
		errs := err.(validator.ValidationErrors)
		messages := ValidationErrors(errs, &a)
		assert.Contains(t, messages["account"], "min")
	}

	// max length
	a.Account = "testtesttesttesttestt"
	if err := binding.Validator.ValidateStruct(a); err != nil {
		errs := err.(validator.ValidationErrors)
		messages := ValidationErrors(errs, &a)
		assert.Contains(t, messages["account"], "max")
	}

	// password
	a.Password = "hogefug"
	if err := binding.Validator.ValidateStruct(a); err != nil {
		errs := err.(validator.ValidationErrors)
		messages := ValidationErrors(errs, &a)
		assert.Contains(t, messages["password"], "invalid")
	}
}
