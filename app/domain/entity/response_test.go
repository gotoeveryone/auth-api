package entity

import "testing"

func TestError(t *testing.T) {
	e := Error{}
	if e.Code != 0 {
		t.Errorf("Code is not default value")
	}
	if e.Message != nil {
		t.Errorf("Message is not default value")
	}
}

func TestStatus(t *testing.T) {
	s := State{}
	if s.Status != "" {
		t.Errorf("Status is not default value")
	}
	if s.Environment != "" {
		t.Errorf("Environment is not default value")
	}
}
