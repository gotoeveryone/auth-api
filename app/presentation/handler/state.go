package handler

import "github.com/gin-gonic/gin"

// State is action handler for application state
type State interface {
	Get(c *gin.Context)
	NoRoute(c *gin.Context)
	NoMethod(c *gin.Context)
}
