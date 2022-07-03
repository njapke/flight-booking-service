package database

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRawPutGet(t *testing.T) {
	db := New()
	val := []byte("test")
	err := db.rawPut("test", "key", val)
	require.NoError(t, err)

	res, err := db.rawGet("test", "key")
	require.NoError(t, err)
	require.Equal(t, val, res.Value)
}

func TestPutGetUser(t *testing.T) {
	db := New()
	u := &User{
		ID: "123", Name: "test", Email: "test@test.com",
	}
	err := db.Put(u)
	require.NoError(t, err)

	var res User
	err = db.Get(u.Key(), &res)
	require.NoError(t, err)

	require.Equal(t, u, &res)
}
