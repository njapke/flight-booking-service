package service

import (
	"errors"
	"net/http"

	"github.com/christophwitzko/flight-booking-service/pkg/database/models"
	"github.com/dgraph-io/badger/v3"
	"github.com/go-chi/chi/v5"
)

func (s *Service) handlerGetFlights(w http.ResponseWriter, r *http.Request) {
	s.contentTypeJSON(w)
	err := s.db.RawValues(w, "flights")
	if err != nil {
		s.log.Errorf("error getting flights: %v", err)
	}
}

func (s *Service) handlerGetFlight(w http.ResponseWriter, r *http.Request) {
	flightId := chi.URLParam(r, "id")
	flightData, err := s.db.RawGet("flights", flightId)
	if err == nil {
		s.contentTypeJSON(w)
		if _, err := w.Write(flightData); err != nil {
			s.log.Errorf("write error: %v", err)
		}
	} else if errors.Is(err, badger.ErrKeyNotFound) {
		s.sendError(w, "flight not found", http.StatusNotFound)
	} else {
		s.sendError(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Service) handlerGetFlightSeats(w http.ResponseWriter, r *http.Request) {
	flightId := chi.URLParam(r, "id")
	allSeats, err := s.db.Values(&models.Seat{}, flightId)
	if err != nil {
		s.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	seats := make([]*models.Seat, 0)
	for _, seat := range allSeats {
		seat := seat.(*models.Seat)
		if seat.Available {
			seats = append(seats, seat)
		}
	}
	if len(seats) == 0 {
		s.sendError(w, "no seats available", http.StatusNotFound)
		return
	}
	s.writeJSON(w, seats)
}
