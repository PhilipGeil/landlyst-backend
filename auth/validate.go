package auth

import (
	"context"
	"database/sql"
	"encoding/base64"
	"log"

	"github.com/PhilipGeil/landlyst-backend/api/resources"
	"github.com/jmoiron/sqlx"
)

func ValidateEmailAndPassword(db *sqlx.DB, ctx context.Context, password, email string) bool {
	salt, err := getSalt(email, db, ctx)
	if err != nil {
		log.Fatal(err)
	}

	hash, err := CreateHashWithSalt(password, []byte(salt))
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	var user resources.User

	err = tx.QueryRowxContext(
		ctx,
		`
			SELECT fName, lName, address, city, zip_code, phone, email, password FROM users
			WHERE email = $1 AND password = $2;
		`,
		email,
		base64.StdEncoding.EncodeToString(hash),
	).Scan(&user)

	if err == sql.ErrNoRows {
		return false
	}

	return true
}

func getSalt(email string, db *sqlx.DB, ctx context.Context) (string, error) {
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
