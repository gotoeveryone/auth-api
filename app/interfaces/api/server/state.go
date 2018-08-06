package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/auth-api/app/config"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"github.com/gotoeveryone/auth-api/app/presentation/handler"
)

type stateHandler struct{}

// NewStateHandler is state action handler
func NewStateHandler() handler.StateHandler {
	return &stateHandler{}
}

// Get is get application state
func (h *stateHandler) Get(c *gin.Context) {
	c.JSON(http.StatusOK, entity.State{
		Status:      "Active",
		Environment: gin.Mode(),
		LogLevel:    config.AppConfig.Log.Level,
		TimeZone:    time.Local.String(),
	})
}
