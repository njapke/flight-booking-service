package seeder

import (
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/christophwitzko/flight-booking-service/pkg/database"
	"github.com/christophwitzko/flight-booking-service/pkg/database/models"
)

func generateFlight() *models.Flight {
	startTime := gofakeit.DateRange(time.Now(), time.Now().Add(time.Hour*48))
	randomFlightDuration := time.Duration(gofakeit.IntRange(30, 300)) * time.Minute
	return &models.Flight{
		ID:        gofakeit.UUID(),
		From:      strings.ToUpper(gofakeit.LetterN(3)),
		To:        strings.ToUpper(gofakeit.LetterN(3)),
		Departure: startTime,
		Arrival:   startTime.Add(randomFlightDuration),
		Status:    gofakeit.RandomString([]string{"scheduled", "cancelled", "delayed"}),
	}
}

func Seed(db *database.Database) error {
	gofakeit.Seed(99)
	for i := 0; i < 100; i++ {
		f := generateFlight()
		err := db.Put(f)
		if err != nil {
			return err
		}
	}
	return nil
}
