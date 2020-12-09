package auth

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func VerifyEmail(ctx context.Context, db *sqlx.DB, user_id int) (uuid string, err error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return
	}

	err = tx.QueryRowxContext(
		ctx,
		`
			INSERT INTO verifications (user_id) VALUES ($1)	RETURNING uuid
		`,
		user_id,
	).Scan(&uuid)
	if err != nil {
		return
	}

	tx.Commit()

	return
}
