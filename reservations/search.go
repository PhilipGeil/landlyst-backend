package reservations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/PhilipGeil/landlyst-backend/api/resources"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func SearchByDate(ctx context.Context, db *sqlx.DB, rs resources.ReservationSearch) ([]resources.ReservationSearchResult, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var additions []int

	for _, i := range rs.Items {
		additions = append(additions, i.ID)
	}

	fmt.Println(rs)

	var rows *sqlx.Rows

	if len(rs.Items) == 0 {
		rows, err = tx.QueryxContext(
			ctx,
			`
			SELECT count(r.room_id) as amount, array_agg(r.room_id) as rooms, b.add_id, b.items, b.price::integer FROM (SELECT room_id FROM rooms_room_additions
				GROUP by room_id) r
				LEFT JOIN
				(SELECT ra.room_id, array_agg(ab.id) as add_id, array_agg(ab.item) as items, sum(ab.price) as price FROM rooms_room_additions ra
					LEFT JOIN room_additions ab
					ON ab.id = ra.room_additions_id
					GROUP BY ra.room_id) as b
				ON b.room_id = r.room_id
				WHERE r.room_id NOT IN (
				SELECT room_id FROM reservations_rooms
				JOIN reservations
				ON reservations.id = reservations_rooms.reservation_id 
				WHERE tsrange($1, $2) @> reservations.start_date::timestamp
				OR tsrange($1, $2) @> reservations.end_date::timestamp
				AND confirmed = true
				)	
				GROUP BY (b.add_id, b.items, b.price)
				ORDER BY (b.price)
			`,
			rs.Dates.StartDate,
			rs.Dates.EndDate,
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
			SELECT count(r.room_id) as amount, array_agg(r.room_id) as "rooms", b.add_id, b.items, b.price::integer FROM (SELECT room_id FROM rooms_room_additions
				GROUP by room_id) r
				LEFT JOIN
				(SELECT ra.room_id, array_agg(ab.id) as "add_id", array_agg(ab.item) as "items", sum(ab.price) as "price" FROM rooms_room_additions ra
					LEFT JOIN room_additions ab
					ON ab.id = ra.room_additions_id
					GROUP BY ra.room_id) as b
				ON b.room_id = r.room_id
				WHERE r.room_id NOT IN (
				SELECT room_id FROM reservations_rooms
				JOIN reservations
				ON reservations.id = reservations_rooms.reservation_id 
				WHERE tsrange($1, $2) @> reservations.start_date::timestamp
				OR tsrange($1, $2) @> reservations.end_date::timestamp
				AND confirmed = true
				)	
				AND
				r.room_id IN (
					SELECT room_id FROM rooms_room_additions
					JOIN room_additions
					ON room_additions.id = rooms_room_additions.room_additions_id
					WHERE room_additions.id = ANY($3)
					GROUP BY room_id
					HAVING count(room_id) = $4
					)
				GROUP BY (b.add_id, b.items, b.price)
				ORDER BY (b.price)
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

	var result []resources.ReservationSearchResult

	var r resources.ReservationSearchResult

	for rows.Next() {
		rows.Scan(&r.Amount, &r.Rooms, &r.AddID, &r.Items, &r.Price)
		result = append(result, r)
	}

	return result, nil
}
