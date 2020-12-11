package auth

import (
	"context"
	"encoding/base64"
	"log"

	"github.com/PhilipGeil/landlyst-backend/api/resources"
	"github.com/jmoiron/sqlx"
)

//CreateUser Creates a new user
func CreateUser(ctx context.Context, user resources.User, db *sqlx.DB) (userID int, err error) {
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

	_, err = CheckIfUserExists(ctx, user, db)
	if err != nil {
		return
	}

	err = tx.QueryRowxContext(
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
			) RETURNING id;
		`,
		user.Fname, user.Lname, user.Address, user.Zip_code, user.Phone, user.Email, base64.StdEncoding.EncodeToString(hash), base64.StdEncoding.EncodeToString(salt),
	).Scan(&userID)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
	return
}
