package server

import (
	"net/http"

	gmux "github.com/gorilla/mux"
)

type Request struct {
	w        http.ResponseWriter
	r        *http.Request
	userAuth *cachedUserAuthentication

	srv *Server

	Vars   map[string]string
	Method string
}

type cachedUserAuthentication struct {
	userID        string
	authenticated bool
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
	req.w = w
	req.r = r
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
