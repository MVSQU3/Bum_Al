package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SECRET = []byte("SECRET_KEY")

func GenerateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 10).Unix(),
	})

	return token.SignedString(SECRET)
}