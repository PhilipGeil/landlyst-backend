package main

import (
	"net/http"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	dbConnectionStringEnv = "postgres://landlyst:landlyst123@localhost:5432/landlyst"
	tokenSecretEnv        = "TOKEN_SECRET"
	tokenSecretPathEnv    = "TOKEN_SECRET_PATH"
	logDebugEnv           = "LOG_DEBUG"
	logJSONEnv            = "LOG_JSON"
	smtpUsernameEnv       = "SMTP_USERNAME"
	smtpPasswordEnv       = "SMTP_PASSWORD"
	encryptionKeyEnv      = "ENCRYPTION_KEY"
)

var (
	dbConnectionString string
	tokenSecret        string

	db               *sqlx.DB
	userAuthTokenAlg *jwt.HMACSHA
	// mailer           email.Mailer

	httpSrv *http.Server
)

func main() {
	makeDatabaseConnection()
	makeUserAuthTokenAlg()
	makeAndStartHTTPServer()
}

func makeDatabaseConnection() {
	// var ok bool
	// dbConnectionString, ok = os.LookupEnv(dbConnectionStringEnv)
	// if !ok {
	// 	logrus.Fatalf("Missing environment variable %s", dbConnectionStringEnv)
	// }

	var err error
	dbConnectionString = dbConnectionStringEnv
	db, err = sqlx.Connect("postgres", dbConnectionString)
	if err != nil {
		panic("open database connection failed: " + err.Error())
	}
}

func makeUserAuthTokenAlg() {
	var tokenSecret []byte
	tokenSecret = []byte(tokenSecretEnv)
	// tokenSecretString, ok := os.LookupEnv(tokenSecretEnv)
	// if ok {
	// 	tokenSecret = []byte(tokenSecretString)
	// } else if path, ok := os.LookupEnv(tokenSecretPathEnv); ok {
	// 	var err error
	// 	tokenSecret, err = ioutil.ReadFile(path)
	// 	if err != nil {
	// 		logrus.Fatalf("Read token secret file error: %s", err)
	// 	}
	// } else {
	// 	logrus.Fatalf("Missing environment variable %s or %s", tokenSecretEnv, tokenSecretPathEnv)
	// }

	userAuthTokenAlg = jwt.NewHS256(tokenSecret)
}
