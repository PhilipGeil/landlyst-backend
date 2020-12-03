package auth

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

//CreateSession creates a new session
func CreateSession(ctx context.Context, db *sqlx.DB, userID int, ipAddress string) (sessionID int, err error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return
	}

	err = tx.QueryRowxContext(
		ctx,
		`
			INSERT INTO sessions (
				user_id,
				ip_address
			) VALUES (
				$1,
				$2
			) RETURNING id
		`,
		userID,
		ipAddress,
	).Scan(&sessionID)
	if err != nil {
		return
	}

	tx.Commit()

	return
}

func ExtendSession(ctx context.Context, db *sqlx.DB, sessionID int) (err error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return
	}

	_, err = tx.ExecContext(
		ctx,
		`
			UPDATE sessions 
			SET last_seen = $1
			WHERE id = $2	
		`,
		time.Now(),
		sessionID,
	)
	if err != nil {
		return
	}

	tx.Commit()

	return
}
