package reservations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/PhilipGeil/landlyst-backend/api/resources"
	"github.com/PhilipGeil/landlyst-backend/email"
	_ "github.com/jackc/pgtype"
	"github.com/jmoiron/sqlx"
)

func SetReservation(ctx context.Context, db *sqlx.DB, r resources.Reservation) (res resources.ReservationResponse, err error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return
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
		return
	}

	customerID, err := CreateCustomer(ctx, tx, r.Customer)
	if err != nil {
		return
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

	res.Reservation = r
	res.Id = id

	var discount resources.Discount

	var discountType []uint8

	err = tx.QueryRowxContext(
		ctx,
		`
			SELECT id, type, number::integer, description FROM reservations_discounts
			JOIN discounts ON discounts.id = reservations_discounts.discount_id
			WHERE reservations_discounts.reservation_id = $1 
		`,
		id,
	).Scan(&discount.Id, &discountType, &discount.Number, &discount.Description)
	if err == sql.ErrNoRows {
		return res, nil
	} else if err != nil {
		return
	}

	discount.Type = B2S(discountType)

	tx.Commit()

	res.Discount = discount

	return
}

func B2S(bs []uint8) string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		b[i] = byte(v)
	}
	return string(b)
}

func ConfirmReservation(ctx context.Context, db *sqlx.DB, res resources.ReservationResponse) (err error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return
	}

	_, err = tx.ExecContext(
		ctx,
		`
			UPDATE reservations
			SET confirmed = true
			WHERE id = $1	
		`,
		res.Id,
	)
	if err != nil {
		return
	}

	tx.Commit()

	email.SendConfirmEmail(res)

	return
}

func GetReservation(ctx context.Context, db *sqlx.DB, id int, email string) (res []resources.GetReservation, err error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return
	}

	rows, err := tx.QueryxContext(
		ctx,
		`
			SELECT reservation_id::integer, start_date, end_date FROM customer_reservations
			JOIN reservations ON reservations.id = customer_reservations.reservation_id
			JOIN customers ON customers.id = customer_reservations.customer_id
			WHERE customer_id IN (
				SELECT customer_id FROM customers_users
				JOIN users ON users.id = customers_users.user_id
				WHERE user_id = $1
			) OR email IN (
				SELECT email FROM customers_users
				JOIN users ON users.id = customers_users.user_id
				WHERE email = $2
			) 
			AND confirmed = true
		`,
		id,
		email,
	)
	if err == sql.ErrNoRows {
		return res, nil
	}
	if err != nil {
		fmt.Println("It's here")
		return
	}

	var r resources.GetReservation

	for rows.Next() {
		rows.Scan(&r.Id, &r.ReservationDates.StartDate, &r.ReservationDates.EndDate)
		res = append(res, r)
	}

	return
}
