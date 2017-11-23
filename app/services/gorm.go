package services

import (
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/gotoeveryone/general-api/app/models"
	"github.com/gotoeveryone/golib"
	"github.com/jinzhu/gorm"
)

var (
	// DbManager データベース接続用インスタンス
	dbManager *gorm.DB
)

// InitDB テーブル初期化
func InitDB(dbConfig golib.DB) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
		url.QueryEscape(dbConfig.Timezone),
	)

	dbManager, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	if gin.Mode() == gin.DebugMode {
		dbManager.LogMode(true)
	}

	// マイグレーション実行
	if err := dbManager.AutoMigrate(models.Token{}, models.User{}).Error; err != nil {
		panic(err)
	}
}
