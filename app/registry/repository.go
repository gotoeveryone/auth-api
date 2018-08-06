package registry

import (
	"github.com/gotoeveryone/auth-api/app/config"
	"github.com/gotoeveryone/auth-api/app/domain/repository"
	"github.com/gotoeveryone/auth-api/app/infrastructure/database"
	"github.com/gotoeveryone/auth-api/app/infrastructure/redis"
)

// NewUserRepository is create user management repository.
func NewUserRepository() repository.UserRepository {
	return database.NewUserRepository()
}

// NewTokenRepository is create token management repository.
func NewTokenRepository() repository.TokenRepository {
	if config.AppConfig.Cache.Use {
		c := redis.Client{
			Config: config.AppConfig.Cache,
		}
		return redis.NewTokenRepository(c)
	}
	return database.NewTokenRepository()
}
