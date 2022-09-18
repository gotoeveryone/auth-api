package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateUnmarshalJSON(t *testing.T) {
	d := Date{}
	assert.NotNil(t, d.UnmarshalJSON([]byte("\"hogefuga\"")))
	assert.NotNil(t, d.UnmarshalJSON([]byte("2022-09-01")))
	assert.Nil(t, d.UnmarshalJSON([]byte("\"2022-09-01\"")))
	tm, _ := time.Parse("2006-01-02", "2022-09-01")
	assert.Equal(t, tm, d.Time)
}

func TestDateMarshalJSON(t *testing.T) {
	tm, _ := time.Parse("2006-01-02", "2022-09-01")
	d := Date{Time: tm}
	b, err := d.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, "\"2022-09-01\"", string(b))
}

func TestDateScan(t *testing.T) {
	tm, _ := time.Parse("2006-01-02", "2022-09-01")
	d := Date{}
	assert.Nil(t, d.Scan(tm))
	assert.Equal(t, tm, d.Time)
}

func TestDateValue(t *testing.T) {
	tm, _ := time.Parse("2006-01-02", "2022-09-01")
	d := Date{Time: tm}
	v, err := d.Value()
	assert.Nil(t, err)
	assert.Equal(t, tm, v)
}

func TestValidUser(t *testing.T) {
	u := User{}
	assert.False(t, u.Valid())
	u.Account = "testuser"
	assert.False(t, u.Valid())
	u.IsEnable = true
	assert.True(t, u.Valid())
}

func TestUserDefaultRole(t *testing.T) {
	u := User{}
	assert.Equal(t, u.DefaultRole(), RoleGeneral)
}
