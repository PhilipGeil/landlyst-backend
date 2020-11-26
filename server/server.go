package server

import (
	"net/http"

	"github.com/gbrlsnchs/jwt/v3"
	gmux "github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"gocloud.dev/blob"
)

type Server struct {
	Database         *sqlx.DB
	UserAuthTokenAlg *jwt.HMACSHA
	mux              *gmux.Router
	buckets          map[string]*blob.Bucket
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.mux.ServeHTTP(w, r)
}

func New(db *sqlx.DB, tokenAlg *jwt.HMACSHA) *Server {
	mux := gmux.NewRouter()
	buckets := make(map[string]*blob.Bucket)
	return &Server{
		Database:         db,
		UserAuthTokenAlg: tokenAlg,
		mux:              mux,
		buckets:          buckets,
	}
}

type httpStatusHijacker struct {
	code int
	http.ResponseWriter
}

type metricsHandler struct {
	Path    string
	Handler http.Handler
}

func (m *metricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hijacker := httpStatusHijacker{code: 200, ResponseWriter: w}
	// startTime := time.Now()
	m.Handler.ServeHTTP(&hijacker, r)
	// duration := time.Now().Sub(startTime)
	// log(hijacker.code, r, duration)
}
