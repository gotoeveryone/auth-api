package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/auth-api/app/config"
	"github.com/gotoeveryone/auth-api/app/registry"
	_ "github.com/gotoeveryone/auth-api/docs"
	"github.com/sirupsen/logrus"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// @title    General authentication API
// @version  1.0
// @license.name Kazuki Kamizuru
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Initialize logger
	logrus.SetFormatter(&logrus.JSONFormatter{})

	c := config.App{
		DB: config.DB{
			Host:     getEnv("DATABASE_HOST", "127.0.0.1"),
			Port:     getEnv("DATABASE_PORT", "3306"),
			Name:     getEnv("DATABASE_NAME", "auth_api"),
			User:     getEnv("DATABASE_USER", "auth_api"),
			Password: getEnv("DATABASE_PASSWORD", ""),
		},
	}

	if getEnv("APP_ENV", "dev") == "dev" {
		c.Debug = true
	}

	// Set timezone
	var err error
	time.Local, err = time.LoadLocation(getEnv("TZ", "Asia/Tokyo"))
	if err != nil {
		logrus.Error(fmt.Sprintf("Get location error: %s", err))
		// continue with default timezone.
	}

	// Set release mode
	if !c.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize datastore
	if err := registry.InitDatastore(c.Debug, c.DB); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	// Initialize application
	r := registry.NewRouter(c)

	host := getEnv("APP_HOST", "0.0.0.0")
	port := getEnv("APP_PORT", "8080")
	if err := r.Run(fmt.Sprintf("%s:%s", host, port)); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
