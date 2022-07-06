package database

import (
	"fmt"
	"testing"

	"github.com/christophwitzko/flight-booking-service/pkg/database/models"
	"github.com/stretchr/testify/require"
)

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

func TestValuesWithPrefixes(t *testing.T) {
	db, err := New()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, db.Close())
	}()

	expectedValues := make([]Model, 0)
	for i := 0; i < 3; i++ {
		u := &models.Flight{ID: fmt.Sprintf("%d", i)}
		require.NoError(t, db.Put(u))
		for s := 0; s < 100; s++ {
			seat := &models.Seat{FlightID: u.ID, Seat: fmt.Sprintf("%d", s)}
			require.NoError(t, db.Put(seat))
			if i == 2 {
				expectedValues = append(expectedValues, seat)
			}
		}
	}

	for _, p := range []string{"0", "1"} {
		values, err := db.Values(&models.Seat{}, p)
		require.NoError(t, err)
		require.Len(t, values, 100)
	}

	values, err := db.Values(&models.Seat{}, "2")
	require.NoError(t, err)

	require.ElementsMatch(t, expectedValues, values)
}
