package auth

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func UserVerification(ctx context.Context, db *sqlx.DB, uuid string) (err error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return
	}

	var userID int

	err = db.QueryRowxContext(
		ctx,
		`
			SELECT user_id FROM verifications WHERE uuid = $1 
		`,
		uuid,
	).Scan(&userID)

	_, err = tx.ExecContext(
		ctx,
		`
			UPDATE users
			SET verified = true
			WHERE id = $1	
		`,
		userID,
	)

	tx.Commit()
	return
}
