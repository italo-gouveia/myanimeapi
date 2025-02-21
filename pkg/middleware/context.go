package middleware

import "context"

// Define a custom type for context keys
type contextKey string

const (
	userContextKey    contextKey = "user"
	isAdminContextKey contextKey = "is_admin"
)

// GetUserFromContext retrieves the user ID from the context
func GetUserFromContext(ctx context.Context) uint {
	if userID, ok := ctx.Value(userContextKey).(uint); ok {
		return userID
	}
	return 0
}

// GetIsAdminFromContext retrieves the admin status from the context
func GetIsAdminFromContext(ctx context.Context) bool {
	if isAdmin, ok := ctx.Value(isAdminContextKey).(bool); ok {
		return isAdmin
	}
	return false
}
