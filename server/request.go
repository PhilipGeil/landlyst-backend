package server

import (
	"context"
	"net/http"

	"github.com/PhilipGeil/landlyst-backend/auth"
	gmux "github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type Request struct {
	W        http.ResponseWriter
	R        *http.Request
	userAuth *cachedUserAuthentication

	srv *Server

	Vars   map[string]string
	Method string
}

type cachedUserAuthentication struct {
	userID        string
	authenticated bool
}

//UserAuthentication authenticates the user
func (r *Request) UserAuthentication(ctx context.Context, db *sqlx.DB) (ok bool, email string, id int, err error) {
	c, err := r.R.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			r.W.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		r.W.WriteHeader(http.StatusBadRequest)
		return
	}

	tknStr := c.Value
	if tknStr == "" {
		return
	}

	ok, email, id = auth.ValidateToken(tknStr)
	var sessionID int
	if ok {
		sessionID, err = auth.RenewToken(r.W, tknStr, email, id)
		if err != nil {
			return
		}
		auth.ExtendSession(ctx, db, sessionID)
	}
	return
}

func request(w http.ResponseWriter, r *http.Request, srv *Server) *Request {
	req := new(Request)
	req.W = w
	req.R = r
	req.Vars = gmux.Vars(r)
	req.Method = r.Method
	req.srv = srv
	return req
}
