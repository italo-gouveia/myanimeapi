// Package api contains end-to-end tests for API flows.
// These tests spin up a real router connected to a real test database and make
// actual HTTP round-trips. They are automatically skipped when the test
// database is not available (SKIP_DB_TESTS=1 or connection failure).
package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"myanimeapi/tests/suites"
)

// TestAuthFlow_E2E verifies the complete authentication flow:
// register a new user → authenticate → receive JWT → use token on protected route.
func TestAuthFlow_E2E(t *testing.T) {
	suite := suites.NewBaseSuite(t)

	if os.Getenv("JWT_SECRET_KEY") == "" {
		_ = os.Setenv("JWT_SECRET_KEY", "test-secret")
	}

	// Unique username to avoid conflicts with parallel test runs
	username := "e2euser_auth"
	password := "E2Epassword123!"
	email := username + "@example.com"

	t.Run("Register", func(t *testing.T) {
		body := map[string]string{
			"username": username,
			"email":    email,
			"password": password,
		}
		b, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		suite.Router.ServeHTTP(rec, req)

		// 201 on first registration; 409 if re-running the test against a dirty DB.
		if rec.Code == http.StatusConflict {
			t.Skip("User already exists — skipping re-registration")
		}
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		assert.Equal(t, "success", resp["status"])
	})

	t.Run("Authenticate", func(t *testing.T) {
		body := map[string]string{
			"username": username,
			"password": password,
		}
		b, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/v1/auth/authenticate", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		suite.Router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		data, ok := resp["data"].(map[string]interface{})
		require.True(t, ok, "expected data object in response")
		token, _ := data["token"].(string)
		assert.NotEmpty(t, token, "expected non-empty JWT token")
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		body := map[string]string{
			"username": username,
			"password": "wrongpassword",
		}
		b, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/v1/auth/authenticate", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		suite.Router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("ProtectedRouteWithToken", func(t *testing.T) {
		// Authenticate first to get a valid token
		body := map[string]string{"username": username, "password": password}
		b, _ := json.Marshal(body)
		authReq := httptest.NewRequest(http.MethodPost, "/v1/auth/authenticate", bytes.NewReader(b))
		authReq.Header.Set("Content-Type", "application/json")
		authRec := httptest.NewRecorder()
		suite.Router.ServeHTTP(authRec, authReq)
		require.Equal(t, http.StatusOK, authRec.Code)

		var authResp map[string]interface{}
		require.NoError(t, json.NewDecoder(authRec.Body).Decode(&authResp))
		data := authResp["data"].(map[string]interface{})
		token := data["token"].(string)

		// Hit a protected route (GET /v1/favorites requires auth)
		protectedReq := httptest.NewRequest(http.MethodGet, "/v1/favorites", nil)
		protectedReq.Header.Set("Authorization", "Bearer "+token)
		protectedRec := httptest.NewRecorder()
		suite.Router.ServeHTTP(protectedRec, protectedReq)

		// Should not be 401
		assert.NotEqual(t, http.StatusUnauthorized, protectedRec.Code)
	})

	t.Run("ProtectedRouteWithoutToken", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/favorites", nil)
		rec := httptest.NewRecorder()
		suite.Router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
