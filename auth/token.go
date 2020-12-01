package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(userEmail string, expiration time.Time) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":   "landlyst",
		"sub":   "users",
		"aud":   "any",
		"email": userEmail,
		"exp":   expiration.Unix(),
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

func RenewToken(w http.ResponseWriter) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := jwt.MapClaims{
		"exp": expirationTime.Unix(),
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
}
