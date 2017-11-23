package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/general-api/app/controllers"
	"github.com/gotoeveryone/general-api/app/middlewares"
	"github.com/gotoeveryone/general-api/app/models"
	"github.com/gotoeveryone/general-api/app/services"
	"github.com/gotoeveryone/golib"
	"github.com/gotoeveryone/golib/logs"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	// 設定ファイル読み出し
	if err := golib.LoadConfig(&services.AppConfig, ""); err != nil {
		panic(fmt.Errorf("LoadConfig error: %s", err))
	}
	config := services.AppConfig

	// ログ設定
	if err := logs.Init(config.Log.Prefix, config.Log.Path, config.Log.Level); err != nil {
		panic(fmt.Errorf("LogConfig error: %s", err))
	}

	// タイムゾーンの設定
	loc, err := time.LoadLocation(config.AppTimezone)
	if err != nil {
		// UTCから9時間後
		loc = time.FixedZone(config.AppTimezone, 9*60*60)
	}
	time.Local = loc

	// DB設定初期化
	services.InitDB(config.DB)

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
	r.GET("/", controllers.GetState)
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

	r.Run(fmt.Sprintf(":%d", config.Port))
}
