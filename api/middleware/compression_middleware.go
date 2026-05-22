// api/middleware/compression_middleware.go
// Package middleware provides HTTP middleware utilities for handling requests.
// This file defines a gzip compression middleware that compresses HTTP responses
// when the client signals support via the Accept-Encoding header.
package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
	"sync"
)

// gzipResponseWriter wraps http.ResponseWriter to transparently compress the response body.
// If the inner handler sets Content-Encoding itself (e.g. "br"), compression is skipped
// and writes go directly to the underlying ResponseWriter.
type gzipResponseWriter struct {
	http.ResponseWriter
	gz          *gzip.Writer
	wroteHeader bool
	skip        bool // true when inner handler already set Content-Encoding
}

func (grw *gzipResponseWriter) WriteHeader(code int) {
	if !grw.wroteHeader {
		grw.wroteHeader = true
		// If the inner handler set Content-Encoding, respect it and skip gzip.
		if grw.Header().Get("Content-Encoding") != "" {
			grw.skip = true
		} else {
			grw.Header().Set("Content-Encoding", "gzip")
			grw.Header().Set("Vary", "Accept-Encoding")
			// Remove Content-Length — it's no longer valid after compression.
			grw.Header().Del("Content-Length")
		}
		grw.ResponseWriter.WriteHeader(code)
	}
}

func (grw *gzipResponseWriter) Write(b []byte) (int, error) {
	if !grw.wroteHeader {
		grw.WriteHeader(http.StatusOK)
	}
	if grw.skip {
		return grw.ResponseWriter.Write(b)
	}
	return grw.gz.Write(b)
}

// Flush flushes both the gzip writer and the underlying response writer.
func (grw *gzipResponseWriter) Flush() {
	_ = grw.gz.Flush()
	if f, ok := grw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// gzipPool reuses gzip.Writer instances to reduce allocations.
var gzipPool = sync.Pool{
	New: func() interface{} {
		gz, _ := gzip.NewWriterLevel(nil, gzip.BestSpeed)
		return gz
	},
}

// GzipMiddleware compresses HTTP responses with gzip when the client supports it.
// It checks the Accept-Encoding request header; if "gzip" is present the response
// body is compressed and Content-Encoding / Vary headers are set accordingly.
// Responses that already carry a Content-Encoding header are left untouched.
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip if the client doesn't accept gzip.
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz := gzipPool.Get().(*gzip.Writer)
		gz.Reset(w)

		grw := &gzipResponseWriter{ResponseWriter: w, gz: gz}
		next.ServeHTTP(grw, r)

		if !grw.skip {
			_ = gz.Close()
		}
		gzipPool.Put(gz)
	})
}
