package entity

import (
	"time"
)

type Role string

const (
	// RoleAdministrator is administrator user.
	RoleAdministrator = Role("Administrator")
	// RoleGeneral is general user.
	RoleGeneral = Role("General")
)

// User is struct of authenticated user data
type User struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	Account     string     `gorm:"type:varchar(10);not null;unique_index" json:"account" binding:"required,min=6,max=10"`
	Name        string     `gorm:"type:varchar(20);not null" json:"name" binding:"required,max=50"`
	Password    string     `gorm:"type:varchar(255);not null" json:"-"`
	Sex         string     `gorm:"type:enum('Male','Female','Unknown');not null" json:"sex" binding:"required"`
	MailAddress *string    `gorm:"type:varchar(100)" json:"mailAddress" binding:"required,email"`
	Role        Role       `gorm:"type:enum('Administrator','General');not null" json:"role"`
	LastLogged  *time.Time `gorm:"type:datetime" json:"-"`
	IsActive    bool       `gorm:"type:tinyint;not null" json:"-"`
	IsEnable    bool       `gorm:"type:tinyint;not null" json:"-"`
	CreatedAt   time.Time  `gorm:"type:datetime;not null" sql:"default:current_timestamp" json:"-"`
}

func (u *User) Valid() bool {
	return u.Account != "" && u.IsEnable && u.IsActive
}

// ValidRole is valid user role
func (u *User) ValidRole() bool {
	roles := []Role{RoleAdministrator, RoleGeneral}
	for _, role := range roles {
		if u.Role == role {
			return true
		}
	}
	return false
}

// DefaultRole is get user default role
func (u *User) DefaultRole() Role {
	return RoleGeneral
}
