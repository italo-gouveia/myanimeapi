// api/middleware/auth_middleware.go
// Package middleware provides middleware functions for authentication and authorization in the MyAnimeAPI application.
// It includes functions for validating JWT tokens, checking admin privileges, and generating JWT tokens with custom claims.
package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
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
		log := logger.Get()
		log.WithField("path", r.URL.Path).Debug("Authenticate middleware triggered")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Warn("Missing authorization header")
			errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "Missing authorization header", "The 'Authorization' header is required.")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil || !token.Valid {
			log.WithError(err).Warn("Invalid token")
			errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "Invalid token", "The provided token is invalid or expired.")
			return
		}

		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			log.Warn("Invalid token claims")
			errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "Invalid token claims", "The token claims are invalid or malformed.")
			return
		}

		// Store claims in context using the existing contextKey type
		ctx := context.WithValue(r.Context(), UserContextKey, claims.UserID)
		ctx = context.WithValue(ctx, IsAdminContextKey, claims.IsAdmin)
		log.WithFields(logrus.Fields{
			"user_id":  claims.UserID,
			"is_admin": claims.IsAdmin,
		}).Info("User authenticated")
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
		log := logger.Get()
		isAdmin, ok := r.Context().Value(IsAdminContextKey).(bool)
		if !ok || !isAdmin {
			log.Warn("Access denied: user is not an admin")
			errors.WriteErrorResponse(w, http.StatusForbidden, errors.ErrForbidden, "Access denied", "You do not have permission to access this resource.")
			return
		}

		log.Info("User is an admin, granting access")
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
	log := logger.Get()
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
		log.WithError(err).Error("Error generating token")
		return "", err
	}

	log.WithFields(logrus.Fields{
		"user_id":  userID,
		"is_admin": isAdmin,
	}).Info("Token generated successfully")
	return tokenString, nil
}
