package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PhilipGeil/landlyst-backend/email"
	"github.com/PhilipGeil/landlyst-backend/server"
)

func (api *API) SendEmail(ctx context.Context, r *server.APIRequest) error {
	ok, err := r.UserAuthentication(ctx, api.DB)
	if err != nil {
		fmt.Println("The error is here")
		return err
	}
	if !ok {
		http.Error(r.W, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return fmt.Errorf("Unauthorized")
	}

	email.SendVerifyEmail("123123", "pgj.philip@gmail.com", "Philip")

	return nil
}
