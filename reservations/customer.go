package reservations

import (
	"context"

	"github.com/PhilipGeil/landlyst-backend/api/resources"
	"github.com/jmoiron/sqlx"
)

//CreateCustomer creates a new customer in the DB and returns the newly created ID
func CreateCustomer(ctx context.Context, db *sqlx.DB, c resources.Customer) (customerID int, err error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return
	}

	err = tx.QueryRowxContext(
		ctx,
		`
			INSERT INTO customers ("fName", "lName", address, zip_code, phone, email)
			VALUES
			(
				$1,
				$2,
				$3,
				$4,
				$5,
				$6
			) RETURNING id
		`,
		c.FName,
		c.LName,
		c.Address,
		c.Zip_code,
		c.Phone,
		c.Email,
	).Scan(&customerID)
	if err != nil {
		return
	}

	if c.UserID != 0 {
		_, err = tx.ExecContext(
			ctx,
			`
				INSERT INTO customers_users (customer_id, user_id) VALUES ($1, $2)
			`,
			customerID,
			c.UserID,
		)
	}

	tx.Commit()

	return

}
