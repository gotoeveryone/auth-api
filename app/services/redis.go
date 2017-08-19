package services

import (
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/gotoeveryone/golib"
)

var (
	con *redis.Conn
)

// RedisService Redis接続サービス
type RedisService struct{}

// Connect Redisへ接続
func (s RedisService) Connect() error {
	redisConig := golib.AppConfig.Cache
	// Redisに保存
	newCon, err := redis.Dial("tcp", redisConig.Host+":"+strconv.Itoa(redisConig.Port))
	if err != nil {
		return err
	}

	// AUTHが取得できた場合は認証
	if redisConig.Auth != "" {
		if _, err := newCon.Do("AUTH", redisConig.Auth); err != nil {
			return err
		}
	}

	con = &newCon
	return nil
}

// Get キーの取得
func (s RedisService) Get(key string) (interface{}, error) {
	if con == nil {
		if err := s.Connect(); err != nil {
			return nil, err
		}
	}
	return (*con).Do("GET", key)
}

// Set キーの設定
func (s RedisService) Set(key string, value interface{}) (interface{}, error) {
	if con == nil {
		if err := s.Connect(); err != nil {
			return nil, err
		}
	}
	return (*con).Do("SET", key, value)
}

// Delete キーの削除
func (s RedisService) Delete(key string) (interface{}, error) {
	if con == nil {
		if err := s.Connect(); err != nil {
			return nil, err
		}
	}
	return (*con).Do("DEL", key)
}

// Expire キーの有効期限設定
func (s RedisService) Expire(key string, expire int) (interface{}, error) {
	if con == nil {
		if err := s.Connect(); err != nil {
			return nil, err
		}
	}
	return (*con).Do("EXPIRE", key, expire)
}

// SetWithExpire キーの保存と有効期限を設定
func (s RedisService) SetWithExpire(key string, expire int, value interface{}) error {
	if con == nil {
		if err := s.Connect(); err != nil {
			return err
		}
	}

	if _, err := s.Set(key, value); err != nil {
		return err
	}

	if _, err := s.Expire(key, expire); err != nil {
		return err
	}

	return nil
}
