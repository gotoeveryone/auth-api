package database

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mock sqlmock.Sqlmock
)

func init() {
	mock = initDBMock()
}

func initDBMock() sqlmock.Sqlmock {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	// defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}

	dbManager = gdb
	return mock
}
