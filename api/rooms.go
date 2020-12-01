package api

import (
	"context"

	"github.com/PhilipGeil/landlyst-backend/room"
	"github.com/PhilipGeil/landlyst-backend/server"
)

func (api *API) Rooms(ctx context.Context, r *server.APIRequest) error {
	type temp struct {
		Response string
	}

	ok, err := r.UserAuthentication()
	if err != nil {
		return err
	}

	if ok {
		r.Encode(temp{
			Response: "Shit works",
		})
	} else {
		r.Encode(temp{
			Response: "Shit don't work",
		})
	}
	return nil
}

func (api *API) RoomAdditions(ctx context.Context, r *server.APIRequest) error {
	r.UserAuthentication()
	ra, err := room.Room(ctx, api.DB)
	if err != nil {
		return err
	}
	r.Encode(ra)
	return nil
}
