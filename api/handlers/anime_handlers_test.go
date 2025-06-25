package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"myanimeapi/api/middleware"
	"myanimeapi/api/mocks"
	"myanimeapi/api/models"
	apperrors "myanimeapi/internal/errors"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// createTestContext creates a context with middleware values for testing
func createTestContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.UserContextKey, uint(1))
	return ctx
}

func TestAnimeHandler_GetAnimeHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnimeService := mocks.NewMockAnimeServiceInterface(ctrl)
	handler := NewAnimeHandler(mockAnimeService)

	tests := []struct {
		name           string
		animeID        string
		setupMock      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:    "Success",
			animeID: "1",
			setupMock: func() {
				anime := &models.Anime{
					ID:          1,
					Title:       "Test Anime",
					Description: "Test Description",
					Rating:      8.5,
					Episodes:    12,
					Status:      "Completed",
					Genres: []models.Genre{
						{
							ID:   1,
							Name: "Action",
						},
					},
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
					StartDate: time.Time{},
					EndDate:   time.Time{},
				}
				mockAnimeService.EXPECT().
					GetAnimeByID(gomock.Any(), uint(1)).
					DoAndReturn(func(ctx context.Context, id uint) (*models.Anime, error) {
						// Verify that the context has the user ID
						userID, ok := ctx.Value(middleware.UserContextKey).(uint)
						assert.True(t, ok)
						assert.Equal(t, uint(1), userID)
						return anime, nil
					})
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":          float64(1),
				"title":       "Test Anime",
				"description": "Test Description",
				"rating":      8.5,
				"episodes":    float64(12),
				"status":      "Completed",
				"genres": []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "Action",
					},
				},
				"created_at": "0001-01-01T00:00:00Z",
				"updated_at": "0001-01-01T00:00:00Z",
				"start_date": "0001-01-01T00:00:00Z",
				"end_date":   "0001-01-01T00:00:00Z",
			},
		},
		{
			name:    "Invalid ID",
			animeID: "invalid",
			setupMock: func() {
				// No mock setup needed for invalid ID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-001",
					"message": "Invalid ID format",
					"details": "The provided ID is not a valid unsigned integer.",
					"context": map[string]interface{}{
						"id": "invalid",
					},
				},
			},
		},
		{
			name:    "Not Found",
			animeID: "999",
			setupMock: func() {
				mockAnimeService.EXPECT().
					GetAnimeByID(gomock.Any(), uint(999)).
					Return(nil, apperrors.NewError(apperrors.ErrResourceNotFound, "Anime not found", "The requested anime could not be found", http.StatusNotFound, nil, nil))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-002",
					"message": "Anime not found",
					"details": "The requested anime could not be found",
				},
			},
		},
		{
			name:    "Internal Server Error",
			animeID: "1",
			setupMock: func() {
				mockAnimeService.EXPECT().
					GetAnimeByID(gomock.Any(), uint(1)).
					Return(nil, apperrors.NewError(apperrors.ErrInternalServer, "Database error", "Failed to retrieve anime from database", http.StatusInternalServerError, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-004",
					"message": "Database error",
					"details": "Failed to retrieve anime from database",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			// Create a new request with the test context
			req := httptest.NewRequest(http.MethodGet, "/animes/"+tt.animeID, nil)
			req = req.WithContext(createTestContext())

			// Set up the URL variables for the request
			vars := map[string]string{
				"id": tt.animeID,
			}
			req = mux.SetURLVars(req, vars)

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handler.GetAnimeHandler(rr, req)

			// Verify the response
			assert.Equal(t, tt.expectedStatus, rr.Code)
			var response map[string]interface{}
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

func TestAnimeHandler_GetAllAnimesHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnimeService := mocks.NewMockAnimeServiceInterface(ctrl)
	handler := NewAnimeHandler(mockAnimeService)

	tests := []struct {
		name           string
		query          string
		setupMock      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:  "Success",
			query: "?page=1&limit=10",
			setupMock: func() {
				animes := []*models.Anime{
					{
						ID:          1,
						Title:       "Anime 1",
						Description: "Description 1",
						Rating:      8.5,
						Episodes:    12,
						Status:      "Completed",
						Genres: []models.Genre{
							{
								ID:   1,
								Name: "Action",
							},
						},
					},
					{
						ID:          2,
						Title:       "Anime 2",
						Description: "Description 2",
						Rating:      9.0,
						Episodes:    24,
						Status:      "Ongoing",
						Genres: []models.Genre{
							{
								ID:   2,
								Name: "Drama",
							},
						},
					},
				}
				mockAnimeService.EXPECT().
					GetAllAnimes(gomock.Any(), 1, 10).
					Return(animes, int64(2), nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"id":          float64(1),
						"title":       "Anime 1",
						"description": "Description 1",
						"rating":      8.5,
						"episodes":    float64(12),
						"status":      "Completed",
						"genres": []interface{}{
							map[string]interface{}{
								"id":   float64(1),
								"name": "Action",
							},
						},
					},
					map[string]interface{}{
						"id":          float64(2),
						"title":       "Anime 2",
						"description": "Description 2",
						"rating":      9.0,
						"episodes":    float64(24),
						"status":      "Ongoing",
						"genres": []interface{}{
							map[string]interface{}{
								"id":   float64(2),
								"name": "Drama",
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
			name:  "Invalid Pagination",
			query: "?page=0&limit=10",
			setupMock: func() {
				// No mock setup needed for invalid pagination
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-001",
					"message": "Invalid pagination parameters",
					"details": "invalid page number. Must be a positive integer",
				},
			},
		},
		{
			name:  "Internal Server Error",
			query: "?page=1&limit=10",
			setupMock: func() {
				mockAnimeService.EXPECT().
					GetAllAnimes(gomock.Any(), 1, 10).
					Return(nil, apperrors.NewError(apperrors.ErrInternalServer, "Database error", "Failed to retrieve animes from database", http.StatusInternalServerError, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-004",
					"message": "Database error",
					"details": "Failed to retrieve animes from database",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			// Create a new request
			req := httptest.NewRequest(http.MethodGet, "/animes"+tt.query, nil)
			req = req.WithContext(createTestContext())

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handler.GetAllAnimesHandler(rr, req)

			// Verify the response
			assert.Equal(t, tt.expectedStatus, rr.Code)
			var response map[string]interface{}
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

func TestAnimeHandler_CreateAnimeHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnimeService := mocks.NewMockAnimeServiceInterface(ctrl)
	handler := NewAnimeHandler(mockAnimeService)

	tests := []struct {
		name           string
		requestBody    string
		setupMock      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:        "Success",
			requestBody: `{"title":"New Anime","description":"New Description","rating":8.5,"episodes":12,"status":"Completed","genre_ids":[1,2],"tag_ids":[1,2]}`,
			setupMock: func() {
				mockAnimeService.EXPECT().
					CreateAnime(gomock.Any(), gomock.Any()).
					Return(nil)
				mockAnimeService.EXPECT().
					GetAnimeByID(gomock.Any(), uint(1)).
					Return(&models.Anime{
						ID:          1,
						Title:       "New Anime",
						Description: "New Description",
						Rating:      8.5,
						Episodes:    12,
						Status:      "Completed",
					}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"id":          float64(1),
				"title":       "New Anime",
				"description": "New Description",
				"rating":      8.5,
				"episodes":    float64(12),
				"status":      "Completed",
			},
		},
		{
			name:        "Invalid Request Body",
			requestBody: `{"invalid":"json"`,
			setupMock: func() {
				// No mock setup needed for invalid request body
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-003",
					"message": "Invalid request body",
					"details": "Invalid JSON format",
				},
			},
		},
		{
			name:        "Validation Error",
			requestBody: `{"title":""}`,
			setupMock: func() {
				// No mock setup needed for validation error
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-003",
					"message": "Validation error",
					"details": "Title is required",
				},
			},
		},
		{
			name:        "Internal Server Error",
			requestBody: `{"title":"New Anime","description":"New Description","rating":8.5,"episodes":12,"status":"Completed"}`,
			setupMock: func() {
				mockAnimeService.EXPECT().
					CreateAnime(gomock.Any(), gomock.Any()).
					Return(apperrors.NewError(apperrors.ErrInternalServer, "Database error", "Failed to create anime in database", http.StatusInternalServerError, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-004",
					"message": "Database error",
					"details": "Failed to create anime in database",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			// Create a new request with the request body
			req := httptest.NewRequest(http.MethodPost, "/animes", bytes.NewBufferString(tt.requestBody))
			req = req.WithContext(createTestContext())

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handler.CreateAnimeHandler(rr, req)

			// Verify the response
			assert.Equal(t, tt.expectedStatus, rr.Code)
			var response map[string]interface{}
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

func TestAnimeHandler_UpdateAnimeHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnimeService := mocks.NewMockAnimeServiceInterface(ctrl)
	handler := NewAnimeHandler(mockAnimeService)

	tests := []struct {
		name           string
		animeID        string
		requestBody    string
		setupMock      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:        "Success",
			animeID:     "1",
			requestBody: `{"title":"Updated Anime","description":"Updated Description","rating":9.0,"episodes":24,"status":"Ongoing"}`,
			setupMock: func() {
				mockAnimeService.EXPECT().
					GetAnimeByID(gomock.Any(), uint(1)).
					Return(&models.Anime{
						ID:          1,
						Title:       "Test Anime",
						Description: "Test Description",
						Rating:      8.5,
						Episodes:    12,
						Status:      "Completed",
					}, nil)
				mockAnimeService.EXPECT().
					UpdateAnime(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":          float64(1),
				"title":       "Updated Anime",
				"description": "Updated Description",
				"rating":      9.0,
				"episodes":    float64(24),
				"status":      "Ongoing",
			},
		},
		{
			name:        "Invalid ID",
			animeID:     "invalid",
			requestBody: `{"title":"Updated Anime"}`,
			setupMock: func() {
				// No mock setup needed for invalid ID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-001",
					"message": "Invalid ID format",
					"details": "The provided ID is not a valid unsigned integer.",
					"context": map[string]interface{}{
						"id": "invalid",
					},
				},
			},
		},
		{
			name:        "Not Found",
			animeID:     "999",
			requestBody: `{"title":"Updated Anime"}`,
			setupMock: func() {
				mockAnimeService.EXPECT().
					GetAnimeByID(gomock.Any(), uint(999)).
					Return(nil, apperrors.NewError(apperrors.ErrResourceNotFound, "Anime not found", "The requested anime could not be found", http.StatusNotFound, nil, nil))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-002",
					"message": "Anime not found",
					"details": "The requested anime could not be found",
				},
			},
		},
		{
			name:        "Internal Server Error",
			animeID:     "1",
			requestBody: `{"title":"Updated Anime"}`,
			setupMock: func() {
				mockAnimeService.EXPECT().
					GetAnimeByID(gomock.Any(), uint(1)).
					Return(&models.Anime{
						ID:          1,
						Title:       "Test Anime",
						Description: "Test Description",
						Rating:      8.5,
						Episodes:    12,
						Status:      "Completed",
					}, nil)
				mockAnimeService.EXPECT().
					UpdateAnime(gomock.Any(), gomock.Any()).
					Return(apperrors.NewError(apperrors.ErrInternalServer, "Database error", "Failed to update anime in database", http.StatusInternalServerError, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-004",
					"message": "Database error",
					"details": "Failed to update anime in database",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodPut, "/animes/"+tt.animeID, bytes.NewBufferString(tt.requestBody))
			req = req.WithContext(createTestContext())
			req = mux.SetURLVars(req, map[string]string{"anime_id": tt.animeID})
			rr := httptest.NewRecorder()

			handler.UpdateAnimeHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			var response map[string]interface{}
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

func TestAnimeHandler_DeleteAnimeHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnimeService := mocks.NewMockAnimeServiceInterface(ctrl)
	handler := NewAnimeHandler(mockAnimeService)

	tests := []struct {
		name           string
		animeID        string
		setupMock      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:    "Success",
			animeID: "1",
			setupMock: func() {
				mockAnimeService.EXPECT().
					DeleteAnime(gomock.Any(), uint(1)).
					Return(nil)
			},
			expectedStatus: http.StatusNoContent,
			expectedBody:   nil,
		},
		{
			name:    "Invalid ID",
			animeID: "invalid",
			setupMock: func() {
				// No mock setup needed for invalid ID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-001",
					"message": "Invalid ID format",
					"details": "The provided ID is not a valid unsigned integer.",
					"context": map[string]interface{}{
						"id": "invalid",
					},
				},
			},
		},
		{
			name:    "Not Found",
			animeID: "999",
			setupMock: func() {
				mockAnimeService.EXPECT().
					DeleteAnime(gomock.Any(), uint(999)).
					Return(apperrors.NewError(apperrors.ErrResourceNotFound, "Anime not found", "The requested anime could not be found", http.StatusNotFound, nil, nil))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-002",
					"message": "Anime not found",
					"details": "The requested anime could not be found",
				},
			},
		},
		{
			name:    "Internal Server Error",
			animeID: "1",
			setupMock: func() {
				mockAnimeService.EXPECT().
					DeleteAnime(gomock.Any(), uint(1)).
					Return(apperrors.NewError(apperrors.ErrInternalServer, "Database error", "Failed to delete anime from database", http.StatusInternalServerError, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-004",
					"message": "Database error",
					"details": "Failed to delete anime from database",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			// Create a new request
			req := httptest.NewRequest(http.MethodDelete, "/animes/"+tt.animeID, nil)
			req = req.WithContext(createTestContext())

			// Set up the URL variables for the request
			vars := map[string]string{
				"id": tt.animeID,
			}
			req = mux.SetURLVars(req, vars)

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handler.DeleteAnimeHandler(rr, req)

			// Verify the response
			assert.Equal(t, tt.expectedStatus, rr.Code)
			var response map[string]interface{}
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

func TestAnimeHandler_GetAnimesByTitleHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnimeService := mocks.NewMockAnimeServiceInterface(ctrl)
	handler := NewAnimeHandler(mockAnimeService)

	tests := []struct {
		name           string
		title          string
		page           int
		limit          int
		setupMock      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:  "Success",
			title: "Naruto",
			page:  1,
			limit: 10,
			setupMock: func() {
				animes := []*models.Anime{
					{
						ID:          1,
						Title:       "Naruto",
						Description: "A story about ninjas",
						Rating:      8.5,
						Episodes:    220,
						Status:      "Completed",
						StartDate:   time.Time{},
						EndDate:     time.Time{},
						CreatedAt:   time.Time{},
						UpdatedAt:   time.Time{},
						Genres: []models.Genre{
							{
								ID:   1,
								Name: "Action",
							},
						},
					},
					{
						ID:          2,
						Title:       "Naruto Shippuden",
						Description: "The continuation of Naruto",
						Rating:      9.0,
						Episodes:    500,
						Status:      "Completed",
						StartDate:   time.Time{},
						EndDate:     time.Time{},
						CreatedAt:   time.Time{},
						UpdatedAt:   time.Time{},
						Genres: []models.Genre{
							{
								ID:   2,
								Name: "Drama",
							},
						},
					},
				}
				mockAnimeService.EXPECT().
					GetAnimesByTitle(gomock.Any(), "Naruto", 1, 10).
					Return(animes, int64(2), nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"id":          float64(1),
						"title":       "Naruto",
						"description": "A story about ninjas",
						"rating":      8.5,
						"episodes":    float64(220),
						"status":      "Completed",
						"start_date":  "0001-01-01T00:00:00Z",
						"end_date":    "0001-01-01T00:00:00Z",
						"created_at":  "0001-01-01T00:00:00Z",
						"updated_at":  "0001-01-01T00:00:00Z",
						"genres": []interface{}{
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
						"title":       "Naruto Shippuden",
						"description": "The continuation of Naruto",
						"rating":      9.0,
						"episodes":    float64(500),
						"status":      "Completed",
						"start_date":  "0001-01-01T00:00:00Z",
						"end_date":    "0001-01-01T00:00:00Z",
						"created_at":  "0001-01-01T00:00:00Z",
						"updated_at":  "0001-01-01T00:00:00Z",
						"genres": []interface{}{
							map[string]interface{}{
								"id":         float64(2),
								"name":       "Drama",
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
			name:  "No Results",
			title: "NonExistentAnime",
			page:  1,
			limit: 10,
			setupMock: func() {
				mockAnimeService.EXPECT().
					GetAnimesByTitle(gomock.Any(), "NonExistentAnime", 1, 10).
					Return([]*models.Anime{}, int64(0), nil)
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
			title: "Naruto",
			page:  1,
			limit: 10,
			setupMock: func() {
				mockAnimeService.EXPECT().
					GetAnimesByTitle(gomock.Any(), "Naruto", 1, 10).
					Return(nil, int64(0), apperrors.NewError(apperrors.ErrInternalServer, "Database error", "Failed to search animes", http.StatusInternalServerError, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-004",
					"message": "Database error",
					"details": "Failed to search animes",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/animes/search?title=%s&page=%d&limit=%d", tt.title, tt.page, tt.limit), nil)
			req = req.WithContext(createTestContext())
			rr := httptest.NewRecorder()

			handler.GetAnimesByTitleHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			var response map[string]interface{}
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

func TestAnimeHandler_GetAnimesByGenreHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnimeService := mocks.NewMockAnimeServiceInterface(ctrl)
	handler := NewAnimeHandler(mockAnimeService)

	tests := []struct {
		name           string
		genreID        string
		page           int
		limit          int
		setupMock      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:    "Success",
			genreID: "1",
			page:    1,
			limit:   10,
			setupMock: func() {
				animes := []*models.Anime{
					{
						ID:          1,
						Title:       "Naruto",
						Description: "A story about ninjas",
						Rating:      8.5,
						Episodes:    220,
						Status:      "Completed",
						StartDate:   time.Time{},
						EndDate:     time.Time{},
						CreatedAt:   time.Time{},
						UpdatedAt:   time.Time{},
						Genres: []models.Genre{
							{
								ID:   1,
								Name: "Action",
							},
						},
					},
				}
				mockAnimeService.EXPECT().
					GetAnimesByGenre(gomock.Any(), "1", 1, 10).
					Return(animes, int64(1), nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"id":          float64(1),
						"title":       "Naruto",
						"description": "A story about ninjas",
						"rating":      8.5,
						"episodes":    float64(220),
						"status":      "Completed",
						"start_date":  "0001-01-01T00:00:00Z",
						"end_date":    "0001-01-01T00:00:00Z",
						"created_at":  "0001-01-01T00:00:00Z",
						"updated_at":  "0001-01-01T00:00:00Z",
						"genres": []interface{}{
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
			name:    "Invalid Genre ID",
			genreID: "",
			page:    1,
			limit:   10,
			setupMock: func() {
				// No mock setup needed for invalid genre ID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-001",
					"message": "Invalid input",
					"details": "Genre parameter is required",
				},
			},
		},
		{
			name:    "Internal Server Error",
			genreID: "1",
			page:    1,
			limit:   10,
			setupMock: func() {
				mockAnimeService.EXPECT().
					GetAnimesByGenre(gomock.Any(), "1", 1, 10).
					Return(nil, int64(0), apperrors.NewError(apperrors.ErrInternalServer, "Database error", "Failed to search animes by genre", http.StatusInternalServerError, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-004",
					"message": "Database error",
					"details": "Failed to search animes by genre",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/animes/genre/%s?page=%d&limit=%d", tt.genreID, tt.page, tt.limit), nil)
			req = req.WithContext(createTestContext())
			req = mux.SetURLVars(req, map[string]string{"genre": tt.genreID})
			rr := httptest.NewRecorder()

			handler.GetAnimesByGenreHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			var response map[string]interface{}
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

func TestAnimeHandler_AddGenresToAnimeHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnimeService := mocks.NewMockAnimeServiceInterface(ctrl)
	handler := NewAnimeHandler(mockAnimeService)

	tests := []struct {
		name           string
		animeID        string
		requestBody    string
		setupMock      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:    "Success",
			animeID: "1",
			requestBody: `{
				"genre_ids": [1, 2, 3]
			}`,
			setupMock: func() {
				mockAnimeService.EXPECT().
					AddGenresToAnime(gomock.Any(), uint(1), []uint{1, 2, 3}).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Genres added successfully",
			},
		},
		{
			name:    "Invalid Anime ID",
			animeID: "invalid",
			requestBody: `{
				"genre_ids": [1, 2, 3]
			}`,
			setupMock: func() {
				// No mock setup needed for invalid anime ID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-001",
					"message": "Invalid ID format",
					"details": "The provided ID is not a valid unsigned integer.",
					"context": map[string]interface{}{
						"id": "invalid",
					},
				},
			},
		},
		{
			name:    "Invalid Request Body",
			animeID: "1",
			requestBody: `{
				"invalid": "json"
			}`,
			setupMock: func() {
				// No mock setup needed for invalid request body
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-003",
					"message": "Invalid request body",
					"details": "Invalid JSON format",
				},
			},
		},
		{
			name:    "Internal Server Error",
			animeID: "1",
			requestBody: `{
				"genre_ids": [1, 2, 3]
			}`,
			setupMock: func() {
				mockAnimeService.EXPECT().
					AddGenresToAnime(gomock.Any(), uint(1), []uint{1, 2, 3}).
					Return(apperrors.NewError(apperrors.ErrInternalServer, "Database error", "Failed to add genres to anime", http.StatusInternalServerError, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-004",
					"message": "Database error",
					"details": "Failed to add genres to anime",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodPost, "/animes/"+tt.animeID+"/genres", bytes.NewBufferString(tt.requestBody))
			req = req.WithContext(createTestContext())
			req = mux.SetURLVars(req, map[string]string{"id": tt.animeID})
			rr := httptest.NewRecorder()

			handler.AddGenresToAnimeHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			var response map[string]interface{}
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

func TestAnimeHandler_RemoveGenresFromAnimeHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnimeService := mocks.NewMockAnimeServiceInterface(ctrl)
	handler := NewAnimeHandler(mockAnimeService)

	tests := []struct {
		name           string
		animeID        string
		requestBody    string
		setupMock      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:    "Success",
			animeID: "1",
			requestBody: `{
				"genre_ids": [1, 2, 3]
			}`,
			setupMock: func() {
				mockAnimeService.EXPECT().
					RemoveGenresFromAnime(gomock.Any(), uint(1), []uint{1, 2, 3}).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Genres removed successfully",
			},
		},
		{
			name:    "Invalid Anime ID",
			animeID: "invalid",
			requestBody: `{
				"genre_ids": [1, 2, 3]
			}`,
			setupMock: func() {
				// No mock setup needed for invalid anime ID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-001",
					"message": "Invalid ID format",
					"details": "The provided ID is not a valid unsigned integer.",
					"context": map[string]interface{}{
						"id": "invalid",
					},
				},
			},
		},
		{
			name:    "Invalid Request Body",
			animeID: "1",
			requestBody: `{
				"invalid": "json"
			}`,
			setupMock: func() {
				// No mock setup needed for invalid request body
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-003",
					"message": "Invalid request body",
					"details": "Invalid JSON format",
				},
			},
		},
		{
			name:    "Internal Server Error",
			animeID: "1",
			requestBody: `{
				"genre_ids": [1, 2, 3]
			}`,
			setupMock: func() {
				mockAnimeService.EXPECT().
					RemoveGenresFromAnime(gomock.Any(), uint(1), []uint{1, 2, 3}).
					Return(apperrors.NewError(apperrors.ErrInternalServer, "Database error", "Failed to remove genres from anime", http.StatusInternalServerError, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-004",
					"message": "Database error",
					"details": "Failed to remove genres from anime",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodDelete, "/animes/"+tt.animeID+"/genres", bytes.NewBufferString(tt.requestBody))
			req = req.WithContext(createTestContext())
			req = mux.SetURLVars(req, map[string]string{"id": tt.animeID})
			rr := httptest.NewRecorder()

			handler.RemoveGenresFromAnimeHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			var response map[string]interface{}
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

func TestAnimeHandler_AddTagsToAnimeHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnimeService := mocks.NewMockAnimeServiceInterface(ctrl)
	handler := NewAnimeHandler(mockAnimeService)

	tests := []struct {
		name           string
		animeID        string
		requestBody    string
		setupMock      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:    "Success",
			animeID: "1",
			requestBody: `{
				"tag_ids": [1, 2, 3]
			}`,
			setupMock: func() {
				mockAnimeService.EXPECT().
					AddTagsToAnime(gomock.Any(), uint(1), []uint{1, 2, 3}).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Tags added successfully",
			},
		},
		{
			name:    "Invalid Anime ID",
			animeID: "invalid",
			requestBody: `{
				"tag_ids": [1, 2, 3]
			}`,
			setupMock: func() {
				// No mock setup needed for invalid anime ID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-001",
					"message": "Invalid ID format",
					"details": "The provided ID is not a valid unsigned integer.",
					"context": map[string]interface{}{
						"id": "invalid",
					},
				},
			},
		},
		{
			name:    "Invalid Request Body",
			animeID: "1",
			requestBody: `{
				"invalid": "json"
			}`,
			setupMock: func() {
				// No mock setup needed for invalid request body
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-003",
					"message": "Invalid request body",
					"details": "Invalid JSON format",
				},
			},
		},
		{
			name:    "Internal Server Error",
			animeID: "1",
			requestBody: `{
				"tag_ids": [1, 2, 3]
			}`,
			setupMock: func() {
				mockAnimeService.EXPECT().
					AddTagsToAnime(gomock.Any(), uint(1), []uint{1, 2, 3}).
					Return(apperrors.NewError(apperrors.ErrInternalServer, "Database error", "Failed to add tags to anime", http.StatusInternalServerError, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-004",
					"message": "Database error",
					"details": "Failed to add tags to anime",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodPost, "/animes/"+tt.animeID+"/tags", bytes.NewBufferString(tt.requestBody))
			req = req.WithContext(createTestContext())
			req = mux.SetURLVars(req, map[string]string{"id": tt.animeID})
			rr := httptest.NewRecorder()

			handler.AddTagsToAnimeHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			var response map[string]interface{}
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

func TestAnimeHandler_RemoveTagsFromAnimeHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnimeService := mocks.NewMockAnimeServiceInterface(ctrl)
	handler := NewAnimeHandler(mockAnimeService)

	tests := []struct {
		name           string
		animeID        string
		requestBody    string
		setupMock      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:    "Success",
			animeID: "1",
			requestBody: `{
				"tag_ids": [1, 2, 3]
			}`,
			setupMock: func() {
				mockAnimeService.EXPECT().
					RemoveTagsFromAnime(gomock.Any(), uint(1), []uint{1, 2, 3}).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Tags removed successfully",
			},
		},
		{
			name:    "Invalid Anime ID",
			animeID: "invalid",
			requestBody: `{
				"tag_ids": [1, 2, 3]
			}`,
			setupMock: func() {
				// No mock setup needed for invalid anime ID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-001",
					"message": "Invalid ID format",
					"details": "The provided ID is not a valid unsigned integer.",
					"context": map[string]interface{}{
						"id": "invalid",
					},
				},
			},
		},
		{
			name:    "Invalid Request Body",
			animeID: "1",
			requestBody: `{
				"invalid": "json"
			}`,
			setupMock: func() {
				// No mock setup needed for invalid request body
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-003",
					"message": "Invalid request body",
					"details": "Invalid JSON format",
				},
			},
		},
		{
			name:    "Internal Server Error",
			animeID: "1",
			requestBody: `{
				"tag_ids": [1, 2, 3]
			}`,
			setupMock: func() {
				mockAnimeService.EXPECT().
					RemoveTagsFromAnime(gomock.Any(), uint(1), []uint{1, 2, 3}).
					Return(apperrors.NewError(apperrors.ErrInternalServer, "Database error", "Failed to remove tags from anime", http.StatusInternalServerError, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ERR-004",
					"message": "Database error",
					"details": "Failed to remove tags from anime",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodDelete, "/animes/"+tt.animeID+"/tags", bytes.NewBufferString(tt.requestBody))
			req = req.WithContext(createTestContext())
			req = mux.SetURLVars(req, map[string]string{"id": tt.animeID})
			rr := httptest.NewRecorder()

			handler.RemoveTagsFromAnimeHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			var response map[string]interface{}
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}
