package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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

func setupTestRouter(handler *handlers.AnimeHandler) *mux.Router {
	router := mux.NewRouter()

	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/animes").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware)

	// Register routes
	protectedRouter.HandleFunc("/{id}", handler.GetAnimeHandler).Methods("GET")
	protectedRouter.HandleFunc("/search", handler.GetAnimesByTitleHandler).Methods("GET")
	protectedRouter.HandleFunc("/genre/{genre}", handler.GetAnimesByGenreHandler).Methods("GET")

	return router
}

func generateTestToken() string {
	// Set test secret key
	os.Setenv("JWT_SECRET_KEY", "test-secret-key")

	// Generate token
	token, _ := middleware.GenerateToken("1", true)
	return token
}

func TestGetAnimeByID(t *testing.T) {
	mockService := &MockAnimeService{}
	handler := handlers.NewAnimeHandler(mockService)
	router := setupTestRouter(handler)

	tests := []struct {
		name           string
		id             string
		mockAnime      *models.Anime
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
		withAuth       bool
	}{
		{
			name: "Success",
			id:   "1",
			mockAnime: &models.Anime{
				ID:          1,
				Title:       "Test Anime",
				Description: "Test Description",
				Genres:      []models.Genre{{ID: 1, Name: "Action"}},
				Tags:        []models.Tag{{ID: 1, Name: "Action"}},
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":          float64(1),
				"title":       "Test Anime",
				"description": "Test Description",
				"genres": []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "Action",
					},
				},
				"tags": []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "Action",
					},
				},
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-01-01T00:00:00Z",
			},
			withAuth: true,
		},
		{
			name:           "Invalid ID Format",
			id:             "invalid",
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
			name:           "Not Found",
			id:             "999",
			mockError:      errors.NewError(errors.ErrResourceNotFound, "Not found", "Anime not found", http.StatusNotFound, nil, nil),
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrResourceNotFound,
					"message": "Not found",
					"details": "Anime not found",
				},
			},
			withAuth: true,
		},
		{
			name:           "Internal Server Error",
			id:             "1",
			mockError:      errors.NewError(errors.ErrInternalServer, "Database error", "Failed to retrieve anime", http.StatusInternalServerError, nil, nil),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInternalServer,
					"message": "Database error",
					"details": "Failed to retrieve anime",
				},
			},
			withAuth: true,
		},
		{
			name:           "Unauthorized",
			id:             "1",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrUnauthorized,
					"message": "Unauthorized",
					"details": "Missing or invalid authorization header",
				},
			},
			withAuth: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService.On("GetAnimeByID", mock.Anything, tc.id).Return(tc.mockAnime, tc.mockError)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/animes/%s", tc.id), nil)
			if tc.withAuth {
				req.Header.Set("Authorization", "Bearer test-token")
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, response)
		})
	}
}

func TestGetAnimesByTitle(t *testing.T) {
	mockService := &MockAnimeService{}
	handler := handlers.NewAnimeHandler(mockService)
	router := setupTestRouter(handler)

	tests := []struct {
		name           string
		title          string
		mockAnimes     []models.Anime
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
		withAuth       bool
	}{
		{
			name:  "Success",
			title: "test",
			mockAnimes: []models.Anime{
				{
					ID:          1,
					Title:       "Test Anime 1",
					Description: "Test Description 1",
					Genres:      []models.Genre{{ID: 1, Name: "Action"}},
					Tags:        []models.Tag{{ID: 1, Name: "Action"}},
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				{
					ID:          2,
					Title:       "Test Anime 2",
					Description: "Test Description 2",
					Genres:      []models.Genre{{ID: 1, Name: "Action"}},
					Tags:        []models.Tag{{ID: 1, Name: "Action"}},
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"id":          float64(1),
						"title":       "Test Anime 1",
						"description": "Test Description 1",
						"genres": []interface{}{
							map[string]interface{}{
								"id":   float64(1),
								"name": "Action",
							},
						},
						"tags": []interface{}{
							map[string]interface{}{
								"id":   float64(1),
								"name": "Action",
							},
						},
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z",
					},
					map[string]interface{}{
						"id":          float64(2),
						"title":       "Test Anime 2",
						"description": "Test Description 2",
						"genres": []interface{}{
							map[string]interface{}{
								"id":   float64(1),
								"name": "Action",
							},
						},
						"tags": []interface{}{
							map[string]interface{}{
								"id":   float64(1),
								"name": "Action",
							},
						},
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z",
					},
				},
				"total": float64(2),
				"page":  float64(1),
				"limit": float64(10),
			},
			withAuth: true,
		},
		{
			name:           "No Animes Found",
			title:          "nonexistent",
			mockAnimes:     []models.Anime{},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data":  []interface{}{},
				"total": float64(0),
				"page":  float64(1),
				"limit": float64(10),
			},
			withAuth: true,
		},
		{
			name:           "Internal Server Error",
			title:          "test",
			mockError:      errors.NewError(errors.ErrInternalServer, "Database error", "Failed to search animes", http.StatusInternalServerError, nil, nil),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInternalServer,
					"message": "Database error",
					"details": "Failed to search animes",
				},
			},
			withAuth: true,
		},
		{
			name:           "Unauthorized",
			title:          "test",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrUnauthorized,
					"message": "Unauthorized",
					"details": "Missing or invalid authorization header",
				},
			},
			withAuth: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService.On("GetAnimesByTitle", mock.Anything, tc.title).Return(tc.mockAnimes, tc.mockError)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/animes/search/%s", tc.title), nil)
			if tc.withAuth {
				req.Header.Set("Authorization", "Bearer test-token")
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, response)
		})
	}
}

func TestGetAnimesByGenre(t *testing.T) {
	mockService := &MockAnimeService{}
	handler := handlers.NewAnimeHandler(mockService)
	router := setupTestRouter(handler)

	tests := []struct {
		name           string
		genre          string
		mockAnimes     []models.Anime
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
		withAuth       bool
	}{
		{
			name:  "Success",
			genre: "action",
			mockAnimes: []models.Anime{
				{
					ID:          1,
					Title:       "Test Anime 1",
					Description: "Test Description 1",
					Genres:      []models.Genre{{ID: 1, Name: "Action"}},
					Tags:        []models.Tag{{ID: 1, Name: "Action"}},
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				{
					ID:          2,
					Title:       "Test Anime 2",
					Description: "Test Description 2",
					Genres:      []models.Genre{{ID: 1, Name: "Action"}},
					Tags:        []models.Tag{{ID: 1, Name: "Action"}},
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"id":          float64(1),
						"title":       "Test Anime 1",
						"description": "Test Description 1",
						"genres": []interface{}{
							map[string]interface{}{
								"id":   float64(1),
								"name": "Action",
							},
						},
						"tags": []interface{}{
							map[string]interface{}{
								"id":   float64(1),
								"name": "Action",
							},
						},
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z",
					},
					map[string]interface{}{
						"id":          float64(2),
						"title":       "Test Anime 2",
						"description": "Test Description 2",
						"genres": []interface{}{
							map[string]interface{}{
								"id":   float64(1),
								"name": "Action",
							},
						},
						"tags": []interface{}{
							map[string]interface{}{
								"id":   float64(1),
								"name": "Action",
							},
						},
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z",
					},
				},
				"total": float64(2),
				"page":  float64(1),
				"limit": float64(10),
			},
			withAuth: true,
		},
		{
			name:           "No Animes Found",
			genre:          "nonexistent",
			mockAnimes:     []models.Anime{},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data":  []interface{}{},
				"total": float64(0),
				"page":  float64(1),
				"limit": float64(10),
			},
			withAuth: true,
		},
		{
			name:           "Internal Server Error",
			genre:          "action",
			mockError:      errors.NewError(errors.ErrInternalServer, "Database error", "Failed to get animes by genre", http.StatusInternalServerError, nil, nil),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInternalServer,
					"message": "Database error",
					"details": "Failed to get animes by genre",
				},
			},
			withAuth: true,
		},
		{
			name:           "Unauthorized",
			genre:          "action",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrUnauthorized,
					"message": "Unauthorized",
					"details": "Missing or invalid authorization header",
				},
			},
			withAuth: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService.On("GetAnimesByGenre", mock.Anything, tc.genre).Return(tc.mockAnimes, tc.mockError)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/animes/genre/%s", tc.genre), nil)
			if tc.withAuth {
				req.Header.Set("Authorization", "Bearer test-token")
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, response)
		})
	}
}
