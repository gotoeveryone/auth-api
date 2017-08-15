package models

import (
	"time"
)

// Error エラー
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Login ログイン時の入力情報
type Login struct {
	Account  string `json:"account" binding:"required,min=6"`
	Password string `json:"password" binding:"required,min=8"`
}

// User ユーザ
type User struct {
	ID          uint    `gorm:"primary_key" json:"id"`
	Account     string  `json:"userId"`
	Name        string  `json:"userName"`
	Password    string  `json:"-"`
	Sex         string  `json:"sex"`
	MailAddress *string `json:"mailAddress"`
	Role        string  `json:"role"`
}

// Token トークン
type Token struct {
	ID          uint      `gorm:"primary_key" json:"-"`
	UserID      uint      `json:"id"`
	Token       string    `json:"access_token"`
	Environment string    `json:"environment"`
	User        User      `json:"-"`
	Expire      int       `json:"-"`
	CreatedAt   time.Time `gorm:"column:created;type:datetime" json:"-"`
}
