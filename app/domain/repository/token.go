package repository

import "github.com/gotoeveryone/general-api/app/domain/entity"

type (
	// TokenRepository is operate of token data.
	TokenRepository interface {
		FindToken(token string, t *entity.Token) error
		Create(u *entity.User, t *entity.Token) error
		Delete(token string) error
		DeleteExpired() (int64, error)
		CanAutoDeleteExpired() bool
	}
)
