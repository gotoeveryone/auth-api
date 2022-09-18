package handler

import "github.com/gin-gonic/gin"

// User is action handler about user data
type User interface {
	Register(c *gin.Context)
	Activate(c *gin.Context)
	Identity(c *gin.Context)
}
