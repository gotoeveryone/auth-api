package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/general-api/app/models"
	"github.com/gotoeveryone/golib/logs"
)

const (
	// TokenKey is authenticated token key
	TokenKey = "authenticated-token"
)

// Bad request
func errorBadRequest(c *gin.Context, message string) {
	errorJSON(c, models.Error{
		Code:    http.StatusBadRequest,
		Message: message,
		Error:   nil,
	})
}

// Unauthorized
func errorUnauthorized(c *gin.Context, message string) {
	errorJSON(c, models.Error{
		Code:    http.StatusUnauthorized,
		Message: message,
		Error:   nil,
	})
}

// Internal server error
func errorInternalServerError(c *gin.Context, err error) {
	logs.Error(fmt.Errorf("Error: %s", err))
	errorJSON(c, models.Error{
		Code:    http.StatusInternalServerError,
		Message: "",
		Error:   err,
	})
}

// Outputting error with JSON format
func errorJSON(c *gin.Context, err models.Error) {
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
	}
	c.AbortWithStatusJSON(err.Code, err)
}
