package api

import (
	"context"

	"github.com/PhilipGeil/landlyst-backend/api/resources"
	"github.com/PhilipGeil/landlyst-backend/auth"
	"github.com/PhilipGeil/landlyst-backend/server"
)

func (api *API) CreateUser(ctx context.Context, r *server.APIRequest) error {
	var u resources.User
	if err := r.Decode(&u); err != nil {
		return err
	}
	auth.CreateUser(ctx, u, api.DB)

	// type rows struct {
	// 	RowsAffected int64
	// }

	// ro := rows{
	// 	RowsAffected: rowsAffected,
	// }

	// if rowsAffected > 0 {
	// 	r.Encode(ro)
	// }

	return nil
}
