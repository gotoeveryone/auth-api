package database

import (
	"fmt"
	"net/url"

	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"github.com/gotoeveryone/golib/config"
	"github.com/jinzhu/gorm"
)

var (
	// Instance of connected database
	dbManager *gorm.DB
)

// Init is execute database connection initial setting
func Init(debug bool, dbConfig config.DB) error {
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
		return err
	}

	dbManager.LogMode(debug)

	// マイグレーション実行
	if err := dbManager.AutoMigrate(entity.Token{}, entity.User{}).Error; err != nil {
		return err
	}

	return nil
}
