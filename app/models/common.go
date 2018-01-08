package models

import (
	"time"

	"github.com/gotoeveryone/golib"
	"golang.org/x/crypto/bcrypt"
)

// AppConfig is sturct of application configuration
type AppConfig struct {
	golib.Config
	Port        int    `json:"port"`
	AppTimezone string `json:"appTimezone"`
}

// Error is struct of error object
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   error  `json:"-"`
}

// State is struct of Application state
type State struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	LogLevel    string `json:"logLevel"`
	TimeZone    string `json:"timeZone"`
}

// Activate is validation struct of using during activate user
type Activate struct {
	Authenticate
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

// Authenticate is validation struct of using during authentication
type Authenticate struct {
	Account  string `json:"account" binding:"required,min=6,max=10"`
	Password string `json:"password" binding:"required,min=8"`
}

// User is struct of authenticated user data
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
	CreatedAt   time.Time  `gorm:"type:datetime;not null" sql:"default:current_timestamp" json:"-"`
}

// Token is struct of authenticated token data
type Token struct {
	ID          uint      `gorm:"primary_key" json:"-"`
	UserID      uint      `gorm:"type:int unsigned;not null" json:"id"`
	Token       string    `gorm:"type:varchar(64);not null;unique_index" json:"accessToken"`
	Environment string    `gorm:"type:varchar(20);not null" json:"environment"`
	CreatedAt   time.Time `gorm:"type:datetime;not null" json:"-"`
	ExpiredAt   time.Time `gorm:"type:datetime;not null" json:"-"`
	User        User      `json:"-"`
}

// MatchPassword is check whether password match
func (u *User) MatchPassword(input string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input))
}

// GetDefaultRole is get user default role
func (u *User) GetDefaultRole() string {
	return "General"
}
