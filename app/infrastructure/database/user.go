package database

import (
	"errors"
	"time"

	"github.com/gotoeveryone/auth-api/app/config"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"github.com/gotoeveryone/auth-api/app/domain/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userRepository struct{}

// NewUserRepository is create user management repository
func NewUserRepository() repository.User {
	return &userRepository{}
}

// Exists is confirm to account already exists
func (r userRepository) Exists(account string) (bool, error) {
	var count int64
	err := dbManager.Model(&entity.User{}).Where(&entity.User{Account: account}).Count(&count).Error
	if err != nil {
		return false, err
	}
	return (count > 0), nil
}

// Find is execute user data finding
func (r userRepository) Find(id uint) (*entity.User, error) {
	var u entity.User
	err := dbManager.Where(&entity.User{ID: id}).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// FindByAccount is find user data from account and password
func (r userRepository) FindByAccount(account string) (*entity.User, error) {
	var u entity.User
	err := dbManager.Where(&entity.User{Account: account, IsEnable: true}).Find(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// MatchPassword is check password matching from user has password
func (r userRepository) MatchPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Create is create user data and return generate password
func (r userRepository) Create(u *entity.User) (string, error) {
	// Issue initial password
	password := config.Generate(16)
	hashPassword, err := r.hashedPassword(password)
	if err != nil {
		return "", err
	}
	u.Password = hashPassword

	// If not specify role, use default role
	if u.Role == "" {
		u.Role = u.DefaultRole()
	}

	u.IsEnable = true
	return password, dbManager.Create(u).Error
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
