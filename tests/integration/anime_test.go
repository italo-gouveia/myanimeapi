package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"myanimeapi/api/handlers"
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

func (m *MockAnimeService) GetAllAnimes(ctx context.Context, page, limit int) ([]*models.Anime, int64, error) {
	args := m.Called(ctx, page, limit)
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

func TestGetAnimeByID(t *testing.T) {
	// Create a new mock service
	mockService := new(MockAnimeService)

	// Create a new handler with the mock service
	handler := handlers.NewAnimeHandler(mockService)

	// Create a new router
	router := mux.NewRouter()
	router.HandleFunc("/anime/{id}", handler.GetAnimeHandler).Methods("GET")

	tests := []struct {
		name           string
		id             string
		mockAnime      *models.Anime
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Success",
			id:   "1",
			mockAnime: &models.Anime{
				ID:    1,
				Title: "Test Anime",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":    float64(1),
				"title": "Test Anime",
			},
		},
		{
			name:           "Invalid ID",
			id:             "invalid",
			mockAnime:      nil,
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
			id:             "999",
			mockAnime:      nil,
			mockError:      errors.NewError(errors.ErrResourceNotFound, "Anime not found", "The requested anime could not be found", http.StatusNotFound, nil, nil),
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"code":    errors.ErrResourceNotFound,
				"message": "Anime not found",
				"details": "The requested anime could not be found",
			},
		},
		{
			name:           "Internal Server Error",
			id:             "1",
			mockAnime:      nil,
			mockError:      fmt.Errorf("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"code":    errors.ErrInternalServer,
				"message": "Database error",
				"details": "Failed to get anime",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock expectations
			if tt.mockError == nil && tt.mockAnime != nil {
				mockService.On("GetAnimeByID", mock.Anything, tt.mockAnime.ID).Return(tt.mockAnime, nil)
			} else if tt.mockError != nil {
				mockService.On("GetAnimeByID", mock.Anything, mock.Anything).Return(nil, tt.mockError)
			}

			// Create a new request
			req := httptest.NewRequest("GET", fmt.Sprintf("/anime/%s", tt.id), nil)
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

func TestGetAnimesByTitle(t *testing.T) {
	// Create a new mock service
	mockService := new(MockAnimeService)

	// Create a new handler with the mock service
	handler := handlers.NewAnimeHandler(mockService)

	// Create a new router
	router := mux.NewRouter()
	router.HandleFunc("/animes/search", handler.GetAnimesByTitleHandler).Methods("GET")

	tests := []struct {
		name           string
		query          string
		mockAnimes     []*models.Anime
		mockTotal      int64
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:  "Success",
			query: "Naruto",
			mockAnimes: []*models.Anime{
				{ID: 1, Title: "Naruto"},
				{ID: 2, Title: "Naruto Shippuden"},
			},
			mockTotal:      2,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{"id": float64(1), "title": "Naruto"},
					map[string]interface{}{"id": float64(2), "title": "Naruto Shippuden"},
				},
				"total": float64(2),
				"page":  float64(1),
				"limit": float64(10),
			},
		},
		{
			name:           "Empty Query",
			query:          "",
			mockAnimes:     nil,
			mockTotal:      0,
			mockError:      nil,
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
			name:           "Internal Server Error",
			query:          "Naruto",
			mockAnimes:     nil,
			mockTotal:      0,
			mockError:      errors.NewError(errors.ErrInternalServer, "Database error", "Failed to search animes", http.StatusInternalServerError, nil, nil),
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
			// Set up mock expectations
			if tt.mockError == nil {
				mockService.On("GetAnimesByTitle", mock.Anything, tt.query, 1, 10).Return(tt.mockAnimes, tt.mockTotal, nil)
			} else {
				mockService.On("GetAnimesByTitle", mock.Anything, tt.query, 1, 10).Return(nil, int64(0), tt.mockError)
			}

			// Create a new request
			req := httptest.NewRequest("GET", fmt.Sprintf("/animes/search?title=%s", tt.query), nil)
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse the response body
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Compare the response body
			assert.Equal(t, tt.expectedBody, response)

			// Verify that all expectations were met
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetAnimesByGenre(t *testing.T) {
	// Create a new mock service
	mockService := new(MockAnimeService)

	// Create a new handler with the mock service
	handler := handlers.NewAnimeHandler(mockService)

	// Create a new router
	router := mux.NewRouter()
	router.HandleFunc("/animes/genre/{genre}", handler.GetAnimesByGenreHandler).Methods("GET")

	tests := []struct {
		name           string
		query          string
		mockAnimes     []*models.Anime
		mockTotal      int64
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:  "Success",
			query: "1",
			mockAnimes: []*models.Anime{
				{ID: 1, Title: "Action Anime 1"},
				{ID: 2, Title: "Action Anime 2"},
			},
			mockTotal:      2,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{"id": float64(1), "title": "Action Anime 1"},
					map[string]interface{}{"id": float64(2), "title": "Action Anime 2"},
				},
				"total": float64(2),
				"page":  float64(1),
				"limit": float64(10),
			},
		},
		{
			name:           "Empty Query",
			query:          "",
			mockAnimes:     nil,
			mockTotal:      0,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInvalidInput,
					"message": "Invalid input",
					"details": "Genre parameter is required",
				},
			},
		},
		{
			name:           "Invalid Genre ID",
			query:          "invalid",
			mockAnimes:     nil,
			mockTotal:      0,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInvalidInput,
					"message": "Invalid genre ID format",
					"details": "The provided genre ID is not a valid unsigned integer.",
					"context": map[string]interface{}{
						"genre_id": "invalid",
					},
				},
			},
		},
		{
			name:           "Internal Server Error",
			query:          "1",
			mockAnimes:     nil,
			mockTotal:      0,
			mockError:      errors.NewError(errors.ErrInternalServer, "Database error", "Failed to search animes by genre", http.StatusInternalServerError, nil, nil),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInternalServer,
					"message": "Database error",
					"details": "Failed to search animes by genre",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock expectations
			if tt.mockError == nil && tt.query != "" && tt.query != "invalid" {
				mockService.On("GetAnimesByGenre", mock.Anything, tt.query, 1, 10).Return(tt.mockAnimes, tt.mockTotal, nil)
			} else if tt.mockError != nil {
				mockService.On("GetAnimesByGenre", mock.Anything, tt.query, 1, 10).Return(nil, int64(0), tt.mockError)
			}

			// Create a new request
			req := httptest.NewRequest("GET", fmt.Sprintf("/animes/genre/%s", tt.query), nil)
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse the response body
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Compare the response body
			assert.Equal(t, tt.expectedBody, response)

			// Verify that all expectations were met
			mockService.AssertExpectations(t)
		})
	}
}
