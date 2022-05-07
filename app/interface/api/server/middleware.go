package server

import (
	"net/http"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/auth-api/app/config"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"github.com/gotoeveryone/auth-api/app/domain/repository"
	"github.com/gotoeveryone/auth-api/app/presentation/middleware"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	timeout time.Duration = time.Hour * 2
)

type jwtAuth struct {
	repo repository.User
}

// NewAuthMiddleware is create middleware for auth
func NewAuthMiddleware(ur repository.User) middleware.Auth {
	return &jwtAuth{
		repo: ur,
	}
}

// @Summary Execute authentication for user
// @Tags Authenticate
// @Produce json
// @Param data body entity.Authenticate true "request data"
// @Success 200 {object} entity.Claim
// @Failure 404 {object} entity.Error
// @Failure 405 {object} entity.Error
// @Router /v1/auth [post]
func loginResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(code, entity.Claim{
		Token:  token,
		Expire: expire.Format(time.RFC3339),
	})
}

// @Summary Publish refresh token for user
// @Tags Authenticate
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} entity.Claim
// @Failure 404 {object} entity.Error
// @Failure 405 {object} entity.Error
// @Router /v1/refresh_token [get]
func refreshResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(code, entity.Claim{
		Token:  token,
		Expire: expire.Format(time.RFC3339),
	})
}

// @Summary Execute deauthentication for user
// @Tags Authenticate
// @Security ApiKeyAuth
// @Produce json
// @Success 204
// @Failure 404 {object} entity.Error
// @Failure 405 {object} entity.Error
// @Router /v1/deauth [delete]
func logoutResponse(c *gin.Context, code int) {
	c.JSON(http.StatusNoContent, gin.H{})
}

// Create is create auth middleware
func (m jwtAuth) Create() (*jwt.GinJWTMiddleware, error) {
	identityKey := config.IdentityKey
	middleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "auth-api",
		Key:         []byte(os.Getenv("SECRET_KEY")),
		Timeout:     timeout,
		MaxRefresh:  timeout,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*entity.User); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			key, ok := claims[identityKey]
			if !ok {
				return nil
			}
			var user entity.User
			if err := m.repo.Find(uint(key.(float64)), &user); err != nil {
				logrus.Error(err)
				return nil
			}
			if &user == nil {
				return nil
			}
			return &user
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var p entity.Authenticate
			if err := c.ShouldBind(&p); err != nil {
				return nil, errUnauthorized
			}

			user, err := m.repo.FindByAccount(p.Account)
			if err != nil {
				logrus.Error(err)
				return nil, errUnauthorized
			}

			if user == nil {
				return nil, errUnauthorized
			}

			if !user.Valid() {
				return nil, errInvalidAccount
			}

			if !user.IsActive {
				return nil, errMustChangePassword
			}

			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p.Password)); err != nil {
				logrus.Error(err)
				return nil, errUnauthorized
			}

			return user, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*entity.User); ok {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, entity.Error{
				Code:    code,
				Message: message,
			})
		},

		TokenLookup: "header: Authorization",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,

		// Response
		LoginResponse:   loginResponse,
		LogoutResponse:  logoutResponse,
		RefreshResponse: refreshResponse,
	})

	if err != nil {
		return nil, err
	}

	if err := middleware.MiddlewareInit(); err != nil {
		return nil, err
	}

	return middleware, nil
}
