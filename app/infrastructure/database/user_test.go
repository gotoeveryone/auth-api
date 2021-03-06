package database

import (
	"reflect"
	"testing"
	"time"

	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestNewUserRepository(t *testing.T) {
	a := reflect.TypeOf(&userRepository{})
	e := reflect.TypeOf(NewUserRepository())
	if a != e {
		t.Errorf("NewTokenRepository type is mismatch, actual: %s, expected: %s", a, e)
	}
}

func TestExists(t *testing.T) {
	mock.ExpectQuery("SELECT count").
		WithArgs("test").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	r := userRepository{}

	v := "test"
	e, err := r.Exists(v)
	if err != nil {
		t.Error(err)
	}
	if !e {
		t.Errorf("User %s is not exists", v)
	}

	mock.ExpectQuery("SELECT count").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	e, err = r.Exists(v)
	if err != nil {
		t.Error(err)
	}
	if e {
		t.Errorf("User %s is exists", v)
	}
}

func TestFindByAccount(t *testing.T) {
	mock.ExpectQuery("SELECT *").
		WillReturnRows(sqlmock.NewRows([]string{"account"}))

	r := userRepository{}

	v := "test"
	u, err := r.FindByAccount(v)
	if err != nil {
		t.Error(err)
	}
	if u.Account != "" {
		t.Errorf("actual: not exists, expected: exists [%s]", u.Account)
	}

	mock.ExpectQuery("SELECT *").
		WillReturnRows(sqlmock.NewRows([]string{"account"}).AddRow("test"))

	u, err = r.FindByAccount(v)
	if err != nil {
		t.Error(err)
	}
	if u.Account == "" {
		t.Errorf("actual: exists [%s], expected: not exists", v)
	}
}

func TestFindByToken(t *testing.T) {
	mock.ExpectQuery("SELECT *").
		WillReturnRows(sqlmock.NewRows([]string{"account"}))

	r := userRepository{}

	v := "test"
	u, err := r.FindByToken(v)
	if err != nil {
		t.Error(err)
	}
	if u.Account != "" {
		t.Errorf("actual: not exists, expected: exists [%s]", u.Account)
	}

	mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"account"}).AddRow("test"))

	u, err = r.FindByToken(v)
	if err != nil {
		t.Error(err)
	}
	if u.Account == "" {
		t.Errorf("actual: exists [%s], expected: not exists", v)
	}
}

func TestValidUser(t *testing.T) {
	r := userRepository{}

	if r.ValidUser(nil) {
		t.Errorf("Valid user is not nil")
	}

	u := entity.User{}
	if r.ValidUser(&u) {
		t.Errorf("Valid user is empty key: %v", u)
	}

	u.Account = "test"
	if r.ValidUser(&u) {
		t.Errorf("Valid user is disable: %v", u)
	}

	u.IsEnable = true
	if !r.ValidUser(&u) {
		t.Errorf("Invalid user: %v", u)
	}
}

func TestValidRole(t *testing.T) {
	r := userRepository{}

	role := "test"
	if r.ValidRole(role) {
		t.Errorf("Invalid role %s", role)
	}

	if !r.ValidRole(entity.RoleAdministrator) {
		t.Errorf("Invalid role")
	}

	if !r.ValidRole(entity.RoleGeneral) {
		t.Errorf("Invalid role")
	}
}

func TestMatchPassword(t *testing.T) {
	r := userRepository{}

	s := "testtest"
	d, err := r.hashedPassword(s)
	if err != nil {
		t.Error(err)
	}

	if err := r.MatchPassword(d, "testtest1"); err == nil {
		t.Errorf("MatchPassword test failed: %s", d)
	}

	if err := r.MatchPassword(d, s); err != nil {
		t.Error(err)
	}
}

func TestCreateUser(t *testing.T) {
	mock.ExpectExec("INSERT INTO").
		WillReturnResult(sqlmock.NewResult(1, 1))

	r := userRepository{}

	u := entity.User{}
	if pass, err := r.Create(&u); err != nil {
		t.Error(err)
	} else if pass == "" {
		t.Errorf("Generated password is empty")
	}

	if u.Role == "" {
		t.Errorf("Role is empty")
	}
	if u.Role != entity.RoleGeneral {
		t.Errorf("Role is not matched, actual: %s, expected: %s", entity.RoleGeneral, u.Role)
	}
	if !u.IsEnable {
		t.Errorf("User is disable")
	}

	mock.ExpectExec("INSERT INTO").
		WillReturnResult(sqlmock.NewResult(1, 1))

	u = entity.User{
		Role: entity.RoleAdministrator,
	}
	if pass, err := r.Create(&u); err != nil {
		t.Error(err)
	} else if pass == "" {
		t.Errorf("Generated password is empty")
	}

	if u.Role != entity.RoleAdministrator {
		t.Errorf("Role is not matched, actual: %s, expected: %s", entity.RoleAdministrator, u.Role)
	}
}

func TestUpdatePassword(t *testing.T) {
	mock.ExpectExec("UPDATE").
		WillReturnResult(sqlmock.NewResult(1, 1))

	r := userRepository{}

	np := "newpassword"
	u := entity.User{
		ID: 1,
	}
	if err := r.UpdatePassword(&u, np); err != nil {
		t.Error(err)
	}
	if err := r.MatchPassword(u.Password, np); err != nil {
		t.Error(err)
	}

	if !u.IsActive {
		t.Errorf("User is not active")
	}
}

func TestUpdateAuthed(t *testing.T) {
	mock.ExpectExec("UPDATE").
		WillReturnResult(sqlmock.NewResult(1, 1))

	r := userRepository{}

	s := time.Now()
	u := entity.User{
		ID:         1,
		LastLogged: &s,
	}
	if err := r.UpdateAuthed(&u); err != nil {
		t.Error(err)
	}

	if s.Equal(*u.LastLogged) {
		t.Errorf("LastLogged not updated")
	}
}
