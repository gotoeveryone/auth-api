package main

import (
	"fmt"
	"general-api/app/controllers"
	"general-api/app/middlewares"
	"general-api/app/models"
	"general-api/app/services"
	"time"

	"github.com/gotoeveryone/golib/logs"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/golib"
)

func main() {
	// 設定ファイル読み出し
	golib.LoadConfig()

	// タイムゾーンの設定
	location := "Asia/Tokyo"
	loc, err := time.LoadLocation(location)
	if err != nil {
		// UTCから9時間後
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc

	// DB設定初期化
	services.InitDB()

	// Route初期化
	r := gin.Default()

	// 404
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, models.Error{
			Code:    404,
			Message: "Not Found",
		})
	})

	// 405
	r.NoMethod(func(c *gin.Context) {
		c.JSON(405, models.Error{
			Code:    405,
			Message: "Method Not Allowed",
		})
	})

	// ルーティング
	r.GET("/", controllers.RouteRedirect)
	v1 := r.Group("v1")
	{
		v1.GET("/", controllers.GetState)
		v1.POST("/auth", controllers.Authenticate)
		auth := v1.Group("")
		{
			auth.Use(middlewares.HasToken())
			auth.GET("/users", controllers.GetUser)
			auth.DELETE("/deauth", controllers.Deauthenticate)
		}
	}

	// トークンの削除を定期的に実施するためのバッチ処理
	// キャッシュサーバを利用している場合はそもそも動作させない
	var ts services.TokensService
	if !ts.UseCached() {
		go func(ts services.TokensService) {
			for {
				cnt, err := ts.DeleteExpired()
				if err != nil {
					logs.Error(err)
				}
				if cnt > 0 {
					logs.Info(fmt.Sprintf("トークンを%d件削除しました。", cnt))
				}
				time.Sleep(60 * time.Second)
			}
		}(ts)
	}

	r.Run()
}
