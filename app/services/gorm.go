package services

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gotoeveryone/golang/common"
	"github.com/jinzhu/gorm"
)

var (
	// DbManager データベース接続用インスタンス
	dbManager *gorm.DB
)

// InitDB テーブル初期化
func InitDB(config common.Config) {

	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=%s",
		config.DB.User,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name,
		"Asia%2FTokyo",
	)

	dbManager, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	dbManager.LogMode(true)
}
