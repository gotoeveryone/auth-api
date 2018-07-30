package infrastructure

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/auth-api/app/application/client"
	"github.com/gotoeveryone/auth-api/app/config"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"github.com/gotoeveryone/auth-api/app/domain/repository"
)

type (
	// Token control common behavior.
	tokenRepository struct{}

	// Token control by database.
	dbTokenRepository struct {
		tokenRepository
	}

	// Token control by Redis.
	redisTokenRepository struct {
		tokenRepository
		Client client.RedisClient
	}
)

// NewTokenRepository is create token management repository.
func NewTokenRepository() repository.TokenRepository {
	if config.AppConfig.Cache.Use {
		redisClient := client.RedisClient{
			Config: config.AppConfig.Cache,
		}
		return redisTokenRepository{
			Client: redisClient,
		}
	}
	return dbTokenRepository{}
}

// Create entity from user.
func (r tokenRepository) createFromUser(u *entity.User, t *entity.Token) {
	key := []byte(u.Account + time.Now().Format("20060102150405000"))
	bytes := sha512.Sum512_256(key)
	t.Token = hex.EncodeToString(bytes[:])
	t.UserID = u.ID
	t.Environment = gin.Mode()
	t.ExpiredAt = time.Now().Add(time.Duration(r.getExpire()) * time.Second)
}

// Get expire seconds.
func (r tokenRepository) getExpire() int {
	return 600
}

// FindToken is execute token data finding
func (r dbTokenRepository) FindToken(token string, t *entity.Token) error {
	return dbManager.Where(&entity.Token{Token: token}).
		Where("expired_at >= ?", time.Now()).First(t).Error
}

// Create is execute token data creating
func (r dbTokenRepository) Create(u *entity.User, t *entity.Token) error {
	r.createFromUser(u, t)
	return dbManager.Create(t).Error
}

// Delete is execute token data deleting
func (r dbTokenRepository) Delete(token string) error {
	return dbManager.Where(&entity.Token{Token: token}).
		Delete(entity.Token{}).Error
}

// DeleteExpired is execute expired token data deleting
func (r dbTokenRepository) DeleteExpired() (int64, error) {
	cnt := dbManager.Where("expired_at < ?", time.Now()).
		Delete(entity.Token{}).RowsAffected
	return cnt, dbManager.Error
}

// CanAutoDeleteExpired is returns whether token data can be automatically deleted.
func (r dbTokenRepository) CanAutoDeleteExpired() bool {
	return false
}

// FindToken is execute token data finding
func (r redisTokenRepository) FindToken(token string, t *entity.Token) error {
	o, err := r.Client.Get(token)
	if err != nil {
		return err
	} else if o == nil {
		return errors.New("Token is invalid")
	}

	return json.Unmarshal(o.([]byte), t)
}

// Create is execute token data creating
func (r redisTokenRepository) Create(u *entity.User, t *entity.Token) error {
	r.createFromUser(u, t)
	// Conver to JSON
	o, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return r.Client.SetWithExpire(t.Token, r.getExpire(), o)
}

// Delete is execute token data deleting
func (r redisTokenRepository) Delete(token string) error {
	_, err := r.Client.Delete(token)
	return err
}

// DeleteExpired is execute expired token data deleting
func (r redisTokenRepository) DeleteExpired() (int64, error) {
	return 0, nil
}

// CanAutoDeleteExpired is returns whether token data can be automatically deleted.
func (r redisTokenRepository) CanAutoDeleteExpired() bool {
	return true
}
