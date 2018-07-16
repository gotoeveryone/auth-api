package repository

import (
	"github.com/gotoeveryone/general-api/app/domain/entity"
)

// TokenRepository is operate of token data.
type TokenRepository interface {
	FindToken(token string, t *entity.Token) error
	Create(u *entity.User, t *entity.Token) error
	Delete(token string) error
	DeleteExpired() (int64, error)
	CanAutoDeleteExpired() bool
}
