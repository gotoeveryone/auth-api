package registry

import (
	"github.com/gotoeveryone/auth-api/app/domain/repository"
	"github.com/gotoeveryone/auth-api/app/infrastructure/database"
)

// NewUserRepository is create user management repository.
func NewUserRepository() repository.User {
	return database.NewUserRepository()
}
