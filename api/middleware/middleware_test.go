// api/middleware/middleware_test.go
// Package middleware provides HTTP middleware utilities for handling requests.
// This file contains tests for the middleware functions defined in the package.
package middleware

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"myanimeapi/internal/errors"
	apperrors "myanimeapi/internal/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const (
	testSecretKey = "test-secret-key" // Secret key for testing JWT
)

// setup initializes the environment for testing.
func setup() {
	os.Setenv("JWT_SECRET_KEY", testSecretKey)
}

// teardown cleans up the environment after testing.
func teardown() {
	os.Unsetenv("JWT_SECRET_KEY")
}

// TestRequestIDMiddleware tests the RequestID middleware functionality
func TestRequestIDMiddleware(t *testing.T) {
	t.Run("New Request ID", func(t *testing.T) {
		// Create a request
		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err)

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Create a handler to use the middleware
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := GetRequestIDFromContext(r.Context())
			assert.NotEmpty(t, requestID, "Request ID should not be empty")
			assert.Equal(t, requestID, w.Header().Get(RequestIDHeader), "Request ID in header should match context")
		})

		// Apply the middleware
		middleware := RequestIDMiddleware()(handler)
		middleware.ServeHTTP(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200")
	})

	t.Run("Existing Request ID", func(t *testing.T) {
		// Create a request with an existing request ID
		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err)
		existingID := "test-request-id"
		req.Header.Set(RequestIDHeader, existingID)

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Create a handler to use the middleware
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := GetRequestIDFromContext(r.Context())
			assert.Equal(t, existingID, requestID, "Request ID should match existing ID")
			assert.Equal(t, existingID, w.Header().Get(RequestIDHeader), "Request ID in header should match existing ID")
		})

		// Apply the middleware
		middleware := RequestIDMiddleware()(handler)
		middleware.ServeHTTP(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200")
	})
}

// TestHTTPSMiddleware tests the HTTPS middleware functionality
func TestHTTPSMiddleware(t *testing.T) {
	t.Run("HTTP Request", func(t *testing.T) {
		// Create an HTTP request
		req, err := http.NewRequest("GET", "http://example.com", nil)
		assert.NoError(t, err)

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Create a handler to use the middleware
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Should not reach this point")
		})

		// Apply the middleware
		middleware := HTTPSMiddleware(handler)
		middleware.ServeHTTP(rr, req)

		// Check the status code and location header
		assert.Equal(t, http.StatusPermanentRedirect, rr.Code, "Status code should be 308")
		assert.Equal(t, "https://example.com", rr.Header().Get("Location"), "Location header should point to HTTPS")
	})

	t.Run("HTTPS Request", func(t *testing.T) {
		// Create an HTTPS request
		req, err := http.NewRequest("GET", "https://example.com", nil)
		assert.NoError(t, err)
		req.TLS = &tls.ConnectionState{} // Simulate HTTPS

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Create a handler to use the middleware
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Apply the middleware
		middleware := HTTPSMiddleware(handler)
		middleware.ServeHTTP(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200")
	})
}

// TestRateLimiter tests the rate limiter functionality
func TestRateLimiter(t *testing.T) {
	rl := NewRateLimiter()

	t.Run("Normal Rate Limit", func(t *testing.T) {
		// Create a request
		req, err := http.NewRequest("GET", "/api/test", nil)
		assert.NoError(t, err)

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Create a handler to use the middleware
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Apply the middleware
		middleware := rl.RateLimitMiddleware(handler)

		// Make requests up to the limit
		for i := 0; i < 50; i++ {
			rr := httptest.NewRecorder()
			middleware.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200")
		}

		// Make one more request that should be rate limited
		rr = httptest.NewRecorder()
		middleware.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusTooManyRequests, rr.Code, "Status code should be 429")
	})

	t.Run("Auth Rate Limit", func(t *testing.T) {
		// Create a new rate limiter for auth endpoints
		rl := NewRateLimiter()

		// Create a request to an auth endpoint
		req, err := http.NewRequest("POST", "/auth/authenticate", nil)
		assert.NoError(t, err)

		// Create a handler to use the middleware
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Apply the middleware
		middleware := rl.RateLimitMiddleware(handler)

		// Make requests up to the limit (5 requests per minute for auth endpoints)
		for i := 0; i < 5; i++ {
			rr := httptest.NewRecorder()
			middleware.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200")
		}

		// Make one more request that should be rate limited
		rr := httptest.NewRecorder()
		middleware.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusTooManyRequests, rr.Code, "Status code should be 429")
	})

	t.Run("Rate Limit Reset", func(t *testing.T) {
		// Create a request
		req, err := http.NewRequest("GET", "/api/test", nil)
		assert.NoError(t, err)

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Create a handler to use the middleware
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Apply the middleware
		middleware := rl.RateLimitMiddleware(handler)

		// Set a custom clock that advances time
		now := time.Now()
		rl.SetClock(func() time.Time {
			return now.Add(time.Minute + time.Second)
		})

		// Make a request after the time window has passed
		middleware.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200")
	})
}

// TestErrorHandlingMiddleware tests the error handling middleware functionality
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

		// Check the status code and error response
		assert.Equal(t, http.StatusInternalServerError, rr.Code, "Status code should be 500")
		var errResp errors.ErrorResponse
		err = json.NewDecoder(rr.Body).Decode(&errResp)
		assert.NoError(t, err, "Error response should be valid JSON")
		assert.Equal(t, errors.ErrInternalServer, errResp.Error.Code, "Error code should match")
	})

	t.Run("Error Metrics", func(t *testing.T) {
		// Create a request
		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err)

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Create a handler that returns an error
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apperrors.WriteErrorResponse(w, http.StatusBadRequest, apperrors.ErrInvalidInput, "Invalid input", "The provided input is invalid", nil)
		})

		// Apply the middleware
		middleware := ErrorHandlingMiddleware(handler)
		middleware.ServeHTTP(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code should be 400")
		var errResp apperrors.ErrorResponse
		err = json.NewDecoder(rr.Body).Decode(&errResp)
		assert.NoError(t, err, "Error response should be valid JSON")
		assert.Equal(t, apperrors.ErrInvalidInput, errResp.Error.Code, "Error code should match")
	})
}

// TestValidateAndSanitizePayloadMiddleware tests the payload validation and sanitization middleware
func TestValidateAndSanitizePayloadMiddleware(t *testing.T) {
	type Payload struct {
		Name  string `json:"name" validate:"required,min=3,max=50"`
		Email string `json:"email" validate:"required,email"`
		HTML  string `json:"html"`
	}

	t.Run("Valid Payload", func(t *testing.T) {
		// Create a request with a valid payload
		payload := Payload{
			Name:  "John Doe",
			Email: "john.doe@example.com",
			HTML:  "<p>This is a <strong>test</strong> paragraph with <a href='https://example.com'>a link</a>.</p>",
		}
		payloadJSON, _ := json.Marshal(payload)
		req, err := http.NewRequest("POST", "/", strings.NewReader(string(payloadJSON)))
		assert.NoError(t, err)

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Create a handler to use the middleware
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			validatedPayload := GetValidatedPayloadFromContext(r.Context())
			assert.NotNil(t, validatedPayload, "Validated payload should not be nil")

			// Check if HTML was sanitized
			p, ok := validatedPayload.(*Payload)
			assert.True(t, ok, "Payload should be of type *Payload")
			assert.Equal(t, "<p>This is a <strong>test</strong> paragraph with <a href=\"https://example.com\" rel=\"nofollow\">a link</a>.</p>", p.HTML, "HTML should be sanitized")

			w.WriteHeader(http.StatusOK)
		})

		// Apply the middleware
		middleware := ValidateAndSanitizePayload(Payload{})(handler)
		middleware.ServeHTTP(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200")
	})

	t.Run("Invalid Payload", func(t *testing.T) {
		// Create a request with an invalid payload
		payload := Payload{
			Name:  "Jo", // Too short
			Email: "invalid-email",
			HTML:  "<script>alert('xss')</script>",
		}
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
		middleware := ValidateAndSanitizePayload(Payload{})(handler)
		middleware.ServeHTTP(rr, req)

		// Check the status code and error response
		assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code should be 400")
		var resp map[string]interface{}
		err = json.NewDecoder(rr.Body).Decode(&resp)
		assert.NoError(t, err, "Error response should be valid JSON")
		errorsMap, ok := resp["errors"].(map[string]interface{})
		assert.True(t, ok, "Response should contain 'errors' field")
		assert.Contains(t, errorsMap, "Name", "Errors should contain 'Name' field")
		assert.Contains(t, errorsMap, "Email", "Errors should contain 'Email' field")
		assert.Equal(t, "Name must be at least 3 characters long.", errorsMap["Name"], "Name error message should match")
		assert.Equal(t, "Validation failed for Email.", errorsMap["Email"], "Email error message should match")
	})
}

// TestAuthenticateMiddleware tests the authentication middleware
func TestAuthenticateMiddleware(t *testing.T) {
	setup()
	defer teardown()

	t.Run("Valid Token", func(t *testing.T) {
		// Generate a valid token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
			UserID:  "1",
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
		middleware := AuthMiddleware(handler)
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
		middleware := AuthMiddleware(handler)
		middleware.ServeHTTP(rr, req)

		// Check the status code and error response
		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Status code should be 401")
		var errResp errors.ErrorResponse
		err = json.NewDecoder(rr.Body).Decode(&errResp)
		assert.NoError(t, err, "Error response should be valid JSON")
		assert.Equal(t, errors.ErrUnauthorized, errResp.Error.Code, "Error code should match")
	})

	t.Run("Expired Token", func(t *testing.T) {
		// Generate an expired token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
			UserID:  "1",
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
		middleware := AuthMiddleware(handler)
		middleware.ServeHTTP(rr, req)

		// Check the status code and error response
		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Status code should be 401")
		var errResp errors.ErrorResponse
		err = json.NewDecoder(rr.Body).Decode(&errResp)
		assert.NoError(t, err, "Error response should be valid JSON")
		assert.Equal(t, errors.ErrUnauthorized, errResp.Error.Code, "Error code should match")
	})
}

// TestRequireAdmin tests the admin requirement middleware
func TestRequireAdmin(t *testing.T) {
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
		middleware := RequireAdmin(handler)
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
		middleware := RequireAdmin(handler)
		middleware.ServeHTTP(rr, req)

		// Check the status code and error response
		assert.Equal(t, http.StatusForbidden, rr.Code, "Status code should be 403")
		var errResp errors.ErrorResponse
		err = json.NewDecoder(rr.Body).Decode(&errResp)
		assert.NoError(t, err, "Error response should be valid JSON")
		assert.Equal(t, errors.ErrForbidden, errResp.Error.Code, "Error code should match")
	})
}
