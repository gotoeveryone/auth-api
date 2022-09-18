package config

import (
	"math/rand"
	"os"
	"time"
)

// DB データベース接続設定
type DB struct {
	Name     string
	Host     string
	Port     string
	User     string
	Password string
	Timezone string
}

// App is application configuration
type App struct {
	Debug bool
	DB
}

const (
	IdentityKey = "id"
)

var (
	// Rand for this package.
	r *rand.Rand
)

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomString is generate string at specify length from character and digit
func RandomString(l int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	v := ""
	for i := 0; i < l; i++ {
		idx := r.Intn(len(letters))
		v += letters[idx : idx+1]
	}
	return v
}

// GetenvOrDefault is return got value from env or default value
func GetenvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
