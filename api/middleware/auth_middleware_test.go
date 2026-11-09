package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	apperrors "myanimeapi/internal/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserID(t *testing.T) {
	t.Run("returns user ID when present", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r = r.WithContext(context.WithValue(r.Context(), UserContextKey, uint(42)))

		id, err := GetUserID(r)

		require.NoError(t, err)
		assert.Equal(t, uint(42), id)
	})

	t.Run("returns error when absent", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		_, err := GetUserID(r)

		require.Error(t, err)
		appErr, ok := err.(*apperrors.AppError)
		require.True(t, ok)
		assert.Equal(t, apperrors.ErrUnauthorized, appErr.Code)
	})
}

func TestIsUserAdmin(t *testing.T) {
	t.Run("returns true for admin", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r = r.WithContext(context.WithValue(r.Context(), IsAdminContextKey, true))

		isAdmin, err := IsUserAdmin(r)

		require.NoError(t, err)
		assert.True(t, isAdmin)
	})

	t.Run("returns false for non-admin", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r = r.WithContext(context.WithValue(r.Context(), IsAdminContextKey, false))

		isAdmin, err := IsUserAdmin(r)

		require.NoError(t, err)
		assert.False(t, isAdmin)
	})

	t.Run("returns error when absent", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		_, err := IsUserAdmin(r)

		require.Error(t, err)
		appErr, ok := err.(*apperrors.AppError)
		require.True(t, ok)
		assert.Equal(t, apperrors.ErrUnauthorized, appErr.Code)
	})
}

func TestJwtExpiry(t *testing.T) {
	t.Run("defaults to 24h", func(t *testing.T) {
		_ = os.Unsetenv("JWT_EXPIRY_HOURS")
		assert.Equal(t, 24*time.Hour, jwtExpiry())
	})

	t.Run("reads from env", func(t *testing.T) {
		_ = os.Setenv("JWT_EXPIRY_HOURS", "48")
		defer func() { _ = os.Unsetenv("JWT_EXPIRY_HOURS") }()
		assert.Equal(t, 48*time.Hour, jwtExpiry())
	})

	t.Run("ignores invalid env value", func(t *testing.T) {
		_ = os.Setenv("JWT_EXPIRY_HOURS", "notanumber")
		defer func() { _ = os.Unsetenv("JWT_EXPIRY_HOURS") }()
		assert.Equal(t, 24*time.Hour, jwtExpiry())
	})

	t.Run("ignores zero value", func(t *testing.T) {
		_ = os.Setenv("JWT_EXPIRY_HOURS", "0")
		defer func() { _ = os.Unsetenv("JWT_EXPIRY_HOURS") }()
		assert.Equal(t, 24*time.Hour, jwtExpiry())
	})
}

func TestGenerateToken(t *testing.T) {
	_ = os.Setenv("JWT_SECRET_KEY", testSecretKey)
	defer func() { _ = os.Unsetenv("JWT_SECRET_KEY") }()

	t.Run("generates non-empty token", func(t *testing.T) {
		token, err := GenerateToken("123", false, "user")
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("generates admin token", func(t *testing.T) {
		token, err := GenerateToken("1", true, "admin")
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestGetRoleFromContext(t *testing.T) {
	t.Run("returns role when present", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r = r.WithContext(context.WithValue(r.Context(), RoleContextKey, "admin"))

		assert.Equal(t, "admin", GetRoleFromContext(r))
	})

	t.Run("returns empty string when absent", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		assert.Equal(t, "", GetRoleFromContext(r))
	})
}
