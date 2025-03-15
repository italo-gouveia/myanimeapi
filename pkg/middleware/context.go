// pkg/middleware/context.go
// Package middleware defines a custom context key type and functions to retrieve user and admin status from the context.
// It is used by the authentication middleware to store and retrieve user and admin information in the request context.
package middleware

import "context"

// contextKey is a custom type for context keys to avoid key collisions.
type contextKey string

const (
	// UserContextKey is the context key for storing and retrieving the user ID.
	UserContextKey contextKey = "user"

	// IsAdminContextKey is the context key for storing and retrieving the admin status.
	IsAdminContextKey contextKey = "is_admin"
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
