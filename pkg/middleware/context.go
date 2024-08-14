package middleware

import "context"

// Define a custom type for context keys
type contextKey string

const (
    userContextKey contextKey = "user"
)

// GetUserFromContext retrieves the user ID from the context
func GetUserFromContext(ctx context.Context) uint {
    if userID, ok := ctx.Value(userContextKey).(uint); ok {
        return userID
    }
    return 0
}
