package services

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/general-api/app/models"
	"golang.org/x/crypto/bcrypt"
)

// TokensService is operate of token data.
type TokensService struct{}

// FindUser is judge user has valid token
func (s TokensService) FindUser(token string) (*models.User, error) {
	var t models.Token
	if err := s.FindToken(token, &t); err != nil {
		return nil, err
	}

	var u models.User
	if err := dbManager.Where(&models.User{ID: uint(t.UserID)}).
		Find(&u).Error; err != nil {
		return nil, err
	}

	if u.Account == "" {
		return nil, errors.New("Token is invalid")
	}

	return &u, nil
}

// FindToken is execute token data finding
func (s TokensService) FindToken(token string, t *models.Token) error {
	if !s.UseCached() {
		return dbManager.Where(&models.Token{Token: token}).
			Where("expired_at >= ?", time.Now()).First(t).Error
	}

	var rs RedisService
	o, err := rs.Get(token)
	if err != nil {
		return err
	} else if o == nil {
		return errors.New("Token is invalid")
	}

	return json.Unmarshal(o.([]byte), t)
}

// Create is execute token data creating
func (s TokensService) Create(u *models.User, t *models.Token) error {
	// Generate token value
	key := []byte(u.Account + time.Now().Format("20060102150405000"))
	bytes := sha256.Sum256(key)
	t.Token = hex.EncodeToString(bytes[:])
	t.UserID = u.ID
	t.Environment = gin.Mode()

	expire := 600
	t.ExpiredAt = time.Now().Add(time.Duration(expire) * time.Second)

	if !s.UseCached() {
		return dbManager.Create(t).Error
	}

	// Conver to JSON
	o, err := json.Marshal(t)
	if err != nil {
		return err
	}

	var rs RedisService
	return rs.SetWithExpire(t.Token, expire, o)
}

// Delete is execute token data deleting
func (s TokensService) Delete(token string) error {
	if !s.UseCached() {
		return dbManager.Where(&models.Token{Token: token}).
			Delete(models.Token{}).Error
	}

	var rs RedisService
	_, err := rs.Delete(token)
	return err
}

// DeleteExpired is execute expired token data deleting
func (s TokensService) DeleteExpired() (int64, error) {
	if s.UseCached() {
		return 0, nil
	}
	cnt := dbManager.Where("expired_at < ?", time.Now()).
		Delete(models.Token{}).RowsAffected
	return cnt, dbManager.Error
}

// UseCached is execute using cache service to whether save tokens confirming
func (s TokensService) UseCached() bool {
	return AppConfig.Cache.Use
}

// UsersService is operate of user data.
type UsersService struct{}

// Exists is confirm to account already exists
func (s UsersService) Exists(account string) (bool, error) {
	var count int
	dbManager.Model(&models.User{}).Where(&models.User{Account: account}).Count(&count)
	if err := dbManager.Error; err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

// FindUser is find user data from account and password
func (s UsersService) FindUser(account string, password string) (*models.User, error) {
	var u models.User
	dbManager.Where(&models.User{Account: account}).Find(&u)
	if err := dbManager.Error; err != nil {
		return nil, err
	}

	// Check password matching from user has password
	if err := u.MatchPassword(password); err != nil {
		return nil, err
	}

	return &u, nil
}

// Create is create user data
func (s UsersService) Create(u *models.User, pass string) error {
	pass, err := s.hashedPassword(pass)
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
func (s UsersService) UpdatePassword(u *models.User, pass string) error {
	newpass, err := s.hashedPassword(pass)
	if err != nil {
		return err
	}
	u.Password = newpass
	u.IsActive = true
	return dbManager.Save(u).Error
}

// UpdateAuthed is update authenticated date
func (s UsersService) UpdateAuthed(u *models.User) error {
	now := time.Now()
	u.LastLogged = &now
	return dbManager.Save(u).Error
}

// Get hashed password
func (s UsersService) hashedPassword(pass string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPass), nil
}
