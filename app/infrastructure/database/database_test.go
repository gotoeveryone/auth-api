package database

import (
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
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
	dsn := "sqlmock-user_test"
	db, mock, err := sqlmock.NewWithDSN(dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}))
	if err != nil {
		panic(err)
	}

	dbManager = gdb
	return mock
}
