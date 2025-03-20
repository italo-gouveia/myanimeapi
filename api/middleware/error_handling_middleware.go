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
	"log"
	"net/http"

	"myanimeapi/internal/errors"
)

// ErrorHandlingMiddleware is a middleware that catches and formats errors.
func ErrorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
				errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Internal Server Error", "An unexpected error occurred.")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
