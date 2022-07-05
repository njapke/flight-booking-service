package seeder

import (
	"fmt"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/christophwitzko/flight-booking-service/pkg/database"
	"github.com/christophwitzko/flight-booking-service/pkg/database/models"
)

func generateSeats(flightId string) []*models.Seat {
	seats := make([]*models.Seat, 174)
	i := 0
	for row := 1; row <= 30; row++ {
		if row == 13 {
			continue
		}
		for _, seat := range []string{"A", "B", "C", "D", "E", "F"} {
			seats[i] = &models.Seat{
				ID:        fmt.Sprintf("%s/%d", flightId, i),
				FlightID:  flightId,
				Seat:      fmt.Sprintf("%d%s", row, seat),
				Row:       row,
				Price:     gofakeit.IntRange(20, 500),
				Available: gofakeit.Bool(),
			}
			i++
		}
	}
	return seats
}

func generateFlight() (*models.Flight, []*models.Seat) {
	startTime := gofakeit.DateRange(time.Now(), time.Now().Add(time.Hour*48))
	randomFlightDuration := time.Duration(gofakeit.IntRange(30, 300)) * time.Minute
	flight := &models.Flight{
		ID:        gofakeit.UUID(),
		From:      strings.ToUpper(gofakeit.LetterN(3)),
		To:        strings.ToUpper(gofakeit.LetterN(3)),
		Departure: startTime,
		Arrival:   startTime.Add(randomFlightDuration),
		Status:    gofakeit.RandomString([]string{"scheduled", "cancelled", "delayed"}),
	}

	return flight, generateSeats(flight.ID)
}

func Seed(db *database.Database) error {
	gofakeit.Seed(999)
	for i := 0; i < 100; i++ {
		f, seats := generateFlight()
		if err := db.Put(f); err != nil {
			return err
		}
		for _, s := range seats {
			if err := db.Put(s); err != nil {
				return err
			}
		}
	}
	return nil
}
