package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
)

// TODO: parameterize this value
// JWTSecret - Placeholder for Jwt Secret
var JWTSecret = []byte("!!SECRET!!")

// GenerateJWT - Method to generate jwt token
func GenerateJWT(name string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = name
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, _ := token.SignedString(JWTSecret)
	return t
}
