package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"myanimeapi/api/handlers"
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

func TestAddFavorite(t *testing.T) {
	// Create a new mock service
	mockService := new(MockFavoriteService)

	// Create a new handler with the mock service
	handler := handlers.NewFavoriteHandler(mockService)

	// Create a new router
	router := mux.NewRouter()
	router.HandleFunc("/favorites", handler.AddFavoriteHandler).Methods("POST")

	tests := []struct {
		name           string
		userID         int
		animeID        int
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "Success",
			userID:         1,
			animeID:        1,
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"message": "Anime added to favorites",
			},
		},
		{
			name:           "Anime Already in Favorites",
			userID:         1,
			animeID:        1,
			mockError:      errors.NewError(errors.ErrConflict, "Anime already in favorites", "The anime is already in the user's favorites", http.StatusConflict, nil, nil),
			expectedStatus: http.StatusConflict,
			expectedBody: map[string]interface{}{
				"code":    errors.ErrConflict,
				"message": "Anime already in favorites",
				"details": "The anime is already in the user's favorites",
			},
		},
		{
			name:           "Internal Server Error",
			userID:         1,
			animeID:        1,
			mockError:      fmt.Errorf("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"code":    errors.ErrInternalServer,
				"message": "Database error",
				"details": "Failed to add favorite to database",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request body
			reqBody := map[string]interface{}{
				"user_id":  tt.userID,
				"anime_id": tt.animeID,
			}
			jsonBody, _ := json.Marshal(reqBody)

			// Set up mock expectations
			if tt.mockError == nil {
				mockService.On("AddFavorite", mock.Anything, tt.userID, tt.animeID).Return(nil, nil)
			} else {
				mockService.On("AddFavorite", mock.Anything, tt.userID, tt.animeID).Return(nil, tt.mockError)
			}

			// Create a new request
			req := httptest.NewRequest("POST", "/favorites", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			// Check the response
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse the response body
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)

			// Compare the response body
			assert.Equal(t, tt.expectedBody, response)

			// Verify that all expectations were met
			mockService.AssertExpectations(t)
		})
	}
}

func TestRemoveFavorite(t *testing.T) {
	// Create a new mock service
	mockService := new(MockFavoriteService)

	// Create a new handler with the mock service
	handler := handlers.NewFavoriteHandler(mockService)

	// Create a new router
	router := mux.NewRouter()
	router.HandleFunc("/favorites/{user_id}/{anime_id}", handler.RemoveFavoriteHandler).Methods("DELETE")

	tests := []struct {
		name           string
		userID         string
		animeID        string
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "Success",
			userID:         "1",
			animeID:        "1",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Anime removed from favorites",
			},
		},
		{
			name:           "Invalid User ID",
			userID:         "invalid",
			animeID:        "1",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code":    errors.ErrInvalidInput,
				"message": "Invalid input",
				"details": "Invalid user ID format",
			},
		},
		{
			name:           "Invalid Anime ID",
			userID:         "1",
			animeID:        "invalid",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code":    errors.ErrInvalidInput,
				"message": "Invalid input",
				"details": "Invalid anime ID format",
			},
		},
		{
			name:           "Not Found",
			userID:         "1",
			animeID:        "999",
			mockError:      errors.NewError(errors.ErrResourceNotFound, "Favorite not found", "The requested favorite could not be found", http.StatusNotFound, nil, nil),
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"code":    errors.ErrResourceNotFound,
				"message": "Favorite not found",
				"details": "The requested favorite could not be found",
			},
		},
		{
			name:           "Internal Server Error",
			userID:         "1",
			animeID:        "1",
			mockError:      fmt.Errorf("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"code":    errors.ErrInternalServer,
				"message": "Database error",
				"details": "Failed to remove favorite from database",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock expectations
			if tt.mockError != nil {
				userID, _ := strconv.Atoi(tt.userID)
				animeID, _ := strconv.Atoi(tt.animeID)
				mockService.On("RemoveFavorite", mock.Anything, userID, animeID).Return(tt.mockError)
			}

			// Create a new request
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/favorites/%s/%s", tt.userID, tt.animeID), nil)
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			// Check the response
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse the response body
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)

			// Compare the response body
			assert.Equal(t, tt.expectedBody, response)

			// Verify that all expectations were met
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetFavorites(t *testing.T) {
	// Create a new mock service
	mockService := new(MockFavoriteService)

	// Create a new handler with the mock service
	handler := handlers.NewFavoriteHandler(mockService)

	// Create a new router
	router := mux.NewRouter()
	router.HandleFunc("/favorites/{user_id}", handler.GetFavoritesHandler).Methods("GET")

	tests := []struct {
		name           string
		userID         string
		mockFavorites  []models.Favorite
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:   "Success",
			userID: "1",
			mockFavorites: []models.Favorite{
				{ID: 1, UserID: 1, AnimeID: 1},
				{ID: 2, UserID: 1, AnimeID: 2},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: []interface{}{
				map[string]interface{}{"id": float64(1), "user_id": float64(1), "anime_id": float64(1)},
				map[string]interface{}{"id": float64(2), "user_id": float64(1), "anime_id": float64(2)},
			},
		},
		{
			name:           "Invalid User ID",
			userID:         "invalid",
			mockFavorites:  nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code":    errors.ErrInvalidInput,
				"message": "Invalid input",
				"details": "Invalid user ID format",
			},
		},
		{
			name:           "Internal Server Error",
			userID:         "1",
			mockFavorites:  nil,
			mockError:      fmt.Errorf("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"code":    errors.ErrInternalServer,
				"message": "Database error",
				"details": "Failed to retrieve favorites from database",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock expectations
			if tt.mockError == nil && tt.userID != "invalid" {
				userID, _ := strconv.Atoi(tt.userID)
				mockService.On("GetFavorites", mock.Anything, userID).Return(tt.mockFavorites, nil)
			} else if tt.mockError != nil {
				userID, _ := strconv.Atoi(tt.userID)
				mockService.On("GetFavorites", mock.Anything, userID).Return(nil, tt.mockError)
			}

			// Create a new request
			req := httptest.NewRequest("GET", fmt.Sprintf("/favorites/%s", tt.userID), nil)
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			// Check the response
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse the response body
			var response interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)

			// Compare the response body
			assert.Equal(t, tt.expectedBody, response)

			// Verify that all expectations were met
			mockService.AssertExpectations(t)
		})
	}
}
