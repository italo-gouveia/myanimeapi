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

func setupAnimeTestRouter(handler *handlers.AnimeHandler) *mux.Router {
	router := mux.NewRouter()

	// Public routes (no authentication required) - Order matters!
	router.HandleFunc("/animes/search", handler.GetAnimesByTitleHandler).Methods("GET")
	router.HandleFunc("/animes/genre/{genre}", handler.GetAnimesByGenreHandler).Methods("GET")
	router.HandleFunc("/animes", handler.GetAllAnimesHandler).Methods("GET")
	router.HandleFunc("/animes/{id}", handler.GetAnimeHandler).Methods("GET")

	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/animes").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware)

	// Protected routes
	protectedRouter.HandleFunc("", handler.CreateAnimeHandler).Methods("POST")
	protectedRouter.HandleFunc("/{id}", handler.UpdateAnimeHandler).Methods("PUT")
	protectedRouter.HandleFunc("/{id}", handler.DeleteAnimeHandler).Methods("DELETE")

	return router
}

func generateAnimeTestToken() string {
	// Set test secret key
	os.Setenv("JWT_SECRET_KEY", "test-secret-key")

	// Generate token
	token, _ := middleware.GenerateToken("1", true)
	return token
}

func TestGetAnimeByID(t *testing.T) {
	mockService := &MockAnimeService{}
	handler := handlers.NewAnimeHandler(mockService)
	router := setupAnimeTestRouter(handler)

	// Generate test token
	token := generateAnimeTestToken()

	tests := []struct {
		name           string
		animeID        string
		mockAnime      *models.Anime
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
		withAuth       bool
	}{
		{
			name:    "Success",
			animeID: "1",
			mockAnime: &models.Anime{
				ID:          1,
				Title:       "Test Anime",
				Description: "Test Description",
				Rating:      8.5,
				Episodes:    12,
				Status:      "Completed",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				StartDate:   time.Time{},
				EndDate:     time.Time{},
				Genres: []models.Genre{
					{
						ID:   1,
						Name: "Action",
					},
				},
				Tags: []models.Tag{
					{
						ID:   1,
						Name: "Action",
					},
				},
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
				"genres": []interface{}{
					map[string]interface{}{
						"id":         float64(1),
						"name":       "Action",
						"created_at": "0001-01-01T00:00:00Z",
						"updated_at": "0001-01-01T00:00:00Z",
						"deleted_at": "0001-01-01T00:00:00Z",
					},
				},
				"tags": []interface{}{
					map[string]interface{}{
						"id":         float64(1),
						"name":       "Action",
						"created_at": "0001-01-01T00:00:00Z",
						"updated_at": "0001-01-01T00:00:00Z",
						"deleted_at": "0001-01-01T00:00:00Z",
					},
				},
			},
			withAuth: true,
		},
		{
			name:           "Invalid ID Format",
			animeID:        "invalid",
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
			name:           "Not Found",
			animeID:        "999",
			mockError:      errors.NewError(errors.ErrResourceNotFound, "Anime not found", "The requested anime could not be found", http.StatusNotFound, nil, nil),
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
		{
			name:    "Unauthorized",
			animeID: "1",
			mockAnime: &models.Anime{
				ID:          1,
				Title:       "Test Anime",
				Description: "Test Description",
				Rating:      8.5,
				Episodes:    12,
				Status:      "Completed",
				CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				StartDate:   time.Time{},
				EndDate:     time.Time{},
				Genres: []models.Genre{
					{
						ID:   1,
						Name: "Action",
					},
				},
				Tags: []models.Tag{
					{
						ID:   1,
						Name: "Action",
					},
				},
			},
			expectedStatus: http.StatusOK, // Route is now public
			expectedBody: map[string]interface{}{
				"id":          float64(1),
				"title":       "Test Anime",
				"description": "Test Description",
				"rating":      float64(8.5),
				"episodes":    float64(12),
				"status":      "Completed",
				"start_date":  "0001-01-01T00:00:00Z",
				"end_date":    "0001-01-01T00:00:00Z",
				"created_at":  "2023-01-01T00:00:00Z",
				"updated_at":  "2023-01-01T00:00:00Z",
				"genres": []interface{}{
					map[string]interface{}{
						"id":         float64(1),
						"name":       "Action",
						"created_at": "0001-01-01T00:00:00Z",
						"updated_at": "0001-01-01T00:00:00Z",
						"deleted_at": "0001-01-01T00:00:00Z",
					},
				},
				"tags": []interface{}{
					map[string]interface{}{
						"id":         float64(1),
						"name":       "Action",
						"created_at": "0001-01-01T00:00:00Z",
						"updated_at": "0001-01-01T00:00:00Z",
						"deleted_at": "0001-01-01T00:00:00Z",
					},
				},
			},
			withAuth: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock expectations
			if tt.mockAnime != nil || tt.mockError != nil {
				animeID, _ := strconv.ParseUint(tt.animeID, 10, 32)
				mockService.On("GetAnimeByID", mock.Anything, uint(animeID)).Return(tt.mockAnime, tt.mockError)
			}

			// Create request
			req, _ := http.NewRequest("GET", fmt.Sprintf("/animes/%s", tt.animeID), nil)
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

func TestGetAnimesByTitle(t *testing.T) {
	mockService := &MockAnimeService{}
	handler := handlers.NewAnimeHandler(mockService)
	router := setupAnimeTestRouter(handler)

	tests := []struct {
		name           string
		title          string
		page           string
		limit          string
		mockAnimes     []*models.Anime
		mockTotal      int64
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:  "Success",
			title: "Test",
			page:  "1",
			limit: "10",
			mockAnimes: []*models.Anime{
				{
					ID:          1,
					Title:       "Test Anime",
					Description: "Test Description",
					Rating:      8.5,
					Episodes:    12,
					Status:      "Completed",
					CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					StartDate:   time.Time{},
					EndDate:     time.Time{},
					Genres: []models.Genre{
						{
							ID:   1,
							Name: "Action",
						},
					},
					Tags: []models.Tag{
						{
							ID:   1,
							Name: "Action",
						},
					},
				},
			},
			mockTotal:      1,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"id":          float64(1),
						"title":       "Test Anime",
						"description": "Test Description",
						"rating":      float64(8.5),
						"episodes":    float64(12),
						"status":      "Completed",
						"start_date":  "0001-01-01T00:00:00Z",
						"end_date":    "0001-01-01T00:00:00Z",
						"created_at":  "2023-01-01T00:00:00Z",
						"updated_at":  "2023-01-01T00:00:00Z",
						"genres": []interface{}{
							map[string]interface{}{
								"id":         float64(1),
								"name":       "Action",
								"created_at": "0001-01-01T00:00:00Z",
								"updated_at": "0001-01-01T00:00:00Z",
								"deleted_at": "0001-01-01T00:00:00Z",
							},
						},
						"tags": []interface{}{
							map[string]interface{}{
								"id":         float64(1),
								"name":       "Action",
								"created_at": "0001-01-01T00:00:00Z",
								"updated_at": "0001-01-01T00:00:00Z",
								"deleted_at": "0001-01-01T00:00:00Z",
							},
						},
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
			name:           "Invalid Page",
			title:          "Test",
			page:           "0",
			limit:          "10",
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInvalidInput,
					"message": "Invalid pagination parameters",
					"details": "invalid page number. Must be a positive integer",
				},
			},
		},
		{
			name:           "Invalid Limit",
			title:          "Test",
			page:           "1",
			limit:          "101",
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInvalidInput,
					"message": "Invalid pagination parameters",
					"details": "invalid limit number. Must be a positive integer between 1 and 100",
				},
			},
		},
		{
			name:           "Internal Server Error",
			title:          "Test",
			page:           "1",
			limit:          "10",
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
			if tt.mockAnimes != nil || tt.mockError != nil {
				page, _ := strconv.Atoi(tt.page)
				limit, _ := strconv.Atoi(tt.limit)
				mockService.On("GetAnimesByTitle", mock.Anything, tt.title, page, limit).Return(tt.mockAnimes, tt.mockTotal, tt.mockError)
			}

			// Create request
			req, _ := http.NewRequest("GET", fmt.Sprintf("/animes/search?title=%s&page=%s&limit=%s", tt.title, tt.page, tt.limit), nil)

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

func TestGetAnimesByGenre(t *testing.T) {
	mockService := &MockAnimeService{}
	handler := handlers.NewAnimeHandler(mockService)
	router := setupAnimeTestRouter(handler)

	tests := []struct {
		name           string
		genre          string
		page           string
		limit          string
		mockAnimes     []*models.Anime
		mockTotal      int64
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:  "Success",
			genre: "action",
			page:  "1",
			limit: "10",
			mockAnimes: []*models.Anime{
				{
					ID:          1,
					Title:       "Test Anime 1",
					Description: "Test Description 1",
					Rating:      8.5,
					Episodes:    12,
					Status:      "Completed",
					Genres:      []models.Genre{{ID: 1, Name: "Action"}},
					Tags:        []models.Tag{{ID: 1, Name: "Action"}},
					CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:          2,
					Title:       "Test Anime 2",
					Description: "Test Description 2",
					Rating:      9.0,
					Episodes:    24,
					Status:      "Ongoing",
					Genres:      []models.Genre{{ID: 1, Name: "Action"}},
					Tags:        []models.Tag{{ID: 1, Name: "Action"}},
					CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			mockTotal:      2,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"id":          float64(1),
						"title":       "Test Anime 1",
						"description": "Test Description 1",
						"rating":      float64(8.5),
						"episodes":    float64(12),
						"status":      "Completed",
						"start_date":  "2023-01-01T00:00:00Z",
						"end_date":    "0001-01-01T00:00:00Z",
						"created_at":  "2023-01-01T00:00:00Z",
						"updated_at":  "2023-01-01T00:00:00Z",
						"genres": []interface{}{
							map[string]interface{}{
								"id":         float64(1),
								"name":       "Action",
								"created_at": "0001-01-01T00:00:00Z",
								"updated_at": "0001-01-01T00:00:00Z",
								"deleted_at": "0001-01-01T00:00:00Z",
							},
						},
						"tags": []interface{}{
							map[string]interface{}{
								"id":         float64(1),
								"name":       "Action",
								"created_at": "0001-01-01T00:00:00Z",
								"updated_at": "0001-01-01T00:00:00Z",
								"deleted_at": "0001-01-01T00:00:00Z",
							},
						},
					},
					map[string]interface{}{
						"id":          float64(2),
						"title":       "Test Anime 2",
						"description": "Test Description 2",
						"rating":      float64(9.0),
						"episodes":    float64(24),
						"status":      "Ongoing",
						"start_date":  "0001-01-01T00:00:00Z",
						"end_date":    "0001-01-01T00:00:00Z",
						"created_at":  "2023-01-01T00:00:00Z",
						"updated_at":  "2023-01-01T00:00:00Z",
						"genres": []interface{}{
							map[string]interface{}{
								"id":         float64(1),
								"name":       "Action",
								"created_at": "0001-01-01T00:00:00Z",
								"updated_at": "0001-01-01T00:00:00Z",
								"deleted_at": "0001-01-01T00:00:00Z",
							},
						},
						"tags": []interface{}{
							map[string]interface{}{
								"id":         float64(1),
								"name":       "Action",
								"created_at": "0001-01-01T00:00:00Z",
								"updated_at": "0001-01-01T00:00:00Z",
								"deleted_at": "0001-01-01T00:00:00Z",
							},
						},
					},
				},
				"total": float64(2),
				"page":  float64(1),
				"limit": float64(10),
			},
		},
		{
			name:           "No Animes Found",
			genre:          "nonexistent",
			page:           "1",
			limit:          "10",
			mockAnimes:     []*models.Anime{},
			mockTotal:      0,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data":  []interface{}{},
				"total": float64(0),
				"page":  float64(1),
				"limit": float64(10),
			},
		},
		{
			name:           "Invalid Pagination",
			genre:          "action",
			page:           "0",
			limit:          "10",
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInvalidInput,
					"message": "Invalid pagination parameters",
					"details": "invalid page number. Must be a positive integer",
				},
			},
		},
		{
			name:           "Internal Server Error",
			genre:          "action",
			page:           "1",
			limit:          "10",
			mockAnimes:     nil,
			mockTotal:      0,
			mockError:      errors.NewError(errors.ErrInternalServer, "Database error", "Failed to get animes by genre", http.StatusInternalServerError, nil, nil),
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
			// Set up mock expectations
			if tc.mockAnimes != nil || tc.mockError != nil {
				page, _ := strconv.Atoi(tc.page)
				limit, _ := strconv.Atoi(tc.limit)
				mockService.On("GetAnimesByGenre", mock.Anything, tc.genre, page, limit).Return(tc.mockAnimes, tc.mockTotal, tc.mockError)
			}

			// Create request
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/animes/genre/%s?page=%s&limit=%s", tc.genre, tc.page, tc.limit), nil)

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
