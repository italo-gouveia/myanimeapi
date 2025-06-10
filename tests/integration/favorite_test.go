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
	protectedRouter.HandleFunc("", handler.AddFavoriteHandler).Methods("POST")
	protectedRouter.HandleFunc("/{animeID}", handler.RemoveFavoriteHandler).Methods("DELETE")
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
		expectedBody   map[string]interface{}
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
			expectedBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id":         float64(1),
					"user_id":    float64(1),
					"anime_id":   float64(1),
					"created_at": "2024-01-01T00:00:00Z",
					"updated_at": "2024-01-01T00:00:00Z",
				},
			},
			withAuth: true,
		},
		{
			name:           "Invalid Anime ID",
			animeID:        "invalid",
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInvalidInput,
					"message": "Invalid input",
					"details": "Invalid anime ID format",
				},
			},
			withAuth: true,
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
				},
			},
			withAuth: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock expectations
			if tt.mockFavorite != nil {
				animeID, _ := strconv.ParseUint(tt.animeID, 10, 32)
				mockService.On("AddFavorite", mock.Anything, uint(1), uint(animeID)).Return(tt.mockFavorite, tt.mockError)
			}

			// Create request
			req, _ := http.NewRequest("POST", fmt.Sprintf("/favorites?anime_id=%s", tt.animeID), nil)
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
			assert.NoError(t, err)

			// Compare response with expected
			assert.Equal(t, tt.expectedBody, response)
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
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": map[string]interface{}{
					"message": "Favorite removed successfully",
				},
			},
			withAuth: true,
		},
		{
			name:           "Invalid Anime ID",
			animeID:        "invalid",
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInvalidInput,
					"message": "Invalid input",
					"details": "Invalid anime ID format",
				},
			},
			withAuth: true,
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
				},
			},
			withAuth: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock expectations
			if tt.mockError == nil {
				animeID, _ := strconv.ParseUint(tt.animeID, 10, 32)
				mockService.On("RemoveFavorite", mock.Anything, uint(1), uint(animeID)).Return(tt.mockError)
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

			// Parse response body
			var response map[string]interface{}
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Compare response with expected
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

func TestGetFavorites(t *testing.T) {
	mockService := new(MockFavoriteService)
	handler := handlers.NewFavoriteHandler(mockService)
	router := setupFavoriteTestRouter(handler)

	tests := []struct {
		name           string
		userID         uint
		mockFavorites  []models.Favorite
		mockError      error
		expectedStatus int
		expectedBody   interface{}
		withAuth       bool
	}{
		{
			name:   "Success",
			userID: 1,
			mockFavorites: []models.Favorite{
				{
					ID:      1,
					UserID:  1,
					AnimeID: 1,
					Anime: models.Anime{
						ID:          1,
						Title:       "Test Anime",
						Description: "Test Description",
						Episodes:    12,
						Status:      "ongoing",
						Rating:      4.5,
						StartDate:   time.Now(),
						EndDate:     time.Now().AddDate(0, 1, 0),
					},
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: []interface{}{
				map[string]interface{}{
					"id":       float64(1),
					"user_id":  float64(1),
					"anime_id": float64(1),
					"anime": map[string]interface{}{
						"id":          float64(1),
						"title":       "Test Anime",
						"description": "Test Description",
						"episodes":    float64(12),
						"status":      "ongoing",
						"rating":      float64(4.5),
						"start_date":  time.Now().Format(time.RFC3339),
						"end_date":    time.Now().AddDate(0, 1, 0).Format(time.RFC3339),
					},
				},
			},
			withAuth: true,
		},
		{
			name:           "Empty Favorites",
			userID:         1,
			mockFavorites:  []models.Favorite{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   []interface{}{},
			withAuth:       true,
		},
		{
			name:           "Internal Server Error",
			userID:         1,
			mockFavorites:  nil,
			mockError:      errors.NewError(errors.ErrInternalServer, "Database error", "Failed to retrieve favorites", http.StatusInternalServerError, nil, nil),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInternalServer,
					"message": "Database error",
					"details": "Failed to retrieve favorites",
				},
			},
			withAuth: true,
		},
		{
			name:           "Unauthorized",
			userID:         0,
			mockFavorites:  nil,
			mockError:      nil,
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrUnauthorized,
					"message": "User not authenticated",
					"details": "Authentication required to access this resource",
				},
			},
			withAuth: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.withAuth {
				mockService.On("GetFavorites", mock.Anything, tt.userID).Return(tt.mockFavorites, tt.mockError)
			}

			req := httptest.NewRequest("GET", "/favorites", nil)

			if tt.withAuth {
				token := generateFavoriteTestToken()
				req.Header.Set("Authorization", "Bearer "+token)
			}

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, tt.expectedStatus, recorder.Code)

			var response interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}
