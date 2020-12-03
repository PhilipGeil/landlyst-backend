package reservations

import (
	"context"

	"github.com/PhilipGeil/landlyst-backend/api/resources"
	"github.com/jmoiron/sqlx"
)

func SetReservation(ctx context.Context, db *sqlx.DB, r resources.Reservation) (bool, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return false, err
	}

	var id int

	err = tx.QueryRowxContext(
		ctx,
		`
			INSERT INTO reservations (start_date, end_date) VALUES ($1, $2)	RETURNING id
		`,
		r.Dates.StartDate,
		r.Dates.EndDate,
	).Scan(&id)

	for _, room := range r.Rooms {
		_, err = tx.ExecContext(
			ctx,
			`
			INSERT INTO	reservations_rooms (reservation_id, room_id) VALUES ($1, $2)
			`,
			id,
			room.ID,
		)
	}
	if err != nil {
		return false, err
	}

	customerID, err := CreateCustomer(ctx, db, r.Customer)
	if err != nil {
		return false, err
	}

	_, err = tx.ExecContext(
		ctx,
		`
			INSERT INTO customer_reservations (customer_id, reservation_id)	
			VALUES
			(
				$1,
				$2
			)
		`,
		customerID,
		id,
	)

	tx.Commit()

	return true, nil
}
