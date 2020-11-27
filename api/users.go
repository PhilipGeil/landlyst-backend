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
	createUserResponse := auth.CreateUser(ctx, u, api.DB)

	type response struct {
		Response string
	}

	res := response{
		Response: createUserResponse,
	}

	r.Encode(res)

	return nil
}

func (api *API) Login(ctx context.Context, r *server.APIRequest) error {
	auth.Middleware(r.Request)
	return nil
}
