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
		logFields := buildRequestLogFields(r)
		logger.WithFields(logFields).Info("Incoming request")

		lw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		serveWithRecovery(next, lw, r, w, logFields)

		logFields["status"] = lw.statusCode
		logFields["duration_ms"] = time.Since(start).Milliseconds()
		logFields["size"] = lw.size
		logCompletion(logFields, lw.statusCode)
	})
}

func buildRequestLogFields(r *http.Request) map[string]interface{} {
	requestID := r.Context().Value(RequestIDContextKey)
	if requestID == nil {
		requestID = uuid.New().String()
	}
	fields := map[string]interface{}{
		"request_id":   requestID,
		"method":       r.Method,
		"path":         r.URL.Path,
		"remote_addr":  r.RemoteAddr,
		"user_agent":   r.UserAgent(),
		"content_type": r.Header.Get("Content-Type"),
	}
	if r.Method == "GET" && len(r.URL.RawQuery) > 0 {
		fields["query"] = r.URL.RawQuery
	}
	if r.Method == "POST" || r.Method == "PUT" {
		fields["content_length"] = r.ContentLength
	}
	return fields
}

func serveWithRecovery(next http.Handler, lw *loggingResponseWriter, r *http.Request, w http.ResponseWriter, logFields map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logger.WithFields(logFields).WithField("error", err).Error("Panic recovered in logging middleware")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()
	next.ServeHTTP(lw, r)
}

func logCompletion(logFields map[string]interface{}, statusCode int) {
	switch {
	case statusCode >= 500:
		logger.WithFields(logFields).Error("Request completed")
	case statusCode >= 400:
		logger.WithFields(logFields).Warning("Request completed")
	default:
		logger.WithFields(logFields).Info("Request completed")
	}
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
