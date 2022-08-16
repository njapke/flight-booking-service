package service

import (
	"encoding/json"
	"net/http"

	"github.com/christophwitzko/flight-booking-service/pkg/database"
	"github.com/christophwitzko/flight-booking-service/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
		Auth:   make(map[string]string),
	}
	svc.setupMiddleware()
	svc.setupRoutes()
	return svc
}

func (s *Service) setupMiddleware() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Compress(5))
	s.router.Use(s.log.Middleware)
	s.router.Use(s.recoverMiddleware)

	s.router.Mount("/debug", middleware.Profiler())
}

func (s *Service) sendError(w http.ResponseWriter, err string, code int) {
	s.log.Warnf("error(code=%d): %s", code, err)
	w.WriteHeader(code)
	s.writeJSON(w, map[string]string{"error": err})
}

func (s *Service) contentTypeJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func (s *Service) writeJSON(w http.ResponseWriter, d any) {
	s.contentTypeJSON(w)
	err := json.NewEncoder(w).Encode(d)
	if err != nil {
		s.log.Errorf("json write error: %v", err)
	}
}

func (s *Service) handlerNotFound(w http.ResponseWriter, r *http.Request) {
	s.sendError(w, "not found", http.StatusNotFound)
}

func (s *Service) handlerIndex(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]string{"service": "flight-booking-service"})
}

func (s *Service) setupRoutes() {
	s.router.NotFound(s.handlerNotFound)

	s.router.Get("/", s.handlerIndex)

	s.router.
		With(middleware.CleanPath).
		Route("/flights", func(r chi.Router) {
			r.Get("/", s.handlerGetFlights)
			r.Get("/{id}", s.handlerGetFlight)
			r.Get("/{id}/seats", s.handlerGetFlightSeats)
		})

	s.router.Get("/destinations", s.handlerGetDestinations)

	s.router.
		With(middleware.BasicAuth("auth", s.Auth)).
		Route("/bookings", func(r chi.Router) {
			r.Get("/", s.handlerGetBookings)
			r.Post("/", s.handlerCreateBooking)
		})
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
