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
	srv.API("GET", "/api/login", api.Login)
}

func rollbacktx(tx *sqlx.Tx) {
	if err := tx.Rollback(); err != nil {
		logrus.Errorf("Failed to rollback transaction for patching run: %s", err)
	}
}
