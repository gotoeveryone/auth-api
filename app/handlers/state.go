package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/general-api/app/models"
	"github.com/gotoeveryone/general-api/app/services"
)

// GetState 状態監視
func GetState(c *gin.Context) {
	c.JSON(http.StatusOK, models.State{
		Status:      "Active",
		Environment: gin.Mode(),
		LogLevel:    services.AppConfig.Log.Level,
		TimeZone:    time.Local.String(),
	})
}
