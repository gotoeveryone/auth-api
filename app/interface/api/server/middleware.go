package server

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"github.com/gotoeveryone/auth-api/app/domain/repository"
	"github.com/gotoeveryone/auth-api/app/presentation/middleware"
	"github.com/sirupsen/logrus"
)

const (
	// Authentication HTTP header prefix
	tokenPrefix = "Bearer "
)

type tokenAuth struct {
	userRepo  repository.User
	tokenRepo repository.Token
}

// NewTokenAuth is create middleware use of token
func NewTokenAuth(ur repository.User, tr repository.Token) middleware.Auth {
	return &tokenAuth{
		userRepo:  ur,
		tokenRepo: tr,
	}
}

// Authorized is confirm has exactly token during proccessing request
func (m *tokenAuth) Authorized() gin.HandlerFunc {
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
		var t entity.Token
		if err := m.tokenRepo.Find(token, &t); err != nil {
			logrus.Error(err)
			errorInternalServerError(c, err)
			return
		}
		if t.Token == "" || t.UserID == 0 {
			c.Writer.Header().Set("WWW-Authenticate", "Bearer error=\"invalid_token\"")
			errorUnauthorized(c, errInvalidAccessToken)
		}
		var u entity.User
		if err := m.userRepo.Find(t.UserID, &u); err != nil {
			logrus.Error(err)
			errorInternalServerError(c, err)
			return
		}
		if !m.userRepo.ValidUser(&u) {
			c.Writer.Header().Set("WWW-Authenticate", "Bearer error=\"invalid_token\"")
			errorUnauthorized(c, errInvalidAccessToken)
			return
		}

		// Set token value
		c.Set(TokenKey, token)
		c.Next()
	}
}
