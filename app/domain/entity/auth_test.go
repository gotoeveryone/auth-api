package entity

import "testing"

func TestUser(t *testing.T) {
	u := User{}
	if u.GetDefaultRole() != RoleGeneral {
		t.Errorf("Role is not default value.")
	}
	if u.ID != 0 {
		t.Errorf("ID is not default value")
	}
	if u.Account != "" {
		t.Errorf("Account is not default value")
	}
}

func TestToken(t *testing.T) {
	to := Token{}
	if to.ID != 0 {
		t.Errorf("ID is not default value")
	}
	if to.Environment != "" {
		t.Errorf("Environment is not default value")
	}
}
