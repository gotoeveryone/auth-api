package services

import (
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"github.com/gotoeveryone/general-api/app/models"
)

var r *rand.Rand // Rand for this package.

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// TokensService トークンモデルのサービス
type TokensService struct{}

// FindUser 対象トークンのユーザが有効かどうかを判定
func (s TokensService) FindUser(token string) (*models.User, error) {
	var t models.Token
	if err := s.FindToken(token, &t); err != nil {
		return nil, err
	}

	var u models.User
	if err := dbManager.Where("id = ?", t.UserID).Find(&u).Error; err != nil {
		return nil, err
	}

	if u.Account == "" {
		return nil, errors.New("Token is invalid")
	}

	return &u, nil
}

// FindToken トークンにマッチするログイン情報を取得します。
func (s TokensService) FindToken(token string, t *models.Token) error {
	if !s.UseCached() {
		if err := dbManager.Where("token = ?", token).Order("created desc").First(t).Error; err != nil {
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
		return dbManager.Where("token = ?", token).
			Delete(models.Token{}).Error
	}

	// キャッシュサーバに接続
	var rs RedisService
	_, err := rs.Delete(token)
	return err
}

// GenerateToken トークン生成
func (s TokensService) GenerateToken() string {
	strlen := 50
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := ""
	for i := 0; i < strlen; i++ {
		idx := r.Intn(len(letters))
		result += letters[idx : idx+1]
	}
	return result
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

// Find モデル検索
func (s UsersService) Find(id int) (*models.User, error) {
	var u models.User
	dbManager.Where("id = ?", id).Find(&u)
	return &u, dbManager.Error
}

// FindActiveUser アクティブなユーザを検索
func (s UsersService) FindActiveUser(account string) (*models.User, error) {
	var u models.User
	dbManager.Where(&models.User{Account: account, IsActive: true}).Find(&u)
	return &u, dbManager.Error
}

// UpdateAuthed 認証日時保存処理
func (s UsersService) UpdateAuthed(u *models.User) error {
	now := time.Now()
	u.LastLogged = &now
	return dbManager.Save(u).Error
}
