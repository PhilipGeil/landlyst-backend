package main

import (
	"net/http"
	"time"

	"github.com/PhilipGeil/landlyst-backend/api"
	"github.com/PhilipGeil/landlyst-backend/server"
	"github.com/sirupsen/logrus"
)

type httpHandler struct {
	Server *server.Server
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Server.ServeHTTP(w, r)
}

func makeHTTPHandler() *httpHandler {
	srv := server.New(db, userAuthTokenAlg)
	api.Init(srv, db)

	var handler httpHandler
	handler.Server = srv
	return &handler
}

func makeAndStartHTTPServer() {
	logrus.Debugln("Starting HTTP server")
	httpSrv = &http.Server{
		Addr:         ":8080",
		Handler:      makeHTTPHandler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := httpSrv.ListenAndServe()
	if err == http.ErrServerClosed {
		return
	}
	if err != nil {
		logrus.Fatalf("HTTP server error: %s", err)
	}
}
