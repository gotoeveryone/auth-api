package registry

import (
	"github.com/gotoeveryone/auth-api/app/domain/repository"
	"github.com/gotoeveryone/auth-api/app/interfaces/api/server"
	"github.com/gotoeveryone/auth-api/app/presentation/middleware"
)

// NewAuthenticateMiddleware is create authenticate middleware
func NewAuthenticateMiddleware(ur repository.UserRepository, tr repository.TokenRepository) middleware.Authenticate {
	return server.NewTokenAuthenticate(ur, tr)
}
