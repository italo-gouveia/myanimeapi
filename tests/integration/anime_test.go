package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	httphandler "myanimeapi/api/adapters/http"
	"myanimeapi/api/middleware"
	"myanimeapi/api/models"
	"myanimeapi/internal/errors"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAnimeService is a mock for AnimeServiceInterface using testify/mock
type MockAnimeService struct {
	mock.Mock
}

func (m *MockAnimeService) GetAnimeByID(ctx context.Context, id uint) (*models.Anime, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Anime), args.Error(1)
}

func (m *MockAnimeService) GetAllAnimes(ctx context.Context, page, limit int, filter models.AnimeFilter) ([]*models.Anime, int64, error) {
	args := m.Called(ctx, page, limit, filter)
	return args.Get(0).([]*models.Anime), args.Get(1).(int64), args.Error(2)
}

func (m *MockAnimeService) CreateAnime(ctx context.Context, anime *models.Anime) error {
	args := m.Called(ctx, anime)
	return args.Error(0)
}

func (m *MockAnimeService) UpdateAnime(ctx context.Context, anime *models.Anime) error {
	args := m.Called(ctx, anime)
	return args.Error(0)
}

func (m *MockAnimeService) DeleteAnime(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAnimeService) GetAnimesByTitle(ctx context.Context, title string, page, limit int) ([]*models.Anime, int64, error) {
	args := m.Called(ctx, title, page, limit)
	return args.Get(0).([]*models.Anime), args.Get(1).(int64), args.Error(2)
}

func (m *MockAnimeService) GetAnimesByGenre(ctx context.Context, genre string, page, limit int) ([]*models.Anime, int64, error) {
	args := m.Called(ctx, genre, page, limit)
	return args.Get(0).([]*models.Anime), args.Get(1).(int64), args.Error(2)
}

func (m *MockAnimeService) GetReviewsForAnime(ctx context.Context, animeID uint, page, limit int) ([]*models.Review, int64, error) {
	args := m.Called(ctx, animeID, page, limit)
	return args.Get(0).([]*models.Review), args.Get(1).(int64), args.Error(2)
}

func (m *MockAnimeService) AddGenresToAnime(ctx context.Context, animeID uint, genreIDs []uint) error {
	args := m.Called(ctx, animeID, genreIDs)
	return args.Error(0)
}

func (m *MockAnimeService) RemoveGenresFromAnime(ctx context.Context, animeID uint, genreIDs []uint) error {
	args := m.Called(ctx, animeID, genreIDs)
	return args.Error(0)
}

func (m *MockAnimeService) AddTagsToAnime(ctx context.Context, animeID uint, tagIDs []uint) error {
	args := m.Called(ctx, animeID, tagIDs)
	return args.Error(0)
}

func (m *MockAnimeService) RemoveTagsFromAnime(ctx context.Context, animeID uint, tagIDs []uint) error {
	args := m.Called(ctx, animeID, tagIDs)
	return args.Error(0)
}

func setupAnimeTestRouter(handler *httphandler.AnimeHandler) *mux.Router {
	router := mux.NewRouter()

	// Public routes (no authentication required) - Order matters!
	router.HandleFunc("/animes/search", handler.GetAnimesByTitleHandler).Methods("GET")
	router.HandleFunc("/animes/genre/{genre}", handler.GetAnimesByGenreHandler).Methods("GET")
	router.HandleFunc("/animes", handler.GetAllAnimesHandler).Methods("GET")
	router.HandleFunc("/animes/{id}", handler.GetAnimeHandler).Methods("GET")

	// Create a subrouter for protected routes (authenticated + admin)
	protectedRouter := router.PathPrefix("/animes").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware)

	// Admin-only routes
	protectedRouter.Handle("", middleware.RequireAdmin(http.HandlerFunc(handler.CreateAnimeHandler))).Methods("POST")
	protectedRouter.Handle("/{id}", middleware.RequireAdmin(http.HandlerFunc(handler.UpdateAnimeHandler))).Methods("PUT")
	protectedRouter.Handle("/{id}", middleware.RequireAdmin(http.HandlerFunc(handler.DeleteAnimeHandler))).Methods("DELETE")

	return router
}

func generateAnimeTestToken() string {
	_ = os.Setenv("JWT_SECRET_KEY", "test-secret-key")
	token, _ := middleware.GenerateToken("1", true, "admin")
	return token
}

func generateNonAdminTestToken() string {
	_ = os.Setenv("JWT_SECRET_KEY", "test-secret-key")
	token, _ := middleware.GenerateToken("2", false, "user") // isAdmin = false
	return token
}

func TestGetAnimeByID(t *testing.T) {
	mockService := &MockAnimeService{}
	handler := httphandler.NewAnimeHandler(mockService)
	router := setupAnimeTestRouter(handler)
	token := generateAnimeTestToken()

	tests := []struct {
		name           string
		animeID        string
		setupMock      func()
		expectedStatus int
		expectedBody   map[string]interface{}
		withAuth       bool
	}{
		{
			name:    "Success",
			animeID: "1",
			setupMock: func() {
				mockService.On("GetAnimeByID", mock.Anything, uint(1)).Return(&models.Anime{
					ID:          1,
					Title:       "Test Anime",
					Description: "Test Description",
					Rating:      8.5,
					Episodes:    12,
					Status:      "Completed",
				}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":          float64(1),
				"title":       "Test Anime",
				"description": "Test Description",
				"rating":      float64(8.5),
				"episodes":    float64(12),
				"status":      "Completed",
				"start_date":  "0001-01-01T00:00:00Z",
				"end_date":    "0001-01-01T00:00:00Z",
			},
			withAuth: true,
		},
		{
			name:           "Invalid ID Format",
			animeID:        "invalid",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInvalidInput,
					"message": "Invalid ID format",
					"details": "The provided ID is not a valid unsigned integer.",
					"context": map[string]interface{}{
						"id": "invalid",
					},
				},
			},
			withAuth: true,
		},
		{
			name:    "Not Found",
			animeID: "999",
			setupMock: func() {
				mockService.On("GetAnimeByID", mock.Anything, uint(999)).Return(nil, errors.NewError(errors.ErrResourceNotFound, "Anime not found", "The requested anime could not be found", http.StatusNotFound, nil, nil)).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrResourceNotFound,
					"message": "Anime not found",
					"details": "The requested anime could not be found",
				},
			},
			withAuth: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.Calls = nil
			tt.setupMock()

			req, _ := http.NewRequest("GET", fmt.Sprintf("/animes/%s", tt.animeID), nil)
			if tt.withAuth {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			var response map[string]interface{}
			jsonErr := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, jsonErr)
			delete(response, "created_at")
			delete(response, "updated_at")

			assert.Equal(t, tt.expectedBody, response)
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetAnimesByTitle(t *testing.T) {
	mockService := &MockAnimeService{}
	handler := httphandler.NewAnimeHandler(mockService)
	router := setupAnimeTestRouter(handler)

	tests := []struct {
		name           string
		title          string
		page           string
		limit          string
		setupMock      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:  "Success",
			title: "Test",
			page:  "1",
			limit: "10",
			setupMock: func() {
				mockService.On("GetAllAnimes", mock.Anything, 1, 10, models.AnimeFilter{Title: "Test"}).Return([]*models.Anime{
					{
						ID:    1,
						Title: "Test Anime",
					},
				}, int64(1), nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"id":          float64(1),
						"title":       "Test Anime",
						"description": "",
						"rating":      float64(0),
						"episodes":    float64(0),
						"status":      "",
						"start_date":  "0001-01-01T00:00:00Z",
						"end_date":    "0001-01-01T00:00:00Z",
					},
				},
				"total": float64(1),
				"page":  float64(1),
				"limit": float64(10),
			},
		},
		{
			name:           "Missing Title",
			title:          "",
			page:           "1",
			limit:          "10",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInvalidInput,
					"message": "Invalid input",
					"details": "Title parameter is required",
				},
			},
		},
		{
			name:  "Internal Server Error",
			title: "Test",
			page:  "1",
			limit: "10",
			setupMock: func() {
				mockService.On("GetAllAnimes", mock.Anything, 1, 10, models.AnimeFilter{Title: "Test"}).Return([]*models.Anime{}, int64(0), errors.NewError(errors.ErrInternalServer, "Database error", "Failed to search animes", http.StatusInternalServerError, nil, nil)).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInternalServer,
					"message": "Database error",
					"details": "Failed to search animes",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.Calls = nil
			tt.setupMock()

			req, _ := http.NewRequest("GET", fmt.Sprintf("/animes/search?title=%s&page=%s&limit=%s", tt.title, tt.page, tt.limit), nil)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			var response map[string]interface{}
			jsonErr := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, jsonErr)
			if data, ok := response["data"].([]interface{}); ok {
				for _, item := range data {
					if anime, ok := item.(map[string]interface{}); ok {
						delete(anime, "created_at")
						delete(anime, "updated_at")
					}
				}
			}

			assert.Equal(t, tt.expectedBody, response)
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetAnimesByGenre(t *testing.T) {
	mockService := &MockAnimeService{}
	handler := httphandler.NewAnimeHandler(mockService)
	router := setupAnimeTestRouter(handler)

	tests := []struct {
		name           string
		genre          string
		page           string
		limit          string
		setupMock      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:  "Success",
			genre: "action",
			page:  "1",
			limit: "10",
			setupMock: func() {
				mockService.On("GetAnimesByGenre", mock.Anything, "action", 1, 10).Return([]*models.Anime{
					{ID: 1, Title: "Test Anime 1"},
					{ID: 2, Title: "Test Anime 2"},
				}, int64(2), nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"id": float64(1), "title": "Test Anime 1", "description": "", "rating": float64(0), "episodes": float64(0), "status": "", "start_date": "0001-01-01T00:00:00Z", "end_date": "0001-01-01T00:00:00Z",
					},
					map[string]interface{}{
						"id": float64(2), "title": "Test Anime 2", "description": "", "rating": float64(0), "episodes": float64(0), "status": "", "start_date": "0001-01-01T00:00:00Z", "end_date": "0001-01-01T00:00:00Z",
					},
				},
				"total": float64(2),
				"page":  float64(1),
				"limit": float64(10),
			},
		},
		{
			name:  "No Animes Found",
			genre: "nonexistent",
			page:  "1",
			limit: "10",
			setupMock: func() {
				mockService.On("GetAnimesByGenre", mock.Anything, "nonexistent", 1, 10).Return([]*models.Anime{}, int64(0), nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data":  []interface{}{},
				"total": float64(0),
				"page":  float64(1),
				"limit": float64(10),
			},
		},
		{
			name:  "Internal Server Error",
			genre: "action",
			page:  "1",
			limit: "10",
			setupMock: func() {
				mockService.On("GetAnimesByGenre", mock.Anything, "action", 1, 10).Return([]*models.Anime{}, int64(0), errors.NewError(errors.ErrInternalServer, "Database error", "Failed to get animes by genre", http.StatusInternalServerError, nil, nil)).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInternalServer,
					"message": "Database error",
					"details": "Failed to get animes by genre",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.Calls = nil
			tc.setupMock()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/animes/genre/%s?page=%s&limit=%s", tc.genre, tc.page, tc.limit), nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			jsonErr := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, jsonErr)
			if data, ok := response["data"].([]interface{}); ok {
				for _, item := range data {
					if anime, ok := item.(map[string]interface{}); ok {
						delete(anime, "created_at")
						delete(anime, "updated_at")
					}
				}
			}

			assert.Equal(t, tc.expectedBody, response)
			mockService.AssertExpectations(t)
		})
	}
}

// TestAdminOnlyRoutes_ForbiddenForNonAdmin verifies that authenticated non-admin
// users receive 403 Forbidden when attempting any admin-only anime mutation.
func TestAdminOnlyRoutes_ForbiddenForNonAdmin(t *testing.T) {
	mockService := &MockAnimeService{}
	handler := httphandler.NewAnimeHandler(mockService)
	router := setupAnimeTestRouter(handler)
	nonAdminToken := generateNonAdminTestToken()

	routes := []struct {
		method string
		path   string
		body   string
	}{
		{http.MethodPost, "/animes", `{"title":"test"}`},
		{http.MethodPut, "/animes/1", `{"title":"test"}`},
		{http.MethodDelete, "/animes/1", ""},
	}

	for _, rt := range routes {
		t.Run(rt.method+" "+rt.path, func(t *testing.T) {
			var body *bytes.Buffer
			if rt.body != "" {
				body = bytes.NewBufferString(rt.body)
			} else {
				body = &bytes.Buffer{}
			}

			req := httptest.NewRequest(rt.method, rt.path, body)
			req.Header.Set("Authorization", "Bearer "+nonAdminToken)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusForbidden, w.Code, "expected 403 for non-admin on %s %s", rt.method, rt.path)

			var resp map[string]interface{}
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			errObj, ok := resp["error"].(map[string]interface{})
			assert.True(t, ok, "response should contain error object")
			assert.Equal(t, errors.ErrForbidden, errObj["code"])
		})
	}

	mockService.AssertExpectations(t) // no service calls should have been made
}
