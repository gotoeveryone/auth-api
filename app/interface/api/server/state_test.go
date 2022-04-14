package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
)

func TestNewStateHandler(t *testing.T) {
	a := reflect.TypeOf(&stateHandler{})
	e := reflect.TypeOf(NewStateHandler())
	if a != e {
		t.Errorf("NewStateHandler type is mismatch, actual: %s, expected: %s", a, e)
	}
}

func TestGet(t *testing.T) {
	h := &stateHandler{}
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.GET("/", h.Get)
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("HTTP status code failed, actual: %d, expected: %d", http.StatusOK, w.Code)
	}

	s := entity.State{}
	if err := json.Unmarshal(w.Body.Bytes(), &s); err != nil {
		t.Error(err)
	}
	if s.Environment != gin.Mode() {
		t.Errorf("Environment is invalid, actual: %s, expected: %s", gin.Mode(), s.Environment)
	}
}

func TestNoRoute(t *testing.T) {
	h := &stateHandler{}
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.GET("/", h.NoRoute)
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("HTTP status code failed, actual: %d, expected: %d", http.StatusNotFound, w.Code)
	}

	e := entity.Error{}
	if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
		t.Error(err)
	}
	if e.Code != http.StatusNotFound {
		t.Errorf("HTTP status code failed, actual: %d, expected: %d", http.StatusNotFound, e.Code)
	}
	m := http.StatusText(http.StatusNotFound)
	if e.Message != m {
		t.Errorf("Message is invalid, actual: %s, expected: %s", m, e.Message)
	}
}

func TestNoMethod(t *testing.T) {
	h := &stateHandler{}
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.GET("/", h.NoMethod)
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("HTTP status code failed, actual: %d, expected: %d", http.StatusMethodNotAllowed, w.Code)
	}

	e := entity.Error{}
	if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
		t.Error(err)
	}
	if e.Code != http.StatusMethodNotAllowed {
		t.Errorf("HTTP status code failed, actual: %d, expected: %d", http.StatusMethodNotAllowed, e.Code)
	}
	m := http.StatusText(http.StatusMethodNotAllowed)
	if e.Message != m {
		t.Errorf("Message is invalid, actual: %s, expected: %s", m, e.Message)
	}
}
