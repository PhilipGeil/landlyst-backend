package api

import (
	"github.com/PhilipGeil/landlyst-backend/server"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type API struct {
	DB *sqlx.DB
}

func Init(srv *server.Server, db *sqlx.DB) {
	api := new(API)
	api.DB = db

	srv.API("POST", "/api/createUser", api.CreateUser)
	srv.API("POST", "/api/login", api.Login)
	srv.API("POST", "/api/rooms", api.Rooms)
	srv.API("GET", "/api/room-additions", api.RoomAdditions)
	srv.API("GET", "/api/sign-out", api.SignOut)
	srv.API("POST", "/api/search-reservation", api.SearchForReservation)
	srv.API("POST", "/api/set-reservation", api.SetReservation)
	srv.API("POST", "/api/send-email", api.SendEmail)
	srv.API("GET, POST", "/api/verify/{id}", api.VerifyEmail)
	srv.API("POST", "/api/zip-code", api.ZipCode)
}

func rollbacktx(tx *sqlx.Tx) {
	if err := tx.Rollback(); err != nil {
		logrus.Errorf("Failed to rollback transaction for patching run: %s", err)
	}
}
