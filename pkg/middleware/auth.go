// pkg/middleware/auth.go
// Package middleware provides middleware functions for authentication and authorization in the MyAnimeAPI application.
// It includes functions for validating JWT tokens, checking admin privileges, and generating JWT tokens with custom claims.
package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims defines the JWT claims structure.
// It includes the user ID and admin status.
type CustomClaims struct {
	IsAdmin bool `json:"is_admin"` // Indicates if the user is an admin
	UserID  uint `json:"user_id"`  // User ID
	jwt.RegisteredClaims
}

// Authenticate is a middleware function that checks for a valid JWT token in the request.
// If the token is valid, it stores the user ID and admin status in the request context.
// If the token is invalid or missing, it returns a 401 Unauthorized response.
//
// Example:
//
//	router.Use(middleware.Authenticate)
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Authenticate middleware triggered for: %s", r.URL.Path)
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("Missing authorization header")
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil || !token.Valid {
			log.Printf("Invalid token: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			log.Println("Invalid token claims")
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Store claims in context using the existing contextKey type
		ctx := context.WithValue(r.Context(), UserContextKey, claims.UserID)
		ctx = context.WithValue(ctx, IsAdminContextKey, claims.IsAdmin)
		log.Printf("User %d authenticated, isAdmin: %v", claims.UserID, claims.IsAdmin)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CheckAdmin is a middleware function that checks if the user has admin privileges.
// If the user is not an admin, it returns a 403 Forbidden response.
// If the user is an admin, it allows the request to proceed.
//
// Example:
//
//	router.Use(middleware.CheckAdmin)
func CheckAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value(IsAdminContextKey).(bool)
		if !ok || !isAdmin {
			log.Println("Access denied: user is not an admin")
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		log.Println("User is an admin, granting access")
		next.ServeHTTP(w, r)
	})
}

// GenerateToken generates a JWT token with custom claims.
// It includes the user ID, admin status, and an expiration time of 24 hours.
// The token is signed using the JWT secret key from the environment variables.
//
// Example:
//
//	token, err := middleware.GenerateToken(1, true)
//	if err != nil {
//	    log.Fatalf("Error generating token: %v", err)
//	}
func GenerateToken(userID uint, isAdmin bool) (string, error) {
	claims := CustomClaims{
		UserID:  userID,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    "myanimeapi",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return "", err
	}

	log.Printf("Token generated for user %d, isAdmin: %v", userID, isAdmin)
	return tokenString, nil
}
