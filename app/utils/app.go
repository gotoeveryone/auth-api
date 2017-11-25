package utils

import (
	"math/rand"
	"time"
)

var r *rand.Rand // Rand for this package.

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Generate 指定桁にマッチするランダム文字列を生成
func Generate(l int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	v := ""
	for i := 0; i < l; i++ {
		idx := r.Intn(len(letters))
		v += letters[idx : idx+1]
	}
	return v
}
