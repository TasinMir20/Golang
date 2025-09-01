package utils

import (
	"fmt"
	"os"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
)

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

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	JWT_SECRET := os.Getenv("JWT_SECRET")
	if JWT_SECRET == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is not set")
	}

	// Clean the token string - remove common prefixes and whitespace
	tokenString = strings.TrimSpace(tokenString)
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	tokenString = strings.TrimPrefix(tokenString, "bearer ")

	// Validate basic JWT format (should have exactly 3 dots)
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format: expected 3 parts separated by dots, got %d", len(parts))
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("could not parse token claims")
	}

	return claims, nil
}
