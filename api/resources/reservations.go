package resources

import (
	"time"
)

type ReservationDates struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type ReservationSearch struct {
	Items []RoomAdditions
	Dates ReservationDates
}
