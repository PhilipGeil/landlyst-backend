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
	ok, err := r.UserAuthentication(ctx, api.DB)
	if err != nil {
		fmt.Println("The error is here")
		return err
	}
	if !ok {
		http.Error(r.W, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return fmt.Errorf("Unauthorized")
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

//SetReservation creates a new reservation
func (api *API) SetReservation(ctx context.Context, r *server.APIRequest) error {
	ok, err := r.UserAuthentication(ctx, api.DB)
	if err != nil {
		fmt.Println("The error is here")
		return err
	}
	if !ok {
		http.Error(r.W, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return fmt.Errorf("Unauthorized")
	}

	var res resources.Reservation

	r.Decode(&res)

	resOK, err := reservations.SetReservation(ctx, api.DB, res)
	if err != nil {
		return err
	}

	type response struct {
		Response string
	}

	if resOK {
		r.Encode(response{
			Response: "Success",
		})
	}
	return nil
}

/*

SELECT count(r.room_id) as amount, array_agg(r.room_id) as rooms, b.adds, b.items, b.price FROM (SELECT room_id FROM rooms_room_additions
GROUP by room_id) r
LEFT JOIN
(SELECT ra.room_id, array_agg(ab.id) as adds, array_agg(ab.item) as items, sum(ab.price) as price FROM rooms_room_additions ra
    LEFT JOIN room_additions ab
    ON ab.id = ra.room_additions_id
    GROUP BY ra.room_id) as b
ON b.room_id = r.room_id
GROUP BY (b.adds, b.items, b.price)
*/
