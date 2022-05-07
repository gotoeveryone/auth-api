package mock

import (
	"errors"

	"github.com/gotoeveryone/auth-api/app/domain/entity"
)

type UserRepository struct {
	User            *entity.User
	IsMatchPassword bool
}

func (r *UserRepository) Exists(account string) (bool, error) {
	return r.User != nil, nil
}

func (r *UserRepository) Find(id uint, u *entity.User) error {
	u = r.User
	return nil
}

func (r *UserRepository) FindByAccount(account string) (*entity.User, error) {
	return r.User, nil
}

func (r *UserRepository) MatchPassword(hashedPassword, password string) error {
	if r.IsMatchPassword {
		return nil
	}
	return errors.New("Password not matched")
}

func (r *UserRepository) Create(u *entity.User) (string, error) {
	return "hogefuga", nil
}

func (r *UserRepository) UpdatePassword(u *entity.User, pass string) error {
	return nil
}

func (r *UserRepository) UpdateAuthed(u *entity.User) error {
	return nil
}
