package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	// RequestIDHeader is the header key for the request ID
	RequestIDHeader = "X-Request-ID"
	// RequestIDContextKey is the context key for the request ID
	RequestIDContextKey = "request_id"
)

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get request ID from header if it exists
			requestID := r.Header.Get(RequestIDHeader)

			// Generate new request ID if none exists
			if requestID == "" {
				requestID = uuid.New().String()
			}

			// Add request ID to response headers
			w.Header().Set(RequestIDHeader, requestID)

			// Add request ID to request context
			ctx := r.Context()
			ctx = context.WithValue(ctx, RequestIDContextKey, requestID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
