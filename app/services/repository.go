package services

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gotoeveryone/general-api/app/models"
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
		if err := dbManager.Where(&models.Token{Token: token}).
			Order("created desc").First(t).Error; err != nil {
			return err
		}

		// 有効期限内のもの
		comp := t.CreatedAt.Add(time.Duration(t.Expire) * time.Second)
		if time.Now().Sub(comp).Seconds() > 0 {
			return errors.New("Token is invalid")
		}

		return nil
	}

	// キャッシュサーバに接続
	var rs RedisService
	o, err := rs.Get(token)
	if o == nil {
		return errors.New("Token is invalid")
	}
	if err != nil {
		return err
	}

	return json.Unmarshal(o.([]byte), t)
}

// Create トークン保存処理
func (s TokensService) Create(t models.Token) error {
	if !s.UseCached() {
		return dbManager.Create(&t).Error
	}

	// JSONに変換
	o, err := json.Marshal(t)
	if err != nil {
		return err
	}

	// キャッシュサーバに接続
	var rs RedisService
	return rs.SetWithExpire(t.Token, t.Expire, o)
}

// Delete トークン削除処理
func (s TokensService) Delete(token string) error {
	if !s.UseCached() {
		return dbManager.Where(&models.Token{Token: token}).
			Delete(models.Token{}).Error
	}

	// キャッシュサーバに接続
	var rs RedisService
	_, err := rs.Delete(token)
	return err
}

// DeleteExpired 有効期限切れトークンを削除する
func (s TokensService) DeleteExpired() (int64, error) {
	if s.UseCached() {
		return 0, nil
	}
	cnt := dbManager.Where("DATE_ADD(created, INTERVAL expire second) < ?",
		time.Now()).Delete(models.Token{}).RowsAffected
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
	pass, err := HashedPassword(pass)
	if err != nil {
		return err
	}
	u.Password = pass

	// 指定がなければ一般ユーザ
	if u.Role == "" {
		u.Role = "General"
	}
	return dbManager.Create(u).Error
}

// UpdatePassword パスワード更新処理
func (s UsersService) UpdatePassword(u *models.User, pass string) error {
	newpass, err := HashedPassword(pass)
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
