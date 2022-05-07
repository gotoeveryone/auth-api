package handler

import "github.com/gin-gonic/gin"

// Authenticate is action handler about authentication and authorization
type Authenticate interface {
	Registration(c *gin.Context)
	Activate(c *gin.Context)
	User(c *gin.Context)
}
