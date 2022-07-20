package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	h := &stateHandler{}
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.GET("/", h.Get)
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)

	s := entity.State{}
	if err := json.Unmarshal(w.Body.Bytes(), &s); err != nil {
		t.Error(err)
	}
	assert.Equal(t, s.Environment, gin.Mode())
}

func TestNoRoute(t *testing.T) {
	h := &stateHandler{}
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.GET("/", h.NoRoute)
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusNotFound)

	e := entity.Error{}
	if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
		t.Error(err)
	}
	assert.Equal(t, e.Message, http.StatusText(http.StatusNotFound))
}

func TestNoMethod(t *testing.T) {
	h := &stateHandler{}
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.GET("/", h.NoMethod)
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusMethodNotAllowed)

	e := entity.Error{}
	if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
		t.Error(err)
	}
	assert.Equal(t, e.Message, http.StatusText(http.StatusMethodNotAllowed))
}
