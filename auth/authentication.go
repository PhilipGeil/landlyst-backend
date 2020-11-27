package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PhilipGeil/landlyst-backend/server"
	"github.com/jmoiron/sqlx"
	"github.com/shaj13/go-guardian/auth"
	"github.com/shaj13/go-guardian/auth/strategies/basic"
	"github.com/shaj13/go-guardian/store"
)

var authenticator auth.Authenticator
var cache store.Cache
var DB *sqlx.DB

func setupGoGuardian(db *sqlx.DB) {
	authenticator = auth.New()
	cache = store.NewFIFO(context.Background(), time.Minute*10)
	DB = db

	basicStrategy := basic.New(validateUser, cache)
	// tokenStrategy := bearer.New(verifyToken, cache)

	authenticator.EnableStrategy(basic.StrategyKey, basicStrategy)
	// authenticator.EnableStrategy(bearer.CachedStrategyKey, tokenStrategy)
}

func validateUser(ctx context.Context, r *http.Request, email, password string) (auth.Info, error) {
	// here connect to db or any other service to fetch user and validate it.
	if ValidateEmailAndPassword(DB, ctx, password, email) {
		return auth.NewDefaultUser("medium", "1", nil, nil), nil
	}

	return nil, fmt.Errorf("Invalid credentials")
}

func Middleware(r *server.Request) {
	log.Println("Executing Auth Middleware")
	user, err := authenticator.Authenticate(r.R)
	if err != nil {
		code := http.StatusUnauthorized
		http.Error(r.W, http.StatusText(code), code)
	}
	log.Printf("User %s Authenticated\n", user.UserName())
}
