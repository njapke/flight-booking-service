package service

import (
	"encoding/json"
	"net/http"

	"github.com/christophwitzko/flight-booking-service/pkg/database"
	"github.com/christophwitzko/flight-booking-service/pkg/database/models"
	"github.com/christophwitzko/flight-booking-service/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Service struct {
	router *chi.Mux
	log    *logger.Logger
	db     *database.Database
}

func New(logger *logger.Logger, db *database.Database) *Service {
	svc := &Service{
		router: chi.NewRouter(),
		log:    logger,
		db:     db,
	}
	svc.setupMiddleware()
	svc.setupRoutes()
	return svc
}

func (s *Service) setupMiddleware() {
	s.router.Use(middleware.CleanPath)
	s.router.Use(middleware.RequestID)
	s.router.Use(s.log.Middleware)
	s.router.Use(s.recoverMiddleware)
}

func (s *Service) sendError(w http.ResponseWriter, err string, code int) {
	s.log.Errorf("error(code=%d): %s", code, err)
	w.WriteHeader(code)
	s.writeJSON(w, map[string]string{"error": err})
}

func (s *Service) writeJSON(w http.ResponseWriter, d any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(w).Encode(d)
	if err != nil {
		s.log.Errorf("json write error: %v", err)
	}
}

func (s *Service) setupRoutes() {
	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		s.writeJSON(w, map[string]string{"service": "flight-booking-service"})
	})
	s.router.Get("/flights", func(w http.ResponseWriter, r *http.Request) {
		flights, err := s.db.Values(&models.Flight{})
		if err != nil {
			s.sendError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.writeJSON(w, flights)
	})
	s.router.Get("/flights/{id}/seats", func(w http.ResponseWriter, r *http.Request) {

		flightId := chi.URLParam(r, "id")
		allSeats, err := s.db.Values(&models.Seat{})
		if err != nil {
			s.sendError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		seats := make([]*models.Seat, 0)
		for _, seat := range allSeats {
			seat := seat.(*models.Seat)
			if seat.Available && seat.FlightID == flightId {
				seats = append(seats, seat)
			}
		}
		s.writeJSON(w, seats)
	})
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
