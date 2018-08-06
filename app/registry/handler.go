package registry

import (
	"github.com/gotoeveryone/auth-api/app/domain/repository"
	"github.com/gotoeveryone/auth-api/app/interfaces/api/server"
	"github.com/gotoeveryone/auth-api/app/presentation/handler"
)

// NewStateHandler is create state action handler
func NewStateHandler() handler.StateHandler {
	return server.NewStateHandler()
}

// NewAuthenticateHandler is create authenticate action handler
func NewAuthenticateHandler(ur repository.UserRepository, tr repository.TokenRepository) handler.AuthenticateHandler {
	return server.NewAuthenticateHandler(ur, tr)
}
