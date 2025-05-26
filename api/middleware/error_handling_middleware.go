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
	"context"
	"net/http"
	"runtime/debug"
	"time"

	apperrors "myanimeapi/internal/errors"
	"myanimeapi/internal/logger"

	"github.com/sirupsen/logrus"
)

// ErrorMetrics tracks error-related metrics
type ErrorMetrics struct {
	PanicCount    int64
	ErrorCount    int64
	LastErrorTime time.Time
}

var metrics = &ErrorMetrics{}

// ErrorHandlingMiddleware provides robust error handling for HTTP requests.
// It recovers from panics, logs them, and returns appropriate error responses.
// It ensures that sensitive error details are not exposed to the client while
// maintaining detailed logs for debugging.
func ErrorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.Get()
		start := time.Now()
		requestID := r.Context().Value(RequestIDContextKey)

		// Create a context with timeout
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()
		r = r.WithContext(ctx)

		// Create base log fields
		logFields := logrus.Fields{
			"request_id":  requestID,
			"method":      r.Method,
			"path":        r.URL.Path,
			"remote_addr": r.RemoteAddr,
			"user_agent":  r.UserAgent(),
		}

		defer func() {
			if err := recover(); err != nil {
				metrics.PanicCount++
				metrics.LastErrorTime = time.Now()

				// Get stack trace
				stack := debug.Stack()

				// Log the panic with detailed information
				log.WithFields(logFields).WithFields(logrus.Fields{
					"error":       err,
					"stack_trace": string(stack),
					"duration_ms": time.Since(start).Milliseconds(),
				}).Error("Panic recovered in ErrorHandlingMiddleware")

				// Use the centralized error response writer
				apperrors.WriteErrorResponse(w, http.StatusInternalServerError, apperrors.ErrInternalServer, "An unexpected error occurred. Please try again later.", "", logFields)
			}
		}()

		// Create a custom response writer to capture the status code
		lw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Handle the request
		next.ServeHTTP(lw, r)

		// Log errors based on status code
		if lw.statusCode >= 400 {
			metrics.ErrorCount++
			metrics.LastErrorTime = time.Now()

			log.WithFields(logFields).WithFields(logrus.Fields{
				"status_code": lw.statusCode,
				"duration_ms": time.Since(start).Milliseconds(),
			}).Error("Request completed with error")
		}
	})
}

// HandleError is a helper function to handle different types of errors
func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	log := logger.Get()
	requestID := r.Context().Value(RequestIDContextKey)

	logFields := logrus.Fields{
		"request_id":  requestID,
		"method":      r.Method,
		"path":        r.URL.Path,
		"remote_addr": r.RemoteAddr,
		"user_agent":  r.UserAgent(),
		"error":       err,
	}

	// Handle different types of errors
	switch {
	case apperrors.IsUnauthorized(err):
		log.WithFields(logFields).Warning("Unauthorized access attempt")
		apperrors.WriteErrorResponse(w, http.StatusUnauthorized, apperrors.ErrUnauthorized, "Unauthorized access", "You are not authorized to access this resource.", logFields)
	case apperrors.IsForbidden(err):
		log.WithFields(logFields).Warning("Forbidden access attempt")
		apperrors.WriteErrorResponse(w, http.StatusForbidden, apperrors.ErrForbidden, "Access forbidden", "You do not have permission to access this resource.", logFields)
	case apperrors.IsNotFound(err):
		log.WithFields(logFields).Info("Resource not found")
		apperrors.WriteErrorResponse(w, http.StatusNotFound, apperrors.ErrResourceNotFound, "Resource not found", "The requested resource could not be found.", logFields)
	case apperrors.IsInvalidInput(err):
		log.WithFields(logFields).Info("Invalid input")
		apperrors.WriteErrorResponse(w, http.StatusBadRequest, apperrors.ErrInvalidInput, "Invalid input", "The provided input is invalid.", logFields)
	default:
		log.WithFields(logFields).Error("Unexpected error")
		apperrors.WriteErrorResponse(w, http.StatusInternalServerError, apperrors.ErrInternalServer, "An unexpected error occurred", "Please try again later.", logFields)
	}
}

// GetErrorMetrics returns the current error metrics
func GetErrorMetrics() *ErrorMetrics {
	return metrics
}
