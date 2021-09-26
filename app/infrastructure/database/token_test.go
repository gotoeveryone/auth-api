package database

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"gorm.io/gorm"
)

func TestNewTokenRepository(t *testing.T) {
	a := reflect.TypeOf(&tokenRepository{})
	e := reflect.TypeOf(NewTokenRepository())
	if a != e {
		t.Errorf("NewTokenRepository type is mismatch, actual: %s, expected: %s", a, e)
	}
}

func TestFindToken(t *testing.T) {
	mock.ExpectQuery("SELECT *").
		WillReturnRows(sqlmock.NewRows([]string{"token"}))

	r := tokenRepository{}

	v := "test"
	e := entity.Token{}
	if err := r.Find(v, &e); err != nil {
		if err != gorm.ErrRecordNotFound {
			t.Error(err)
		}
	}
	if e.Token != "" {
		t.Errorf("actual: exists, expected: not exists [%s]", e.Token)
	}

	mock.ExpectQuery("SELECT *").
		WillReturnRows(sqlmock.NewRows([]string{"token"}).AddRow("test"))

	if err := r.Find(v, &e); err != nil {
		t.Error(err)
	}
	if e.Token == "" {
		t.Errorf("actual: not exists, expected: exists [%s]", v)
	}
}

func TestCreateToken(t *testing.T) {
	mock.ExpectExec("INSERT INTO").
		WillReturnResult(sqlmock.NewResult(1, 1))

	r := tokenRepository{}

	u := entity.User{
		ID: 1,
	}
	e := entity.Token{}
	if err := r.Create(&u, &e); err != nil {
		t.Error(err)
	}
	if u.ID != e.UserID {
		t.Errorf("ID not matched, user: %d, token: %d", u.ID, e.UserID)
	}
	if gin.Mode() != e.Environment {
		t.Errorf("Environment not matched, user: %s, token: %s", gin.Mode(), e.Environment)
	}
	if e.Token == "" {
		t.Errorf("Token is empty")
	}
}

func TestDeleteToken(t *testing.T) {
	mock.ExpectExec("DELETE FROM").
		WillReturnResult(sqlmock.NewResult(1, 1))

	r := tokenRepository{}

	v := "test"
	if err := r.Delete(v); err != nil {
		t.Error(err)
	}
}

func TestDeleteExpiredToken(t *testing.T) {
	a := int64(10)
	mock.ExpectExec("DELETE FROM").
		WillReturnResult(sqlmock.NewResult(1, a))

	r := tokenRepository{}

	cnt, err := r.DeleteExpired()
	if err != nil {
		t.Error(err)
	}
	if cnt != a {
		t.Errorf("RowsAffected is not matched, actual: %d, expected: %d", a, cnt)
	}
}

func TestCanAutoDeleteExpired(t *testing.T) {
	r := tokenRepository{}
	if r.CanAutoDeleteExpired() {
		t.Errorf("TestCanAutoDeleteExpired failed")
	}
}
