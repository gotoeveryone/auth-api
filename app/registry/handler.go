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

// NewUserHandler is create action handler for user
func NewUserHandler(r repository.User) handler.User {
	return server.NewUserHandler(r)
}
