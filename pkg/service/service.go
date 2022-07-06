package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/christophwitzko/flight-booking-service/pkg/database"
	"github.com/christophwitzko/flight-booking-service/pkg/database/models"
	"github.com/christophwitzko/flight-booking-service/pkg/logger"
	"github.com/dgraph-io/badger/v3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type Service struct {
	router *chi.Mux
	log    *logger.Logger
	db     *database.Database
	Auth   map[string]string
}

func New(logger *logger.Logger, db *database.Database) *Service {
	svc := &Service{
		router: chi.NewRouter(),
		log:    logger,
		db:     db,
		Auth:   map[string]string{"user": "pw"},
	}
	svc.setupMiddleware()
	svc.setupRoutes()
	return svc
}

func (s *Service) setupMiddleware() {
	s.router.Use(middleware.CleanPath)
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Compress(5))
	s.router.Use(s.log.Middleware)
	s.router.Use(s.recoverMiddleware)
}

func (s *Service) sendError(w http.ResponseWriter, err string, code int) {
	s.log.Warnf("error(code=%d): %s", code, err)
	w.WriteHeader(code)
	s.writeJSON(w, map[string]string{"error": err})
}

func (s *Service) contentTypeJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func (s *Service) writeJSON(w http.ResponseWriter, d any) {
	s.contentTypeJson(w)
	err := json.NewEncoder(w).Encode(d)
	if err != nil {
		s.log.Errorf("json write error: %v", err)
	}
}

func (s *Service) setupRoutes() {
	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		s.sendError(w, "not found", http.StatusNotFound)
	})

	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		s.writeJSON(w, map[string]string{"service": "flight-booking-service"})
	})

	s.router.Route("/flights", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			flights, err := s.db.Values(&models.Flight{})
			if err != nil {
				s.sendError(w, err.Error(), http.StatusInternalServerError)
				return
			}
			s.writeJSON(w, flights)
		})
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			flightId := chi.URLParam(r, "id")
			flightData, err := s.db.RawGet("flights", flightId)
			if err == nil {
				s.contentTypeJson(w)
				if _, err := w.Write(flightData); err != nil {
					s.log.Errorf("write error: %v", err)
				}
			} else if errors.Is(err, badger.ErrKeyNotFound) {
				s.sendError(w, "flight not found", http.StatusNotFound)
			} else {
				s.sendError(w, err.Error(), http.StatusInternalServerError)
			}
		})
		r.Get("/{id}/seats", func(w http.ResponseWriter, r *http.Request) {
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
		})
	})

	s.router.
		With(middleware.BasicAuth("auth", s.Auth)).
		Route("/bookings", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				user, _, _ := r.BasicAuth()
				bookings, err := s.db.Values(&models.Booking{}, user)
				if err != nil {
					s.sendError(w, err.Error(), http.StatusInternalServerError)
					return
				}
				s.writeJSON(w, bookings)
			})
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				userId, _, _ := r.BasicAuth()
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
					s.sendError(w, "could not find flight", http.StatusNotFound)
					return
				}

				price := 0
				updates := make([]database.Model, len(bookingRequest.Passengers)+1)
				for i, passenger := range bookingRequest.Passengers {
					var seat models.Seat
					key := fmt.Sprintf("%s/%s", flight.ID, passenger.Seat)
					if err := s.db.Get(key, &seat); err != nil {
						s.log.Warnf("could not find flight: %v", err)
						s.sendError(w, "could not find seat", http.StatusNotFound)
						return
					}
					if !seat.Available {
						s.sendError(w, "seat not available", http.StatusNotFound)
						return
					}
					price += seat.Price
					seat.Available = false
					updates[i] = &seat
				}

				booking := &models.Booking{
					ID:         uuid.NewString(),
					UserID:     userId,
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
			})
		})
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
