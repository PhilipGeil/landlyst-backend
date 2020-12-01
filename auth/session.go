package auth

import (
	"context"

	"github.com/jmoiron/sqlx"
)

//CreateSession creates a new session
func CreateSession(ctx context.Context, db *sqlx.DB, userID int, token, ipAddress string) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		`
			INSERT INTO sessions (
				user_id,
				auth_token,
				ip_address
			) VALUES (
				$1,
				$2,
				$3
			)
		`,
		userID,
		token,
		ipAddress,
	)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}
