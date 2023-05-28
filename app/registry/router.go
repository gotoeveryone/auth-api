package registry

import (
	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/auth-api/app/config"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(config config.App) (*gin.Engine, error) {
	// Initialize application
	r := gin.Default()
	r.HandleMethodNotAllowed = true

	// Repository
	ur := NewUserRepository()

	// Handler
	sh := NewStateHandler()
	uh := NewUserHandler(ur)

	// Middleware
	m, err := NewAuthMiddleware(ur).Create()
	if err != nil {
		return nil, err
	}

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
		v1.POST("/users", uh.Register)
		v1.POST("/activate", uh.Activate)
		v1.POST("/auth", m.LoginHandler)
		v1.GET("/refresh_token", m.RefreshHandler)
		auth := v1.Group("")
		{
			auth.Use(m.MiddlewareFunc())
			{
				auth.GET("/me", uh.Identity)
				auth.DELETE("/deauth", m.LogoutHandler)
			}
		}
	}

	// show swagger ui to /swagger/index.html
	if config.Debug {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return r, nil
}
