package database

import (
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"gorm.io/gorm"
)

func TestNewUserRepository(t *testing.T) {
	a := reflect.TypeOf(&userRepository{})
	e := reflect.TypeOf(NewUserRepository())
	if a != e {
		t.Errorf("NewTokenRepository type is mismatch, actual: %s, expected: %s", a, e)
	}
}

func TestExists(t *testing.T) {
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users`")).
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

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users`")).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	e, err = r.Exists(v)
	if err != nil {
		t.Error(err)
	}
	if e {
		t.Errorf("User %s is exists", v)
	}
}

func TestFindUser(t *testing.T) {
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	r := userRepository{}

	id := uint(1)
	e := entity.User{}
	if err := r.Find(id, &e); err != nil {
		if err != gorm.ErrRecordNotFound {
			t.Error(err)
		}
	}
	if e.ID != 0 {
		t.Errorf("actual: exists, expected: not exists [%d]", e.ID)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	if err := r.Find(id, &e); err != nil {
		t.Error(err)
	}
	if e.ID != 1 {
		t.Errorf("actual: not exists, expected: exists [%d]", e.ID)
	}
}

func TestFindByAccount(t *testing.T) {
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
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

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
		WillReturnRows(sqlmock.NewRows([]string{"account"}).AddRow("test"))

	u, err = r.FindByAccount(v)
	if err != nil {
		t.Error(err)
	}
	if u.Account == "" {
		t.Errorf("actual: exists [%s], expected: not exists", v)
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
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
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

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
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
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users`")).
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
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users`")).
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
