package server

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"github.com/sirupsen/logrus"
)

const (
	// TokenKey is authenticated token key
	TokenKey = "authenticated-token"
)

var (
	errUnauthorized        = errors.New("Authorization failed")
	errRequiredAccessToken = errors.New("Token is required")
	errInvalidAccessToken  = errors.New("Token is invalid")

	errInvalidAccount = errors.New("Account is invalid")
	errExistsAccount  = errors.New("Account is already exists")
	errInvalidRole    = errors.New("Role is invalid")

	errUpdatePassword     = errors.New("Update password failed")
	errSamePassword       = errors.New("Not allowed changing to same password")
	errMustChangePassword = errors.New("Password must be changed")

	errValidationFailed = errors.New("Validation failed")
)

// Return bad request response.
func errorBadRequest(c *gin.Context, message interface{}) {
	errorJSON(c, entity.Error{
		Code:    http.StatusBadRequest,
		Message: message,
		Error:   nil,
	})
}

// Return unauthorized response.
func errorUnauthorized(c *gin.Context, message interface{}) {
	errorJSON(c, entity.Error{
		Code:    http.StatusUnauthorized,
		Message: message,
		Error:   nil,
	})
}

// Return internal server error response.
func errorInternalServerError(c *gin.Context, err error) {
	logrus.Error(fmt.Errorf("error: %s", err))
	errorJSON(c, entity.Error{
		Code:    http.StatusInternalServerError,
		Message: "",
		Error:   err,
	})
}

// Outputting error with JSON format
func errorJSON(c *gin.Context, err entity.Error) {
	var header string
	switch err.Code {
	case http.StatusBadRequest:
		header = "error=\"invalid_request\""
	case http.StatusUnauthorized:
		header = "error=\"invalid_token\""
	}
	if header != "" && c.GetString(TokenKey) != "" {
		c.Writer.Header().Set("WWW-Authenticate", "Bearer "+header)
	}
	if err.Message == "" {
		err.Message = http.StatusText(err.Code)
	} else {
		et := reflect.TypeOf(err.Message)
		if et.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			err.Message = err.Message.(error).Error()
		}
	}
	c.AbortWithStatusJSON(err.Code, err)
}
