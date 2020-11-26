package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type APIRequest struct {
	*Request
}

type APIHandler func(ctx context.Context, r *APIRequest) error

func (srv *Server) API(methods, path string, handler APIHandler) {
	srv.mux.Handle(path, &metricsHandler{path, &apiHandler{
		handler: handler,
		srv:     srv,
	}}).Methods(strings.Split(methods, ",")...)
}

type apiHandler struct {
	handler APIHandler
	srv     *Server
}

func (h *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	req := &APIRequest{
		Request: request(w, r, h.srv),
	}
	err := h.handler(r.Context(), req)
	if err != nil {
		logrus.Errorf("API handler returned error: %s", err)
	}
}

func (r *APIRequest) Decode(p interface{}) error {
	jsonEncoded := strings.Contains(r.r.Header.Get("Content-Type"), "application/json")
	if !jsonEncoded {
		return nil
	}
	err := json.NewDecoder(r.r.Body).Decode(p)
	if err != nil {
		logrus.Warnf("Decode JSON error: %s", err)
		return err
	}
	return nil
}

func (r *APIRequest) Encode(p interface{}) bool {
	r.r.Header.Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(r.w).Encode(p)
	if err != nil {
		return false
	}
	return true
}
