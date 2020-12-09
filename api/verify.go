package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PhilipGeil/landlyst-backend/auth"
	"github.com/PhilipGeil/landlyst-backend/server"
)

func (api *API) VerifyEmail(ctx context.Context, r *server.APIRequest) (err error) {
	id, ok := r.Vars["id"]
	if !ok {
		return fmt.Errorf("Missing id")
	}
	fmt.Println(id)
	auth.UserVerification(ctx, api.DB, id)

	http.ServeFile(r.W, r.R, "C:\\Users\\phil2643\\development\\landlyst\\api-server\\public\\email-verified.html")

	return

}
