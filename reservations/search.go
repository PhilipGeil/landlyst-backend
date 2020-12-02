package reservations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/PhilipGeil/landlyst-backend/api/resources"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func SearchByDate(ctx context.Context, db *sqlx.DB, rs *resources.ReservationSearch) ([]resources.Room, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var additions []int

	for _, i := range rs.Items {
		additions = append(additions, i.ID)
	}

	var rows *sqlx.Rows

	if len(rs.Items) == 0 {
		rows, err = tx.QueryxContext(
			ctx,
			`
				SELECT room_number, id, cleaned, current_staying FROM rooms
				WHERE id NOT IN (
					SELECT room_id FROM reservations_rooms
					JOIN reservations
					ON reservations.id = reservations_rooms.reservation_id 
					WHERE reservations.start_date BETWEEN $1 AND $2
					OR reservations.end_date BETWEEN $1 AND $2
				)	
			`,
			&rs.Dates.StartDate,
			&rs.Dates.EndDate,
		)
		if err == sql.ErrNoRows {
			fmt.Println("Sorry mate, didn't find anything")
		} else if err != nil {
			return nil, err
		}
	} else {

		rows, err = tx.QueryxContext(
			ctx,
			`
			SELECT room_number, id, cleaned, current_staying FROM rooms
			WHERE id NOT IN (
				SELECT room_id FROM reservations_rooms
				JOIN reservations
				ON reservations.id = reservations_rooms.reservation_id 
				WHERE reservations.start_date BETWEEN $1 AND $2
				OR reservations.end_date BETWEEN $1 AND $2
			)
			AND
			id IN (
				SELECT room_id FROM rooms_room_additions
				JOIN room_additions
				ON room_additions.id = rooms_room_additions.room_additions_id
				WHERE room_additions.id = ANY($3)
				GROUP BY room_id
				HAVING count(room_id) = $4
				)	
				`,
			&rs.Dates.StartDate,
			&rs.Dates.EndDate,
			pq.Array(additions),
			len(additions),
		)
		if err == sql.ErrNoRows {
			fmt.Println("Sorry mate, didn't find anything")
		}
		if err != nil {
			return nil, err
		}
	}

	var rooms []resources.Room

	var r resources.Room

	for rows.Next() {
		rows.Scan(&r.RoomNumber, &r.ID, &r.Cleaned, &r.CurrentStaying)
		rooms = append(rooms, r)
		fmt.Println(r)
	}

	fmt.Println(rooms)

	return rooms, nil
}
