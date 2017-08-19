package services

import (
	"fmt"

	"github.com/gotoeveryone/golib"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	// DbManager データベース接続用インスタンス
	dbManager *gorm.DB
)

// InitDB テーブル初期化
func InitDB() {
	dbConfig := golib.AppConfig.DB

	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
		"Asia%2FTokyo",
	)

	dbManager, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	dbManager.LogMode(true)
}
