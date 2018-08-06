package redis

import (
	"encoding/json"
	"errors"

	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"github.com/gotoeveryone/auth-api/app/domain/repository"
	"github.com/gotoeveryone/auth-api/app/infrastructure"
)

type tokenRepository struct {
	infrastructure.BaseTokenRepository
	client Client
}

// NewTokenRepository is create token management repository using Redis.
func NewTokenRepository(c Client) repository.TokenRepository {
	return tokenRepository{
		client: c,
	}
}

// FindToken is execute token data finding
func (r tokenRepository) FindToken(token string, t *entity.Token) error {
	o, err := r.client.Get(token)
	if err != nil {
		return err
	} else if o == nil {
		return errors.New("Token is invalid")
	}

	return json.Unmarshal(o.([]byte), t)
}

// Create is execute token data creating
func (r tokenRepository) Create(u *entity.User, t *entity.Token) error {
	r.CreateFromUser(u, t)
	// Conver to JSON
	o, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return r.client.SetWithExpire(t.Token, r.Expire(), o)
}

// Delete is execute token data deleting
func (r tokenRepository) Delete(token string) error {
	_, err := r.client.Delete(token)
	return err
}

// DeleteExpired is execute expired token data deleting
func (r tokenRepository) DeleteExpired() (int64, error) {
	return 0, nil
}

// CanAutoDeleteExpired is returns whether token data can be automatically deleted.
func (r tokenRepository) CanAutoDeleteExpired() bool {
	return true
}
