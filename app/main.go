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

	// Set timezone
	var err error
	time.Local, err = time.LoadLocation(config.GetenvOrDefault("TZ", "Asia/Tokyo"))
	if err != nil {
		logrus.Error(fmt.Sprintf("Get location error: %s", err))
		// continue with default timezone.
	}

	c := config.App{
		DB: config.DB{
			Host:     config.GetenvOrDefault("DATABASE_HOST", "127.0.0.1"),
			Port:     config.GetenvOrDefault("DATABASE_PORT", "3306"),
			Name:     config.GetenvOrDefault("DATABASE_NAME", "auth_api"),
			User:     config.GetenvOrDefault("DATABASE_USER", "auth_api"),
			Password: config.GetenvOrDefault("DATABASE_PASSWORD", ""),
			Timezone: time.Local,
		},
	}

	if config.GetenvOrDefault("APP_ENV", "dev") == "dev" {
		c.Debug = true
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

	host := config.GetenvOrDefault("APP_HOST", "0.0.0.0")
	port := config.GetenvOrDefault("APP_PORT", "8080")
	if err := r.Run(fmt.Sprintf("%s:%s", host, port)); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
