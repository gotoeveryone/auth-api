package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/auth-api/app/config"
	"github.com/gotoeveryone/auth-api/app/registry"
	_ "github.com/gotoeveryone/auth-api/docs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// @title    General authentication API
// @version  1.0
// @license.name Kazuki Kamizuru
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	isDebug := false
	if config.GetenvOrDefault("APP_ENV", "dev") == "dev" {
		isDebug = true
	}

	// Set log level
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// Set release mode
	if !isDebug {
		gin.SetMode(gin.ReleaseMode)
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Set timezone
	var err error
	time.Local, err = time.LoadLocation(config.GetenvOrDefault("TZ", "Asia/Tokyo"))
	if err != nil {
		log.Error().Msgf("Get location error: %s", err)
		// continue with default timezone.
	}

	c := config.App{
		Debug: isDebug,
		DB: config.DB{
			Host:     config.GetenvOrDefault("DATABASE_HOST", "127.0.0.1"),
			Port:     config.GetenvOrDefault("DATABASE_PORT", "3306"),
			Name:     config.GetenvOrDefault("DATABASE_NAME", "auth_api"),
			User:     config.GetenvOrDefault("DATABASE_USER", "auth_api"),
			Password: config.GetenvOrDefault("DATABASE_PASSWORD", ""),
			Timezone: time.Local,
		},
	}

	// Initialize datastore
	if err := registry.InitDatastore(c.Debug, c.DB); err != nil {
		log.Fatal().Err(err)
	}

	// Initialize router
	r, err := registry.NewRouter(c)
	if err != nil {
		log.Fatal().Err(err)
	}

	host := config.GetenvOrDefault("APP_HOST", "0.0.0.0")
	port := config.GetenvOrDefault("APP_PORT", "8080")
	if err := r.Run(fmt.Sprintf("%s:%s", host, port)); err != nil {
		log.Fatal().Err(err)
	}
}
