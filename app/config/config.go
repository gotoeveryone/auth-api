package config

import (
	"math/rand"
	"time"

	"github.com/gotoeveryone/golib"
)

// configuration
type appConfig struct {
	golib.Config
	App app
}

// application configuration
type app struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Timezone string `json:"timezone"`
}

var (
	// AppConfig is configuration data read from JSON file
	AppConfig appConfig

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
