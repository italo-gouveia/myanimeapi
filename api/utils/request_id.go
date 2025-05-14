package utils

import (
	"context"
	"myanimeapi/api/middleware"
)

// GetRequestID retrieves the request ID from the context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(middleware.RequestIDContextKey).(string); ok {
		return requestID
	}
	return ""
}
