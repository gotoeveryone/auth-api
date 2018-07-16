package handler

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/general-api/app/domain/entity"
	"github.com/gotoeveryone/golib/logs"
)

const (
	// TokenKey is authenticated token key
	TokenKey = "authenticated-token"
)

var (
	// ErrUnauthorized is unauthorized message.
	ErrUnauthorized = errors.New("Authorization failed")
	// ErrRequiredAccessToken is access token is required message.
	ErrRequiredAccessToken = errors.New("Token is required")
	// ErrInvalidAccessToken is access token is invalid message.
	ErrInvalidAccessToken = errors.New("Token is invalid")
)

// ErrorBadRequest is return bad request response.
func ErrorBadRequest(c *gin.Context, message interface{}) {
	errorJSON(c, entity.Error{
		Code:    http.StatusBadRequest,
		Message: message,
		Error:   nil,
	})
}

// ErrorUnauthorized is return unauthorized response.
func ErrorUnauthorized(c *gin.Context, message interface{}) {
	errorJSON(c, entity.Error{
		Code:    http.StatusUnauthorized,
		Message: message,
		Error:   nil,
	})
}

// ErrorInternalServerError is return internal server error response.
func ErrorInternalServerError(c *gin.Context, err error) {
	logs.Error(fmt.Errorf("Error: %s", err))
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
