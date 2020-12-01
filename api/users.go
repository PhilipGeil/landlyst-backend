package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

	var c resources.Credentials

	if err := r.Decode(&c); err != nil {
		http.Error(r.W, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return err
	}

	fmt.Println(c.Email)

	type loginResponse struct {
		Response string
		ID       int
		User     *resources.User
	}

	user, err := auth.ValidateUser(ctx, c, api.DB)
	if err != nil {
		r.Encode(loginResponse{Response: "Error logging in", ID: 100, User: nil})
		return err
	}
	expiration := time.Now().Add(time.Minute * 2)
	token := auth.CreateToken(user.Email, expiration)

	ip := GetIP(r.R)

	err = auth.CreateSession(ctx, api.DB, user.ID, token, ip)

	http.SetCookie(r.W, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expiration,
	})

	r.Encode(loginResponse{
		Response: "Succes",
		ID:       http.StatusOK,
		User:     user,
	})

	return nil
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
