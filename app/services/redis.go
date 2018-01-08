package services

import (
	"strconv"

	"github.com/garyburd/redigo/redis"
)

var (
	con *redis.Conn
)

// RedisService Operationg of redis connection.
type RedisService struct{}

// Connect Connect to Redis.
func (s RedisService) Connect() error {
	redisConig := AppConfig.Cache
	newCon, err := redis.Dial("tcp", redisConig.Host+":"+strconv.Itoa(redisConig.Port))
	if err != nil {
		return err
	}

	// When could got auth property from configuration data, Execute "AUTH" command.
	if redisConig.Auth != "" {
		if _, err := newCon.Do("AUTH", redisConig.Auth); err != nil {
			return err
		}
	}

	con = &newCon
	return nil
}

// Get Execute "GET" command.
func (s RedisService) Get(key string) (interface{}, error) {
	if con == nil {
		if err := s.Connect(); err != nil {
			return nil, err
		}
	}
	return (*con).Do("GET", key)
}

// Set Execute "SET" command.
func (s RedisService) Set(key string, value interface{}) (interface{}, error) {
	if con == nil {
		if err := s.Connect(); err != nil {
			return nil, err
		}
	}
	return (*con).Do("SET", key, value)
}

// Delete Execute "DEL" command.
func (s RedisService) Delete(key string) (interface{}, error) {
	if con == nil {
		if err := s.Connect(); err != nil {
			return nil, err
		}
	}
	return (*con).Do("DEL", key)
}

// Expire Execute "EXPIRE" command.
func (s RedisService) Expire(key string, expire int) (interface{}, error) {
	if con == nil {
		if err := s.Connect(); err != nil {
			return nil, err
		}
	}
	return (*con).Do("EXPIRE", key, expire)
}

// SetWithExpire Execute "SET" and "EXPIRE" command.
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
