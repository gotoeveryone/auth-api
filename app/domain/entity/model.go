package entity

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Role string

type Gender string

const (
	RoleAdministrator = Role("Administrator")
	RoleGeneral       = Role("General")

	GenderMale    = Gender("Male")
	GenderFemale  = Gender("Female")
	GenderUnknown = Gender("Unknown")
)

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(b []byte) error {
	var err error
	d.Time, err = time.Parse("\"2006-01-02\"", string(b))
	return err
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", d.Time.Format("2006-01-02"))), nil
}

func (d *Date) Scan(value any) error {
	d.Time = value.(time.Time)
	return nil
}

func (d Date) Value() (driver.Value, error) {
	return d.Time, nil
}

// User is struct of authenticated user data
type User struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	Account     string     `gorm:"type:varchar(20);not null;unique_index" json:"account"`
	Name        string     `gorm:"type:varchar(50);not null" json:"name"`
	Password    string     `gorm:"type:varchar(255);not null" json:"-"`
	Gender      Gender     `gorm:"type:enum('Male','Female','Unknown');not null" json:"gender"`
	MailAddress string     `gorm:"type:varchar(255);not null" json:"mailAddress"`
	Birthday    Date       `gorm:"type:date;not null" json:"birthday"`
	Role        Role       `gorm:"type:enum('Administrator','General');not null"`
	LastLogged  *time.Time `gorm:"type:datetime" json:"-"`
	IsActive    bool       `gorm:"type:tinyint;not null" json:"-"`
	IsEnable    bool       `gorm:"type:tinyint;not null" json:"-"`
	CreatedAt   time.Time  `gorm:"type:datetime;not null" sql:"default:current_timestamp" json:"-"`
}

// Valid is valid user data
func (u *User) Valid() bool {
	return u.Account != "" && u.IsEnable
}

// DefaultRole is get user default role
func (u *User) DefaultRole() Role {
	return RoleGeneral
}
