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
	if u.Valid() {
		t.Error("User is valid.")
	}
	u.IsActive = true
	if !u.Valid() {
		t.Error("User is not valid.")
	}
}

func TestValidRole(t *testing.T) {
	u := User{Role: Role("test")}
	if u.ValidRole() {
		t.Errorf("Role is valid: %s", u.Role)
	}

	u.Role = RoleAdministrator
	if !u.ValidRole() {
		t.Errorf("Role is invalid: %s", u.Role)
	}

	u.Role = RoleGeneral
	if !u.ValidRole() {
		t.Errorf("Role is invalid: %s", u.Role)
	}
}

func TestUserDefaultRole(t *testing.T) {
	u := User{}
	if u.DefaultRole() != RoleGeneral {
		t.Errorf("Role is not default value, actual: [%s], expected: [%s]", u.DefaultRole(), RoleGeneral)
	}
}
