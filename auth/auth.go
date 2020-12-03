package auth

import (
	"context"
	"encoding/base64"
	"log"
	"strconv"

	"github.com/PhilipGeil/landlyst-backend/api/resources"
	"github.com/jmoiron/sqlx"
)

//CreateUser Creates a new user
func CreateUser(ctx context.Context, user resources.User, db *sqlx.DB) string {
	//Start with hashing password with salt
	salt, err := CreateSalt()
	if err != nil {
		log.Fatal(err)
	}

	hash, err := CreateHashWithSalt(user.Password, salt)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		log.Fatal("Begin DB transaction failed")
	}

	if CheckIfUserExists(ctx, user, db) {
		return "User already exists"
	}

	res, err := tx.ExecContext(
		ctx,
		`
			INSERT INTO users 
			(
				"fName",
				"lName",
				"address",
				"zip_code",
				"phone",
				"email",
				"password",
				"salt",
				"permissions"
			)
			VALUES
			(
				$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7,
				$8,
			'{guest}'
			);
		`,
		user.FName, user.LName, user.Address, user.Zip_code, user.Phone, user.Email, base64.StdEncoding.EncodeToString(hash), base64.StdEncoding.EncodeToString(salt),
	)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	return strconv.FormatInt(rowsAffected, 10) + " Rows was inserted"
}
