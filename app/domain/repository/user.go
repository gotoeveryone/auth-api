package repository

import (
	"github.com/gotoeveryone/auth-api/app/domain/entity"
)

// User is repository for operate about user.
type User interface {
	Exists(account string) (bool, error)
	Find(id uint, u *entity.User) error
	FindByAccount(account string) (*entity.User, error)
	MatchPassword(hashedPassword, password string) error
	Create(u *entity.User) (string, error)
	UpdatePassword(u *entity.User, pass string) error
	UpdateAuthed(u *entity.User) error
}
