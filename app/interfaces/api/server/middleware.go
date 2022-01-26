package server

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/auth-api/app/domain/repository"
	"github.com/gotoeveryone/auth-api/app/presentation/middleware"
	"github.com/sirupsen/logrus"
)

const (
	// Authentication HTTP header prefix
	tokenPrefix = "Bearer "
)

type tokenAuthenticate struct {
	repo repository.UserRepository
}

// NewTokenAuthenticate is create authenticate token middleware
func NewTokenAuthenticate(r repository.UserRepository) middleware.Authenticate {
	return &tokenAuthenticate{
		repo: r,
	}
}

// Authorized is confirm has exactly token during proccessing request
func (m *tokenAuthenticate) Authorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Confirm has "Authorization" to HTTP header
		tokenHeader := c.Request.Header.Get("Authorization")
		if !strings.HasPrefix(tokenHeader, tokenPrefix) {
			c.Writer.Header().Set("WWW-Authenticate", "Bearer realm=\"token_required\"")
			errorUnauthorized(c, errRequiredAccessToken)
			return
		}

		// Confirm can get token value from HTTP header
		token := strings.TrimSpace(strings.Replace(tokenHeader, tokenPrefix, "", 1))
		if token == "" {
			c.Writer.Header().Set("WWW-Authenticate", "Bearer error=\"token_required\"")
			errorUnauthorized(c, errRequiredAccessToken)
			return
		}

		// Confirm has token valid
		user, err := m.repo.FindByToken(token)
		if err != nil {
			logrus.Error(err)
			errorInternalServerError(c, err)
			return
		} else if !m.repo.ValidUser(user) {
			c.Writer.Header().Set("WWW-Authenticate", "Bearer error=\"invalid_token\"")
			errorUnauthorized(c, errInvalidAccessToken)
			return
		}

		// Set token value
		c.Set(TokenKey, token)
		c.Next()
	}
}
