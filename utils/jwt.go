package utils

import (
	"fmt"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("my_secret_key")

func CreateSecretToken(id string, sessionUUID string) (string, error) {
	JWT_SECRET := os.Getenv("JWT_SECRET")

	claims := jwt.MapClaims{
		"sessionId":   id,
		"sessionUUID": sessionUUID,
		"sessionName": "LoginSession",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWT_SECRET))
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
}
