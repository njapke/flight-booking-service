package seeder

import (
	"fmt"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/christophwitzko/flight-booking-service/pkg/database"
	"github.com/christophwitzko/flight-booking-service/pkg/database/models"
)

func generateSeats(flightID string, rows int) []*models.Seat {
	seats := make([]*models.Seat, rows*6)
	i := 0
	for row := 1; row <= rows+1; row++ {
		if row == 13 {
			continue
		}
		for _, seat := range []string{"A", "B", "C", "D", "E", "F"} {
			seats[i] = &models.Seat{
				FlightID:  flightID,
				Seat:      fmt.Sprintf("%d%s", row, seat),
				Row:       row,
				Price:     gofakeit.IntRange(20, 500),
				Available: true,
			}
			i++
		}
	}
	return seats
}

func generateFlight(rows int) (*models.Flight, []*models.Seat) {
	startTime := gofakeit.DateRange(time.Now(), time.Now().Add(time.Hour*48))
	randomFlightDuration := time.Duration(gofakeit.IntRange(30, 300)) * time.Minute
	flight := &models.Flight{
		ID:        gofakeit.UUID(),
		From:      strings.ToUpper(gofakeit.LetterN(3)),
		To:        strings.ToUpper(gofakeit.LetterN(3)),
		Departure: startTime,
		Arrival:   startTime.Add(randomFlightDuration),
		// the scheduled status should be the most common status
		Status: gofakeit.RandomString([]string{"scheduled", "scheduled", "cancelled", "delayed"}),
	}

	return flight, generateSeats(flight.ID, rows)
}

func Seed(db *database.Database) error {
	return SeedWithSize(db, 1000, 29)
}

func SeedWithSize(db *database.Database, flights, seatRowsPerFlight int) error {
	gofakeit.Seed(999)
	for i := 0; i < flights; i++ {
		f, seats := generateFlight(seatRowsPerFlight)
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
