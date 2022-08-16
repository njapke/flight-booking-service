package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/christophwitzko/flight-booking-service/pkg/database"
	"github.com/christophwitzko/flight-booking-service/pkg/database/seeder"
	"github.com/christophwitzko/flight-booking-service/pkg/logger"
	"github.com/christophwitzko/flight-booking-service/pkg/service"
)

func main() {
	level := logger.DebugLevel
	levelName := os.Getenv("LOG_LEVEL")
	switch levelName {
	case "debug":
		level = logger.DebugLevel
	case "info":
		level = logger.InfoLevel
	case "warn":
		level = logger.WarnLevel
	case "error":
		level = logger.ErrorLevel
	}
	log := logger.New(level)
	if levelName != "" {
		log.Infof("log level: %s", levelName)
	}
	if err := run(log); err != nil {
		log.Fatal(err)
	}
}

func getBindAddress() string {
	if bindAddress := os.Getenv("BIND_ADDRESS"); bindAddress != "" {
		return bindAddress
	}
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return "127.0.0.1:3000"
}

func run(log *logger.Logger) error {
	db, err := database.New()
	if err != nil {
		return err
	}
	err = seeder.Seed(db)
	if err != nil {
		return err
	}

	s := service.New(log, db)
	s.Auth["user"] = "pw"

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
		Addr:         getBindAddress(),
		Handler:      s,
	}

	listenErrCh := make(chan error)
	go func() {
		log.Infof("listening on %s", srv.Addr)
		sErr := srv.ListenAndServe()
		if !errors.Is(sErr, http.ErrServerClosed) {
			listenErrCh <- sErr
		}
		close(listenErrCh)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Info("shutting down...")
	case err = <-listenErrCh:
		log.Errorf("error listening on %s: %s", srv.Addr, err)
	}
	stop()

	log.Info("stopping server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); errors.Is(err, context.DeadlineExceeded) {
		log.Info("closing server...")
		if err = srv.Close(); err != nil {
			log.Error(err)
		}
		// finishing pending database writes
		<-time.After(time.Second)
	} else if err != nil {
		log.Error(err)
	}

	log.Info("closing database...")
	return db.Close()
}
