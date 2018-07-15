package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/general-api/app/config"
	"github.com/gotoeveryone/general-api/app/domain/entity"
	"github.com/gotoeveryone/general-api/app/domain/repository"
	"github.com/gotoeveryone/general-api/app/handler"
	"github.com/gotoeveryone/general-api/app/infrastructure"
	"github.com/gotoeveryone/general-api/app/middleware"
	"github.com/gotoeveryone/golib"
	"github.com/gotoeveryone/golib/logs"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	// Load configuration from JSON file
	if err := golib.LoadConfig(&config.AppConfig, ""); err != nil {
		panic(fmt.Errorf("LoadConfig error: %s", err))
	}
	config := config.AppConfig

	// Initial log
	if err := logs.Init(config.Log.Prefix, config.Log.Path, config.Log.Level); err != nil {
		panic(fmt.Errorf("LogConfig error: %s", err))
	}

	// Set timezone
	time.Local, _ = time.LoadLocation(config.AppTimezone)

	// Initial database
	infrastructure.InitDB(config.DB)

	// Initial application
	r := gin.Default()

	// Not found
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, entity.Error{
			Code:    http.StatusNotFound,
			Message: http.StatusText(http.StatusNotFound),
		})
	})

	// Method not allowed
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, entity.Error{
			Code:    http.StatusMethodNotAllowed,
			Message: http.StatusText(http.StatusMethodNotAllowed),
		})
	})

	// Routing
	r.GET("/", handler.GetState)
	v1 := r.Group("v1")
	{
		v1.GET("/", handler.GetState)
		v1.POST("/users", handler.Registration)
		v1.POST("/activate", handler.Activate)
		v1.POST("/auth", handler.Authenticate)
		auth := v1.Group("")
		{
			auth.Use(middleware.HasToken())
			auth.GET("/users", handler.GetUser)
			auth.DELETE("/deauth", handler.Deauthenticate)
		}
	}

	// Deleting expired tokens.
	// When can't auto delete expired tokens, this function is behavior.
	tr := infrastructure.NewTokenRepository()
	if !tr.CanAutoDeleteExpired() {
		go func(repo repository.TokenRepository) {
			for {
				cnt, err := repo.DeleteExpired()
				if err != nil {
					logs.Error(err)
				}
				if cnt > 0 {
					logs.Info(fmt.Sprintf("Expired %d tokens was deleted.", cnt))
				}
				time.Sleep(60 * time.Second)
			}
		}(tr)
	}

	r.Run(fmt.Sprintf(":%d", config.Port))
}
