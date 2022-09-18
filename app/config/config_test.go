package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	s := []int{1, 10, 30, 50, 64, 128}
	for _, v := range s {
		g := RandomString(v)
		assert.Equal(t, len(g), v)
	}
}
