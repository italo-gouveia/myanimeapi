package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"myanimeapi/api/handlers"
	"myanimeapi/api/middleware"
	"myanimeapi/api/models"
	"myanimeapi/internal/errors"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFavoriteService is a mock for FavoriteServiceInterface using testify/mock
type MockFavoriteService struct {
	mock.Mock
}

func (m *MockFavoriteService) AddFavorite(ctx context.Context, userID uint, animeID uint) (*models.Favorite, error) {
	args := m.Called(ctx, userID, animeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Favorite), args.Error(1)
}

func (m *MockFavoriteService) RemoveFavorite(ctx context.Context, userID uint, animeID uint) error {
	args := m.Called(ctx, userID, animeID)
	return args.Error(0)
}

func (m *MockFavoriteService) GetFavorites(ctx context.Context, userID uint) ([]models.Favorite, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Favorite), args.Error(1)
}

func setupFavoriteTestRouter(handler *handlers.FavoriteHandler) *mux.Router {
	router := mux.NewRouter()

	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/favorites").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware)

	// Register routes
	protectedRouter.HandleFunc("/{anime_id:[0-9]+}", handler.AddFavoriteHandler).Methods("POST")
	protectedRouter.HandleFunc("/{anime_id:[0-9]+}", handler.RemoveFavoriteHandler).Methods("DELETE")
	protectedRouter.HandleFunc("", handler.GetFavoritesHandler).Methods("GET")

	return router
}

func generateFavoriteTestToken() string {
	// Set test secret key
	os.Setenv("JWT_SECRET_KEY", "test-secret-key")

	// Generate token
	token, _ := middleware.GenerateToken("1", true)
	return token
}

func TestAddFavorite(t *testing.T) {
	mockService := &MockFavoriteService{}
	handler := handlers.NewFavoriteHandler(mockService)
	router := setupFavoriteTestRouter(handler)

	// Generate test token
	token := generateFavoriteTestToken()

	tests := []struct {
		name           string
		animeID        string
		mockFavorite   *models.Favorite
		mockError      error
		expectedStatus int
		expectedBody   func() map[string]interface{}
		withAuth       bool
	}{
		{
			name:    "Success",
			animeID: "1",
			mockFavorite: &models.Favorite{
				ID:        1,
				UserID:    1,
				AnimeID:   1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedStatus: http.StatusCreated,
			expectedBody: func() map[string]interface{} {
				return map[string]interface{}{
					"id":       float64(1),
					"user_id":  float64(1),
					"anime_id": float64(1),
				}
			},
			withAuth: true,
		},
		{
			name:           "Invalid Anime ID",
			animeID:        "invalid",
			expectedStatus: http.StatusNotFound, // Route will not match
			expectedBody: func() map[string]interface{} {
				return map[string]interface{}{"error": "404 not found"}
			},
			withAuth: true,
		},
		{
			name:           "Unauthorized",
			animeID:        "1",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: func() map[string]interface{} {
				return map[string]interface{}{
					"error": map[string]interface{}{
						"code":    errors.ErrUnauthorized,
						"message": "Authentication required",
						"details": "Missing Authorization header",
						"context": map[string]interface{}{"error": "No authorization token provided"},
					},
				}
			},
			withAuth: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockService.ExpectedCalls = nil
			mockService.Calls = nil

			// Set up mock expectations only for valid ID and withAuth
			if tt.mockFavorite != nil && tt.withAuth && tt.animeID != "invalid" {
				animeID, _ := strconv.ParseUint(tt.animeID, 10, 32)
				mockService.On("AddFavorite", mock.Anything, uint(1), uint(animeID)).Return(tt.mockFavorite, tt.mockError).Once()
			}

			// Create request
			req, _ := http.NewRequest("POST", fmt.Sprintf("/favorites/%s", tt.animeID), nil)
			if tt.withAuth {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve request
			router.ServeHTTP(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Parse response body
			var response map[string]interface{}
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			if err != nil && rr.Body.String() != "" {
				// Handle plain text 404
				assert.Contains(t, rr.Body.String(), "not found")
				return
			}

			// Compare response with expected
			if tt.expectedBody != nil {
				expected := tt.expectedBody()
				// Remove dynamic fields for comparison
				delete(response, "created_at")
				delete(response, "updated_at")
				delete(response, "anime")
				delete(response, "user")
				assert.Equal(t, expected, response)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestRemoveFavorite(t *testing.T) {
	mockService := &MockFavoriteService{}
	handler := handlers.NewFavoriteHandler(mockService)
	router := setupFavoriteTestRouter(handler)

	// Generate test token
	token := generateFavoriteTestToken()

	tests := []struct {
		name           string
		animeID        string
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
		withAuth       bool
	}{
		{
			name:           "Success",
			animeID:        "1",
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
			withAuth:       true,
		},
		{
			name:           "Invalid Anime ID",
			animeID:        "invalid",
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]interface{}{"error": "404 not found"},
			withAuth:       true,
		},
		{
			name:           "Unauthorized",
			animeID:        "1",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrUnauthorized,
					"message": "Authentication required",
					"details": "Missing Authorization header",
					"context": map[string]interface{}{"error": "No authorization token provided"},
				},
			},
			withAuth: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockService.ExpectedCalls = nil
			mockService.Calls = nil

			// Set up mock expectations
			if tt.name == "Success" {
				animeID, _ := strconv.ParseUint(tt.animeID, 10, 32)
				mockService.On("RemoveFavorite", mock.Anything, uint(1), uint(animeID)).Return(tt.mockError).Once()
			}

			// Create request
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/favorites/%s", tt.animeID), nil)
			if tt.withAuth {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve request
			router.ServeHTTP(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			if rr.Body.Len() > 0 {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					assert.Contains(t, rr.Body.String(), "not found")
					return
				}
				assert.Equal(t, tt.expectedBody, response)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestGetFavorites(t *testing.T) {
	mockService := &MockFavoriteService{}
	handler := handlers.NewFavoriteHandler(mockService)
	router := setupFavoriteTestRouter(handler)
	token := generateFavoriteTestToken()

	tests := []struct {
		name           string
		setupMock      func()
		expectedStatus int
		expectedBody   interface{}
		withAuth       bool
	}{
		{
			name: "Success",
			setupMock: func() {
				mockService.On("GetFavorites", mock.Anything, uint(1)).Return([]models.Favorite{
					{
						ID:      1,
						UserID:  1,
						AnimeID: 1,
						Anime:   models.Anime{ID: 1, Title: "Test Anime"},
					},
				}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: []interface{}{
				map[string]interface{}{
					"id":       float64(1),
					"user_id":  float64(1),
					"anime_id": float64(1),
					"anime": map[string]interface{}{
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
			},
			withAuth: true,
		},
		{
			name: "Empty Favorites",
			setupMock: func() {
				mockService.On("GetFavorites", mock.Anything, uint(1)).Return([]models.Favorite{}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   []interface{}{},
			withAuth:       true,
		},
		{
			name: "Internal Server Error",
			setupMock: func() {
				mockService.On("GetFavorites", mock.Anything, uint(1)).Return([]models.Favorite{}, errors.NewError(errors.ErrInternalServer, "DB error", "...", http.StatusInternalServerError, nil, nil)).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInternalServer,
					"message": "DB error",
					"details": "...",
				},
			},
			withAuth: true,
		},
		{
			name:           "Unauthorized",
			setupMock:      func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrUnauthorized,
					"message": "Authentication required",
					"details": "Missing Authorization header",
					"context": map[string]interface{}{"error": "No authorization token provided"},
				},
			},
			withAuth: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockService.ExpectedCalls = nil
			mockService.Calls = nil

			tt.setupMock()

			req, _ := http.NewRequest("GET", "/favorites", nil)
			if tt.withAuth {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			var response interface{}
			json.Unmarshal(rr.Body.Bytes(), &response)

			// Remove dynamic fields for comparison
			if responseArray, ok := response.([]interface{}); ok {
				for _, item := range responseArray {
					if anime, ok := item.(map[string]interface{}); ok {
						delete(anime, "created_at")
						delete(anime, "updated_at")
						delete(anime, "user")
						if animeData, ok := anime["anime"].(map[string]interface{}); ok {
							delete(animeData, "created_at")
							delete(animeData, "updated_at")
						}
					}
				}
			}

			assert.Equal(t, tt.expectedBody, response)

			mockService.AssertExpectations(t)
		})
	}
}
