package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/general-api/app/models"
	"github.com/gotoeveryone/golib/logs"
)

const (
	// TokenKey トークン取得用のキー
	TokenKey = "authenticated-token"
)

// 400エラー
func errorBadRequest(c *gin.Context, message string) {
	errorJSON(c, models.Error{
		Code:    http.StatusBadRequest,
		Message: message,
		Error:   nil,
	})
}

// 401エラー
func errorUnauthorized(c *gin.Context, message string) {
	errorJSON(c, models.Error{
		Code:    http.StatusUnauthorized,
		Message: message,
		Error:   nil,
	})
}

// 500エラー
func errorInternalServerError(c *gin.Context, err error) {
	logs.Error(fmt.Errorf("Error: %s", err))
	errorJSON(c, models.Error{
		Code:    http.StatusInternalServerError,
		Message: "",
		Error:   err,
	})
}

// エラー用JSONの出力
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
