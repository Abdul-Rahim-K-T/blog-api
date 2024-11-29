package middleware

import (
	"blog-api/pkg/jwt"
	"context"
	"log"
	"net/http"
)

// // AuthorMiddleware checks if the user is an author by JWT Tken.
// func AuthorMiddleware(secretKey string) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			role, err := jwt.ExtractRoleFromToken(r, secretKey)
// 			log.Printf("ROLE:%s, %v+", role, err)
// 			if err != nil || role != "author" {
// 				http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 				return
// 			}

// 			// Lost the claims for debugging purposes
// 			log.Printf("Authorization claims: %+v", claims)
// 			log.Println("Author middleware")
// 			// Continue with the request
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

func AuthorMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := jwt.ExtractClaims(r, secretKey)
			if err != nil || claims.Role != "author" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Log the claims for debugging purposes
			log.Printf("AuthorMiddleware claims: %+v", claims)

			// Add user ID to context
			ctx := context.WithValue(r.Context(), jwt.UserIDKey, claims.UserID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
