package middleware

import (
	"context"
	"fmt"
	"moviesapi/internal/shared"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// jwtSecret validates tokens issued by the user-service. It defaults to a dev
// value and must be overridden with SetJWTSecret using the same SECERET_KEY
// the user-service signs with, otherwise every authenticated request 401s.
var jwtSecret = []byte("super-secret-key-for-development")

// SetJWTSecret overrides the development default with the configured secret.
func SetJWTSecret(secret string) {
	if secret != "" {
		jwtSecret = []byte(secret)
	}
}

// Auth middleware validates the JWT token and extracts the User ID
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Bypass Auth for public GET movie endpoints
		if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/v1/movies") {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			shared.Log.Warn("missing authorization header")
			http.Error(w, "Unauthorized: Missing Authorization Header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			shared.Log.Warn("invalid authorization format")
			http.Error(w, "Unauthorized: Invalid Authorization Format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what you expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			shared.Log.Warn("invalid token", zap.Error(err))
			http.Error(w, "Unauthorized: Invalid Token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// The user-service signs tokens with a custom "id" claim (Customclaim.ID);
			// fall back to the standard "sub" claim for other issuers.
			userID, _ := claims["id"].(string)
			if userID == "" {
				userID, _ = claims["sub"].(string)
			}
			if userID == "" {
				shared.Log.Warn("token has no usable subject claim")
				http.Error(w, "Unauthorized: Token Missing Subject", http.StatusUnauthorized)
				return
			}

			// Inject User ID into headers for downstream services
			r.Header.Set("X-User-ID", userID)

			// Optional: store in context if we need it in the gateway itself
			ctx := context.WithValue(r.Context(), "userID", userID)
			r = r.WithContext(ctx)
		} else {
			shared.Log.Warn("invalid token claims")
			http.Error(w, "Unauthorized: Invalid Token Claims", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
