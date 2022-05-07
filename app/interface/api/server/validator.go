package server

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	dateRegex     = regexp.MustCompile("^\\d{4}-\\d{2}-\\d{2}$")
	passwordRegex = regexp.MustCompile("^[A-Za-z0-9]{8,}$")
)

func date(fl validator.FieldLevel) bool {
	return dateRegex.MatchString(fl.Field().String())
}

func password(fl validator.FieldLevel) bool {
	return passwordRegex.MatchString(fl.Field().String())
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("date", date)
		v.RegisterValidation("password", password)
	}
}

// ValidationErrors is create validation error message.
func ValidationErrors(ve validator.ValidationErrors, o interface{}) map[string]string {
	res := map[string]string{}

	for _, err := range ve {
		field, _ := reflect.TypeOf(o).Elem().FieldByName(err.Field())
		key := field.Tag.Get("json")
		if err.Param() != "" {
			res[key] = fmt.Sprintf("Value is %s %s", err.Tag(), err.Param())
		} else if err.Tag() == "required" {
			res[key] = fmt.Sprint("Value is required")
		} else {
			res[key] = fmt.Sprint("Value is invalid")
		}
	}

	return res
}
