package database

import (
	"fmt"

	mysqlDriver "github.com/go-sql-driver/mysql"
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
	fmt.Sprintln(dbConfig.Host, dbConfig.Port)
	c := mysqlDriver.Config{
		User:                 dbConfig.User,
		Passwd:               dbConfig.Password,
		DBName:               dbConfig.Name,
		Addr:                 fmt.Sprintf("%s:%s", dbConfig.Host, dbConfig.Port),
		Net:                  "tcp",
		ParseTime:            true,
		Loc:                  dbConfig.Timezone,
		AllowNativePasswords: true,
	}

	logMode := logger.Warn
	if debug {
		logMode = logger.Info
	}

	var err error
	dbManager, err = gorm.Open(mysql.New(mysql.Config{
		DSN: c.FormatDSN(),
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
