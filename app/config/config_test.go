package config

import "testing"

func TestGenerate(t *testing.T) {
	s := []int{1, 10, 30, 50, 64, 128}
	for _, v := range s {
		g := Generate(v)
		if len(g) != v {
			t.Errorf("length not matched. [%d]", v)
		}
	}
}
