// api/middleware/logging_middleware.go
// Package middleware provides HTTP middleware utilities for handling requests.
// This file defines a logging middleware that logs incoming requests, including the method, path, and duration of the request.
package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware is an HTTP middleware that logs incoming requests.
// It logs the start and completion of each request, including the HTTP method, URL path, and the time taken to process the request.
// The middleware passes the request to the next handler in the chain after logging the start of the request.
// After the next handler completes, it logs the completion of the request along with the duration.
//
// Example usage:
//
//	http.Handle("/path", LoggingMiddleware(myHandler))
//
// This will log all requests to "/path" with their method, path, and duration.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}
