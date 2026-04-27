package httphandler

import (
	"context"
	"myanimeapi/api/middleware"
)

// CreateTestContext creates a context with middleware values for testing
func CreateTestContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.UserContextKey, uint(1))
	return ctx
}
