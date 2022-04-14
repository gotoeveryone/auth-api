package handler

import "github.com/gin-gonic/gin"

// Authenticate is action handler about authentication and authorization
type Authenticate interface {
	Authenticate(c *gin.Context)
	Deauthenticate(c *gin.Context)
	Registration(c *gin.Context)
	Activate(c *gin.Context)
	GetUser(c *gin.Context)
}
