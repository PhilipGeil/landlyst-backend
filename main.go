package main

import (
	"context"
	"fmt"

	"github.com/PhilipGeil/landlyst-backend/auth"
	"github.com/PhilipGeil/landlyst-backend/db"
)

func main() {
	auth.CreateUser()
	db, err := db.ConnectToDB("postgres://landlyst:landlyst123@localhost:5432/landlyst")
	if err != nil {
		fmt.Println(err)
	}

	ctx := context.Background()

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		fmt.Println((err))
	}

	var rooms string
	err = tx.QueryRowContext(
		ctx,
		`
		SELECT room_number FROM rooms
		WHERE id = 1
		`,
	).Scan(&rooms)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rooms)
}
