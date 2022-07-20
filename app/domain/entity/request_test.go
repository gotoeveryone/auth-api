package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActivate(t *testing.T) {
	a := Activate{
		Authenticate: Authenticate{
			Account:  "testtest",
			Password: "testtest",
		},
	}

	assert.Empty(t, a.NewPassword)
}

func TestAuthenticate(t *testing.T) {
	a := Authenticate{}

	assert.Empty(t, a.Account)
	assert.Empty(t, a.Password)
}
