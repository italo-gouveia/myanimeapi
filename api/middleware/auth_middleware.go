// api/middleware/auth_middleware.go
// Package middleware provides HTTP middleware utilities for handling requests in the MyAnimeAPI application.
// It includes middleware functions for error handling, authentication, logging, and more.
package middleware

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	apperrors "myanimeapi/internal/errors"
	"myanimeapi/internal/logger"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims represents the JWT claims structure
type CustomClaims struct {
	UserID  string `json:"user_id"`
	IsAdmin bool   `json:"is_admin"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, apperrors.NewError(apperrors.ErrUnauthorized, "Invalid token",
		"The provided token is invalid", http.StatusUnauthorized, nil, nil)
}

// AuthMiddleware handles JWT token validation and user authentication.
// It extracts the token from the Authorization header, validates it,
// and adds the user information to the request context.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.Get()

		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.WithFields(map[string]interface{}{
				"method":      r.Method,
				"path":        r.URL.Path,
				"remote_addr": r.RemoteAddr,
			}).Warning("Missing Authorization header")

			apperrors.WriteErrorResponse(w, http.StatusUnauthorized, apperrors.ErrUnauthorized,
				"Authentication required",
				"Missing Authorization header",
				map[string]interface{}{
					"error": "No authorization token provided",
				})
			return
		}

		// Check if the Authorization header has the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.WithFields(map[string]interface{}{
				"method":      r.Method,
				"path":        r.URL.Path,
				"remote_addr": r.RemoteAddr,
			}).Warning("Invalid Authorization header format")

			apperrors.WriteErrorResponse(w, http.StatusUnauthorized, apperrors.ErrUnauthorized,
				"Invalid token format",
				"Authorization header must be in the format: Bearer <token>",
				map[string]interface{}{
					"error": "Invalid authorization header format",
				})
			return
		}

		// Extract the token
		tokenString := parts[1]

		// Validate the token
		claims, err := ValidateToken(tokenString)
		if err != nil {
			log.WithFields(map[string]interface{}{
				"method":      r.Method,
				"path":        r.URL.Path,
				"remote_addr": r.RemoteAddr,
				"error":       err,
			}).Warning("Invalid token")

			apperrors.WriteErrorResponse(w, http.StatusUnauthorized, apperrors.ErrUnauthorized,
				"Invalid token",
				"The provided token is invalid or has expired",
				map[string]interface{}{
					"error": err.Error(),
				})
			return
		}

		// Add user info to context
		userID, err := strconv.ParseUint(claims.UserID, 10, 32)
		if err != nil {
			log.WithField("error", err).Error("Failed to parse user ID")
			apperrors.WriteErrorResponse(w, http.StatusUnauthorized, apperrors.ErrUnauthorized,
				"Invalid token",
				"The authentication token contains an invalid user ID",
				map[string]interface{}{
					"error": "Invalid user ID in token",
				})
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, uint(userID))
		ctx = context.WithValue(ctx, IsAdminContextKey, claims.IsAdmin)
		ctx = context.WithValue(ctx, RoleContextKey, claims.Role)

		// Log successful authentication
		log.WithFields(map[string]interface{}{
			"method":      r.Method,
			"path":        r.URL.Path,
			"remote_addr": r.RemoteAddr,
			"user_id":     userID,
			"is_admin":    claims.IsAdmin,
			"role":        claims.Role,
		}).Info("User authenticated successfully")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAdmin is a middleware that ensures the user has admin privileges.
// It should be used after the AuthMiddleware.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.Get()

		// Get admin status from context
		isAdmin, ok := r.Context().Value(IsAdminContextKey).(bool)
		if !ok || !isAdmin {
			log.WithFields(map[string]interface{}{
				"method":      r.Method,
				"path":        r.URL.Path,
				"remote_addr": r.RemoteAddr,
			}).Warning("Admin access required")

			apperrors.WriteErrorResponse(w, http.StatusForbidden, apperrors.ErrForbidden,
				"Admin access required",
				"You do not have permission to access this resource",
				map[string]interface{}{
					"error": "Admin privileges required",
				})
			return
		}

		// Continue with the request
		next.ServeHTTP(w, r)
	})
}

// GetUserID extracts the user ID from the request context.
// It should be used in handlers after the AuthMiddleware.
func GetUserID(r *http.Request) (uint, error) {
	userID, ok := r.Context().Value(UserContextKey).(uint)
	if !ok {
		return 0, apperrors.NewError(apperrors.ErrUnauthorized, "User not authenticated",
			"User ID not found in context", http.StatusUnauthorized, nil, nil)
	}
	return userID, nil
}

// IsUserAdmin checks if the user has admin privileges.
// It should be used in handlers after the AuthMiddleware.
func IsUserAdmin(r *http.Request) (bool, error) {
	isAdmin, ok := r.Context().Value(IsAdminContextKey).(bool)
	if !ok {
		return false, apperrors.NewError(apperrors.ErrUnauthorized, "User not authenticated",
			"Admin status not found in context", http.StatusUnauthorized, nil, nil)
	}
	return isAdmin, nil
}

// jwtExpiry returns the JWT expiry duration from JWT_EXPIRY_HOURS env var, defaulting to 24h.
func jwtExpiry() time.Duration {
	if h := os.Getenv("JWT_EXPIRY_HOURS"); h != "" {
		if hours, err := strconv.Atoi(h); err == nil && hours > 0 {
			return time.Duration(hours) * time.Hour
		}
	}
	return 24 * time.Hour
}

// GenerateToken generates a JWT token with custom claims
func GenerateToken(userID string, isAdmin bool, role string) (string, error) {
	claims := &CustomClaims{
		UserID:  userID,
		IsAdmin: isAdmin,
		Role:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpiry())),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

// GetRoleFromContext extracts the role from the request context.
// It should be used in handlers after the AuthMiddleware.
func GetRoleFromContext(r *http.Request) string {
	role, _ := r.Context().Value(RoleContextKey).(string)
	return role
}
