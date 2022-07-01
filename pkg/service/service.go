package service

import (
	"fmt"
	"net/http"

	"github.com/christophwitzko/flight-booking-service/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Service struct {
	Router *chi.Mux
	Logger *logger.Logger
}

func New(logger *logger.Logger) *Service {
	svc := &Service{
		Router: chi.NewRouter(),
		Logger: logger}
	svc.setupMiddleware()
	svc.setupRoutes()
	return svc
}

func (s *Service) setupMiddleware() {
	s.Router.Use(middleware.CleanPath)
	s.Router.Use(middleware.RequestID)
	s.Router.Use(s.Logger.Middleware)
	s.Router.Use(middleware.Recoverer)
}

func (s *Service) setupRoutes() {
	s.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	s.Router.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("oops")
	})
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}
