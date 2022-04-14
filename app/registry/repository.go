package registry

import (
	"github.com/gotoeveryone/auth-api/app/config"
	"github.com/gotoeveryone/auth-api/app/domain/repository"
	"github.com/gotoeveryone/auth-api/app/infrastructure/database"
	"github.com/gotoeveryone/auth-api/app/infrastructure/redis"
)

// NewUserRepository is create user management repository.
func NewUserRepository() repository.User {
	return database.NewUserRepository()
}

// NewTokenRepository is create token management repository.
func NewTokenRepository(c config.App) repository.Token {
	if c.Cache.Use {
		c := redis.Client{
			Config: c.Cache,
		}
		return redis.NewTokenRepository(c)
	}
	return database.NewTokenRepository()
}
