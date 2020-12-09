package api

import (
	"context"

	"github.com/PhilipGeil/landlyst-backend/server"
)

func (api *API) ZipCode(ctx context.Context, r *server.APIRequest) error {
	type zip struct {
		Zip string `json:"zip_code"`
	}

	var z zip

	r.Decode(&z)

	type city struct {
		City string
	}

	var c city

	tx, err := api.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	err = tx.QueryRowxContext(
		ctx,
		`
			SELECT city FROM zip_code WHERE zip_code = $1	
		`,
		z.Zip,
	).Scan(&c.City)
	if err != nil {
		return err
	}

	r.Encode(c)

	return nil
}
