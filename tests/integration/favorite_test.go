package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func (m *MockFavoriteService) AddFavorite(ctx context.Context, userID, animeID uint) (*models.Favorite, error) {
	args := m.Called(ctx, userID, animeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Favorite), args.Error(1)
}

func (m *MockFavoriteService) RemoveFavorite(ctx context.Context, userID, animeID uint) error {
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
	protectedRouter.HandleFunc("", handler.GetFavoritesHandler).Methods("GET")
	protectedRouter.HandleFunc("/{anime_id}", handler.AddFavoriteHandler).Methods("POST")
	protectedRouter.HandleFunc("/{anime_id}", handler.RemoveFavoriteHandler).Methods("DELETE")

	return router
}

func TestAddFavorite(t *testing.T) {
	mockService := new(MockFavoriteService)
	handler := handlers.NewFavoriteHandler(mockService)
	router := setupFavoriteTestRouter(handler)

	tests := []struct {
		name           string
		userID         uint
		animeID        string
		mockFavorite   *models.Favorite
		mockError      error
		expectedStatus int
		expectedBody   interface{}
		withAuth       bool
	}{
		{
			name:    "Success",
			userID:  1,
			animeID: "1",
			mockFavorite: &models.Favorite{
				ID:        1,
				UserID:    1,
				AnimeID:   1,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				Anime: models.Anime{
					ID:          1,
					Title:       "Test Anime",
					Description: "Test Description",
					Rating:      8.5,
					Episodes:    12,
					Status:      "Completed",
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
					StartDate:   time.Time{},
					EndDate:     time.Time{},
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"id":         float64(1),
				"user_id":    float64(1),
				"anime_id":   float64(1),
				"created_at": "0001-01-01T00:00:00Z",
				"updated_at": "0001-01-01T00:00:00Z",
				"anime": map[string]interface{}{
					"id":          float64(1),
					"title":       "Test Anime",
					"description": "Test Description",
					"rating":      float64(8.5),
					"episodes":    float64(12),
					"status":      "Completed",
					"created_at":  "0001-01-01T00:00:00Z",
					"updated_at":  "0001-01-01T00:00:00Z",
					"start_date":  "0001-01-01T00:00:00Z",
					"end_date":    "0001-01-01T00:00:00Z",
				},
			},
			withAuth: true,
		},
		{
			name:           "Invalid Anime ID",
			userID:         1,
			animeID:        "invalid",
			mockFavorite:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInvalidInput,
					"message": "Invalid anime ID format",
					"details": "The provided anime ID must be a valid unsigned integer",
				},
			},
			withAuth: true,
		},
		{
			name:           "Anime Already in Favorites",
			userID:         1,
			animeID:        "1",
			mockFavorite:   nil,
			mockError:      errors.NewError(errors.ErrConflict, "Anime already in favorites", "User 1 already has anime 1 in favorites", http.StatusConflict, nil, nil),
			expectedStatus: http.StatusConflict,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrConflict,
					"message": "Anime already in favorites",
					"details": "User 1 already has anime 1 in favorites",
				},
			},
			withAuth: true,
		},
		{
			name:           "Internal Server Error",
			userID:         1,
			animeID:        "1",
			mockFavorite:   nil,
			mockError:      errors.NewError(errors.ErrInternalServer, "Database error", "Failed to add favorite", http.StatusInternalServerError, nil, nil),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInternalServer,
					"message": "Database error",
					"details": "Failed to add favorite",
				},
			},
			withAuth: true,
		},
		{
			name:           "Unauthorized",
			userID:         0,
			animeID:        "1",
			mockFavorite:   nil,
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
				mockService.On("AddFavorite", mock.Anything, tt.userID, uint(1)).Return(tt.mockFavorite, tt.mockError)
			}

			req := httptest.NewRequest("POST", fmt.Sprintf("/favorites/%s", tt.animeID), nil)

			if tt.withAuth {
				token := generateTestToken()
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

func TestRemoveFavorite(t *testing.T) {
	mockService := new(MockFavoriteService)
	handler := handlers.NewFavoriteHandler(mockService)
	router := setupFavoriteTestRouter(handler)

	tests := []struct {
		name           string
		userID         uint
		animeID        string
		mockError      error
		expectedStatus int
		expectedBody   interface{}
		withAuth       bool
	}{
		{
			name:           "Success",
			userID:         1,
			animeID:        "1",
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
			expectedBody:   nil,
			withAuth:       true,
		},
		{
			name:           "Invalid Anime ID",
			userID:         1,
			animeID:        "invalid",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInvalidInput,
					"message": "Invalid anime ID format",
					"details": "The provided anime ID must be a valid unsigned integer",
				},
			},
			withAuth: true,
		},
		{
			name:           "Favorite Not Found",
			userID:         1,
			animeID:        "999",
			mockError:      errors.NewError(errors.ErrResourceNotFound, "Favorite not found", "The requested favorite could not be found", http.StatusNotFound, nil, nil),
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrResourceNotFound,
					"message": "Favorite not found",
					"details": "The requested favorite could not be found",
				},
			},
			withAuth: true,
		},
		{
			name:           "Internal Server Error",
			userID:         1,
			animeID:        "1",
			mockError:      errors.NewError(errors.ErrInternalServer, "Database error", "Failed to remove favorite", http.StatusInternalServerError, nil, nil),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInternalServer,
					"message": "Database error",
					"details": "Failed to remove favorite",
				},
			},
			withAuth: true,
		},
		{
			name:           "Unauthorized",
			userID:         0,
			animeID:        "1",
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
				mockService.On("RemoveFavorite", mock.Anything, tt.userID, uint(1)).Return(tt.mockError)
			}

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/favorites/%s", tt.animeID), nil)

			if tt.withAuth {
				token := generateTestToken()
				req.Header.Set("Authorization", "Bearer "+token)
			}

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, tt.expectedStatus, recorder.Code)

			if tt.expectedBody != nil {
				var response interface{}
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			} else {
				assert.Empty(t, recorder.Body.String())
			}
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
				token := generateTestToken()
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
