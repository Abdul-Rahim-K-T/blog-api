package jwt

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// ParseJWTToken parses the JWT token and returns the claims
func ParseJWTToken(tokenString, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil
}

// ExtractRoleFromToken extracts the role from the JWT token.
func ExtractRoleFromToken(r *http.Request, secretKey string) (string, error) {
	// Get the token from the Authorization header
	tokenString := r.Header.Get("Authorization")
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	// Parse the token
	claims, err := ParseJWTToken(tokenString, secretKey)
	if err != nil {
		return "", err
	}

	// Extract the role
	role, ok := claims["role"].(string)
	if !ok {
		return "", errors.New("role not found in token")
	}

	return role, nil
}

func ExtractClaims(r *http.Request, secretKey string) (*Claims, error) {
	tokenString := getTokenFromHeader(r)
	if tokenString == "" {
		return nil, errors.New("no token found")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func getTokenFromHeader(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if strings.HasPrefix(bearerToken, "Bearer ") {
		return strings.TrimPrefix(bearerToken, "Bearer ")
	}
	return ""
}
