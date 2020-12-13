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

// func (r *Request) UserAuthentication() (userID string, authenticated bool, err error) {
// 	if r.userAuth != nil {
// 		return r.userAuth.userID, r.userAuth.authenticated, nil
// 	}
// 	userID, authenticated, err = requestUserAuth(r.r, r.srv.UserAuthTokenAlg)
// 	if err != nil {
// 		return
// 	}
// 	r.userAuth = &cachedUserAuthentication{
// 		userID:        userID,
// 		authenticated: authenticated,
// 	}
// 	return
// }

func request(w http.ResponseWriter, r *http.Request, srv *Server) *Request {
	req := new(Request)
	req.W = w
	req.R = r
	req.Vars = gmux.Vars(r)
	req.Method = r.Method
	req.srv = srv
	return req
}

// func requestUserAuth(r *http.Request, userAuthTokenAlg *jwt.HMACSHA) (string, bool, error) {
// 	var (
// 		token string
// 		ok    bool
// 	)
// 	token, ok = headerAuth(r)
// 	if !ok {
// 		var err error
// 		token, ok, err = cookieUserAuth(r)
// 		if err != nil {
// 			return "", false, err
// 		} else if !ok {
// 			return "", false, nil
// 		}
// 	}

// 	// Decode auth token
// 	userID, err := decodeToken(userAuthTokenAlg, []byte(token))
// 	if err != nil {
// 		return "", false, ErrInvalidAuthToken
// 	}

// 	return userID, true, nil
// }
