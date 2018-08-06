package handler

import "github.com/gin-gonic/gin"

// AuthenticateHandler is authenticate action handler
type AuthenticateHandler interface {
	Authenticate(c *gin.Context)
	Deauthenticate(c *gin.Context)
	Registration(c *gin.Context)
	Activate(c *gin.Context)
	GetUser(c *gin.Context)
}
