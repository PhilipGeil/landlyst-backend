package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(userEmail string, expiration time.Time, sessionID, userID int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":       "landlyst",
		"sub":       "users",
		"aud":       "any",
		"email":     userEmail,
		"exp":       expiration.Unix(),
		"sessionID": sessionID,
		"id":        userID,
	})

	//TODO Create a real secret at some point
	jwtToken, _ := token.SignedString([]byte("secret"))

	fmt.Println(jwtToken)

	return jwtToken
}

func ValidateToken(tokenString string) (ok bool, email string, id int) {
	claims := jwt.MapClaims{}
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		fmt.Println(err)
	}
	ok = tkn.Valid
	email = claims["email"].(string)
	id = int(claims["id"].(float64))
	return

}

func RenewToken(w http.ResponseWriter, tokenStr, userEmail string, userID int) (sessionID int, err error) {

	tokenClaim, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return
	}

	claims := tokenClaim.Claims.(jwt.MapClaims)

	fmt.Println("Here comes the sessionID:")
	sessionID = int(claims["sessionID"].(float64))

	expirationTime := time.Now().Add(24 * time.Hour)
	claims = jwt.MapClaims{
		"exp":       expirationTime.Unix(),
		"sessionID": sessionID,
		"id":        userID,
		"email":     userEmail,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	return
}
