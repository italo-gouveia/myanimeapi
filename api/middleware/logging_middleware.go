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
	"github.com/sirupsen/logrus"
)

// LoggingMiddleware logs incoming requests and their processing time.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.Get()
		start := time.Now()
		requestID := r.Context().Value(RequestIDContextKey)
		if requestID == nil {
			requestID = uuid.New().String()
		}

		lw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Log the incoming request details
		log.WithFields(logrus.Fields{
			"request_id":  requestID,
			"method":      r.Method,
			"path":        r.URL.Path,
			"remote_addr": r.RemoteAddr,
			"user_agent":  r.UserAgent(),
		}).Info("Incoming request")

		next.ServeHTTP(lw, r)

		duration := time.Since(start)

		// Determine log level based on status code
		level := logrus.InfoLevel
		if lw.statusCode >= 500 {
			level = logrus.ErrorLevel
		} else if lw.statusCode >= 400 {
			level = logrus.WarnLevel
		}

		// Log the request completion with appropriate level
		log.WithFields(logrus.Fields{
			"request_id": requestID,
			"method":     r.Method,
			"path":       r.URL.Path,
			"status":     lw.statusCode,
			"duration":   duration.String(),
			"size":       lw.size,
		}).Log(level, "Request completed")
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
