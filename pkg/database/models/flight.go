package models

import "time"

type Flight struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Departure time.Time `json:"departure"`
	Arrival   time.Time `json:"arrival"`
	Status    string    `json:"status"`
}

func (u *Flight) Collection() string {
	return "flights"
}

func (u *Flight) Key() string {
	return u.ID
}
