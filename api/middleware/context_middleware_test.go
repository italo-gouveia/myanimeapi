package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequireRole(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("admin always passes", func(t *testing.T) {
		handler := RequireRole("reviewer")(next)
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r = r.WithContext(context.WithValue(r.Context(), IsAdminContextKey, true))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("allowed role passes", func(t *testing.T) {
		handler := RequireRole("reviewer")(next)
		ctx := context.WithValue(context.Background(), IsAdminContextKey, false)
		ctx = context.WithValue(ctx, RoleContextKey, "reviewer")
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r = r.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("disallowed role returns 403", func(t *testing.T) {
		handler := RequireRole("reviewer")(next)
		ctx := context.WithValue(context.Background(), IsAdminContextKey, false)
		ctx = context.WithValue(ctx, RoleContextKey, "user")
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r = r.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, r)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("no role in context returns 403", func(t *testing.T) {
		handler := RequireRole("reviewer")(next)
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, r)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
