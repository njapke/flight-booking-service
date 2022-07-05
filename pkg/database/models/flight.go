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

func (f *Flight) Collection() string {
	return "flights"
}

func (f *Flight) Key() string {
	return f.ID
}
