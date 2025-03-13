// pkg/middleware/middleware_test.go
package middleware

import (
	"context"
	"encoding/json"
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
	testSecretKey = "test-secret-key"
)

func setup() {
	os.Setenv("JWT_SECRET_KEY", testSecretKey)
}

func teardown() {
	os.Unsetenv("JWT_SECRET_KEY")
}

func TestAuthenticateMiddleware(t *testing.T) {
	setup()
	defer teardown()

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
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+tokenString)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a handler to use the middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserFromContext(r.Context())
		isAdmin := GetIsAdminFromContext(r.Context())
		assert.Equal(t, uint(1), userID)
		assert.Equal(t, true, isAdmin)
		w.WriteHeader(http.StatusOK)
	})

	// Apply the middleware
	middleware := Authenticate(handler)
	middleware.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthenticateMiddleware_InvalidToken(t *testing.T) {
	setup()
	defer teardown()

	// Create a request with an invalid token
	req, err := http.NewRequest("GET", "/protected", nil)
	if err != nil {
		t.Fatal(err)
	}
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
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestCheckAdminMiddleware(t *testing.T) {
	// Create a request with admin context
	req, err := http.NewRequest("GET", "/admin", nil)
	if err != nil {
		t.Fatal(err)
	}
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
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestCheckAdminMiddleware_NonAdmin(t *testing.T) {
	// Create a request with non-admin context
	req, err := http.NewRequest("GET", "/admin", nil)
	if err != nil {
		t.Fatal(err)
	}
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
	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestGenerateToken(t *testing.T) {
	setup()
	defer teardown()

	tokenString, err := GenerateToken(1, true)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Parse the token to verify its contents
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(testSecretKey), nil
	})
	assert.NoError(t, err)
	assert.True(t, token.Valid)

	claims, ok := token.Claims.(*CustomClaims)
	assert.True(t, ok)
	assert.Equal(t, uint(1), claims.UserID)
	assert.True(t, claims.IsAdmin)
}

func TestLoggingMiddleware(t *testing.T) {
	// Create a request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

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
	assert.Equal(t, http.StatusOK, rr.Code)
}

/*func TestRateLimiterMiddleware(t *testing.T) {
	rl := NewRateLimiter()

	// Create a custom clock to control time in the test
	now := time.Now()
	rl.SetClock(func() time.Time {
		return now
	})

	// Create a request for an authentication endpoint (stricter rate limit)
	req, err := http.NewRequest("GET", "/auth/authenticate", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Simulate a client IP address
	req.RemoteAddr = "127.0.0.1:12345"

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a handler to use the middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Apply the middleware
	middleware := rl.RateLimitMiddleware(handler)

	// Make requests within the limit (5 requests allowed)
	for i := 0; i < 5; i++ {
		middleware.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	}

	// Make one more request to exceed the limit
	middleware.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)

	// Advance the clock by 1 minute to reset the rate limit
	now = now.Add(time.Minute)

	// Make another request (should be allowed again)
	middleware.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}*/

func TestValidateAndSanitizePayloadMiddleware(t *testing.T) {
	type Payload struct {
		Name  string `json:"name" validate:"required,min=3,max=50"`
		Email string `json:"email" validate:"required,email"`
	}

	// Create a request with a valid payload
	payload := Payload{Name: "John Doe", Email: "john.doe@example.com"}
	payloadJSON, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/", strings.NewReader(string(payloadJSON)))
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a handler to use the middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validatedPayload := r.Context().Value(ValidatedPayloadKey)
		assert.NotNil(t, validatedPayload)
		w.WriteHeader(http.StatusOK)
	})

	// Apply the middleware
	middleware := ValidateAndSanitizePayload(handler, Payload{})
	middleware.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestValidateAndSanitizePayloadMiddleware_InvalidPayload(t *testing.T) {
	type Payload struct {
		Name  string `json:"name" validate:"required,min=3,max=50"`
		Email string `json:"email" validate:"required,email"`
	}

	// Create a request with an invalid payload
	payload := Payload{Name: "Jo", Email: "invalid-email"}
	payloadJSON, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/", strings.NewReader(string(payloadJSON)))
	if err != nil {
		t.Fatal(err)
	}

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
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
