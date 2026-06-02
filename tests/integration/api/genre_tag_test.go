// Package api contains integration tests for Genre and Tag endpoints.
// Tests wire the real handlers against mock services and make actual HTTP
// round-trips, verifying routing, JSON encoding and error responses together.
// No database is required.
package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	httphandler "myanimeapi/api/adapters/http"
	"myanimeapi/api/mocks"
	"myanimeapi/api/models"
	"myanimeapi/api/middleware"
	apperrors "myanimeapi/internal/errors"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Router helpers
// ---------------------------------------------------------------------------

func newGenreRouter(svc *mocks.MockGenreServiceInterface) *mux.Router {
	router := mux.NewRouter()
	httphandler.NewGenreHandler(svc).RegisterGenreRoutes(router)
	return router
}

func newTagRouter(svc *mocks.MockTagServiceInterface) *mux.Router {
	router := mux.NewRouter()
	httphandler.NewTagHandler(svc).RegisterTagRoutes(router)
	return router
}

func generateAdminToken(t *testing.T) string {
	t.Helper()
	_ = os.Setenv("JWT_SECRET_KEY", "test-secret-key")
	token, err := middleware.GenerateToken("1", true, "admin")
	require.NoError(t, err)
	return token
}

func doRequest(t *testing.T, router http.Handler, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	t.Helper()
	var reqBody *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewReader(b)
	} else {
		reqBody = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

// ---------------------------------------------------------------------------
// Genre — GetAll
// ---------------------------------------------------------------------------

func TestGenre_GetAll_HappyPath(t *testing.T) {
	svc := new(mocks.MockGenreServiceInterface)
	svc.EXPECT().
		GetAllGenres(mock.Anything, 1, 10).
		Return([]models.Genre{
			{ID: 1, Name: "Action"},
			{ID: 2, Name: "Adventure"},
		}, int64(2), nil).Once()

	rr := doRequest(t, newGenreRouter(svc), http.MethodGet, "/genres", nil, "")

	assert.Equal(t, http.StatusOK, rr.Code)
	var list []interface{}
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&list))
	assert.Len(t, list, 2)
	svc.AssertExpectations(t)
}

func TestGenre_GetAll_Empty(t *testing.T) {
	svc := new(mocks.MockGenreServiceInterface)
	svc.EXPECT().
		GetAllGenres(mock.Anything, 1, 10).
		Return([]models.Genre{}, int64(0), nil).Once()

	rr := doRequest(t, newGenreRouter(svc), http.MethodGet, "/genres", nil, "")

	assert.Equal(t, http.StatusOK, rr.Code)
	var list []interface{}
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&list))
	assert.Empty(t, list)
	svc.AssertExpectations(t)
}

func TestGenre_GetAll_ServiceError(t *testing.T) {
	svc := new(mocks.MockGenreServiceInterface)
	svc.EXPECT().
		GetAllGenres(mock.Anything, 1, 10).
		Return(nil, int64(0), apperrors.NewError(
			apperrors.ErrInternalServer,
			"DB error",
			"database unavailable",
			http.StatusInternalServerError,
			nil, nil,
		)).Once()

	rr := doRequest(t, newGenreRouter(svc), http.MethodGet, "/genres", nil, "")

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	svc.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Genre — GetByID
// ---------------------------------------------------------------------------

func TestGenre_GetByID_HappyPath(t *testing.T) {
	svc := new(mocks.MockGenreServiceInterface)
	svc.EXPECT().
		GetGenreByID(mock.Anything, uint(1)).
		Return(&models.Genre{ID: 1, Name: "Action"}, nil).Once()

	rr := doRequest(t, newGenreRouter(svc), http.MethodGet, "/genres/1", nil, "")

	assert.Equal(t, http.StatusOK, rr.Code)
	var body map[string]interface{}
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&body))
	assert.Equal(t, float64(1), body["id"])
	assert.Equal(t, "Action", body["name"])
	svc.AssertExpectations(t)
}

func TestGenre_GetByID_NotFound(t *testing.T) {
	svc := new(mocks.MockGenreServiceInterface)
	svc.EXPECT().
		GetGenreByID(mock.Anything, uint(999)).
		Return(nil, apperrors.NewError(
			apperrors.ErrResourceNotFound,
			"Genre not found",
			"No genre with that ID",
			http.StatusNotFound,
			nil, nil,
		)).Once()

	rr := doRequest(t, newGenreRouter(svc), http.MethodGet, "/genres/999", nil, "")

	assert.Equal(t, http.StatusNotFound, rr.Code)
	svc.AssertExpectations(t)
}

func TestGenre_GetByID_InvalidID(t *testing.T) {
	svc := new(mocks.MockGenreServiceInterface)
	rr := doRequest(t, newGenreRouter(svc), http.MethodGet, "/genres/abc", nil, "")
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	svc.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Genre — Create (protected)
// ---------------------------------------------------------------------------

func TestGenre_Create_HappyPath(t *testing.T) {
	svc := new(mocks.MockGenreServiceInterface)
	svc.EXPECT().
		CreateGenre(mock.Anything, mock.MatchedBy(func(g *models.Genre) bool {
			return g.Name == "Isekai"
		})).
		Return(nil).Once()

	token := generateAdminToken(t)
	rr := doRequest(t, newGenreRouter(svc), http.MethodPost, "/genres",
		map[string]string{"name": "Isekai"}, token)

	assert.Equal(t, http.StatusCreated, rr.Code)
	svc.AssertExpectations(t)
}

func TestGenre_Create_Unauthorized(t *testing.T) {
	svc := new(mocks.MockGenreServiceInterface)
	rr := doRequest(t, newGenreRouter(svc), http.MethodPost, "/genres",
		map[string]string{"name": "Isekai"}, "")
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	svc.AssertExpectations(t)
}

func TestGenre_Create_InvalidJSON(t *testing.T) {
	svc := new(mocks.MockGenreServiceInterface)
	token := generateAdminToken(t)

	req := httptest.NewRequest(http.MethodPost, "/genres", bytes.NewReader([]byte("{bad")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	newGenreRouter(svc).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	svc.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Genre — Delete (protected)
// ---------------------------------------------------------------------------

func TestGenre_Delete_HappyPath(t *testing.T) {
	svc := new(mocks.MockGenreServiceInterface)
	svc.EXPECT().
		DeleteGenre(mock.Anything, uint(1)).
		Return(nil).Once()

	token := generateAdminToken(t)
	rr := doRequest(t, newGenreRouter(svc), http.MethodDelete, "/genres/1", nil, token)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	svc.AssertExpectations(t)
}

func TestGenre_Delete_NotFound(t *testing.T) {
	svc := new(mocks.MockGenreServiceInterface)
	svc.EXPECT().
		DeleteGenre(mock.Anything, uint(99)).
		Return(apperrors.NewError(
			apperrors.ErrResourceNotFound,
			"Genre not found",
			"No genre with that ID",
			http.StatusNotFound,
			nil, nil,
		)).Once()

	token := generateAdminToken(t)
	rr := doRequest(t, newGenreRouter(svc), http.MethodDelete, "/genres/99", nil, token)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	svc.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Tag — GetAll
// ---------------------------------------------------------------------------

func TestTag_GetAll_HappyPath(t *testing.T) {
	svc := new(mocks.MockTagServiceInterface)
	svc.EXPECT().
		GetAllTags(mock.Anything, 1, 100).
		Return([]models.Tag{
			{ID: 1, Name: "Shounen"},
			{ID: 2, Name: "Isekai"},
		}, int64(2), nil).Once()

	rr := doRequest(t, newTagRouter(svc), http.MethodGet, "/tags", nil, "")

	assert.Equal(t, http.StatusOK, rr.Code)
	var body map[string]interface{}
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&body))
	data, ok := body["data"].([]interface{})
	require.True(t, ok)
	assert.Len(t, data, 2)
	assert.Equal(t, float64(2), body["total"])
	svc.AssertExpectations(t)
}

func TestTag_GetAll_WithPagination(t *testing.T) {
	svc := new(mocks.MockTagServiceInterface)
	svc.EXPECT().
		GetAllTags(mock.Anything, 2, 5).
		Return([]models.Tag{{ID: 6, Name: "Mecha"}}, int64(10), nil).Once()

	rr := doRequest(t, newTagRouter(svc), http.MethodGet, "/tags?page=2&limit=5", nil, "")

	assert.Equal(t, http.StatusOK, rr.Code)
	var body map[string]interface{}
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&body))
	assert.Equal(t, float64(10), body["total"])
	assert.Equal(t, float64(2), body["page"])
	assert.Equal(t, float64(5), body["limit"])
	svc.AssertExpectations(t)
}

func TestTag_GetAll_InvalidPagination(t *testing.T) {
	svc := new(mocks.MockTagServiceInterface)
	rr := doRequest(t, newTagRouter(svc), http.MethodGet, "/tags?page=abc", nil, "")
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	svc.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Tag — GetByID
// ---------------------------------------------------------------------------

func TestTag_GetByID_HappyPath(t *testing.T) {
	svc := new(mocks.MockTagServiceInterface)
	svc.EXPECT().
		GetTagByID(mock.Anything, uint(3)).
		Return(&models.Tag{ID: 3, Name: "Mecha"}, nil).Once()

	rr := doRequest(t, newTagRouter(svc), http.MethodGet, "/tags/3", nil, "")

	assert.Equal(t, http.StatusOK, rr.Code)
	var body map[string]interface{}
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&body))
	assert.Equal(t, float64(3), body["id"])
	assert.Equal(t, "Mecha", body["name"])
	svc.AssertExpectations(t)
}

func TestTag_GetByID_NotFound(t *testing.T) {
	svc := new(mocks.MockTagServiceInterface)
	svc.EXPECT().
		GetTagByID(mock.Anything, uint(999)).
		Return(nil, apperrors.NewError(
			apperrors.ErrResourceNotFound,
			"Tag not found",
			"No tag with that ID",
			http.StatusNotFound,
			nil, nil,
		)).Once()

	rr := doRequest(t, newTagRouter(svc), http.MethodGet, "/tags/999", nil, "")

	assert.Equal(t, http.StatusNotFound, rr.Code)
	svc.AssertExpectations(t)
}

func TestTag_GetByID_InvalidID(t *testing.T) {
	svc := new(mocks.MockTagServiceInterface)
	rr := doRequest(t, newTagRouter(svc), http.MethodGet, "/tags/notanumber", nil, "")
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	svc.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Tag — Create (protected)
// ---------------------------------------------------------------------------

func TestTag_Create_HappyPath(t *testing.T) {
	svc := new(mocks.MockTagServiceInterface)
	svc.EXPECT().
		CreateTag(mock.Anything, mock.MatchedBy(func(tg *models.Tag) bool {
			return tg.Name == "Mahou Shoujo"
		})).
		Return(nil).Once()

	token := generateAdminToken(t)
	rr := doRequest(t, newTagRouter(svc), http.MethodPost, "/tags",
		map[string]string{"name": "Mahou Shoujo"}, token)

	assert.Equal(t, http.StatusCreated, rr.Code)
	svc.AssertExpectations(t)
}

func TestTag_Create_Unauthorized(t *testing.T) {
	svc := new(mocks.MockTagServiceInterface)
	rr := doRequest(t, newTagRouter(svc), http.MethodPost, "/tags",
		map[string]string{"name": "Mahou Shoujo"}, "")
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	svc.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Tag — Delete (protected)
// ---------------------------------------------------------------------------

func TestTag_Delete_HappyPath(t *testing.T) {
	svc := new(mocks.MockTagServiceInterface)
	svc.EXPECT().
		DeleteTag(mock.Anything, uint(1)).
		Return(nil).Once()

	token := generateAdminToken(t)
	rr := doRequest(t, newTagRouter(svc), http.MethodDelete, fmt.Sprintf("/tags/%d", 1), nil, token)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	svc.AssertExpectations(t)
}

func TestTag_Delete_Unauthorized(t *testing.T) {
	svc := new(mocks.MockTagServiceInterface)
	rr := doRequest(t, newTagRouter(svc), http.MethodDelete, "/tags/1", nil, "")
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	svc.AssertExpectations(t)
}
