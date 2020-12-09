package api

import (
	"context"

	"github.com/PhilipGeil/landlyst-backend/api/resources"
	"github.com/PhilipGeil/landlyst-backend/room"
	"github.com/PhilipGeil/landlyst-backend/server"
)

func (api *API) Rooms(ctx context.Context, r *server.APIRequest) error {
	// ok, err := r.UserAuthentication(ctx, api.DB)
	// if err != nil {
	// 	fmt.Println("The error is here")
	// 	return err
	// }
	// if !ok {
	// 	http.Error(r.W, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	// 	return fmt.Errorf("Unauthorized")
	// }

	type roomAdditions struct {
		additions []resources.RoomAdditions
	}

	var ra []resources.RoomAdditions

	r.Decode(&ra)

	rooms, err := room.Room(ctx, api.DB, ra)
	if err != nil {
		return err
	}

	r.Encode(rooms)
	return nil
}

func (api *API) RoomAdditions(ctx context.Context, r *server.APIRequest) error {
	// ok, err := r.UserAuthentication(ctx, api.DB)
	// if err != nil {
	// 	return err
	// }
	// if !ok {
	// 	http.Error(r.W, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	// 	return fmt.Errorf("Unauthorized")
	// }
	r.W.Header().Set("Access-Control-Allow-Origin", "*")
	ra, err := room.RoomAdditions(ctx, api.DB)
	if err != nil {
		return err
	}
	r.Encode(ra)
	return nil
}
