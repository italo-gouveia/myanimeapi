package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStrictRateLimitMiddleware(t *testing.T) {
	t.Run("allows requests within limit", func(t *testing.T) {
		rl := NewRateLimiter()
		fixed := time.Now()
		rl.SetClock(func() time.Time { return fixed })

		handler := rl.StrictRateLimitMiddleware(5, time.Minute)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		for i := 0; i < 5; i++ {
			req := httptest.NewRequest(http.MethodPost, "/auth/reset", nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code, "request %d should succeed", i+1)
		}
	})

	t.Run("blocks requests exceeding limit", func(t *testing.T) {
		rl := NewRateLimiter()
		fixed := time.Now()
		rl.SetClock(func() time.Time { return fixed })

		handler := rl.StrictRateLimitMiddleware(2, time.Minute)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		// First 2 OK
		for i := 0; i < 2; i++ {
			req := httptest.NewRequest(http.MethodPost, "/auth/reset", nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}

		// 3rd blocked
		req := httptest.NewRequest(http.MethodPost, "/auth/reset", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})

	t.Run("window resets after expiry", func(t *testing.T) {
		rl := NewRateLimiter()
		now := time.Now()
		rl.SetClock(func() time.Time { return now })

		handler := rl.StrictRateLimitMiddleware(1, time.Minute)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		// First request OK
		req := httptest.NewRequest(http.MethodPost, "/reset", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// Advance clock past window
		now = now.Add(2 * time.Minute)
		req2 := httptest.NewRequest(http.MethodPost, "/reset", nil)
		w2 := httptest.NewRecorder()
		handler.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)
	})
}

func TestClientIP(t *testing.T) {
	t.Run("uses X-Forwarded-For", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("X-Forwarded-For", "1.2.3.4, 5.6.7.8")
		assert.Equal(t, "1.2.3.4", clientIP(r))
	})

	t.Run("uses X-Real-IP", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("X-Real-IP", "9.8.7.6")
		assert.Equal(t, "9.8.7.6", clientIP(r))
	})

	t.Run("falls back to RemoteAddr", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.RemoteAddr = "192.168.1.1:4321"
		ip := clientIP(r)
		assert.Equal(t, "192.168.1.1", ip)
	})
}
