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

	h := NewAuthHandler(&mock.UserRepository{})
	r.POST("/v1/users", h.Registration)

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

	if w.Code != http.StatusBadRequest {
		t.Errorf("Failed: HTTP status code is not matched, actual: %d, expected: %d", w.Code, http.StatusBadRequest)
	}
}

func TestRegistrationSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	h := NewAuthHandler(&mock.UserRepository{})
	r.POST("/v1/users", h.Registration)

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

	if w.Code != http.StatusCreated {
		t.Errorf("Failed: HTTP status code is not matched, actual: %d, expected: %d", w.Code, http.StatusCreated)
	}

	e := entity.GeneratedPassword{}
	if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
		t.Error(err)
	}
	if e.Password == "" {
		t.Error("Failed: Password is empty")
	}
}

func TestActivateFailedInvalidParam(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	h := NewAuthHandler(&mock.UserRepository{})
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

	if w.Code != http.StatusBadRequest {
		t.Errorf("Failed: HTTP status code is not matched, actual: %d, expected: %d", w.Code, http.StatusBadRequest)
	}
}

func TestActivateFailedSamePassword(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	h := NewAuthHandler(&mock.UserRepository{})
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

	if w.Code != http.StatusBadRequest {
		t.Errorf("Failed: HTTP status code is not matched, actual: %d, expected: %d", w.Code, http.StatusBadRequest)
	}
}

func TestActivateFailedAccountNotExist(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	h := NewAuthHandler(&mock.UserRepository{})
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

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Failed: HTTP status code is not matched, actual: %d, expected: %d", w.Code, http.StatusUnauthorized)
	}
}

func TestActivateFailedPasswordNotMatched(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	h := NewAuthHandler(&mock.UserRepository{User: &entity.User{
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

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Failed: HTTP status code is not matched, actual: %d, expected: %d", w.Code, http.StatusUnauthorized)
	}
}

func TestActivateSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	h := NewAuthHandler(&mock.UserRepository{User: &entity.User{
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

	if w.Code != http.StatusOK {
		t.Errorf("Failed: HTTP status code is not matched, actual: %d, expected: %d", w.Code, http.StatusCreated)
	}
}

func TestGetUser(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	u := entity.User{Account: "testuser"}
	r.Use(setIdentity(&u))

	h := NewAuthHandler(&mock.UserRepository{})
	r.GET("/v1/me", h.User)

	req, _ := http.NewRequest("GET", "/v1/me", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Failed: HTTP status code is not matched, actual: %d, expected: %d", w.Code, http.StatusOK)
	}

	e := entity.User{}
	if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
		t.Error(err)
	}
	if e.Account != u.Account {
		t.Errorf("Failed: Account is not matched, actual: %s, expected: %s", e.Account, u.Account)
	}
}
