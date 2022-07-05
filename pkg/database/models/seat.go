package models

import "fmt"

type Seat struct {
	FlightID  string `json:"flightId"`
	Seat      string `json:"seat"`
	Row       int    `json:"row"`
	Price     int    `json:"price"`
	Available bool   `json:"available"`
}

func (s *Seat) Collection() string {
	return "seats"
}

func (s *Seat) Key() string {
	return fmt.Sprintf("%s/%s", s.FlightID, s.Seat)
}
