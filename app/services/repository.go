package services

import (
	"encoding/json"
	"errors"
	"general-api/app/models"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gotoeveryone/golang/common"
)

var (
	// Config 設定
	Config common.Config
)

// TokensService トークンモデルのサービス
type TokensService struct {
}

// FindUser 対象トークンのユーザが有効かどうかを判定
func (s TokensService) FindUser(token string) (*models.User, error) {
	var model models.Token
	if err := dbManager.Where("token = ?", token).First(&model).Error; err != nil {
		return nil, err
	}

	// 有効期限内のもの
	comp := model.CreatedAt.Add(time.Duration(model.Expire) * time.Second)
	if time.Now().Sub(comp).Seconds() > 0 {
		return nil, errors.New("Token is invalid")
	}

	if err := dbManager.Model(&model).Related(&model.User).Error; err != nil {
		return nil, err
	}

	if model.User.Account == "" {
		return nil, errors.New("ユーザが不正")
	}

	return &model.User, nil
}

// Save トークン保存処理
func (s TokensService) Create(token models.Token) error {

	// Redisを利用しない場合はデータベースへ
	if !Config.Redis.Use {
		// データベースに保存
		if err := dbManager.Create(&token).Error; err != nil {
			return err
		}

		return nil
	}

	// Redisに保存
	con, err := redis.Dial("tcp", Config.Redis.Host+":"+strconv.Itoa(Config.Redis.Port))
	if err != nil {
		// データベースに保存
		if err := dbManager.Create(&token).Error; err != nil {
			panic(err)
		}
	}
	defer con.Close()

	// AUTHが取得できた場合は認証
	if Config.Redis.Auth != "" {
		if _, err := con.Do("AUTH", Config.Redis.Auth); err != nil {
			return err
		}
	}

	o, err := json.Marshal(token)
	if err != nil {
		return err
	}

	if _, err := con.Do("SET", token.Token, o); err != nil {
		return err
	}

	if _, err := con.Do("EXPIRE", token.Token, token.Expire); err != nil {
		return err
	}

	return nil
}

// UsersService ユーザモデルのサービス
type UsersService struct {
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
