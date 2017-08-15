package main

import (
	"general-api/app/controllers"
	"general-api/app/models"
	"general-api/app/services"
	"net/http"
	"time"

	"github.com/gotoeveryone/golang/common"

	"github.com/gin-gonic/gin"
)

func main() {
	// 設定ファイル読み出し
	common.LoadConfig(&services.Config)

	// タイムゾーンの設定
	location := "Asia/Tokyo"
	loc, err := time.LoadLocation(location)
	if err != nil {
		// UTCから9時間後
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc

	// DB設定初期化
	services.InitDB(services.Config)

	r := gin.Default()

	// ミドルウェア
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 404
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, models.Error{
			Code:    http.StatusNotFound,
			Message: "Not Found",
		})
	})

	// 405
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, models.Error{
			Code:    http.StatusMethodNotAllowed,
			Message: "Method Not Allowed",
		})
	})

	// ルーティング
	r.GET("/", controllers.GetState)
	r.GET("/users", controllers.GetUser)
	r.POST("/login", controllers.Authenticate)

	r.Run()
}
