package database

import (
	"fmt"
	"testing"

	"github.com/christophwitzko/flight-booking-service/pkg/database/models"
	"github.com/stretchr/testify/require"
)

func TestRawPutGet(t *testing.T) {
	db, err := New()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, db.Close())
	}()

	val := []byte("test")
	err = db.rawPut("test", "key", val)
	require.NoError(t, err)

	res, err := db.rawGet("test", "key")
	require.NoError(t, err)
	require.Equal(t, val, res)
}

func TestPutGet(t *testing.T) {
	db, err := New()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, db.Close())
	}()

	u := &models.Flight{
		ID:     "123",
		From:   "AAA",
		To:     "BBB",
		Status: "test",
	}
	err = db.Put(u)
	require.NoError(t, err)

	var res models.Flight
	err = db.Get(u.Key(), &res)
	require.NoError(t, err)

	require.Equal(t, u, &res)
}

func TestValues(t *testing.T) {
	db, err := New()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, db.Close())
	}()

	expectedValues := make([]Model, 0)
	for i := 0; i < 100; i++ {
		u := &models.Flight{ID: fmt.Sprintf("%d", i)}
		require.NoError(t, db.Put(u))
		expectedValues = append(expectedValues, u)
	}

	values, err := db.Values(&models.Flight{})
	require.NoError(t, err)

	require.ElementsMatch(t, expectedValues, values)
}
