package entity

import "testing"

func TestValidUser(t *testing.T) {
	u := User{}
	if u.Valid() {
		t.Error("User is valid.")
	}
	u.Account = "testuser"
	if u.Valid() {
		t.Error("User is valid.")
	}
	u.IsEnable = true
	if !u.Valid() {
		t.Error("User is not valid.")
	}
}

func TestUserDefaultRole(t *testing.T) {
	u := User{}
	if u.DefaultRole() != RoleGeneral {
		t.Errorf("Role is not default value, actual: [%s], expected: [%s]", u.DefaultRole(), RoleGeneral)
	}
}
