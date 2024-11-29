package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWTToken(userID int, role string, secretKey string) (string, error) {
	// Define claims for the jwt
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expiration
	}

	// Create a new JWT token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and return the token
	return token.SignedString([]byte(secretKey))
}
