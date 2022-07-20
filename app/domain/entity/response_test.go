package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	e := Error{}
	assert.Empty(t, e.Code)
	assert.Empty(t, e.Message)
}

func TestStatus(t *testing.T) {
	s := State{}
	assert.Empty(t, s.Status)
	assert.Empty(t, s.Environment)
}

func TestGeneratedPassword(t *testing.T) {
	s := GeneratedPassword{}
	assert.Empty(t, s.Password)
}
