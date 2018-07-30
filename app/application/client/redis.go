package client

import (
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/gotoeveryone/golib"
)

var (
	con *redis.Conn
)

// RedisClient Operationg of redis connection.
type RedisClient struct {
	Config golib.Cache
}

// Connect Connect to Redis.
func (s RedisClient) Connect() error {
	c := s.Config
	newCon, err := redis.Dial("tcp", c.Host+":"+strconv.Itoa(c.Port))
	if err != nil {
		return err
	}

	// When could got auth property from configuration data, Execute "AUTH" command.
	if c.Auth != "" {
		if _, err := newCon.Do("AUTH", c.Auth); err != nil {
			return err
		}
	}

	con = &newCon
	return nil
}

// Get Execute "GET" command.
func (s RedisClient) Get(key string) (interface{}, error) {
	if con == nil {
		if err := s.Connect(); err != nil {
			return nil, err
		}
	}
	return (*con).Do("GET", key)
}

// Set Execute "SET" command.
func (s RedisClient) Set(key string, value interface{}) (interface{}, error) {
	if con == nil {
		if err := s.Connect(); err != nil {
			return nil, err
		}
	}
	return (*con).Do("SET", key, value)
}

// Delete Execute "DEL" command.
func (s RedisClient) Delete(key string) (interface{}, error) {
	if con == nil {
		if err := s.Connect(); err != nil {
			return nil, err
		}
	}
	return (*con).Do("DEL", key)
}

// Expire Execute "EXPIRE" command.
func (s RedisClient) Expire(key string, expire int) (interface{}, error) {
	if con == nil {
		if err := s.Connect(); err != nil {
			return nil, err
		}
	}
	return (*con).Do("EXPIRE", key, expire)
}

// SetWithExpire Execute "SET" and "EXPIRE" command.
func (s RedisClient) SetWithExpire(key string, expire int, value interface{}) error {
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
