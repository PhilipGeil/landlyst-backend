package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PhilipGeil/landlyst-backend/api/resources"
	"github.com/PhilipGeil/landlyst-backend/reservations"
	"github.com/PhilipGeil/landlyst-backend/server"
)

func (api *API) SearchForReservation(ctx context.Context, r *server.APIRequest) error {
	ok, err := r.UserAuthentication()
	if err != nil {
		fmt.Println("The error is here")
		return err
	}
	if !ok {
		http.Error(r.W, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return fmt.Errorf("Unauthorized")
	}

	type roomAdditions struct {
		additions []resources.RoomAdditions
	}

	var rs *resources.ReservationSearch

	r.Decode(&rs)

	search, err := reservations.SearchByDate(ctx, api.DB, rs)
	if err != nil {
		return err
	}

	r.Encode(search)
	return nil
}
