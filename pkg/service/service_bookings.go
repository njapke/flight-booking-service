package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/christophwitzko/flight-booking-service/pkg/database"
	"github.com/christophwitzko/flight-booking-service/pkg/database/models"
	"github.com/google/uuid"
)

func (s *Service) handlerGetBookings(w http.ResponseWriter, r *http.Request) {
	user, _, _ := r.BasicAuth()
	bookings, err := s.db.Values(&models.Booking{}, user)
	if err != nil {
		s.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, bookings)
}

func (s *Service) handlerCreateBooking(w http.ResponseWriter, r *http.Request) {
	userID, _, _ := r.BasicAuth()
	var bookingRequest models.Booking
	if err := json.NewDecoder(r.Body).Decode(&bookingRequest); err != nil {
		s.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(bookingRequest.Passengers) == 0 {
		s.sendError(w, "no passengers", http.StatusBadRequest)
		return
	}
	var flight models.Flight
	if err := s.db.Get(bookingRequest.FlightID, &flight); err != nil {
		s.log.Warnf("could not get flight: %v", err)
		s.sendError(w, "could not find flight", http.StatusBadRequest)
		return
	}

	price := 0
	updates := make([]database.Model, len(bookingRequest.Passengers)+1)
	for i, passenger := range bookingRequest.Passengers {
		var seat models.Seat
		key := fmt.Sprintf("%s/%s", flight.ID, passenger.Seat)
		if err := s.db.Get(key, &seat); err != nil {
			s.log.Warnf("could not find flight: %v", err)
			s.sendError(w, "could not find seat", http.StatusBadRequest)
			return
		}
		if !seat.Available {
			s.sendError(w, "seat not available", http.StatusBadRequest)
			return
		}
		price += seat.Price
		seat.Available = false
		updates[i] = &seat
	}

	booking := &models.Booking{
		ID:         uuid.NewString(),
		UserID:     userID,
		FlightID:   flight.ID,
		Price:      price,
		Status:     "confirmed",
		Passengers: bookingRequest.Passengers,
	}
	updates[len(updates)-1] = booking

	if err := s.db.Put(updates...); err != nil {
		s.sendError(w, err.Error(), http.StatusInternalServerError)
	}

	s.writeJSON(w, booking)
}
