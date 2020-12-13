package auth

import (
	"context"
	"fmt"

	"github.com/PhilipGeil/landlyst-backend/api/resources"
	"github.com/jmoiron/sqlx"
	"github.com/shaj13/go-guardian/auth"
	"github.com/shaj13/go-guardian/store"
)

var authenticator auth.Authenticator
var cache store.Cache

//ValidateUser validates the user by checking if the email already exist in the db
func ValidateUser(ctx context.Context, c resources.Credentials, db *sqlx.DB) (*resources.User, error) {
	user, err := ValidateEmailAndPassword(ctx, db, c.Password, c.Email)
	if err != nil {
		fmt.Println("The error is here auth/ValidateUser ln 37")
		return nil, err
	}

	return user, nil
}
