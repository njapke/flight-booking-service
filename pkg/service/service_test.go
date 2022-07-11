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

func sendRequest(s http.Handler, method, path string, body io.Reader, modReqFns ...func(req *http.Request)) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	for _, f := range modReqFns {
		f(req)
	}
	rr := httptest.NewRecorder()
	s.ServeHTTP(rr, req)
	return rr
}

var testUser = []string{"user", "pw"}

func setBasicAuth(req *http.Request) {
	req.SetBasicAuth(testUser[0], testUser[1])
}

func initService(t *testing.T) *Service {
	db, err := database.New()
	require.NoError(t, err)
	require.NoError(t, seeder.Seed(db))
	s := New(logger.NewNop(), db)
	s.Auth[testUser[0]] = testUser[1]
	return s
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

func putBookingRequestData(s *Service) error {
	seats := []database.Model{
		&models.Flight{ID: "123", From: "AAA", To: "BBB", Status: "test"},
		&models.Seat{FlightID: "123", Seat: "A1", Row: 1, Price: 10, Available: false},
		&models.Seat{FlightID: "123", Seat: "B1", Row: 1, Price: 10, Available: true},
		&models.Seat{FlightID: "123", Seat: "C1", Row: 1, Price: 10, Available: true},
		&models.Seat{FlightID: "123", Seat: "F3", Row: 3, Price: 10, Available: false},
	}
	return s.db.Put(seats...)
}

func TestCreateBooking(t *testing.T) {
	s := initService(t)
	defer func() {
		require.NoError(t, s.db.Close())
	}()
	require.NoError(t, putBookingRequestData(s))

	// check if seats are correctly stored in database
	res := sendRequest(s, "GET", "/flights/123/seats", nil)
	require.Equal(t, http.StatusOK, res.Code)
	var resSeats []*models.Seat
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &resSeats))
	require.Len(t, resSeats, 2)

	bookingRequest := &models.Booking{
		FlightID: "123",
		Passengers: []models.Passenger{
			{Name: "John", Seat: "B1"},
			{Name: "Jane", Seat: "C1"},
		},
	}
	buf := &bytes.Buffer{}
	require.NoError(t, json.NewEncoder(buf).Encode(bookingRequest))
	res = sendRequest(s, "POST", "/bookings", buf, setBasicAuth)
	require.Equal(t, http.StatusOK, res.Code)

	var bookingResponse models.Booking
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &bookingResponse))
	require.Equal(t, bookingRequest.FlightID, bookingResponse.FlightID)
	require.Equal(t, bookingResponse.Passengers, bookingResponse.Passengers)
	require.Equal(t, 20, bookingResponse.Price)
	require.Equal(t, "confirmed", bookingResponse.Status)

	var bookings []*models.Booking
	res = sendRequest(s, "GET", "/bookings", nil, setBasicAuth)
	require.Equal(t, http.StatusOK, res.Code)
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &bookings))
	require.Len(t, bookings, 1)
	require.Equal(t, bookingResponse, *bookings[0])
}

func TestNotFound(t *testing.T) {
	s := initService(t)
	defer func() {
		require.NoError(t, s.db.Close())
	}()
	res := sendRequest(s, "GET", "/not-found", nil)
	require.Equal(t, http.StatusNotFound, res.Code)
}

func TestCreateInvalidBooking(t *testing.T) {
	s := initService(t)
	defer func() {
		require.NoError(t, s.db.Close())
	}()
	require.NoError(t, putBookingRequestData(s))

	bookingRequest := &models.Booking{
		FlightID: "123",
		Passengers: []models.Passenger{
			{Name: "John", Seat: "A1"},
			{Name: "Jane", Seat: "B1"},
		},
	}
	buf := &bytes.Buffer{}
	require.NoError(t, json.NewEncoder(buf).Encode(bookingRequest))
	res := sendRequest(s, "POST", "/bookings", buf, setBasicAuth)
	require.Equal(t, http.StatusBadRequest, res.Code)
}

func BenchmarkFlights(b *testing.B) {
	db, _ := database.New()
	_ = seeder.Seed(db)
	s := New(logger.NewNop(), db)

	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/flights", nil)
	request.Header.Set("X-Forwarded-For", "127.0.0.1")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.ServeHTTP(responseRecorder, request)
	}
}
