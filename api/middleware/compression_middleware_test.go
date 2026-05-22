package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// echoHandler writes a fixed JSON payload so tests have a known response body.
var echoHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message":"hello"}`))
})

func TestGzipMiddleware_CompressesWhenAccepted(t *testing.T) {
	handler := GzipMiddleware(echoHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "gzip", rr.Header().Get("Content-Encoding"))
	assert.Equal(t, "Accept-Encoding", rr.Header().Get("Vary"))

	// Body must be valid gzip that decompresses to the original payload.
	gz, err := gzip.NewReader(rr.Body)
	require.NoError(t, err)
	defer gz.Close()

	body, err := io.ReadAll(gz)
	require.NoError(t, err)
	assert.Equal(t, `{"message":"hello"}`, string(body))
}

func TestGzipMiddleware_SkipsWhenNotAccepted(t *testing.T) {
	handler := GzipMiddleware(echoHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	// No Accept-Encoding header.
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Empty(t, rr.Header().Get("Content-Encoding"))
	assert.Equal(t, `{"message":"hello"}`, rr.Body.String())
}

func TestGzipMiddleware_SkipsAlreadyEncoded(t *testing.T) {
	// Handler that pre-sets Content-Encoding (simulates upstream proxy).
	preEncoded := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "br")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("brotli-encoded-data"))
	})

	handler := GzipMiddleware(preEncoded)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Middleware must leave Content-Encoding unchanged.
	assert.Equal(t, "br", rr.Header().Get("Content-Encoding"))
	assert.Equal(t, "brotli-encoded-data", rr.Body.String())
}

func TestGzipMiddleware_AcceptsMultipleEncodings(t *testing.T) {
	handler := GzipMiddleware(echoHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "deflate, gzip, br")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, "gzip", rr.Header().Get("Content-Encoding"))
}

func TestGzipMiddleware_NonSuccessStatusCode(t *testing.T) {
	errorHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"not found"}`))
	})

	handler := GzipMiddleware(errorHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "gzip", rr.Header().Get("Content-Encoding"))

	gz, err := gzip.NewReader(rr.Body)
	require.NoError(t, err)
	defer gz.Close()

	body, err := io.ReadAll(gz)
	require.NoError(t, err)
	assert.Equal(t, `{"error":"not found"}`, string(body))
}

func TestGzipMiddleware_LargePayload(t *testing.T) {
	large := strings.Repeat("a", 100_000)
	largeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(large))
	})

	handler := GzipMiddleware(largeHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, "gzip", rr.Header().Get("Content-Encoding"))

	// Compressed body must be smaller than the raw payload.
	assert.Less(t, rr.Body.Len(), len(large))

	gz, err := gzip.NewReader(rr.Body)
	require.NoError(t, err)
	defer gz.Close()

	body, err := io.ReadAll(gz)
	require.NoError(t, err)
	assert.Equal(t, large, string(body))
}
