package handlers

import (
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

func TestGetFavoritesHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		setupMock      func(*mocks.MockFavoriteServiceInterface)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Success",
			userID: 1,
			setupMock: func(mockService *mocks.MockFavoriteServiceInterface) {
				favorites := []models.Favorite{
					{
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
				}
				mockService.EXPECT().
					GetFavorites(gomock.Any(), uint(1)).
					Return(favorites, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":1,"user_id":1,"anime_id":1,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","anime":{"id":1,"title":"Test Anime","description":"Test Description","rating":8.5,"episodes":12,"status":"Completed","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","start_date":"0001-01-01T00:00:00Z","end_date":"0001-01-01T00:00:00Z"},"user":{"bio":"","created_at":"0001-01-01T00:00:00Z","deleted_at":"0001-01-01T00:00:00Z","email":"","id":0,"is_active":false,"is_admin":false,"profile_pic":"","social_links":null,"updated_at":"0001-01-01T00:00:00Z","username":""}}]`,
		},
		{
			name:   "No Favorites",
			userID: 1,
			setupMock: func(mockService *mocks.MockFavoriteServiceInterface) {
				mockService.EXPECT().
					GetFavorites(gomock.Any(), uint(1)).
					Return([]models.Favorite{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[]`,
		},
		{
			name:   "Database Error",
			userID: 1,
			setupMock: func(mockService *mocks.MockFavoriteServiceInterface) {
				mockService.EXPECT().
					GetFavorites(gomock.Any(), uint(1)).
					Return(nil, apperrors.NewError(apperrors.ErrInternalServer, "Database error", "Failed to retrieve favorites from database", http.StatusInternalServerError, map[string]interface{}{"user_id": 1}, nil))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":{"code":"ERR-004","message":"Database error","details":"Failed to retrieve favorites from database","context":{"user_id":1}}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create a mock service
			mockService := mocks.NewMockFavoriteServiceInterface(ctrl)
			tt.setupMock(mockService)

			// Create a new handler with the mock service
			handler := NewFavoriteHandler(mockService)

			// Create a new request
			req := httptest.NewRequest(http.MethodGet, "/favorites", nil)
			req = req.WithContext(CreateTestContext())

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handler.GetFavoritesHandler(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Check the response body
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestAddFavoriteHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		animeID        string
		setupMock      func(*mocks.MockFavoriteServiceInterface)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:    "Success",
			userID:  1,
			animeID: "1",
			setupMock: func(mockService *mocks.MockFavoriteServiceInterface) {
				favorite := &models.Favorite{
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
				}
				mockService.EXPECT().
					AddFavorite(gomock.Any(), uint(1), uint(1)).
					Return(favorite, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":1,"user_id":1,"anime_id":1,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","anime":{"id":1,"title":"Test Anime","description":"Test Description","rating":8.5,"episodes":12,"status":"Completed","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","start_date":"0001-01-01T00:00:00Z","end_date":"0001-01-01T00:00:00Z"},"user":{"bio":"","created_at":"0001-01-01T00:00:00Z","deleted_at":"0001-01-01T00:00:00Z","email":"","id":0,"is_active":false,"is_admin":false,"profile_pic":"","social_links":null,"updated_at":"0001-01-01T00:00:00Z","username":""}}`,
		},
		{
			name:    "Invalid Anime ID",
			userID:  1,
			animeID: "invalid",
			setupMock: func(mockService *mocks.MockFavoriteServiceInterface) {
				// No mock setup needed for invalid anime ID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":"ERR-001","message":"Invalid anime ID format","details":"The provided anime ID must be a valid unsigned integer"}}`,
		},
		{
			name:    "Anime Not Found",
			userID:  1,
			animeID: "999",
			setupMock: func(mockService *mocks.MockFavoriteServiceInterface) {
				mockService.EXPECT().
					AddFavorite(gomock.Any(), uint(1), uint(999)).
					Return(nil, apperrors.NewError(apperrors.ErrResourceNotFound, "Anime not found", "The requested anime could not be found", http.StatusNotFound, map[string]interface{}{"user_id": 1, "anime_id": 999}, nil))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":{"code":"ERR-002","message":"Anime not found","details":"The requested anime could not be found","context":{"user_id":1,"anime_id":999}}}`,
		},
		{
			name:    "Favorite Already Exists",
			userID:  1,
			animeID: "1",
			setupMock: func(mockService *mocks.MockFavoriteServiceInterface) {
				mockService.EXPECT().
					AddFavorite(gomock.Any(), uint(1), uint(1)).
					Return(nil, apperrors.NewError(apperrors.ErrConflict, "Favorite already exists", "This anime is already in your favorites list", http.StatusConflict, map[string]interface{}{"user_id": 1, "anime_id": 1}, nil))
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"error":{"code":"ERR-008","message":"Favorite already exists","details":"This anime is already in your favorites list","context":{"user_id":1,"anime_id":1}}}`,
		},
		{
			name:    "Database Error",
			userID:  1,
			animeID: "1",
			setupMock: func(mockService *mocks.MockFavoriteServiceInterface) {
				mockService.EXPECT().
					AddFavorite(gomock.Any(), uint(1), uint(1)).
					Return(nil, apperrors.NewError(apperrors.ErrInternalServer, "Database error", "Failed to add favorite to database", http.StatusInternalServerError, map[string]interface{}{"user_id": 1, "anime_id": 1}, nil))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":{"code":"ERR-004","message":"Database error","details":"Failed to add favorite to database","context":{"user_id":1,"anime_id":1}}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create a mock service
			mockService := mocks.NewMockFavoriteServiceInterface(ctrl)
			tt.setupMock(mockService)

			// Create a new handler with the mock service
			handler := NewFavoriteHandler(mockService)

			// Create a new request
			req := httptest.NewRequest(http.MethodPost, "/favorites/"+tt.animeID, nil)
			req = req.WithContext(CreateTestContext())

			// Set up the URL variables for the request
			vars := map[string]string{
				"anime_id": tt.animeID,
			}
			req = mux.SetURLVars(req, vars)

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handler.AddFavoriteHandler(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Check the response body
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestRemoveFavoriteHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		animeID        string
		setupMock      func(*mocks.MockFavoriteServiceInterface)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:    "Success",
			userID:  1,
			animeID: "1",
			setupMock: func(mockService *mocks.MockFavoriteServiceInterface) {
				mockService.EXPECT().
					RemoveFavorite(gomock.Any(), uint(1), uint(1)).
					Return(nil)
			},
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
		{
			name:    "Invalid Anime ID",
			userID:  1,
			animeID: "invalid",
			setupMock: func(mockService *mocks.MockFavoriteServiceInterface) {
				// No mock setup needed for invalid anime ID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":"ERR-001","message":"Invalid anime ID format","details":"The provided anime ID must be a valid unsigned integer"}}`,
		},
		{
			name:    "Favorite Not Found",
			userID:  1,
			animeID: "999",
			setupMock: func(mockService *mocks.MockFavoriteServiceInterface) {
				mockService.EXPECT().
					RemoveFavorite(gomock.Any(), uint(1), uint(999)).
					Return(apperrors.NewError(apperrors.ErrResourceNotFound, "Favorite not found", "The requested favorite could not be found", http.StatusNotFound, map[string]interface{}{"user_id": 1, "anime_id": 999}, nil))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":{"code":"ERR-002","message":"Favorite not found","details":"The requested favorite could not be found","context":{"user_id":1,"anime_id":999}}}`,
		},
		{
			name:    "Database Error",
			userID:  1,
			animeID: "1",
			setupMock: func(mockService *mocks.MockFavoriteServiceInterface) {
				mockService.EXPECT().
					RemoveFavorite(gomock.Any(), uint(1), uint(1)).
					Return(apperrors.NewError(apperrors.ErrInternalServer, "Database error", "Failed to remove favorite from database", http.StatusInternalServerError, map[string]interface{}{"user_id": 1, "anime_id": 1}, nil))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":{"code":"ERR-004","message":"Database error","details":"Failed to remove favorite from database","context":{"user_id":1,"anime_id":1}}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create a mock service
			mockService := mocks.NewMockFavoriteServiceInterface(ctrl)
			tt.setupMock(mockService)

			// Create a new handler with the mock service
			handler := NewFavoriteHandler(mockService)

			// Create a new request
			req := httptest.NewRequest(http.MethodDelete, "/favorites/"+tt.animeID, nil)
			req = req.WithContext(CreateTestContext())

			// Set up the URL variables for the request
			vars := map[string]string{
				"anime_id": tt.animeID,
			}
			req = mux.SetURLVars(req, vars)

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handler.RemoveFavoriteHandler(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Check the response body
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			} else {
				assert.Empty(t, rr.Body.String())
			}
		})
	}
}
