package service

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/christophwitzko/flight-booking-service/pkg/database"
	"github.com/christophwitzko/flight-booking-service/pkg/database/models"
	"github.com/christophwitzko/flight-booking-service/pkg/logger"
	"github.com/dgraph-io/badger/v3"
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
		With(middleware.BasicAuth("auth", map[string]string{"user": "pw"})).
		Route("/bookings", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				s.writeJSON(w, "ok")
			})
		})
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
