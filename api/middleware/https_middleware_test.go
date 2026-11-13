package middleware

import (
	"crypto/tls"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSecureRequest(t *testing.T) {
	t.Run("returns true for TLS request", func(t *testing.T) {
		r := httptest.NewRequest("GET", "https://example.com/", nil)
		r.TLS = &tls.ConnectionState{}
		assert.True(t, IsSecureRequest(r))
	})

	t.Run("returns false for non-TLS request", func(t *testing.T) {
		r := httptest.NewRequest("GET", "http://example.com/", nil)
		assert.False(t, IsSecureRequest(r))
	})
}

func TestGetProtocol(t *testing.T) {
	t.Run("returns https for TLS", func(t *testing.T) {
		r := httptest.NewRequest("GET", "https://example.com/", nil)
		r.TLS = &tls.ConnectionState{}
		assert.Equal(t, "https", GetProtocol(r))
	})

	t.Run("returns http for non-TLS", func(t *testing.T) {
		r := httptest.NewRequest("GET", "http://example.com/", nil)
		assert.Equal(t, "http", GetProtocol(r))
	})
}
