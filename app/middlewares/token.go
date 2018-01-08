package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/general-api/app/handlers"
	"github.com/gotoeveryone/general-api/app/models"
	"github.com/gotoeveryone/general-api/app/services"
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
			c.AbortWithStatusJSON(401, models.Error{
				Code:    401,
				Message: "Token is required",
			})
			return
		}

		// Confirm can get token value from HTTP header
		token := strings.TrimSpace(strings.Replace(tokenHeader, tokenPrefix, "", 1))
		if token == "" {
			c.Writer.Header().Set("WWW-Authenticate", "Bearer error=\"token_required\"")
			c.AbortWithStatusJSON(401, models.Error{
				Code:    401,
				Message: "Token is required",
			})
			return
		}

		// Confirm has token valid
		var ts services.TokensService
		var m models.Token
		if err := ts.FindToken(token, &m); err != nil {
			logs.Error(err)
			c.Writer.Header().Set("WWW-Authenticate", "Bearer error=\"invalid_token\"")
			c.AbortWithStatusJSON(401, models.Error{
				Code:    401,
				Message: "Token is invalid",
			})
			return
		}

		// Set token value
		c.Set(handlers.TokenKey, m.Token)
		c.Next()
	}
}
