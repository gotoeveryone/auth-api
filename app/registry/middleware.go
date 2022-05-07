package registry

import (
	"github.com/gotoeveryone/auth-api/app/domain/repository"
	"github.com/gotoeveryone/auth-api/app/interface/api/server"
	"github.com/gotoeveryone/auth-api/app/presentation/middleware"
)

// NewAuthMiddleware is create middleware about auth
func NewAuthMiddleware(ur repository.User) middleware.Auth {
	return server.NewAuthMiddleware(ur)
}
