package redis

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
)

// ClientMock is client mock
type ClientMock struct {
	stored  map[string]interface{}
	expired map[string]int
}

func (c ClientMock) Connect() error {
	return nil
}

func (c ClientMock) Get(key string) (interface{}, error) {
	v, ok := c.stored[key]
	if ok {
		return v, nil
	}
	return v, fmt.Errorf("%s not stored", key)
}

func (c ClientMock) Set(key string, value interface{}) (interface{}, error) {
	c.stored[key] = value
	return value, nil
}

func (c ClientMock) Delete(key string) (interface{}, error) {
	v, err := c.Get(key)
	if err != nil {
		return nil, err
	}
	c.stored[key] = nil
	return v, nil
}

func (c ClientMock) Expire(key string, expire int) (interface{}, error) {
	v, err := c.Get(key)
	if err != nil {
		return nil, err
	}
	c.expired[key] = expire
	return v, nil
}

func (c ClientMock) SetWithExpire(key string, expire int, value interface{}) error {
	if _, err := c.Set(key, value); err != nil {
		return err
	}
	if _, err := c.Expire(key, expire); err != nil {
		return err
	}
	return nil
}

var client *ClientMock

func init() {
	client = &ClientMock{
		stored:  map[string]interface{}{},
		expired: map[string]int{},
	}
}

func TestNewTokenRepository(t *testing.T) {
	a := reflect.TypeOf(&tokenRepository{
		client: client,
	})
	e := reflect.TypeOf(NewTokenRepository(client))
	if a != e {
		t.Errorf("NewTokenRepository type is mismatch, actual: %s, expected: %s", a, e)
	}
}

func TestFindToken(t *testing.T) {
	r := tokenRepository{
		client: client,
	}

	v := "test"
	e := entity.Token{}
	if err := r.Find(v, &e); err != nil {
		if !strings.Contains(err.Error(), "not stored") {
			t.Error(err)
		}
	}
	if e.Token != "" {
		t.Errorf("actual: exists, expected: not exists [%s]", e.Token)
	}

	e.Token = v
	b, err := json.Marshal(e)
	if err != nil {
		t.Error(err)
	}
	if _, err := r.client.Set(v, b); err != nil {
		t.Error(err)
	}

	if err := r.Find(v, &e); err != nil {
		t.Error(err)
	}
	if e.Token == "" {
		t.Errorf("actual: not exists, expected: exists [%s]", v)
	}
}

func TestCreateToken(t *testing.T) {
	r := tokenRepository{
		client: client,
	}

	u := entity.User{
		ID: 1,
	}
	e := entity.Token{}
	if err := r.Create(&u, &e); err != nil {
		t.Error(err)
	}
	if u.ID != e.UserID {
		t.Errorf("ID not matched, user: %d, token: %d", u.ID, e.UserID)
	}
	if gin.Mode() != e.Environment {
		t.Errorf("Environment not matched, user: %s, token: %s", gin.Mode(), e.Environment)
	}
	if e.Token == "" {
		t.Errorf("Token is empty")
	}
}

func TestDeleteToken(t *testing.T) {
	r := tokenRepository{
		client: client,
	}

	v := "test"
	if err := r.Delete(v); err != nil {
		t.Error(err)
	}
}

func TestDeleteExpiredToken(t *testing.T) {
	r := tokenRepository{
		client: client,
	}

	cnt, err := r.DeleteExpired()
	if err != nil {
		t.Error(err)
	}
	if cnt != 0 {
		t.Errorf("RowsAffected is not matched, actual: %d, expected: %d", 0, cnt)
	}
}

func TestCanAutoDeleteExpired(t *testing.T) {
	r := tokenRepository{
		client: client,
	}
	if !r.CanAutoDeleteExpired() {
		t.Errorf("TestCanAutoDeleteExpired failed")
	}
}
