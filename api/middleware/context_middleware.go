// api/middleware/context_middleware.go
// Package middleware defines a custom context key type and functions to retrieve user and admin status from the context.
// It is used by the authentication middleware to store and retrieve user and admin information in the request context.
package middleware

import (
	"context"
	"net/http"

	"myanimeapi/internal/errors"
)

// contextKey is a custom type for context keys to avoid key collisions.
type contextKey string

const (
	// UserContextKey is the context key for storing and retrieving the user ID.
	UserContextKey contextKey = "user"

	// IsAdminContextKey is the context key for storing and retrieving the admin status.
	IsAdminContextKey contextKey = "is_admin"

	// RoleContextKey is the context key for storing and retrieving the user role.
	RoleContextKey contextKey = "role"

	// RequestIDContextKey is the context key for storing and retrieving the request ID.
	RequestIDContextKey contextKey = "request_id"

	// ValidatedPayloadKey is the context key for storing the validated and sanitized payload.
	ValidatedPayloadKey contextKey = "validated_payload"
)

// GetUserFromContext retrieves the user ID from the context.
// It returns the user ID as a uint if it exists in the context; otherwise, it returns 0.
func GetUserFromContext(ctx context.Context) uint {
	if userID, ok := ctx.Value(UserContextKey).(uint); ok {
		return userID
	}
	return 0
}

// GetIsAdminFromContext retrieves the admin status from the context.
// It returns a boolean indicating whether the user is an admin if the status exists in the context; otherwise, it returns false.
func GetIsAdminFromContext(ctx context.Context) bool {
	if isAdmin, ok := ctx.Value(IsAdminContextKey).(bool); ok {
		return isAdmin
	}
	return false
}

// GetRequestIDFromContext retrieves the request ID from the context.
// It returns the request ID as a string if it exists in the context; otherwise, it returns an empty string.
func GetRequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDContextKey).(string); ok {
		return requestID
	}
	return ""
}

// GetValidatedPayloadFromContext retrieves the validated payload from the context.
// It returns the payload as an interface{} if it exists in the context; otherwise, it returns nil.
func GetValidatedPayloadFromContext(ctx context.Context) interface{} {
	return ctx.Value(ValidatedPayloadKey)
}

// RequireRole returns a middleware that allows access only to users with one of the given roles.
// Admins (is_admin=true) always pass regardless of role parameter.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isAdmin, _ := r.Context().Value(IsAdminContextKey).(bool)
			if isAdmin {
				next.ServeHTTP(w, r)
				return
			}
			role, _ := r.Context().Value(RoleContextKey).(string)
			if !allowed[role] {
				errors.WriteErrorResponse(w, http.StatusForbidden, errors.ErrForbidden, "Forbidden", "Insufficient role", nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
