package repository

import (
	"github.com/gotoeveryone/auth-api/app/domain/entity"
)

// UserRepository is operate of user data.
type UserRepository interface {
	Exists(account string) (bool, error)
	Find(id uint, u *entity.User) error
	FindByAccount(account string) (*entity.User, error)
	ValidUser(u *entity.User) bool
	ValidRole(role string) bool
	MatchPassword(hashedPassword, password string) error
	Create(u *entity.User) (string, error)
	UpdatePassword(u *entity.User, pass string) error
	UpdateAuthed(u *entity.User) error
}
