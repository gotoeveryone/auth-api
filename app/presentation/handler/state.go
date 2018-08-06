package handler

import "github.com/gin-gonic/gin"

// StateHandler is state action handler
type StateHandler interface {
	Get(c *gin.Context)
}
