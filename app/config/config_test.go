package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	s := []int{1, 10, 30, 50, 64, 128}
	for _, v := range s {
		g := RandomString(v)
		assert.Equal(t, v, len(g))
	}
}

func TestGetenvOrDefault(t *testing.T) {
	os.Setenv("HOGE", "fuga")
	assert.Equal(t, "fuga", GetenvOrDefault("HOGE", "piyo"))
	assert.Equal(t, "piyo", GetenvOrDefault("HOGE1", "piyo"))
}
