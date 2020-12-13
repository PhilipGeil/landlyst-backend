package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/PhilipGeil/landlyst-backend/api/resources"
	"github.com/PhilipGeil/landlyst-backend/auth"
	"github.com/PhilipGeil/landlyst-backend/email"
	"github.com/PhilipGeil/landlyst-backend/server"
)

func (api *API) CreateUser(ctx context.Context, r *server.APIRequest) (err error) {
	var u resources.User
	if err := r.Decode(&u); err != nil {
		return err
	}
	userID, err := auth.CreateUser(ctx, u, api.DB)
	if err != nil {
		return
	}

	uuid, err := auth.VerifyEmail(ctx, api.DB, userID)
	if err != nil {
		return
	}

	email.SendVerifyEmail(uuid, u.Email, u.Fname)

	type response struct {
		Response string
	}

	r.Encode(response{
		Response: "Ok",
	})

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
	ip := GetIP(r.R)

	sessionID, err := auth.CreateSession(ctx, api.DB, user.ID, ip)
	if err != nil {
		return err
	}

	expiration := time.Now().Add(time.Hour * 24)
	token := auth.CreateToken(user.Email, expiration, sessionID, user.ID)

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

func (api *API) SignOut(ctx context.Context, r *server.APIRequest) error {
	type signOutResponse struct {
		Response string
	}

	http.SetCookie(r.W, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Unix(0, 0),
	})

	r.Encode(signOutResponse{
		Response: "Success",
	})

	return nil
}
