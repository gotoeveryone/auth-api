package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/auth-api/app/config"
	"github.com/gotoeveryone/auth-api/app/domain/repository"
	"github.com/gotoeveryone/auth-api/app/registry"
	_ "github.com/gotoeveryone/auth-api/docs"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
		Cache: config.Cache{
			Host: getEnv("CACHE_HOST", "127.0.0.1"),
			Port: getEnv("CACHE_PORT", "6379"),
			Auth: getEnv("CACHE_AUTH", ""),
		},
	}

	if getEnv("APP_ENV", "dev") == "dev" {
		c.Debug = true
	}

	if getEnv("USE_CACHE", "") != "" {
		c.Cache.Use = true
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
	r := gin.Default()
	r.HandleMethodNotAllowed = true

	// Repository
	ur := registry.NewUserRepository()
	tr := registry.NewTokenRepository(c)

	// Handler
	sh := registry.NewStateHandler()
	ah := registry.NewAuthenticateHandler(ur, tr)

	// Middleware
	m := registry.NewAuthenticateMiddleware(ur, tr)

	// Routing
	// Root
	r.GET("/", sh.Get)
	// Not Found
	r.NoRoute(sh.NoRoute)
	// Method Not Allowed
	r.NoMethod(sh.NoMethod)
	// Application
	v1 := r.Group("v1")
	{
		v1.GET("/", sh.Get)
		v1.POST("/users", ah.Registration)
		v1.POST("/activate", ah.Activate)
		v1.POST("/auth", ah.Authenticate)
		auth := v1.Group("")
		{
			auth.Use(m.Authorized())
			auth.GET("/users", ah.GetUser)
			auth.DELETE("/deauth", ah.Deauthenticate)
		}
	}

	// show swagger ui to /swagger/index.html
	if c.Debug {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Deleting expired tokens.
	// When can't auto delete expired tokens, this function is behavior.
	if !tr.CanAutoDeleteExpired() {
		go func(repo repository.TokenRepository) {
			for {
				cnt, err := repo.DeleteExpired()
				if err != nil {
					logrus.Error(err)
				}
				if cnt > 0 {
					logrus.Info(fmt.Sprintf("Expired %d tokens was deleted.", cnt))
				}
				time.Sleep(60 * time.Second)
			}
		}(tr)
	}

	host := getEnv("APP_HOST", "0.0.0.0")
	port := getEnv("APP_PORT", "8080")
	if err := r.Run(fmt.Sprintf("%s:%s", host, port)); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
