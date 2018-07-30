package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/auth-api/app/application/handler"
	"github.com/gotoeveryone/auth-api/app/infrastructure"
	"github.com/gotoeveryone/golib/logs"
)

const (
	// Authentication HTTP header prefix
	tokenPrefix = "Bearer "
)

// HasToken is confirm has exactly token during proccessing request
func HasToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Confirm has "Authorization" to HTTP header
		tokenHeader := c.Request.Header.Get("Authorization")
		if !strings.HasPrefix(tokenHeader, tokenPrefix) {
			c.Writer.Header().Set("WWW-Authenticate", "Bearer realm=\"token_required\"")
			handler.ErrorUnauthorized(c, handler.ErrRequiredAccessToken)
			return
		}

		// Confirm can get token value from HTTP header
		token := strings.TrimSpace(strings.Replace(tokenHeader, tokenPrefix, "", 1))
		if token == "" {
			c.Writer.Header().Set("WWW-Authenticate", "Bearer error=\"token_required\"")
			handler.ErrorUnauthorized(c, handler.ErrRequiredAccessToken)
			return
		}

		// Confirm has token valid
		ur := infrastructure.NewUserRepository()
		user, err := ur.FindByToken(token)
		if err != nil {
			logs.Error(err)
			handler.ErrorInternalServerError(c, err)
			return
		} else if !ur.ValidUser(user) {
			c.Writer.Header().Set("WWW-Authenticate", "Bearer error=\"invalid_token\"")
			handler.ErrorUnauthorized(c, handler.ErrInvalidAccessToken)
			return
		}

		// Set token value
		c.Set(handler.TokenKey, token)
		c.Next()
	}
}
