package database

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users`")).
		WithArgs("test").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	r := userRepository{}

	v := "test"
	e, err := r.Exists(v)
	assert.Nil(t, err)
	assert.True(t, e)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users`")).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	e, err = r.Exists(v)
	assert.Nil(t, err)
	assert.False(t, e)
}

func TestFindUser(t *testing.T) {
	r := userRepository{}
	id := uint(1)

	{
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		e, err := r.Find(id)
		assert.Nil(t, err)
		assert.Nil(t, e)
	}
	{
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		e, err := r.Find(id)
		assert.Nil(t, err)
		assert.NotNil(t, e)
	}
}

func TestFindByAccount(t *testing.T) {
	r := userRepository{}
	v := "test"

	{
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
			WillReturnRows(sqlmock.NewRows([]string{"account"}))

		u, err := r.FindByAccount(v)
		assert.Nil(t, err)
		assert.Nil(t, u)
	}
	{
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
			WillReturnRows(sqlmock.NewRows([]string{"account"}).AddRow("test"))

		u, err := r.FindByAccount(v)
		assert.Nil(t, err)
		assert.NotNil(t, u)
	}
}

func TestMatchPassword(t *testing.T) {
	r := userRepository{}

	s := "testtest"
	d, err := r.hashedPassword(s)
	assert.Nil(t, err)
	assert.NotNil(t, r.MatchPassword(d, "testtest1"))
	assert.Nil(t, r.MatchPassword(d, s))
}

func TestCreateUser(t *testing.T) {
	r := userRepository{}

	{
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
			WillReturnResult(sqlmock.NewResult(1, 1))

		u := entity.User{}
		pass, err := r.Create(&u)
		assert.Nil(t, err)
		assert.NotEmpty(t, pass)
		assert.Equal(t, u.Role, entity.RoleGeneral)
		assert.True(t, u.IsEnable)
	}

	{
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
			WillReturnResult(sqlmock.NewResult(1, 1))

		u := entity.User{
			Role: entity.RoleAdministrator,
		}
		pass, err := r.Create(&u)
		assert.Nil(t, err)
		assert.NotEmpty(t, pass)
		assert.Equal(t, u.Role, entity.RoleAdministrator)
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
	assert.Nil(t, r.UpdatePassword(&u, np))
	assert.Nil(t, r.MatchPassword(u.Password, np))
	assert.True(t, u.IsActive)
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
	assert.Nil(t, r.UpdateAuthed(&u))
	assert.False(t, s.Equal(*u.LastLogged))
}
