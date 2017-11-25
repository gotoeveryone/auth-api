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

// TokensService トークンモデルのサービス
type TokensService struct{}

// FindUser 対象トークンのユーザが有効かどうかを判定
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

// FindToken トークン取得処理
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

// Create トークン保存処理
func (s TokensService) Create(u *models.User, t *models.Token) error {
	// トークンの生成
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

	// JSONに変換
	o, err := json.Marshal(t)
	if err != nil {
		return err
	}

	var rs RedisService
	return rs.SetWithExpire(t.Token, expire, o)
}

// Delete トークン削除処理
func (s TokensService) Delete(token string) error {
	if !s.UseCached() {
		return dbManager.Where(&models.Token{Token: token}).
			Delete(models.Token{}).Error
	}

	var rs RedisService
	_, err := rs.Delete(token)
	return err
}

// DeleteExpired 有効期限切れトークンを削除する
func (s TokensService) DeleteExpired() (int64, error) {
	if s.UseCached() {
		return 0, nil
	}
	cnt := dbManager.Where("expired_at < ?", time.Now()).
		Delete(models.Token{}).RowsAffected
	return cnt, dbManager.Error
}

// UseCached キャッシュサーバを利用するかどうかを判定
func (s TokensService) UseCached() bool {
	return AppConfig.Cache.Use
}

// UsersService ユーザモデルのサービス
type UsersService struct{}

// Exists アカウントが存在するかを確認
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

// FindUser ユーザ取得処理
func (s UsersService) FindUser(account string, password string) (*models.User, error) {
	var u models.User
	dbManager.Where(&models.User{Account: account}).Find(&u)
	if err := dbManager.Error; err != nil {
		return nil, err
	}

	// パスワードの一致確認
	if err := u.MatchPassword(password); err != nil {
		return nil, err
	}

	return &u, nil
}

// Create ユーザ登録処理
func (s UsersService) Create(u *models.User, pass string) error {
	pass, err := s.hashedPassword(pass)
	if err != nil {
		return err
	}
	u.Password = pass

	// 指定がなければデフォルトの権限
	if u.Role == "" {
		u.Role = u.GetDefaultRole()
	}
	return dbManager.Create(u).Error
}

// UpdatePassword パスワード更新処理
func (s UsersService) UpdatePassword(u *models.User, pass string) error {
	newpass, err := s.hashedPassword(pass)
	if err != nil {
		return err
	}
	u.Password = newpass
	u.IsActive = true
	return dbManager.Save(u).Error
}

// UpdateAuthed 認証日時保存処理
func (s UsersService) UpdateAuthed(u *models.User) error {
	now := time.Now()
	u.LastLogged = &now
	return dbManager.Save(u).Error
}

// パスワードをハッシュ化する
func (s UsersService) hashedPassword(pass string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPass), nil
}
