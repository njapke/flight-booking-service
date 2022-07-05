package models

type Seat struct {
	ID        string `json:"id"`
	FlightID  string `json:"flightId"`
	Seat      string `json:"seat"`
	Row       int    `json:"row"`
	Price     int    `json:"price"`
	Available bool   `json:"available"`
}

func (u *Seat) Collection() string {
	return "seats"
}

func (u *Seat) Key() string {
	return u.ID
}
