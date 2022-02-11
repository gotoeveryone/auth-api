package database

import (
	"errors"
	"time"

	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"github.com/gotoeveryone/auth-api/app/domain/repository"
	"github.com/gotoeveryone/auth-api/app/infrastructure"
	"gorm.io/gorm"
)

type tokenRepository struct {
	infrastructure.BaseTokenRepository
}

// NewTokenRepository is create token management repository using Database.
func NewTokenRepository() repository.TokenRepository {
	return &tokenRepository{}
}

// Find is execute token data finding
func (r tokenRepository) Find(token string, t *entity.Token) error {
	err := dbManager.Where(&entity.Token{Token: token}).
		Where("expired_at >= ?", time.Now()).First(t).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return err
}

// Create is execute token data creating
func (r tokenRepository) Create(u *entity.User, t *entity.Token) error {
	r.CreateFromUser(u, t)
	return dbManager.Create(t).Error
}

// Delete is execute token data deleting
func (r tokenRepository) Delete(token string) error {
	return dbManager.Where(&entity.Token{Token: token}).
		Delete(entity.Token{}).Error
}

// DeleteExpired is execute expired token data deleting
func (r tokenRepository) DeleteExpired() (int64, error) {
	cnt := dbManager.Where("expired_at < ?", time.Now()).
		Delete(entity.Token{}).RowsAffected
	return cnt, dbManager.Error
}

// CanAutoDeleteExpired is returns whether token data can be automatically deleted.
func (r tokenRepository) CanAutoDeleteExpired() bool {
	return false
}
