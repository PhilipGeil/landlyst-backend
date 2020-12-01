package room

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/PhilipGeil/landlyst-backend/api/resources"
	"github.com/jmoiron/sqlx"
)

func Room(ctx context.Context, db *sqlx.DB) ([]resources.RoomAdditions, error) {
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
	fmt.Println(r)

	return roomAdditions, nil
}
