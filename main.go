package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/auth-api/app/application/handler"
	"github.com/gotoeveryone/auth-api/app/application/middleware"
	"github.com/gotoeveryone/auth-api/app/config"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"github.com/gotoeveryone/auth-api/app/domain/repository"
	"github.com/gotoeveryone/auth-api/app/infrastructure"
	"github.com/gotoeveryone/golib"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	// Load configuration from JSON file
	if err := golib.LoadConfig(&config.AppConfig, ""); err != nil {
		log.Fatal(fmt.Sprintf("LoadConfig error: %s", err))
	}
	c := config.AppConfig

	// Initial log
	var err error
	config.Logger, err = golib.NewLogger(c.Log)
	if err != nil {
		log.Fatal(fmt.Sprintf("Log initialize error: %s", err))
	}

	// Set timezone
	time.Local, err = time.LoadLocation(c.App.Timezone)
	if err != nil {
		config.Logger.Error(fmt.Sprintf("Get location error: %s", err))
		// continue with default timezone.
	}

	// Initial database
	infrastructure.InitDB(c.DB)

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
					config.Logger.Error(err)
				}
				if cnt > 0 {
					config.Logger.Info(fmt.Sprintf("Expired %d tokens was deleted.", cnt))
				}
				time.Sleep(60 * time.Second)
			}
		}(tr)
	}

	r.Run(fmt.Sprintf("%s:%d", c.App.Host, c.App.Port))
}
