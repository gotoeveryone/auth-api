package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"github.com/gotoeveryone/auth-api/app/presentation/handler"
	"github.com/rs/zerolog"
)

type stateHandler struct{}

// NewStateHandler is create action handler for state
func NewStateHandler() handler.State {
	return &stateHandler{}
}

// Get is get application state
// @Summary Return application state
// @Tags State
// @Produce json
// @Success 200 {object} entity.State
// @Failure 404 {object} entity.Error
// @Failure 405 {object} entity.Error
// @Router /v1 [get]
func (h *stateHandler) Get(c *gin.Context) {
	c.JSON(http.StatusOK, entity.State{
		Status:      "Active",
		Environment: gin.Mode(),
		LogLevel:    zerolog.GlobalLevel().String(),
		TimeZone:    time.Local.String(),
	})
}

// NoRoute is not found response
func (h *stateHandler) NoRoute(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusNotFound, entity.Error{
		Code:    http.StatusNotFound,
		Message: http.StatusText(http.StatusNotFound),
	})
}

// NoMethod is method not allowed response
func (h *stateHandler) NoMethod(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusMethodNotAllowed, entity.Error{
		Code:    http.StatusMethodNotAllowed,
		Message: http.StatusText(http.StatusMethodNotAllowed),
	})
}
