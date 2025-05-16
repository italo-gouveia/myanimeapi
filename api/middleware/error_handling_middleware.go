// api/middleware/error_handling_middleware.go
// Package middleware provides HTTP middleware utilities for handling requests in the MyAnimeAPI application.
// It includes middleware functions for error handling, authentication, logging, and more.
//
// Middleware functions in this package are designed to wrap HTTP handlers and provide additional
// functionality such as panic recovery, request validation, and access control.
//
// Example Usage:
//
//	router := mux.NewRouter()
//	router.Use(middleware.ErrorHandlingMiddleware)
//	router.Use(middleware.LoggingMiddleware)
//	router.HandleFunc("/endpoint", myHandler)
//
// This package is essential for ensuring consistent behavior across all HTTP endpoints in the application.
package middleware

import (
	"net/http"

	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"

	"github.com/sirupsen/logrus"
)

// ErrorHandlingMiddleware provides robust error handling for HTTP requests.
// It recovers from panics, logs them, and returns a generic 500 Internal Server Error.
// It ensures that sensitive error details are not exposed to the client.
func ErrorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.Get()
		defer func() {
			if err := recover(); err != nil {
				requestID := r.Context().Value(RequestIDContextKey)
				log.WithFields(logrus.Fields{
					"request_id":  requestID,
					"error":       err,
					"method":      r.Method,
					"path":        r.URL.Path,
					"remote_addr": r.RemoteAddr,
					"user_agent":  r.UserAgent(),
				}).Error("Panic recovered in ErrorHandlingMiddleware")

				// Use the centralized error response writer
				errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "An unexpected error occurred. Please try again later.", "")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
