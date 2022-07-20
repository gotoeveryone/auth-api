package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidUser(t *testing.T) {
	u := User{}
	assert.False(t, u.Valid())
	u.Account = "testuser"
	assert.False(t, u.Valid())
	u.IsEnable = true
	assert.True(t, u.Valid())
}

func TestUserDefaultRole(t *testing.T) {
	u := User{}
	assert.Equal(t, u.DefaultRole(), RoleGeneral)
}
