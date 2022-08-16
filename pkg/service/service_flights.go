package service

import (
	"errors"
	"net/http"

	"github.com/christophwitzko/flight-booking-service/pkg/database"
	"github.com/christophwitzko/flight-booking-service/pkg/database/models"
	"github.com/dgraph-io/badger/v3"
	"github.com/go-chi/chi/v5"
)

func filterFlights(flights []*models.Flight, queryFrom, queryTo, queryStatus string) []*models.Flight {
	foundFlights := make([]*models.Flight, 0)
	for _, flight := range flights {
		if (queryFrom == "" || flight.From == queryFrom) &&
			(queryTo == "" || flight.To == queryTo) &&
			(queryStatus == "" || flight.Status == queryStatus) {
			foundFlights = append(foundFlights, flight)
		}
	}
	return foundFlights
}

func (s *Service) handlerGetFlights(w http.ResponseWriter, r *http.Request) {
	queryFrom := r.URL.Query().Get("from")
	queryTo := r.URL.Query().Get("to")
	queryStatus := r.URL.Query().Get("status")

	if queryFrom == "" && queryTo == "" && queryStatus == "" {
		s.contentTypeJSON(w)
		err := s.db.RawValues(w, "flights")
		if err != nil {
			s.log.Errorf("error getting flights: %v", err)
		}
		return
	}

	allFlights, err := database.Values[*models.Flight](s.db)
	if err != nil {
		s.sendError(w, "could not get flights", http.StatusInternalServerError)
		return
	}

	foundFlights := filterFlights(allFlights, queryFrom, queryTo, queryStatus)
	if len(foundFlights) == 0 {
		s.sendError(w, "no flights found", http.StatusBadRequest)
		return
	}
	s.writeJSON(w, foundFlights)
}

func (s *Service) handlerGetDestinations(w http.ResponseWriter, r *http.Request) {
	allFlights, err := database.Values[*models.Flight](s.db)
	if err != nil {
		s.sendError(w, "could not get flights", http.StatusInternalServerError)
		return
	}

	from := make(map[string]bool)
	to := make(map[string]bool)
	for _, flight := range allFlights {
		from[flight.From] = true
		to[flight.To] = true
	}

	ret := struct {
		From []string `json:"from"`
		To   []string `json:"to"`
	}{
		From: make([]string, 0, len(from)),
		To:   make([]string, 0, len(to)),
	}
	for k := range from {
		ret.From = append(ret.From, k)
	}
	for k := range to {
		ret.To = append(ret.To, k)
	}
	s.writeJSON(w, ret)
}

func (s *Service) handlerGetFlight(w http.ResponseWriter, r *http.Request) {
	flightID := chi.URLParam(r, "id")
	flightData, err := s.db.RawGet("flights", flightID)
	if err == nil {
		s.contentTypeJSON(w)
		if _, err = w.Write(flightData); err != nil {
			s.log.Errorf("write error: %v", err)
		}
		return
	} else if errors.Is(err, badger.ErrKeyNotFound) {
		s.sendError(w, "flight not found", http.StatusNotFound)
		return
	}
	s.sendError(w, err.Error(), http.StatusInternalServerError)
}

func (s *Service) handlerGetFlightSeats(w http.ResponseWriter, r *http.Request) {
	flightID := chi.URLParam(r, "id")
	allSeats, err := database.Values[*models.Seat](s.db, flightID)
	if err != nil {
		s.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	availableSeats := make([]*models.Seat, 0)
	for _, seat := range allSeats {
		if seat.Available {
			availableSeats = append(availableSeats, seat)
		}
	}
	if len(availableSeats) == 0 {
		s.sendError(w, "no seats available", http.StatusNotFound)
		return
	}
	s.writeJSON(w, availableSeats)
}
