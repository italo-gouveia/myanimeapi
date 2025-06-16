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

	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/animes").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware)

	// Register routes
	protectedRouter.HandleFunc("/{id}", handler.GetAnimeHandler).Methods("GET")
	protectedRouter.HandleFunc("", handler.GetAnimesByTitleHandler).Methods("GET")
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
				"data": map[string]interface{}{
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

			// For success case, we need to handle dynamic timestamps
			if tt.name == "Success" {
				// Verify that created_at and updated_at are present and valid timestamps
				assert.Contains(t, response, "created_at")
				assert.Contains(t, response, "updated_at")
				createdAt, err := time.Parse(time.RFC3339, response["created_at"].(string))
				assert.NoError(t, err)
				assert.True(t, createdAt.After(time.Time{}))
				updatedAt, err := time.Parse(time.RFC3339, response["updated_at"].(string))
				assert.NoError(t, err)
				assert.True(t, updatedAt.After(time.Time{}))

				// Remove dynamic fields for comparison
				delete(response, "created_at")
				delete(response, "updated_at")
			}

			// Compare response with expected
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

func TestGetAnimesByTitle(t *testing.T) {
	mockService := &MockAnimeService{}
	handler := handlers.NewAnimeHandler(mockService)
	router := setupAnimeTestRouter(handler)

	// Generate test token
	token := generateAnimeTestToken()

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
		withAuth       bool
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
			},
			mockTotal:      1,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": map[string]interface{}{
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
				"total": float64(1),
				"page":  float64(1),
				"limit": float64(10),
			},
			withAuth: true,
		},
		{
			name:           "Invalid Page",
			title:          "Test",
			page:           "invalid",
			limit:          "10",
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInvalidInput,
					"message": "Invalid input",
					"details": "Invalid page number",
					"context": map[string]interface{}{
						"error": "Invalid page number",
					},
				},
			},
			withAuth: true,
		},
		{
			name:           "Invalid Limit",
			title:          "Test",
			page:           "1",
			limit:          "invalid",
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errors.ErrInvalidInput,
					"message": "Invalid input",
					"details": "Invalid limit number",
					"context": map[string]interface{}{
						"error": "Invalid limit number",
					},
				},
			},
			withAuth: true,
		},
		{
			name:           "Unauthorized",
			title:          "Test",
			page:           "1",
			limit:          "10",
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
			if tt.mockAnimes != nil {
				page, _ := strconv.Atoi(tt.page)
				limit, _ := strconv.Atoi(tt.limit)
				mockService.On("GetAnimesByTitle", mock.Anything, tt.title, page, limit).Return(tt.mockAnimes, tt.mockTotal, tt.mockError)
			}

			// Create request
			req, _ := http.NewRequest("GET", fmt.Sprintf("/animes?title=%s&page=%s&limit=%s", tt.title, tt.page, tt.limit), nil)
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

func TestGetAnimesByGenre(t *testing.T) {
	mockService := &MockAnimeService{}
	handler := handlers.NewAnimeHandler(mockService)
	router := setupAnimeTestRouter(handler)

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
					"context": map[string]interface{}{
						"error": "Failed to get animes by genre",
					},
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
					"message": "Authentication required",
					"details": "Missing Authorization header",
					"context": map[string]interface{}{
						"error": "No authorization token provided",
					},
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
