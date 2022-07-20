package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"github.com/gotoeveryone/auth-api/app/mock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginFailedInvalidParam(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	m := NewAuthMiddleware(&mock.UserRepository{})
	middleware, err := m.Create()
	if err != nil {
		t.Fatal(err)
	}

	r.POST("/v1/auth", middleware.LoginHandler)

	e := entity.Authenticate{}
	j, err := json.Marshal(e)
	if err != nil {
		t.Error(err)
	}
	body := bytes.NewBuffer(j)
	req, _ := http.NewRequest("POST", "/v1/auth", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

func TestLoginFailedAccountNotExist(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	m := NewAuthMiddleware(&mock.UserRepository{})
	middleware, err := m.Create()
	if err != nil {
		t.Fatal(err)
	}

	r.POST("/v1/auth", middleware.LoginHandler)

	e := entity.Authenticate{
		Account:  "testuser",
		Password: "password",
	}
	j, err := json.Marshal(e)
	if err != nil {
		t.Error(err)
	}
	body := bytes.NewBuffer(j)
	req, _ := http.NewRequest("POST", "/v1/auth", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

func TestLoginFailedInvalidAccount(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	cryptedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), 10)
	m := NewAuthMiddleware(&mock.UserRepository{User: &entity.User{
		Account:  "testuser",
		Password: string(cryptedPassword),
		IsEnable: false,
	}})
	middleware, err := m.Create()
	if err != nil {
		t.Fatal(err)
	}

	r.POST("/v1/auth", middleware.LoginHandler)

	e := entity.Authenticate{
		Account:  "testuser",
		Password: "Invalid001",
	}
	j, err := json.Marshal(e)
	if err != nil {
		t.Error(err)
	}
	body := bytes.NewBuffer(j)
	req, _ := http.NewRequest("POST", "/v1/auth", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

func TestLoginFaileNotActiveAccount(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	cryptedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), 10)
	m := NewAuthMiddleware(&mock.UserRepository{User: &entity.User{
		Account:  "testuser",
		Password: string(cryptedPassword),
		IsEnable: true,
		IsActive: false,
	}})
	middleware, err := m.Create()
	if err != nil {
		t.Fatal(err)
	}

	r.POST("/v1/auth", middleware.LoginHandler)

	e := entity.Authenticate{
		Account:  "testuser",
		Password: "password",
	}
	j, err := json.Marshal(e)
	if err != nil {
		t.Error(err)
	}
	body := bytes.NewBuffer(j)
	req, _ := http.NewRequest("POST", "/v1/auth", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

func TestLoginFailedPasswordNotMatched(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	cryptedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), 10)
	m := NewAuthMiddleware(&mock.UserRepository{User: &entity.User{
		Account:  "testuser",
		Password: string(cryptedPassword),
		IsEnable: true,
		IsActive: true,
	}})
	middleware, err := m.Create()
	if err != nil {
		t.Fatal(err)
	}

	r.POST("/v1/auth", middleware.LoginHandler)

	e := entity.Authenticate{
		Account:  "testuser",
		Password: "invalid_password",
	}
	j, err := json.Marshal(e)
	if err != nil {
		t.Error(err)
	}
	body := bytes.NewBuffer(j)
	req, _ := http.NewRequest("POST", "/v1/auth", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

func TestLoginSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	password := "password"
	cryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	m := NewAuthMiddleware(&mock.UserRepository{User: &entity.User{
		Account:  "testuser",
		Password: string(cryptedPassword),
		IsEnable: true,
		IsActive: true,
	}})
	middleware, err := m.Create()
	if err != nil {
		t.Fatal(err)
	}

	r.POST("/v1/auth", middleware.LoginHandler)

	e := entity.Authenticate{
		Account:  "testuser",
		Password: password,
	}
	j, err := json.Marshal(e)
	if err != nil {
		t.Error(err)
	}
	body := bytes.NewBuffer(j)
	req, _ := http.NewRequest("POST", "/v1/auth", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)

	c := entity.Claim{}
	if err := json.Unmarshal(w.Body.Bytes(), &c); err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, c.Token)
}
