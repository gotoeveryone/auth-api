package models

import (
	"time"

	"github.com/gotoeveryone/golib"
)

// AppConfig アプリケーション設定1
type AppConfig struct {
	golib.Config
	Port        int    `json:"port"`
	AppTimezone string `json:"appTimezone"`
}

// Error エラー
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// State 状態
type State struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	LogLevel    string `json:"logLevel"`
	TimeZone    string `json:"timeZone"`
}

// Login ログイン時の入力情報
type Login struct {
	Account  string `json:"account" binding:"required,min=6,max=10"`
	Password string `json:"password" binding:"required,min=8"`
}

// User ユーザ
type User struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	Account     string     `gorm:"type:varchar(10);not null;unique_index" json:"userId" binding:"required,min=6,max=10"`
	Name        string     `gorm:"type:varchar(20);not null" json:"userName" binding:"required,max=50"`
	Password    string     `gorm:"type:varchar(255);not null" json:"-"`
	Sex         string     `gorm:"type:enum('男性','女性');not null" json:"sex" binding:"required"`
	MailAddress *string    `gorm:"type:varchar(100)" json:"mailAddress" binding:"required"`
	Role        string     `gorm:"type:enum('管理者','一般');not null" json:"role"`
	LastLogged  *time.Time `gorm:"type:datetime" json:"-"`
	IsActive    bool       `gorm:"type:tinyint;not null" json:"-"`
	IsEnable    bool       `gorm:"type:tinyint;not null" json:"-"`
}

// Token トークン
type Token struct {
	ID          uint      `gorm:"primary_key" json:"-"`
	UserID      uint      `gorm:"type:int unsigned;not null" json:"id"`
	Token       string    `gorm:"type:varchar(50);not null;unique_index" json:"access_token"`
	Expire      int       `gorm:"type:smallint unsigned" json:"-"`
	Environment string    `gorm:"type:varchar(20);not null" json:"environment"`
	CreatedAt   time.Time `gorm:"column:created;type:datetime" json:"-"`
	User        User      `json:"-"`
}
