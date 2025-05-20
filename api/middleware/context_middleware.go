// api/middleware/context_middleware.go
// Package middleware defines a custom context key type and functions to retrieve user and admin status from the context.
// It is used by the authentication middleware to store and retrieve user and admin information in the request context.
package middleware

import (
	"context"
	"strconv"
)

// contextKey is a custom type for context keys to avoid key collisions.
type contextKey string

const (
	// UserContextKey is the context key for storing and retrieving the user ID.
	UserContextKey contextKey = "user"

	// IsAdminContextKey is the context key for storing and retrieving the admin status.
	IsAdminContextKey contextKey = "is_admin"

	// RequestIDContextKey is the context key for storing and retrieving the request ID.
	RequestIDContextKey contextKey = "request_id"

	// ValidatedPayloadKey is the context key for storing the validated and sanitized payload.
	ValidatedPayloadKey contextKey = "validated_payload"
)

// GetUserFromContext retrieves the user ID from the context.
// It returns the user ID as a uint if it exists in the context; otherwise, it returns 0.
// The user ID is stored as a string in the context and converted to uint.
func GetUserFromContext(ctx context.Context) uint {
	if userIDStr, ok := ctx.Value(UserContextKey).(string); ok {
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err == nil {
			return uint(userID)
		}
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
