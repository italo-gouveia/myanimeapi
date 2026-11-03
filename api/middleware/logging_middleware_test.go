package middleware

import (
	"bufio"
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildRequestLogFields(t *testing.T) {
	t.Run("GET with query and request_id in context", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/animes?page=1", nil)
		r = r.WithContext(context.WithValue(r.Context(), RequestIDContextKey, "test-id"))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("User-Agent", "go-test/1.0")

		fields := buildRequestLogFields(r)

		assert.Equal(t, "test-id", fields["request_id"])
		assert.Equal(t, http.MethodGet, fields["method"])
		assert.Equal(t, "/animes", fields["path"])
		assert.Equal(t, "page=1", fields["query"])
		assert.Equal(t, "application/json", fields["content_type"])
		assert.NotContains(t, fields, "content_length")
	})

	t.Run("POST with content_length", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/reviews", nil)
		r.ContentLength = 42

		fields := buildRequestLogFields(r)

		assert.Equal(t, int64(42), fields["content_length"])
		assert.NotContains(t, fields, "query")
	})

	t.Run("PUT with content_length", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPut, "/reviews/1", nil)
		r.ContentLength = 10

		fields := buildRequestLogFields(r)

		assert.Equal(t, int64(10), fields["content_length"])
	})

	t.Run("generates request_id when none in context", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		fields := buildRequestLogFields(r)

		reqID, ok := fields["request_id"].(string)
		require.True(t, ok)
		assert.NotEmpty(t, reqID)
	})

	t.Run("GET without query does not include query field", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/animes", nil)

		fields := buildRequestLogFields(r)

		assert.NotContains(t, fields, "query")
	})
}

func TestLogCompletion(t *testing.T) {
	// logCompletion writes to logger; we verify it doesn't panic for each branch.
	t.Run("2xx logs info", func(t *testing.T) {
		assert.NotPanics(t, func() {
			logCompletion(map[string]interface{}{"path": "/"}, http.StatusOK)
		})
	})

	t.Run("4xx logs warning", func(t *testing.T) {
		assert.NotPanics(t, func() {
			logCompletion(map[string]interface{}{"path": "/"}, http.StatusBadRequest)
		})
	})

	t.Run("5xx logs error", func(t *testing.T) {
		assert.NotPanics(t, func() {
			logCompletion(map[string]interface{}{"path": "/"}, http.StatusInternalServerError)
		})
	})
}

func TestServeWithRecovery(t *testing.T) {
	t.Run("passes request to next handler", func(t *testing.T) {
		called := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		lw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		serveWithRecovery(next, lw, r, w, map[string]interface{}{})

		assert.True(t, called)
	})

	t.Run("recovers from panic", func(t *testing.T) {
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		lw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		assert.NotPanics(t, func() {
			serveWithRecovery(next, lw, r, w, map[string]interface{}{})
		})
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestLoggingMiddleware(t *testing.T) {
	t.Run("calls next handler and logs", func(t *testing.T) {
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		handler := LoggingMiddleware(next)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/test", nil)

		assert.NotPanics(t, func() { handler.ServeHTTP(w, r) })
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("logs 500 for panicking handler", func(t *testing.T) {
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("boom")
		})

		handler := LoggingMiddleware(next)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/test", nil)
		r.ContentLength = 5

		assert.NotPanics(t, func() { handler.ServeHTTP(w, r) })
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestLoggingResponseWriter(t *testing.T) {
	t.Run("WriteHeader captures status code", func(t *testing.T) {
		w := httptest.NewRecorder()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		lrw.WriteHeader(http.StatusCreated)

		assert.Equal(t, http.StatusCreated, lrw.statusCode)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Write accumulates size", func(t *testing.T) {
		w := httptest.NewRecorder()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		n, err := lrw.Write([]byte("hello"))

		assert.NoError(t, err)
		assert.Equal(t, 5, n)
		assert.Equal(t, 5, lrw.size)
	})

	t.Run("Hijack returns error when wrapped writer is not hijacker", func(t *testing.T) {
		w := httptest.NewRecorder()
		lrw := &loggingResponseWriter{ResponseWriter: w}

		conn, rw, err := lrw.Hijack()
		assert.Nil(t, conn)
		assert.Nil(t, rw)
		assert.Error(t, err)
	})

	t.Run("Hijack succeeds when wrapped writer supports it", func(t *testing.T) {
		w := &hijackableRecorder{ResponseRecorder: httptest.NewRecorder()}
		lrw := &loggingResponseWriter{ResponseWriter: w}

		conn, _, err := lrw.Hijack()
		assert.NoError(t, err)
		assert.NotNil(t, conn)
	})

	t.Run("Flush calls underlying flusher", func(t *testing.T) {
		w := httptest.NewRecorder()
		lrw := &loggingResponseWriter{ResponseWriter: w}
		assert.NotPanics(t, func() { lrw.Flush() })
	})

	t.Run("Flush with flusher calls Flush", func(t *testing.T) {
		fw := &flushableRecorder{ResponseRecorder: httptest.NewRecorder()}
		lrw := &loggingResponseWriter{ResponseWriter: fw}
		lrw.Flush()
		assert.True(t, fw.flushed)
	})
}

// hijackableRecorder is a test double that implements http.Hijacker.
type hijackableRecorder struct {
	*httptest.ResponseRecorder
}

func (h *hijackableRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	server, client := net.Pipe()
	go func() { _ = server.Close() }()
	br := bufio.NewReadWriter(bufio.NewReader(client), bufio.NewWriter(client))
	return client, br, nil
}

// flushableRecorder is a test double that implements http.Flusher.
type flushableRecorder struct {
	*httptest.ResponseRecorder
	flushed bool
}

func (f *flushableRecorder) Flush() {
	f.flushed = true
}
