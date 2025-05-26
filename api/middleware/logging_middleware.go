// api/middleware/logging_middleware.go
// Package middleware provides HTTP middleware utilities for handling requests.
// This file defines a logging middleware that logs incoming requests, including the method, path, and duration of the request.
package middleware

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"time"

	"myanimeapi/internal/logger"

	"github.com/google/uuid"
)

// LoggingMiddleware logs incoming requests and their processing time.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := r.Context().Value(RequestIDContextKey)
		if requestID == nil {
			requestID = uuid.New().String()
		}

		// Create base log fields
		logFields := map[string]interface{}{
			"request_id":   requestID,
			"method":       r.Method,
			"path":         r.URL.Path,
			"remote_addr":  r.RemoteAddr,
			"user_agent":   r.UserAgent(),
			"content_type": r.Header.Get("Content-Type"),
		}

		// Add query parameters for GET requests
		if r.Method == "GET" && len(r.URL.RawQuery) > 0 {
			logFields["query"] = r.URL.RawQuery
		}

		// Add request body size for POST/PUT requests
		if r.Method == "POST" || r.Method == "PUT" {
			logFields["content_length"] = r.ContentLength
		}

		// Log the incoming request details
		logger.WithFields(logFields).Info("Incoming request")

		lw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Wrap the handler in a panic recovery
		func() {
			defer func() {
				if err := recover(); err != nil {
					logger.WithFields(logFields).WithField("error", err).Error("Panic recovered in logging middleware")
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(lw, r)
		}()

		duration := time.Since(start)

		// Determine log level based on status code
		var level string
		if lw.statusCode >= 500 {
			level = "ERROR"
		} else if lw.statusCode >= 400 {
			level = "WARNING"
		} else {
			level = "INFO"
		}

		// Add response fields
		logFields["status"] = lw.statusCode
		logFields["duration_ms"] = duration.Milliseconds()
		logFields["size"] = lw.size

		// Log the request completion with appropriate level
		switch level {
		case "ERROR":
			logger.WithFields(logFields).Error("Request completed")
		case "WARNING":
			logger.WithFields(logFields).Warning("Request completed")
		default:
			logger.WithFields(logFields).Info("Request completed")
		}
	})
}

// loggingResponseWriter wraps http.ResponseWriter to capture status code and response size.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	if err != nil {
		logger.WithField("error", err).Error("Failed to write response")
	}
	lrw.size += size
	return size, err
}

// Ensure loggingResponseWriter implements http.Hijacker if the wrapped writer does.
func (lrw *loggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := lrw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("http.Hijacker not supported")
	}
	return h.Hijack()
}

// Ensure loggingResponseWriter implements http.Flusher if the wrapped writer does.
func (lrw *loggingResponseWriter) Flush() {
	f, ok := lrw.ResponseWriter.(http.Flusher)
	if ok {
		f.Flush()
	}
}
