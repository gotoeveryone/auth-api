package infrastructure

import (
	"time"

	"github.com/gotoeveryone/general-api/app/config"
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

	return (count > 0), nil
}

// FindByAccount is find user data from account and password
func (r userRepository) FindByAccount(account string) (*entity.User, error) {
	var u entity.User
	dbManager.Where(&entity.User{Account: account, IsEnable: true}).Find(&u)

	if dbManager.RecordNotFound() {
		return nil, nil
	}

	if err := dbManager.Error; err != nil {
		return nil, err
	}

	return &u, nil
}

// FindByToken is judge user has valid token
func (r userRepository) FindByToken(token string) (*entity.User, error) {
	var u entity.User
	dbManager.Joins("INNER JOIN tokens ON users.id = tokens.user_id").
		Where(&entity.User{IsEnable: true}).
		Where("tokens.token = ?", token).Find(&u)

	if dbManager.RecordNotFound() {
		return nil, nil
	}

	if err := dbManager.Error; err != nil {
		return nil, err
	}

	return &u, nil
}

// ValidUser is valid user
func (r userRepository) ValidUser(u *entity.User) bool {
	return u != nil && u.Account != "" && u.IsEnable
}

// ValidRole is valid user role
func (r userRepository) ValidRole(role string) bool {
	roles := []string{entity.RoleAdministrator, entity.RoleGeneral}
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
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
		u.Role = u.GetDefaultRole()
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
