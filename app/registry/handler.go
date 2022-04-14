package registry

import (
	"github.com/gotoeveryone/auth-api/app/domain/repository"
	"github.com/gotoeveryone/auth-api/app/interface/api/server"
	"github.com/gotoeveryone/auth-api/app/presentation/handler"
)

// NewStateHandler is create action handler for state
func NewStateHandler() handler.State {
	return server.NewStateHandler()
}

// NewAuthHandler is create action handler for auth
func NewAuthHandler(ur repository.User, tr repository.Token) handler.Authenticate {
	return server.NewAuthHandler(ur, tr)
}
