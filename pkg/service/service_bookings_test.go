package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/christophwitzko/flight-booking-service/pkg/database"
	"github.com/christophwitzko/flight-booking-service/pkg/database/models"
	"github.com/christophwitzko/flight-booking-service/pkg/database/seeder"
	"github.com/christophwitzko/flight-booking-service/pkg/logger"
)

func BenchmarkHandlerGetBookings(b *testing.B) {
	db, _ := database.New()
	_ = seeder.Seed(db)
	s := New(logger.NewNop(), db)

	for i := 0; i < 100; i++ {
		_ = db.Put(&models.Booking{UserID: "user", ID: fmt.Sprintf("%d", i)})
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.SetBasicAuth("user", "password")

	resWriter := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.handlerGetBookings(resWriter, req)
	}
}

func BenchmarkHandlerCreateBooking(b *testing.B) {
	db, _ := database.New()
	_ = seeder.Seed(db)
	s := New(logger.NewNop(), db)

	amountOfBookingRequests := 300000
	bookingRequests := make([]io.ReadCloser, amountOfBookingRequests)
	for i := 0; i < amountOfBookingRequests; i++ {
		flightId := strconv.Itoa(i)
		_ = db.Put(&models.Flight{ID: flightId}, &models.Seat{FlightID: flightId, Seat: "A1", Available: true})
		bookingRequest := &models.Booking{
			FlightID:   flightId,
			Passengers: []models.Passenger{{Name: "user", Seat: "A1"}},
		}
		payload, _ := json.Marshal(bookingRequest)
		bookingRequests[i] = io.NopCloser(bytes.NewReader(payload))
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.SetBasicAuth("user", "password")
	req.ContentLength = -1

	resWriter := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.Body = bookingRequests[i]
		s.handlerCreateBooking(resWriter, req)
	}
}
