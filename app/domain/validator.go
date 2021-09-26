package domain

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

// ValidationErrors is create validation error message.
func ValidationErrors(ve validator.ValidationErrors, o interface{}) map[string]string {
	res := map[string]string{}

	for _, err := range ve {
		field, _ := reflect.TypeOf(o).Elem().FieldByName(err.Field())
		key := field.Tag.Get("json")
		if err.Param() != "" {
			res[key] = fmt.Sprintf("Value is %s %s", err.Tag(), err.Param())
		} else {
			res[key] = fmt.Sprintf("Value is %s", err.Tag())
		}
	}

	return res
}
