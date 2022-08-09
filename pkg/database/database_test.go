package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/christophwitzko/flight-booking-service/pkg/database/models"
	"github.com/stretchr/testify/require"
)

func TestRawGet(t *testing.T) {
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
	resData, err := db.RawGet(res.Collection(), "123")
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(resData, &res))
	require.Equal(t, u, &res)
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

	res2, err := Get[*models.Flight](db, u.Key())
	require.NoError(t, err)
	require.Equal(t, u, res2)
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

	values2, err := Values[*models.Flight](db)
	require.NoError(t, err)
	require.ElementsMatch(t, expectedValues, values2)
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
		values, vErr := db.Values(&models.Seat{}, p)
		require.NoError(t, vErr)
		require.Len(t, values, 100)
	}

	values, err := db.Values(&models.Seat{}, "2")
	require.NoError(t, err)

	require.ElementsMatch(t, expectedValues, values)
}

func TestValuesRawValues(t *testing.T) {
	db, err := New()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, db.Close())
	}()

	originalFlights := []Model{&models.Flight{ID: "A"}, &models.Flight{ID: "B"}, &models.Flight{ID: "C"}}
	require.NoError(t, db.Put(originalFlights...))

	buf := &bytes.Buffer{}
	require.NoError(t, db.RawValues(buf, "flights"))

	var flights []*models.Flight
	require.NoError(t, json.NewDecoder(buf).Decode(&flights))
	require.ElementsMatch(t, originalFlights, flights)
}

func BenchmarkPut(b *testing.B) {
	db, _ := New()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = db.Put(&models.Flight{ID: "123"})
	}
}

func BenchmarkGet(b *testing.B) {
	db, _ := New()
	_ = db.Put(&models.Flight{ID: "123"})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var flight models.Flight
		_ = db.Get("123", &flight)
	}
}

func BenchmarkGetGenerics(b *testing.B) {
	db, _ := New()
	_ = db.Put(&models.Flight{ID: "123"})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Get[*models.Flight](db, "123")
	}
}

func BenchmarkRawGet(b *testing.B) {
	db, _ := New()
	_ = db.Put(&models.Flight{ID: "123"})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = db.RawGet("flights", "123")
	}
}

func BenchmarkValues(b *testing.B) {
	db, _ := New()
	for i := 0; i < 1000; i++ {
		_ = db.Put(&models.Flight{ID: fmt.Sprintf("%d", i)})
	}

	emptyFlight := &models.Flight{}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = db.Values(emptyFlight)
	}
}

func BenchmarkValuesGenerics(b *testing.B) {
	db, _ := New()
	for i := 0; i < 1000; i++ {
		_ = db.Put(&models.Flight{ID: fmt.Sprintf("%d", i)})
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Values[*models.Flight](db)
	}
}

func BenchmarkRawValues(b *testing.B) {
	db, _ := New()
	for i := 0; i < 1000; i++ {
		_ = db.Put(&models.Flight{ID: fmt.Sprintf("%d", i)})
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = db.RawValues(io.Discard, "flights")
	}
}
