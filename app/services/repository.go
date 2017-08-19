package services

import (
	"encoding/json"
	"errors"
	"general-api/app/models"
	"math/rand"
	"time"

	"github.com/gotoeveryone/golib"
)

var r *rand.Rand // Rand for this package.

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// TokensService トークンモデルのサービス
type TokensService struct{}

// FindUser 対象トークンのユーザが有効かどうかを判定
func (s TokensService) FindUser(token string) (*models.User, error) {
	var m models.Token
	if err := s.FindToken(token, &m); err != nil {
		return nil, err
	}

	if err := dbManager.Preload("User", "is_active = ?", 1).Find(&m).Error; err != nil {
		return nil, err
	}

	if m.User.Account == "" {
		return nil, errors.New("ユーザが不正")
	}

	return &m.User, nil
}

// FindToken トークンにマッチするログイン情報を取得します。
func (s TokensService) FindToken(token string, model *models.Token) error {

	if !s.UseCached() {
		if err := dbManager.Where("token = ?", token).Order("created desc").First(model).Error; err != nil {
			return err
		}

		// 有効期限内のもの
		comp := model.CreatedAt.Add(time.Duration(model.Expire) * time.Second)
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

	if err := json.Unmarshal(o.([]byte), model); err != nil {
		return err
	}
	return nil
}

// Create トークン保存処理
func (s TokensService) Create(token models.Token) error {

	if !s.UseCached() {
		if err := dbManager.Create(&token).Error; err != nil {
			return err
		}

		return nil
	}

	// JSONに変換
	o, err := json.Marshal(token)
	if err != nil {
		return err
	}

	// キャッシュサーバに接続
	var rs RedisService
	if err := rs.SetWithExpire(token.Token, token.Expire, o); err != nil {
		return err
	}
	return nil
}

// Has トークンを保持しているかどうか
func (s TokensService) Has(token string, model *models.Token) (bool, error) {
	if err := s.FindToken(token, model); err != nil {
		return false, err
	}

	if model.Token == "" {
		return false, errors.New("Token is invalid")
	}

	return true, nil
}

// Delete トークン削除処理
func (s TokensService) Delete(token string) error {

	if !s.UseCached() {
		if err := dbManager.Where("token = ?", token).Delete(models.Token{}).Error; err != nil {
			return err
		}

		return nil
	}

	// キャッシュサーバに接続
	var rs RedisService
	if _, err := rs.Delete(token); err != nil {
		return err
	}

	return nil
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

// UseCached キャッシュサーバを利用するかどうかを判定
func (s TokensService) UseCached() bool {
	return golib.AppConfig.Cache.Use
}

// UsersService ユーザモデルのサービス
type UsersService struct{}

// Find モデル検索
func (s UsersService) Find(id int) (*models.User, error) {
	var u models.User
	if err := dbManager.Where("id = ?", id).Find(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// FindActiveUser アクティブなユーザを検索
func (s UsersService) FindActiveUser(account string) (*models.User, error) {
	var user models.User
	dbManager.Where("account = ?", account).Where("is_active = ?", 1).Find(&user)

	if err := dbManager.Error; err != nil {
		return nil, err
	}

	return &user, nil
}
