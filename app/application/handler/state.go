package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/general-api/app/config"
	"github.com/gotoeveryone/general-api/app/domain/entity"
)

// GetState is get application state
func GetState(c *gin.Context) {
	c.JSON(http.StatusOK, entity.State{
		Status:      "Active",
		Environment: gin.Mode(),
		LogLevel:    config.AppConfig.Log.Level,
		TimeZone:    time.Local.String(),
	})
}
