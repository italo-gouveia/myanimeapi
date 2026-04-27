// Package auth contains integration tests for the authentication endpoints.
// These tests wire the real AuthHandler against a mock AuthService and make
// actual HTTP round-trips so routing, JSON decoding and response formatting
// are all verified together — no database required.
package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"myanimeapi/api/handlers"
	"myanimeapi/api/mocks"
	"myanimeapi/api/models"
	apperrors "myanimeapi/internal/errors"
	"myanimeapi/internal/logger"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newAuthRouter(svc *mocks.MockAuthServiceInterface) *mux.Router {
	log := logger.New()
	router := mux.NewRouter()
	handlers.NewAuthHandler(svc, log).RegisterAuthRoutes(router)
	return router
}

func postJSON(t *testing.T, router http.Handler, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	b, err := json.Marshal(body)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func decodeMap(t *testing.T, rr *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&m))
	return m
}

// ---------------------------------------------------------------------------
// Register
// ---------------------------------------------------------------------------

func TestRegister_HappyPath(t *testing.T) {
	svc := new(mocks.MockAuthServiceInterface)
	svc.EXPECT().
		RegisterUser(mock.Anything, mock.MatchedBy(func(u *models.User) bool {
			return u.Username == "alice" && u.Email == "alice@example.com"
		})).
		Return(nil).Once()

	rr := postJSON(t, newAuthRouter(svc), "/auth/register", map[string]string{
		"username": "alice",
		"email":    "alice@example.com",
		"password": "Secret123!",
	})

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")

	body := decodeMap(t, rr)
	assert.Equal(t, "success", body["status"])
	assert.Equal(t, "User registered successfully", body["message"])
	svc.AssertExpectations(t)
}

func TestRegister_InvalidJSON(t *testing.T) {
	svc := new(mocks.MockAuthServiceInterface)
	router := newAuthRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	svc.AssertExpectations(t) // no service calls expected
}

func TestRegister_ServiceError(t *testing.T) {
	svc := new(mocks.MockAuthServiceInterface)
	svc.EXPECT().
		RegisterUser(mock.Anything, mock.Anything).
		Return(apperrors.NewError(
			apperrors.ErrConflict,
			"Username already exists",
			"A user with that username already exists",
			http.StatusConflict,
			nil, nil,
		)).Once()

	rr := postJSON(t, newAuthRouter(svc), "/auth/register", map[string]string{
		"username": "alice",
		"email":    "alice@example.com",
		"password": "Secret123!",
	})

	assert.Equal(t, http.StatusConflict, rr.Code)
	svc.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Authenticate
// ---------------------------------------------------------------------------

func TestAuthenticate_HappyPath(t *testing.T) {
	svc := new(mocks.MockAuthServiceInterface)
	svc.EXPECT().
		AuthenticateUser(mock.Anything, mock.MatchedBy(func(c *models.UserCredentials) bool {
			return c.Username == "alice" && c.Password == "Secret123!"
		})).
		Return("jwt-token-abc", nil).Once()

	rr := postJSON(t, newAuthRouter(svc), "/auth/authenticate", map[string]string{
		"username": "alice",
		"password": "Secret123!",
	})

	assert.Equal(t, http.StatusOK, rr.Code)

	body := decodeMap(t, rr)
	assert.Equal(t, "success", body["status"])
	data, ok := body["data"].(map[string]interface{})
	require.True(t, ok, "data field should be an object")
	assert.Equal(t, "jwt-token-abc", data["token"])
	svc.AssertExpectations(t)
}

func TestAuthenticate_InvalidJSON(t *testing.T) {
	svc := new(mocks.MockAuthServiceInterface)
	router := newAuthRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/auth/authenticate", bytes.NewReader([]byte("{bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	svc.AssertExpectations(t)
}

func TestAuthenticate_InvalidCredentials(t *testing.T) {
	svc := new(mocks.MockAuthServiceInterface)
	svc.EXPECT().
		AuthenticateUser(mock.Anything, mock.Anything).
		Return("", apperrors.NewError(
			apperrors.ErrUnauthorized,
			"Invalid credentials",
			"Username or password is incorrect",
			http.StatusUnauthorized,
			nil, nil,
		)).Once()

	rr := postJSON(t, newAuthRouter(svc), "/auth/authenticate", map[string]string{
		"username": "alice",
		"password": "wrongpass",
	})

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	svc.AssertExpectations(t)
}

func TestAuthenticate_InternalError(t *testing.T) {
	svc := new(mocks.MockAuthServiceInterface)
	svc.EXPECT().
		AuthenticateUser(mock.Anything, mock.Anything).
		Return("", apperrors.NewError(
			apperrors.ErrInternalServer,
			"Internal server error",
			"An unexpected error occurred",
			http.StatusInternalServerError,
			nil, nil,
		)).Once()

	rr := postJSON(t, newAuthRouter(svc), "/auth/authenticate", map[string]string{
		"username": "alice",
		"password": "Secret123!",
	})

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	svc.AssertExpectations(t)
}
