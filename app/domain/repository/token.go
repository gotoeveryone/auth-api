package repository

import (
	"github.com/gotoeveryone/auth-api/app/domain/entity"
)

// TokenRepository is operate of token data.
type TokenRepository interface {
	Find(token string, t *entity.Token) error
	Create(u *entity.User, t *entity.Token) error
	Delete(token string) error
	DeleteExpired() (int64, error)
	CanAutoDeleteExpired() bool
}
