package database

import (
	"fmt"
	"net/url"

	"github.com/gotoeveryone/auth-api/app/config"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// Instance of connected database
	dbManager *gorm.DB
)

// Init is execute database connection initial setting
func Init(debug bool, dbConfig config.DB) error {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
		url.QueryEscape(dbConfig.Timezone),
	)

	logMode := logger.Warn
	if debug {
		logMode = logger.Info
	}

	dbManager, err = gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})

	if err != nil {
		return err
	}

	// マイグレーション実行
	if err := dbManager.AutoMigrate(entity.User{}); err != nil {
		return err
	}

	return nil
}
