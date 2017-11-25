package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/general-api/app/handlers"
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
	time.Local, _ = time.LoadLocation(config.AppTimezone)

	// DB設定初期化
	services.InitDB(config.DB)

	// Route初期化
	r := gin.Default()

	// Not found
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, models.Error{
			Code:    http.StatusNotFound,
			Message: http.StatusText(http.StatusNotFound),
		})
	})

	// Method not allowed
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, models.Error{
			Code:    http.StatusMethodNotAllowed,
			Message: http.StatusText(http.StatusMethodNotAllowed),
		})
	})

	// ルーティング
	r.GET("/", handlers.GetState)
	v1 := r.Group("v1")
	{
		v1.GET("/", handlers.GetState)
		v1.POST("/users", handlers.Publish)
		v1.POST("/activate", handlers.Activate)
		v1.POST("/auth", handlers.Authenticate)
		auth := v1.Group("")
		{
			auth.Use(middlewares.HasToken())
			auth.GET("/users", handlers.GetUser)
			auth.DELETE("/deauth", handlers.Deauthenticate)
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
					logs.Info(fmt.Sprintf("Expired %d tokens was deleted.", cnt))
				}
				time.Sleep(60 * time.Second)
			}
		}(ts)
	}

	r.Run(fmt.Sprintf(":%d", config.Port))
}
