package models

type Passenger struct {
	Name string `json:"name"`
	Seat string `json:"seat"`
}

type Booking struct {
	ID         string      `json:"id"`
	UserID     string      `json:"userId"`
	FlightID   string      `json:"flightId"`
	Price      int         `json:"price"`
	Status     string      `json:"status"`
	Passengers []Passenger `json:"passengers"`
}

func (b *Booking) Collection() string {
	return "bookings"
}

func (b *Booking) Key() string {
	return b.ID
}
