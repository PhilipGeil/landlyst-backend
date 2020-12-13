package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/PhilipGeil/landlyst-backend/api/resources"
	"github.com/PhilipGeil/landlyst-backend/email"
	"github.com/PhilipGeil/landlyst-backend/server"
)

func (api *API) SendEmail(ctx context.Context, r *server.APIRequest) error {
	ok, _, _, err := r.UserAuthentication(ctx, api.DB)
	if err != nil {
		fmt.Println("The error is here")
		return err
	}
	if !ok {
		http.Error(r.W, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return fmt.Errorf("Unauthorized")
	}

	// email.SendVerifyEmail("123123", "pgj.philip@gmail.com", "Philip")

	var c = resources.Customer{
		FName:    "Philip",
		LName:    "Jensen",
		Address:  "Søkrogen 55",
		Phone:    "27581220",
		Email:    "pgj.philip@gmail.com",
		Zip_code: "4281",
	}

	var d = resources.ReservationDates{
		StartDate: time.Now(),
		EndDate:   time.Now(),
	}

	var re = resources.Reservation{
		Customer: c,
		Dates:    d,
	}

	var res = resources.ReservationResponse{
		Reservation: re,
	}
	email.SendConfirmEmail(res)

	return nil
}
