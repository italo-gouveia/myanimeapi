package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"myanimeapi/api/mocks"
	"myanimeapi/api/models"
	apperrors "myanimeapi/internal/errors"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

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
		expectedBody   string
	}{
		{
			name:    "Success",
			animeID: "1",
			setupMock: func() {
				ctx := context.Background()
				anime := &models.Anime{
					ID:          1,
					Title:       "Test Anime",
					Description: "Test Description",
					Rating:      8.5,
					Episodes:    12,
					Status:      "Completed",
					StartDate:   time.Now(),
					EndDate:     time.Now().AddDate(0, 3, 0),
					Genres: []models.Genre{
						{ID: 1, Name: "Action"},
					},
					Tags: []models.Tag{
						{ID: 1, Name: "Fantasy"},
					},
				}
				mockAnimeService.EXPECT().
					GetAnimeByID(ctx, uint(1)).
					Return(anime, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"title":"Test Anime","description":"Test Description","rating":8.5,"episodes":12,"status":"Completed","genres":[{"id":1,"name":"Action"}],"tags":[{"id":1,"name":"Fantasy"}]}`,
		},
		{
			name:    "Invalid ID",
			animeID: "invalid",
			setupMock: func() {
				// No mock setup needed for invalid ID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":"ERR-003","message":"Invalid ID format","details":"The provided ID is not a valid unsigned integer."}}`,
		},
		{
			name:    "Not Found",
			animeID: "999",
			setupMock: func() {
				ctx := context.Background()
				mockAnimeService.EXPECT().
					GetAnimeByID(ctx, uint(999)).
					Return(nil, apperrors.NewError(apperrors.ErrResourceNotFound, "Anime not found", "The requested anime could not be found", http.StatusNotFound, nil, nil))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":{"code":"ERR-002","message":"Anime not found","details":"The requested anime could not be found"}}`,
		},
		{
			name:    "Internal Server Error",
			animeID: "1",
			setupMock: func() {
				ctx := context.Background()
				mockAnimeService.EXPECT().
					GetAnimeByID(ctx, uint(1)).
					Return(nil, apperrors.NewError(apperrors.ErrInternalServer, "Database error", "Failed to retrieve anime from database", http.StatusInternalServerError, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":{"code":"ERR-004","message":"Database error","details":"Failed to retrieve anime from database"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodGet, "/anime/"+tt.animeID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.animeID})
			w := httptest.NewRecorder()

			handler.GetAnimeHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
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
		expectedBody   string
	}{
		{
			name:  "Success",
			query: "?page=1&limit=10",
			setupMock: func() {
				ctx := context.Background()
				animes := []models.Anime{
					{
						ID:          1,
						Title:       "Anime 1",
						Description: "Description 1",
						Rating:      8.5,
						Episodes:    12,
						Status:      "Completed",
						Genres: []models.Genre{
							{ID: 1, Name: "Action"},
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
							{ID: 2, Name: "Drama"},
						},
					},
				}
				mockAnimeService.EXPECT().
					GetAllAnimes(ctx, 1, 10).
					Return(animes, int64(2), nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":[{"id":1,"title":"Anime 1","description":"Description 1","rating":8.5,"episodes":12,"status":"Completed","genres":[{"id":1,"name":"Action"}]},{"id":2,"title":"Anime 2","description":"Description 2","rating":9.0,"episodes":24,"status":"Ongoing","genres":[{"id":2,"name":"Drama"}]}],"total":2,"page":1,"limit":10}`,
		},
		{
			name:  "Invalid Pagination",
			query: "?page=invalid&limit=10",
			setupMock: func() {
				// No mock setup needed for invalid pagination
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":"ERR-003","message":"Invalid pagination parameters","details":"Invalid page number"}}`,
		},
		{
			name:  "Internal Server Error",
			query: "?page=1&limit=10",
			setupMock: func() {
				ctx := context.Background()
				mockAnimeService.EXPECT().
					GetAllAnimes(ctx, 1, 10).
					Return(nil, int64(0), apperrors.NewError(apperrors.ErrInternalServer, "Database error", "Failed to retrieve animes from database", http.StatusInternalServerError, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":{"code":"ERR-004","message":"Database error","details":"Failed to retrieve animes from database"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodGet, "/anime"+tt.query, nil)
			w := httptest.NewRecorder()

			handler.GetAllAnimesHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
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
		expectedBody   string
	}{
		{
			name:        "Success",
			requestBody: `{"title":"New Anime","description":"New Description","rating":8.5,"episodes":12,"status":"Completed","genre_ids":[1,2],"tag_ids":[1,2]}`,
			setupMock: func() {
				ctx := context.Background()
				expectedAnime := &models.Anime{
					Title:       "New Anime",
					Description: "New Description",
					Rating:      8.5,
					Episodes:    12,
					Status:      "Completed",
				}
				mockAnimeService.EXPECT().
					CreateAnime(ctx, expectedAnime).
					Return(nil)
				mockAnimeService.EXPECT().
					GetAnimeByID(ctx, uint(1)).
					Return(expectedAnime, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":1,"title":"New Anime","description":"New Description","rating":8.5,"episodes":12,"status":"Completed"}`,
		},
		{
			name:        "Invalid Request Body",
			requestBody: `{"invalid":"json"`,
			setupMock: func() {
				// No mock setup needed for invalid request body
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":"ERR-003","message":"Invalid request body","details":"Invalid JSON format"}}`,
		},
		{
			name:        "Validation Error",
			requestBody: `{"title":""}`,
			setupMock: func() {
				// No mock setup needed for validation error
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":"ERR-003","message":"Validation error","details":"Title is required"}}`,
		},
		{
			name:        "Internal Server Error",
			requestBody: `{"title":"New Anime","description":"New Description","rating":8.5,"episodes":12,"status":"Completed"}`,
			setupMock: func() {
				ctx := context.Background()
				expectedAnime := &models.Anime{
					Title:       "New Anime",
					Description: "New Description",
					Rating:      8.5,
					Episodes:    12,
					Status:      "Completed",
				}
				mockAnimeService.EXPECT().
					CreateAnime(ctx, expectedAnime).
					Return(apperrors.NewError(apperrors.ErrInternalServer, "Database error", "Failed to create anime in database", http.StatusInternalServerError, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":{"code":"ERR-004","message":"Database error","details":"Failed to create anime in database"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodPost, "/anime", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateAnimeHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
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
		expectedBody   string
	}{
		{
			name:        "Success",
			animeID:     "1",
			requestBody: `{"title":"Updated Anime","description":"Updated Description","rating":9.0,"episodes":24,"status":"Ongoing"}`,
			setupMock: func() {
				ctx := context.Background()
				existingAnime := &models.Anime{
					ID:          1,
					Title:       "Original Anime",
					Description: "Original Description",
					Rating:      8.5,
					Episodes:    12,
					Status:      "Completed",
				}
				updatedAnime := &models.Anime{
					ID:          1,
					Title:       "Updated Anime",
					Description: "Updated Description",
					Rating:      9.0,
					Episodes:    24,
					Status:      "Ongoing",
				}
				mockAnimeService.EXPECT().
					GetAnimeByID(ctx, uint(1)).
					Return(existingAnime, nil)
				mockAnimeService.EXPECT().
					UpdateAnime(ctx, updatedAnime).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"title":"Updated Anime","description":"Updated Description","rating":9.0,"episodes":24,"status":"Ongoing"}`,
		},
		{
			name:        "Invalid ID",
			animeID:     "invalid",
			requestBody: `{"title":"Updated Anime"}`,
			setupMock: func() {
				// No mock setup needed for invalid ID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":"ERR-003","message":"Invalid ID format","details":"The provided ID is not a valid unsigned integer."}}`,
		},
		{
			name:        "Not Found",
			animeID:     "999",
			requestBody: `{"title":"Updated Anime"}`,
			setupMock: func() {
				ctx := context.Background()
				mockAnimeService.EXPECT().
					GetAnimeByID(ctx, uint(999)).
					Return(nil, apperrors.NewError(apperrors.ErrResourceNotFound, "Anime not found", "The requested anime could not be found", http.StatusNotFound, nil, nil))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":{"code":"ERR-002","message":"Anime not found","details":"The requested anime could not be found"}}`,
		},
		{
			name:        "Internal Server Error",
			animeID:     "1",
			requestBody: `{"title":"Updated Anime"}`,
			setupMock: func() {
				ctx := context.Background()
				existingAnime := &models.Anime{
					ID:    1,
					Title: "Original Anime",
				}
				mockAnimeService.EXPECT().
					GetAnimeByID(ctx, uint(1)).
					Return(existingAnime, nil)
				mockAnimeService.EXPECT().
					UpdateAnime(ctx, gomock.Any()).
					Return(apperrors.NewError(apperrors.ErrInternalServer, "Database error", "Failed to update anime in database", http.StatusInternalServerError, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":{"code":"ERR-004","message":"Database error","details":"Failed to update anime in database"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodPut, "/anime/"+tt.animeID, bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{"id": tt.animeID})
			w := httptest.NewRecorder()

			handler.UpdateAnimeHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
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
		expectedBody   string
	}{
		{
			name:    "Success",
			animeID: "1",
			setupMock: func() {
				ctx := context.Background()
				mockAnimeService.EXPECT().
					DeleteAnime(ctx, uint(1)).
					Return(nil)
			},
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
		{
			name:    "Invalid ID",
			animeID: "invalid",
			setupMock: func() {
				// No mock setup needed for invalid ID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":"ERR-003","message":"Invalid ID format","details":"The provided ID is not a valid unsigned integer."}}`,
		},
		{
			name:    "Not Found",
			animeID: "999",
			setupMock: func() {
				ctx := context.Background()
				mockAnimeService.EXPECT().
					DeleteAnime(ctx, uint(999)).
					Return(apperrors.NewError(apperrors.ErrResourceNotFound, "Anime not found", "The requested anime could not be found", http.StatusNotFound, nil, nil))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":{"code":"ERR-002","message":"Anime not found","details":"The requested anime could not be found"}}`,
		},
		{
			name:    "Internal Server Error",
			animeID: "1",
			setupMock: func() {
				ctx := context.Background()
				mockAnimeService.EXPECT().
					DeleteAnime(ctx, uint(1)).
					Return(apperrors.NewError(apperrors.ErrInternalServer, "Database error", "Failed to delete anime from database", http.StatusInternalServerError, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":{"code":"ERR-004","message":"Database error","details":"Failed to delete anime from database"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodDelete, "/anime/"+tt.animeID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.animeID})
			w := httptest.NewRecorder()

			handler.DeleteAnimeHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			} else {
				assert.Empty(t, w.Body.String())
			}
		})
	}
}
