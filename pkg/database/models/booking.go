package models

type Booking struct {
	ID       string `json:"id"`
	UserID   string `json:"userId"`
	FlightID string `json:"flightId"`
	Price    int    `json:"price"`
	Status   string `json:"status"`
}

func (b *Booking) Collection() string {
	return "bookings"
}

func (b *Booking) Key() string {
	return b.ID
}
