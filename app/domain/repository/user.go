package repository

import "github.com/gotoeveryone/general-api/app/domain/entity"

type (
	// UserRepository is operate of user data.
	UserRepository interface {
		Exists(account string) (bool, error)
		FindByUserAndPassword(account string, password string) (*entity.User, error)
		FindByToken(token string) (*entity.User, error)
		Create(u *entity.User, pass string) error
		UpdatePassword(u *entity.User, pass string) error
		UpdateAuthed(u *entity.User) error
	}
)
