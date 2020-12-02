package room

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/PhilipGeil/landlyst-backend/api/resources"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func Room(ctx context.Context, db *sqlx.DB, ra []resources.RoomAdditions) ([]resources.Room, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var additions []int

	for _, i := range ra {
		additions = append(additions, i.ID)
	}

	rows, err := tx.QueryxContext(
		ctx,
		`
		SELECT rooms.id, rooms.room_number, rooms.cleaned, rooms.current_staying FROM rooms_room_additions
		JOIN rooms
		ON rooms.id = rooms_room_additions.room_id
		JOIN room_additions
		ON room_additions.id = rooms_room_additions.room_additions_id
		WHERE room_additions.id = ANY($1)
		GROUP BY (rooms.id, rooms.room_number, rooms.cleaned, rooms.current_staying)
		HAVING count(rooms.id) = $2
		`,
		pq.Array(additions),
		len(additions),
	)
	if err != nil {
		return nil, err
	}
	var r resources.Room

	var rooms []resources.Room

	for rows.Next() {
		rows.StructScan(&r)
		rooms = append(rooms, r)
	}

	return rooms, nil
}

func RoomAdditions(ctx context.Context, db *sqlx.DB) ([]resources.RoomAdditions, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var r resources.RoomAdditions

	rows, err := tx.QueryxContext(
		ctx,
		`
			SELECT id, item, price, description FROM room_additions;
		`,
	)
	if err == sql.ErrNoRows {
		fmt.Println("No rows to be found")
	} else if err != nil {
		return nil, err
	}

	var roomAdditions []resources.RoomAdditions

	for rows.Next() {
		err = rows.StructScan(&r)
		roomAdditions = append(roomAdditions, r)
	}
	if err != nil {
		return nil, err
	}

	return roomAdditions, nil
}
