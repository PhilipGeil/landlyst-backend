package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(userEmail string, expiration time.Time, sessionID int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":       "landlyst",
		"sub":       "users",
		"aud":       "any",
		"email":     userEmail,
		"exp":       expiration.Unix(),
		"sessionID": sessionID,
	})

	//TODO Create a real secret at some point
	jwtToken, _ := token.SignedString([]byte("secret"))

	fmt.Println(jwtToken)

	return jwtToken
}

func ValidateToken(tokenString string) bool {
	claims := jwt.MapClaims{}
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tkn.Valid)
	return tkn.Valid

}

func RenewToken(w http.ResponseWriter, tokenStr string) (sessionID int, err error) {

	tokenClaim, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return
	}

	claims := tokenClaim.Claims.(jwt.MapClaims)

	fmt.Println("Here comes the sessionID:")
	sessionID = int(claims["sessionID"].(float64))

	expirationTime := time.Now().Add(2 * time.Hour)
	claims = jwt.MapClaims{
		"exp":       expirationTime.Unix(),
		"sessionID": sessionID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	return
}
