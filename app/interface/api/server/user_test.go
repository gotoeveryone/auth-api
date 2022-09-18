package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/auth-api/app/config"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"github.com/gotoeveryone/auth-api/app/mock"
	"github.com/stretchr/testify/assert"
)

func setIdentity(user *entity.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(config.IdentityKey, user)
		c.Next()
	}
}

func TestRegistrationFailed(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	h := NewUserHandler(&mock.UserRepository{})
	r.POST("/v1/users", h.Register)

	p := entity.User{
		Account: "testuser",
	}
	j, err := json.Marshal(p)
	if err != nil {
		t.Error(err)
	}
	body := bytes.NewBuffer(j)
	req, _ := http.NewRequest("POST", "/v1/users", body)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusBadRequest)
}

func TestRegistrationSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	h := NewUserHandler(&mock.UserRepository{})
	r.POST("/v1/users", h.Register)

	role := "General"
	p := entity.RegistrationUser{
		Account:     "testuser",
		Name:        "Test User",
		Role:        &role,
		Gender:      "Unknown",
		MailAddress: "hoge@example.com",
		Birthday:    "2006-01-02",
	}
	j, err := json.Marshal(p)
	if err != nil {
		t.Error(err)
	}
	body := bytes.NewBuffer(j)
	req, _ := http.NewRequest("POST", "/v1/users", body)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusCreated)

	e := entity.GeneratedPassword{}
	if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, e.Password)
}

func TestActivateFailedInvalidParam(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	h := NewUserHandler(&mock.UserRepository{})
	r.POST("/v1/activate", h.Activate)

	p := entity.Activate{
		Authenticate: entity.Authenticate{
			Account:  "testuser",
			Password: "password",
		},
	}
	j, err := json.Marshal(p)
	if err != nil {
		t.Error(err)
	}
	body := bytes.NewBuffer(j)
	req, _ := http.NewRequest("POST", "/v1/activate", body)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusBadRequest)
}

func TestActivateFailedSamePassword(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	h := NewUserHandler(&mock.UserRepository{})
	r.POST("/v1/activate", h.Activate)

	p := entity.Activate{
		Authenticate: entity.Authenticate{
			Account:  "testuser",
			Password: "password",
		},
		NewPassword: "password",
	}
	j, err := json.Marshal(p)
	if err != nil {
		t.Error(err)
	}
	body := bytes.NewBuffer(j)
	req, _ := http.NewRequest("POST", "/v1/activate", body)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusBadRequest)
}

func TestActivateFailedAccountNotExist(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	h := NewUserHandler(&mock.UserRepository{})
	r.POST("/v1/activate", h.Activate)

	p := entity.Activate{
		Authenticate: entity.Authenticate{
			Account:  "testuser",
			Password: "HogeFuga001",
		},
		NewPassword: "HogeFuga001New",
	}
	j, err := json.Marshal(p)
	if err != nil {
		t.Error(err)
	}
	body := bytes.NewBuffer(j)
	req, _ := http.NewRequest("POST", "/v1/activate", body)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

func TestActivateFailedPasswordNotMatched(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	h := NewUserHandler(&mock.UserRepository{User: &entity.User{
		Account: "testuser",
	}, IsMatchPassword: false})
	r.POST("/v1/activate", h.Activate)

	p := entity.Activate{
		Authenticate: entity.Authenticate{
			Account:  "testuser",
			Password: "HogeFuga001",
		},
		NewPassword: "HogeFuga001New",
	}
	j, err := json.Marshal(p)
	if err != nil {
		t.Error(err)
	}
	body := bytes.NewBuffer(j)
	req, _ := http.NewRequest("POST", "/v1/activate", body)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

func TestActivateSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	h := NewUserHandler(&mock.UserRepository{User: &entity.User{
		Account: "testuser",
	}, IsMatchPassword: true})
	r.POST("/v1/activate", h.Activate)

	p := entity.Activate{
		Authenticate: entity.Authenticate{
			Account:  "testuser",
			Password: "HogeFuga001",
		},
		NewPassword: "HogeFuga001New",
	}
	j, err := json.Marshal(p)
	if err != nil {
		t.Error(err)
	}
	body := bytes.NewBuffer(j)
	req, _ := http.NewRequest("POST", "/v1/activate", body)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)
}

func TestIdentity(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	u := entity.User{Account: "testuser"}
	r.Use(setIdentity(&u))

	h := NewUserHandler(&mock.UserRepository{})
	r.GET("/v1/me", h.Identity)

	req, _ := http.NewRequest("GET", "/v1/me", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)

	e := entity.User{}
	if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
		t.Error(err)
	}
	assert.Equal(t, e.Account, u.Account)
}
