package registry

import (
	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/auth-api/app/config"
	"github.com/gotoeveryone/auth-api/app/domain/repository"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func NewRouter(config config.App, userRepo repository.UserRepository, tokenRepo repository.TokenRepository) *gin.Engine {
	// Initialize application
	r := gin.Default()
	r.HandleMethodNotAllowed = true

	// Handler
	sh := NewStateHandler()
	ah := NewAuthenticateHandler(userRepo, tokenRepo)

	// Middleware
	m := NewAuthenticateMiddleware(userRepo, tokenRepo)

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
	if config.Debug {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return r
}
