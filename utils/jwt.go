package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TODO: parameterize this value
var JWTSecret = []byte("!!SECRET!!")

func GenerateJWT(name string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = name
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, _ := token.SignedString(JWTSecret)
	return t
}
