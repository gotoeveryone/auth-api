package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/general-api/app/handlers"
	"github.com/gotoeveryone/general-api/app/middlewares"
	"github.com/gotoeveryone/general-api/app/models"
	"github.com/gotoeveryone/general-api/app/services"
	"github.com/gotoeveryone/golib"
	"github.com/gotoeveryone/golib/logs"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	// Load configuration from JSON file
	if err := golib.LoadConfig(&services.AppConfig, ""); err != nil {
		panic(fmt.Errorf("LoadConfig error: %s", err))
	}
	config := services.AppConfig

	// Initial log
	if err := logs.Init(config.Log.Prefix, config.Log.Path, config.Log.Level); err != nil {
		panic(fmt.Errorf("LogConfig error: %s", err))
	}

	// Set timezone
	time.Local, _ = time.LoadLocation(config.AppTimezone)

	// Initial database
	services.InitDB(config.DB)

	// Initial application
	r := gin.Default()

	// Not found
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, models.Error{
			Code:    http.StatusNotFound,
			Message: http.StatusText(http.StatusNotFound),
		})
	})

	// Method not allowed
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, models.Error{
			Code:    http.StatusMethodNotAllowed,
			Message: http.StatusText(http.StatusMethodNotAllowed),
		})
	})

	// Routing
	r.GET("/", handlers.GetState)
	v1 := r.Group("v1")
	{
		v1.GET("/", handlers.GetState)
		v1.POST("/users", handlers.Registration)
		v1.POST("/activate", handlers.Activate)
		v1.POST("/auth", handlers.Authenticate)
		auth := v1.Group("")
		{
			auth.Use(middlewares.HasToken())
			auth.GET("/users", handlers.GetUser)
			auth.DELETE("/deauth", handlers.Deauthenticate)
		}
	}

	// Delete expire token from database.
	// When use cache `false` at configuration file, this function is behavior.
	var ts services.TokensService
	if !ts.UseCached() {
		go func(ts services.TokensService) {
			for {
				cnt, err := ts.DeleteExpired()
				if err != nil {
					logs.Error(err)
				}
				if cnt > 0 {
					logs.Info(fmt.Sprintf("Expired %d tokens was deleted.", cnt))
				}
				time.Sleep(60 * time.Second)
			}
		}(ts)
	}

	r.Run(fmt.Sprintf(":%d", config.Port))
}
