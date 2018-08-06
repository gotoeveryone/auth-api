package infrastructure

import (
	"crypto/sha512"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
)

// BaseTokenRepository is token control base repository
type BaseTokenRepository struct{}

// CreateFromUser is create entity from user.
func (r BaseTokenRepository) CreateFromUser(u *entity.User, t *entity.Token) {
	key := []byte(u.Account + time.Now().Format("20060102150405000"))
	bytes := sha512.Sum512_256(key)
	t.Token = hex.EncodeToString(bytes[:])
	t.UserID = u.ID
	t.Environment = gin.Mode()
	t.ExpiredAt = time.Now().Add(time.Duration(r.Expire()) * time.Second)
}

// Expire is get expire seconds.
func (r BaseTokenRepository) Expire() int {
	return 600
}
