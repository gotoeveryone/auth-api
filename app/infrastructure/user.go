package infrastructure

import (
	"errors"
	"time"

	"github.com/gotoeveryone/general-api/app/domain/entity"
	"github.com/gotoeveryone/general-api/app/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

type (
	// User control repository.
	userRepository struct{}
)

// NewUserRepository is create user management repository.
func NewUserRepository() repository.UserRepository {
	return userRepository{}
}

// Exists is confirm to account already exists
func (r userRepository) Exists(account string) (bool, error) {
	var count int
	dbManager.Model(&entity.User{}).Where(&entity.User{Account: account}).Count(&count)
	if err := dbManager.Error; err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

// FindByUserAndPassword is find user data from account and password
func (r userRepository) FindByUserAndPassword(account string, password string) (*entity.User, error) {
	var u entity.User
	dbManager.Where(&entity.User{Account: account}).Find(&u)
	if err := dbManager.Error; err != nil {
		return nil, err
	}

	// Check password matching from user has password
	if err := u.MatchPassword(password); err != nil {
		return nil, err
	}

	return &u, nil
}

// FindByToken is judge user has valid token
func (r userRepository) FindByToken(token string) (*entity.User, error) {
	var u entity.User
	if err := dbManager.Joins("INNER JOIN tokens ON users.id = tokens.user_id").Where(&entity.Token{Token: token}).
		Find(&u).Error; err != nil {
		return nil, err
	}

	if u.Account == "" {
		return nil, errors.New("Access token is invalid")
	}

	return &u, nil
}

// Create is create user data
func (r userRepository) Create(u *entity.User, pass string) error {
	pass, err := r.hashedPassword(pass)
	if err != nil {
		return err
	}
	u.Password = pass

	// If not specify role, use default role
	if u.Role == "" {
		u.Role = u.GetDefaultRole()
	}
	return dbManager.Create(u).Error
}

// UpdatePassword is update new password
func (r userRepository) UpdatePassword(u *entity.User, pass string) error {
	newpass, err := r.hashedPassword(pass)
	if err != nil {
		return err
	}
	u.Password = newpass
	u.IsActive = true
	return dbManager.Save(u).Error
}

// UpdateAuthed is update authenticated date
func (r userRepository) UpdateAuthed(u *entity.User) error {
	now := time.Now()
	u.LastLogged = &now
	return dbManager.Save(u).Error
}

// Get hashed password
func (r userRepository) hashedPassword(pass string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPass), nil
}
