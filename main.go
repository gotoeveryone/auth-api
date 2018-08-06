package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/auth-api/app/config"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"github.com/gotoeveryone/auth-api/app/domain/repository"
	"github.com/gotoeveryone/auth-api/app/registry"
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

	// Set application mode
	if !c.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Set timezone
	time.Local, err = time.LoadLocation(c.App.Timezone)
	if err != nil {
		config.Logger.Error(fmt.Sprintf("Get location error: %s", err))
		// continue with default timezone.
	}

	// Initial datastore
	if err := registry.InitDatastore(); err != nil {
		config.Logger.Error(err)
		os.Exit(1)
	}

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

	// Repository
	ur := registry.NewUserRepository()
	tr := registry.NewTokenRepository()

	// Handler
	sh := registry.NewStateHandler()
	ah := registry.NewAuthenticateHandler(ur, tr)

	// Middleware
	m := registry.NewAuthenticateMiddleware(ur)

	// Routing
	r.GET("/", sh.Get)
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

	// Deleting expired tokens.
	// When can't auto delete expired tokens, this function is behavior.
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

	if err := r.Run(fmt.Sprintf("%s:%d", c.App.Host, c.App.Port)); err != nil {
		config.Logger.Error(err)
		os.Exit(1)
	}
}
