package services

import (
	"math/rand"
	"time"

	"github.com/gotoeveryone/general-api/app/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	// AppConfig アプリケーション設定
	AppConfig models.AppConfig
)

var r *rand.Rand // Rand for this package.

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// GenerateToken トークン生成
func GenerateToken(l int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	v := ""
	for i := 0; i < l; i++ {
		idx := r.Intn(len(letters))
		v += letters[idx : idx+1]
	}
	return v
}

// HashedPassword パスワードをハッシュ化する
func HashedPassword(pass string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPass), nil
}
