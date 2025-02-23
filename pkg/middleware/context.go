// pkg/middleware/context.go
// This package defines a custom context key type and functions to retrieve user and admin status from the context.
// It is used by the authentication middleware to store and retrieve user and admin information in the request context.
package middleware

import "context"

// Define a custom type for context keys
type contextKey string

const (
	UserContextKey    contextKey = "user"     // Exported context key for user ID
	IsAdminContextKey contextKey = "is_admin" // Exported context key for admin status
)

// GetUserFromContext retrieves the user ID from the context
func GetUserFromContext(ctx context.Context) uint {
	if userID, ok := ctx.Value(UserContextKey).(uint); ok {
		return userID
	}
	return 0
}

// GetIsAdminFromContext retrieves the admin status from the context
func GetIsAdminFromContext(ctx context.Context) bool {
	if isAdmin, ok := ctx.Value(IsAdminContextKey).(bool); ok {
		return isAdmin
	}
	return false
}
