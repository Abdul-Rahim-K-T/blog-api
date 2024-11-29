package middleware

import (
	"blog-api/pkg/jwt"
	"context"
	"net/http"
	"strings"
)

// AuthMiddleware checks for valid JWT token in the Authorization header
func AuthMiddleware(secretKey string, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		// Extract the token from the Authorization header
		tokenString := strings.Split(authHeader, "Bearer ")[1]
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// Parse and validate the JWT token
		claims, err := jwt.ParseJWTToken(tokenString, secretKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Add user ID from claims to request context
		userID := claims["user_id"].(float64) // Assuming user_id is in the claims as a float64
		ctx := context.WithValue(r.Context(), "user_id", int(userID))

		// Pass the context to the next handler
		r = r.WithContext(ctx)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	}
}
