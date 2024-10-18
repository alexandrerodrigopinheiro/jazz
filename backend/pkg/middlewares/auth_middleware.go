// backend/pkg/middlewares/auth_middleware.go
package middlewares

import (
	"context"
	"net/http"
	"strings"

	"jazz/backend/pkg/auth"
)

type contextKey string

const UserContextKey = contextKey("userID")

// AuthMiddleware checks if the request has a valid JWT token.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userID := claims["user_id"].(float64)
		ctx := context.WithValue(r.Context(), UserContextKey, uint(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
