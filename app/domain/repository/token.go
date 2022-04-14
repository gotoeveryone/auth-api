package repository

import (
	"github.com/gotoeveryone/auth-api/app/domain/entity"
)

// Token is repository for operate about token.
type Token interface {
	Find(token string, t *entity.Token) error
	Create(u *entity.User, t *entity.Token) error
	Delete(token string) error
	DeleteExpired() (int64, error)
	CanAutoDeleteExpired() bool
}
