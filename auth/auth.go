package auth

import (
	"context"
	"encoding/base64"
	"log"

	"github.com/jackc/pgtype"

	"github.com/PhilipGeil/landlyst-backend/api/resources"
	"github.com/jmoiron/sqlx"
)

func CToGoString(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}

//CreateUser Creates a new user
func CreateUser(ctx context.Context, user resources.User, db *sqlx.DB) {
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

	var permissionTypes pgtype.EnumArray
	permissionTypes.Status = pgtype.Present
	permissionTypes.Dimensions = []pgtype.ArrayDimension{
		{
			Length:     int32(len(user.Permissions)),
			LowerBound: 1,
		},
	}
	permissionTypes.Elements = make([]pgtype.GenericText, len(user.Permissions))
	for idx, permissionType := range user.Permissions {
		permissionTypes.Elements[idx] = pgtype.GenericText{
			String: permissionType,
		}
		permissionTypes.Elements[idx].Status = pgtype.Present
	}

	_, err = tx.ExecContext(
		ctx,
		`
			INSERT INTO users 
			(
				"fName",
				"lName",
				"address",
				"city",
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
				$9,
				$10
			);
		`,
		user.FName, user.LName, user.Address, user.City, user.Zip_code, user.Phone, user.Email, base64.StdEncoding.EncodeToString(hash), base64.StdEncoding.EncodeToString(salt), permissionTypes,
	)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
}
