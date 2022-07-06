package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/christophwitzko/flight-booking-service/pkg/database"
	"github.com/christophwitzko/flight-booking-service/pkg/database/models"
	"github.com/christophwitzko/flight-booking-service/pkg/database/seeder"
	"github.com/christophwitzko/flight-booking-service/pkg/logger"
	"github.com/stretchr/testify/require"
)

func sendRequest(s http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	rr := httptest.NewRecorder()
	s.ServeHTTP(rr, req)
	return rr
}

func initService(t *testing.T) *Service {
	db, err := database.New()
	require.NoError(t, err)
	require.NoError(t, seeder.Seed(db))
	return New(logger.NewNop(), db)
}

func TestIndex(t *testing.T) {
	s := initService(t)
	defer func() {
		require.NoError(t, s.db.Close())
	}()
	res := sendRequest(s, "GET", "/", nil)
	require.Equal(t, http.StatusOK, res.Code)

	res = sendRequest(s, "POST", "/", bytes.NewReader([]byte("{}")))
	require.Equal(t, http.StatusMethodNotAllowed, res.Code)
}

func TestGetFlights(t *testing.T) {
	s := initService(t)
	defer func() {
		require.NoError(t, s.db.Close())
	}()
	res := sendRequest(s, "GET", "/flights", nil)
	require.Equal(t, http.StatusOK, res.Code)
	var flights []*models.Flight
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &flights))
	require.Len(t, flights, 100)
}

func TestGetFlight(t *testing.T) {
	s := initService(t)
	defer func() {
		require.NoError(t, s.db.Close())
	}()

	flight := &models.Flight{ID: "123", From: "AAA", To: "BBB", Status: "test"}
	require.NoError(t, s.db.Put(flight))

	res := sendRequest(s, "GET", "/flights/123", nil)
	require.Equal(t, http.StatusOK, res.Code)
	var flightRes models.Flight
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &flightRes))
	require.Equal(t, flight, &flightRes)
}

func BenchmarkFlights(b *testing.B) {
	db, _ := database.New()
	_ = seeder.Seed(db)
	s := New(logger.NewNop(), db)

	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/flights", nil)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.ServeHTTP(responseRecorder, request)
	}
}
