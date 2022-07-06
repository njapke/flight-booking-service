package main

import (
	"context"
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
	log := logger.New()
	if err := run(log); err != nil {
		log.Fatal(err)
	}
}

func run(log *logger.Logger) error {
	defer log.Sync()

	db, err := database.New()
	if err != nil {
		return err
	}
	if err := seeder.Seed(db); err != nil {
		return err
	}
	srv := &http.Server{
		Addr:    "127.0.0.1:3000",
		Handler: service.New(log, db),
	}
	go func() {
		log.Infof("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Error(err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	<-ctx.Done()
	stop()

	log.Info("stopping server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err == context.DeadlineExceeded {
		log.Info("closing server...")
		if err := srv.Close(); err != nil {
			log.Error(err)
		}
		// finishing pending database writes
		<-time.After(time.Second)
	} else if err != nil {
		log.Error(err)
	}

	log.Info("closing database...")
	if err := db.Close(); err != nil {
		return err
	}
	return nil
}
