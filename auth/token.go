package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken() []byte {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "landlyst",
		"sub": "users",
		"aud": "any",
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	})

	//TODO Create a real secret at some point
	jwtToken, _ := token.SignedString([]byte("secret"))

	return []byte(jwtToken)
}
