package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

//ConnectToDB creates a connection to the database
func ConnectToDB(conn string) (*sqlx.DB, error) {

	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
