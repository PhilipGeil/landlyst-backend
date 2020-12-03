package auth

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/PhilipGeil/landlyst-backend/api/resources"
	"github.com/jmoiron/sqlx"
)

func ValidateEmailAndPassword(ctx context.Context, db *sqlx.DB, password, email string) (*resources.User, error) {
	salt, err := getSalt(ctx, email, db)
	if err != nil {
		log.Fatal(err)
	}

	saltBytes, err := base64.StdEncoding.DecodeString(salt)

	hash, err := CreateHashWithSalt(password, saltBytes)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(base64.StdEncoding.EncodeToString(hash))

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	var user resources.User

	err = tx.QueryRowxContext(
		ctx,
		`
			SELECT id, "fName", "lName", address, zip_code, phone, email FROM users
			WHERE email = $1 AND password = $2;
		`,
		email,
		base64.StdEncoding.EncodeToString(hash),
	).Scan(&user.ID, &user.FName, &user.LName, &user.Address, &user.Zip_code, &user.Phone, &user.Email)
	if err == sql.ErrNoRows {
		return nil, err
	}

	return &user, nil
}

func getSalt(ctx context.Context, email string, db *sqlx.DB) (string, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	var salt string

	err = tx.QueryRowxContext(
		ctx,
		`
			SELECT salt FROM users WHERE email = $1;
		`,
		email,
	).Scan(&salt)
	if err == sql.ErrNoRows {
		return "", err
	}

	return salt, nil
}

func CheckIfUserExists(ctx context.Context, user resources.User, db *sqlx.DB) bool {
	tx, err := db.BeginTxx(ctx, nil)

	var id string

	err = tx.QueryRowxContext(
		ctx,
		`
			SELECT id FROM users WHERE email = $1
		`,
		user.Email,
	).Scan(&id)
	if err == sql.ErrNoRows {
		return false
	}

	tx.Commit()

	return true

}
