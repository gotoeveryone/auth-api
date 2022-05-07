package config

import (
	"math/rand"
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

// Generate a random character string matching the specified digit
func Generate(l int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	v := ""
	for i := 0; i < l; i++ {
		idx := r.Intn(len(letters))
		v += letters[idx : idx+1]
	}
	return v
}
