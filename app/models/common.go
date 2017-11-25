package models

import (
	"time"

	"github.com/gotoeveryone/golib"
	"golang.org/x/crypto/bcrypt"
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
	Error   error  `json:"-"`
}

// State 状態
type State struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	LogLevel    string `json:"logLevel"`
	TimeZone    string `json:"timeZone"`
}

// Activate ユーザ有効化
type Activate struct {
	Login
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

// Login ログイン時の入力情報
type Login struct {
	Account  string `json:"account" binding:"required,min=6,max=10"`
	Password string `json:"password" binding:"required,min=8"`
}

// User ユーザ
type User struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	Account     string     `gorm:"type:varchar(10);not null;unique_index" json:"account" binding:"required,min=6,max=10"`
	Name        string     `gorm:"type:varchar(20);not null" json:"name" binding:"required,max=50"`
	Password    string     `gorm:"type:varchar(255);not null" json:"-"`
	Sex         string     `gorm:"type:enum('Male','Female');not null" json:"sex" binding:"required"`
	MailAddress *string    `gorm:"type:varchar(100)" json:"mailAddress" binding:"required"`
	Role        string     `gorm:"type:enum('Administrator','General');not null" json:"role"`
	LastLogged  *time.Time `gorm:"type:datetime" json:"-"`
	IsActive    bool       `gorm:"type:tinyint;not null" json:"-"`
	IsEnable    bool       `gorm:"type:tinyint;not null" json:"-"`
}

// Token トークン
type Token struct {
	ID          uint      `gorm:"primary_key" json:"-"`
	UserID      uint      `gorm:"type:int unsigned;not null" json:"id"`
	Token       string    `gorm:"type:varchar(64);not null;unique_index" json:"accessToken"`
	Environment string    `gorm:"type:varchar(20);not null" json:"environment"`
	CreatedAt   time.Time `gorm:"type:datetime;not null" json:"-"`
	ExpiredAt   time.Time `gorm:"type:datetime;not null" json:"-"`
	User        User      `json:"-"`
}

// MatchPassword パスワードが一致しているかを確認する
func (u *User) MatchPassword(input string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input))
}

// GetDefaultRole デフォルトの権限を取得する
func (u *User) GetDefaultRole() string {
	return "General"
}
