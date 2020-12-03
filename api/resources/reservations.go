package resources

import (
	"time"

	"github.com/lib/pq"
)

type ReservationDates struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type ReservationSearch struct {
	Items []RoomAdditions
	Dates ReservationDates
}

type Reservation struct {
	Rooms    []Room           `json:"rooms"`
	Dates    ReservationDates `json:"dates"`
	Customer Customer         `json:"customer"`
}

type ReservationSearchResult struct {
	Amount int            `db:"amount"`
	Rooms  pq.Int64Array  `db:"rooms"`
	AddID  pq.Int64Array  `db:"add_id"`
	Items  pq.StringArray `db:"items"`
	Price  int            `db:"price"`
}
