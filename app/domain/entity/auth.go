package entity

import (
	"time"
)

var (
	// RoleAdministrator is administrator user.
	RoleAdministrator = "Administrator"
	// RoleGeneral is general user.
	RoleGeneral = "General"
)

// User is struct of authenticated user data
type User struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	Account     string     `gorm:"type:varchar(10);not null;unique_index" json:"account" binding:"required,min=6,max=10"`
	Name        string     `gorm:"type:varchar(20);not null" json:"name" binding:"required,max=50"`
	Password    string     `gorm:"type:varchar(255);not null" json:"-"`
	Sex         string     `gorm:"type:enum('Male','Female');not null" json:"sex" binding:"required"`
	MailAddress *string    `gorm:"type:varchar(100)" json:"mailAddress" binding:"required,email"`
	Role        string     `gorm:"type:enum('Administrator','General');not null" json:"role"`
	LastLogged  *time.Time `gorm:"type:datetime" json:"-"`
	IsActive    bool       `gorm:"type:tinyint;not null" json:"-"`
	IsEnable    bool       `gorm:"type:tinyint;not null" json:"-"`
	CreatedAt   time.Time  `gorm:"type:datetime;not null" sql:"default:current_timestamp" json:"-"`
	Tokens      []Token    `json:"-"`
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

// GetDefaultRole is get user default role
func (u *User) GetDefaultRole() string {
	return RoleGeneral
}
