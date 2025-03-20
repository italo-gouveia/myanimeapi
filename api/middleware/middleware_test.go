// api/middleware/middleware_test.go
// Package middleware provides HTTP middleware utilities for handling requests.
// This file contains tests for the middleware functions defined in the package.
package middleware

import (
	"context"
	"encoding/json"
	"myanimeapi/internal/errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const (
	testSecretKey = "test-secret-key" // Secret key for testing JWT
)

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// setup initializes the environment for testing.
func setup() {
	os.Setenv("JWT_SECRET_KEY", testSecretKey)
}

// teardown cleans up the environment after testing.
func teardown() {
	os.Unsetenv("JWT_SECRET_KEY")
}

// TestAuthenticateMiddleware tests the Authenticate middleware with a valid token.
func TestAuthenticateMiddleware(t *testing.T) {
	setup()
	defer teardown()

	t.Run("Valid Token", func(t *testing.T) {
		// Generate a valid token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
			UserID:  1,
			IsAdmin: true,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
			},
		})
		tokenString, _ := token.SignedString([]byte(testSecretKey))

		// Create a request with the token
		req, err := http.NewRequest("GET", "/protected", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Create a handler to use the middleware
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := GetUserFromContext(r.Context())
			isAdmin := GetIsAdminFromContext(r.Context())
			assert.Equal(t, uint(1), userID, "User ID should match")
			assert.Equal(t, true, isAdmin, "User should be admin")
			w.WriteHeader(http.StatusOK)
		})

		// Apply the middleware
		middleware := Authenticate(handler)
		middleware.ServeHTTP(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200")
	})

	t.Run("Invalid Token", func(t *testing.T) {
		// Create a request with an invalid token
		req, err := http.NewRequest("GET", "/protected", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer invalid-token")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Create a handler to use the middleware
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Should not reach this point")
		})

		// Apply the middleware
		middleware := Authenticate(handler)
		middleware.ServeHTTP(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Status code should be 401")
	})

	t.Run("Expired Token", func(t *testing.T) {
		// Generate an expired token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
			UserID:  1,
			IsAdmin: true,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour * 1)), // Expired 1 hour ago
			},
		})
		tokenString, _ := token.SignedString([]byte(testSecretKey))

		// Create a request with the expired token
		req, err := http.NewRequest("GET", "/protected", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Create a handler to use the middleware
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Should not reach this point")
		})

		// Apply the middleware
		middleware := Authenticate(handler)
		middleware.ServeHTTP(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Status code should be 401")
	})
}

// TestCheckAdminMiddleware tests the CheckAdmin middleware.
func TestCheckAdminMiddleware(t *testing.T) {
	t.Run("Admin User", func(t *testing.T) {
		// Create a request with admin context
		req, err := http.NewRequest("GET", "/admin", nil)
		assert.NoError(t, err)
		ctx := context.WithValue(req.Context(), IsAdminContextKey, true)
		req = req.WithContext(ctx)

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Create a handler to use the middleware
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Apply the middleware
		middleware := CheckAdmin(handler)
		middleware.ServeHTTP(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200")
	})

	t.Run("Non-Admin User", func(t *testing.T) {
		// Create a request with non-admin context
		req, err := http.NewRequest("GET", "/admin", nil)
		assert.NoError(t, err)
		ctx := context.WithValue(req.Context(), IsAdminContextKey, false)
		req = req.WithContext(ctx)

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Create a handler to use the middleware
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Should not reach this point")
		})

		// Apply the middleware
		middleware := CheckAdmin(handler)
		middleware.ServeHTTP(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusForbidden, rr.Code, "Status code should be 403")
	})
}

// TestGenerateToken tests the GenerateToken function.
func TestGenerateToken(t *testing.T) {
	setup()
	defer teardown()

	tokenString, err := GenerateToken(1, true)
	assert.NoError(t, err, "Token generation should not fail")
	assert.NotEmpty(t, tokenString, "Token should not be empty")

	// Parse the token to verify its contents
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(testSecretKey), nil
	})
	assert.NoError(t, err, "Token parsing should not fail")
	assert.True(t, token.Valid, "Token should be valid")

	claims, ok := token.Claims.(*CustomClaims)
	assert.True(t, ok, "Claims should be of type CustomClaims")
	assert.Equal(t, uint(1), claims.UserID, "User ID should match")
	assert.True(t, claims.IsAdmin, "User should be admin")
}

// TestLoggingMiddleware tests the LoggingMiddleware function.
func TestLoggingMiddleware(t *testing.T) {
	// Create a request
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a handler to use the middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Apply the middleware
	middleware := LoggingMiddleware(handler)
	middleware.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200")
}

// TestValidateAndSanitizePayloadMiddleware tests the ValidateAndSanitizePayload middleware.
func TestValidateAndSanitizePayloadMiddleware(t *testing.T) {
	type Payload struct {
		Name  string `json:"name" validate:"required,min=3,max=50"`
		Email string `json:"email" validate:"required,email"`
	}

	t.Run("Valid Payload", func(t *testing.T) {
		// Create a request with a valid payload
		payload := Payload{Name: "John Doe", Email: "john.doe@example.com"}
		payloadJSON, _ := json.Marshal(payload)
		req, err := http.NewRequest("POST", "/", strings.NewReader(string(payloadJSON)))
		assert.NoError(t, err)

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Create a handler to use the middleware
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			validatedPayload := r.Context().Value(ValidatedPayloadKey)
			assert.NotNil(t, validatedPayload, "Validated payload should not be nil")
			w.WriteHeader(http.StatusOK)
		})

		// Apply the middleware
		middleware := ValidateAndSanitizePayload(handler, Payload{})
		middleware.ServeHTTP(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200")
	})

	t.Run("Invalid Payload", func(t *testing.T) {
		// Create a request with an invalid payload
		payload := Payload{Name: "Jo", Email: "invalid-email"}
		payloadJSON, _ := json.Marshal(payload)
		req, err := http.NewRequest("POST", "/", strings.NewReader(string(payloadJSON)))
		assert.NoError(t, err)

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Create a handler to use the middleware
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Should not reach this point")
		})

		// Apply the middleware
		middleware := ValidateAndSanitizePayload(handler, Payload{})
		middleware.ServeHTTP(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code should be 400")
	})
}

// TestErrorHandlingMiddleware tests the ErrorHandlingMiddleware function.
func TestErrorHandlingMiddleware(t *testing.T) {
	t.Run("Panic Recovery", func(t *testing.T) {
		// Create a request
		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err)

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Create a handler that panics
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		})

		// Apply the middleware
		middleware := ErrorHandlingMiddleware(handler)
		middleware.ServeHTTP(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusInternalServerError, rr.Code, "Status code should be 500")

		// Check the error response
		var errResp errors.ErrorResponse
		err = json.NewDecoder(rr.Body).Decode(&errResp)
		assert.NoError(t, err, "Error response should be valid JSON")
		assert.Equal(t, "ERR-004", errResp.Error.Code, "Error code should match")
		assert.Equal(t, "Internal Server Error", errResp.Error.Message, "Error message should match")
	})
}
