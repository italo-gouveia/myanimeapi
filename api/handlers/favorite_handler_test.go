package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"myanimeapi/api/middleware"
	"myanimeapi/api/mocks"
	"myanimeapi/api/models"

	"gorm.io/gorm"
)

// TestGetFavoritesHandler tests the GetFavoritesHandler function
func TestGetFavoritesHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		setupMock      func(*mocks.MockDBInterface)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Success",
			userID: 1,
			setupMock: func(mockDB *mocks.MockDBInterface) {
				// Mock user retrieval
				mockDB.EXPECT().First(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
					user := dest.(*models.User)
					user.ID = 1
					user.Username = "testuser"
					return &gorm.DB{}
				})

				// Create a mock DB object that will be returned by each method
				mockGormDB := &gorm.DB{}

				// Set up the chain of method calls
				mockDB.EXPECT().WithContext(gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().Preload(gomock.Any(), gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().Preload(gomock.Any(), gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().Where(gomock.Any(), gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().Find(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
					favorites := dest.(*[]models.Favorite)
					*favorites = []models.Favorite{
						{
							ID:      1,
							UserID:  1,
							AnimeID: 1,
							Anime: models.Anime{
								ID:          1,
								Title:       "Test Anime",
								Description: "Test Description",
							},
							User: models.User{
								ID:       1,
								Username: "testuser",
							},
						},
					}
					return mockGormDB
				})
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":1,"user_id":1,"anime_id":1,"anime":{"id":1,"title":"Test Anime","description":"Test Description"},"user":{"id":1,"username":"testuser"}}]`,
		},
		{
			name:   "Database Error",
			userID: 1,
			setupMock: func(mockDB *mocks.MockDBInterface) {
				// Mock user retrieval
				mockDB.EXPECT().First(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
					user := dest.(*models.User)
					user.ID = 1
					user.Username = "testuser"
					return &gorm.DB{}
				})

				// Create a mock DB object that will be returned by each method
				mockGormDB := &gorm.DB{}

				// Set up the chain of method calls
				mockDB.EXPECT().WithContext(gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().Preload(gomock.Any(), gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().Preload(gomock.Any(), gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().Where(gomock.Any(), gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().Find(gomock.Any(), gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().GetError().Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":{"code":"ERR-004","message":"Failed to get favorites","details":"database error"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create a mock DB
			mockDB := mocks.NewMockDBInterface(ctrl)
			tt.setupMock(mockDB)

			// Create a new handler with the mock DB
			handler := NewFavoriteHandler(mockDB)

			// Create a new request
			req := httptest.NewRequest(http.MethodGet, "/favorites", nil)
			ctx := context.WithValue(req.Context(), middleware.UserContextKey, tt.userID)
			req = req.WithContext(ctx)

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

// TestAddFavoriteHandler tests the AddFavoriteHandler function
func TestAddFavoriteHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		requestBody    string
		setupMock      func(*mocks.MockDBInterface)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Success",
			userID:      1,
			requestBody: `{"anime_id": 1}`,
			setupMock: func(mockDB *mocks.MockDBInterface) {
				// Create a mock DB object that will be returned by each method
				mockGormDB := &gorm.DB{}

				// Mock user retrieval
				mockDB.EXPECT().WithContext(gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().First(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
					user := dest.(*models.User)
					user.ID = 1
					user.Username = "testuser"
					return mockGormDB
				})

				// Mock anime existence check
				mockDB.EXPECT().First(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
					anime := dest.(*models.Anime)
					anime.ID = 1
					anime.Title = "Test Anime"
					return mockGormDB
				})

				// Mock favorite existence check
				mockDB.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().First(gomock.Any(), gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().GetError().Return(errors.New("record not found"))

				// Mock favorite creation
				mockDB.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, value interface{}) *gorm.DB {
					favorite := value.(*models.Favorite)
					favorite.ID = 1
					favorite.UserID = 1
					favorite.AnimeID = 1
					return mockGormDB
				})

				// Mock preloading related data
				mockDB.EXPECT().Preload(gomock.Any(), gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().Preload(gomock.Any(), gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().First(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
					favorite := dest.(*models.Favorite)
					favorite.ID = 1
					favorite.UserID = 1
					favorite.AnimeID = 1
					favorite.Anime = models.Anime{
						ID:          1,
						Title:       "Test Anime",
						Description: "Test Description",
					}
					favorite.User = models.User{
						ID:       1,
						Username: "testuser",
					}
					return mockGormDB
				})
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":1,"user_id":1,"anime_id":1,"anime":{"id":1,"title":"Test Anime","description":"Test Description"},"user":{"id":1,"username":"testuser"}}`,
		},
		{
			name:        "User Not Found",
			userID:      1,
			requestBody: `{"anime_id": 1}`,
			setupMock: func(mockDB *mocks.MockDBInterface) {
				// Create a mock DB object that will be returned by each method
				mockGormDB := &gorm.DB{}

				// Mock user retrieval
				mockDB.EXPECT().WithContext(gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().First(gomock.Any(), gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().GetError().Return(errors.New("user not found"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":{"code":"ERR-001","message":"Unauthorized","details":"User not found"}}`,
		},
		{
			name:        "Anime Not Found",
			userID:      1,
			requestBody: `{"anime_id": 1}`,
			setupMock: func(mockDB *mocks.MockDBInterface) {
				// Create a mock DB object that will be returned by each method
				mockGormDB := &gorm.DB{}

				// Mock user retrieval
				mockDB.EXPECT().WithContext(gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().First(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
					user := dest.(*models.User)
					user.ID = 1
					return mockGormDB
				})

				// Mock anime existence check
				mockDB.EXPECT().First(gomock.Any(), gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().GetError().Return(errors.New("record not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":{"code":"ERR-002","message":"Anime not found","details":"record not found"}}`,
		},
		{
			name:        "Favorite Already Exists",
			userID:      1,
			requestBody: `{"anime_id": 1}`,
			setupMock: func(mockDB *mocks.MockDBInterface) {
				// Create a mock DB object that will be returned by each method
				mockGormDB := &gorm.DB{}

				// Mock user retrieval
				mockDB.EXPECT().WithContext(gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().First(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
					user := dest.(*models.User)
					user.ID = 1
					return mockGormDB
				})

				// Mock anime existence check
				mockDB.EXPECT().First(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
					anime := dest.(*models.Anime)
					anime.ID = 1
					return mockGormDB
				})

				// Mock favorite existence check
				mockDB.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().First(gomock.Any(), gomock.Any()).Return(mockGormDB)
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"error":{"code":"ERR-008","message":"Anime already in favorites","details":"This anime is already in your favorites list"}}`,
		},
		{
			name:        "Invalid Request Body",
			userID:      1,
			requestBody: `{"anime_id": "invalid"}`,
			setupMock: func(mockDB *mocks.MockDBInterface) {
				// No mock setup needed for invalid request body
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":"ERR-003","message":"Invalid request body","details":"json: cannot unmarshal string into Go struct field FavoriteCreateRequest.anime_id of type uint"}}`,
		},
		{
			name:        "Database Error",
			userID:      1,
			requestBody: `{"anime_id": 1}`,
			setupMock: func(mockDB *mocks.MockDBInterface) {
				// Create a mock DB object that will be returned by each method
				mockGormDB := &gorm.DB{}

				// Mock user retrieval
				mockDB.EXPECT().WithContext(gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().First(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
					user := dest.(*models.User)
					user.ID = 1
					return mockGormDB
				})

				// Mock anime existence check
				mockDB.EXPECT().First(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
					anime := dest.(*models.Anime)
					anime.ID = 1
					return mockGormDB
				})

				// Mock favorite existence check
				mockDB.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().First(gomock.Any(), gomock.Any()).Return(mockGormDB)

				// Mock favorite creation
				mockDB.EXPECT().Create(gomock.Any(), gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().GetError().Return(errors.New("database error"))

				// Mock preloading related data
				mockDB.EXPECT().Preload(gomock.Any(), gomock.Any()).Return(mockGormDB)
				mockDB.EXPECT().Preload(gomock.Any(), gomock.Any()).Return(mockGormDB)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":{"code":"ERR-004","message":"Failed to get favorites","details":"database error"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create a mock DB
			mockDB := mocks.NewMockDBInterface(ctrl)
			tt.setupMock(mockDB)

			// Create a new handler with the mock DB
			handler := NewFavoriteHandler(mockDB)

			// Create a new request with the test body
			req := httptest.NewRequest(http.MethodPost, "/favorites", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			ctx := context.WithValue(req.Context(), middleware.UserContextKey, tt.userID)
			req = req.WithContext(ctx)

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

// TestRemoveFavoriteHandler tests the RemoveFavoriteHandler function
func TestRemoveFavoriteHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		animeID        string
		setupMock      func(*mocks.MockDBInterface)
		expectedStatus int
	}{
		{
			name:    "Success",
			userID:  1,
			animeID: "1",
			setupMock: func(mockDB *mocks.MockDBInterface) {
				// Mock user retrieval
				mockDB.EXPECT().WithContext(gomock.Any()).Return(&gorm.DB{})
				mockDB.EXPECT().First(gomock.Any(), uint(1)).DoAndReturn(func(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
					user := dest.(*models.User)
					user.ID = 1
					return &gorm.DB{}
				})

				// Mock favorite deletion
				mockDB.EXPECT().Where("user_id = ? AND anime_id = ?", uint(1), uint(1)).Return(&gorm.DB{})
				mockDB.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(&gorm.DB{})
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:    "User Not Found",
			userID:  1,
			animeID: "1",
			setupMock: func(mockDB *mocks.MockDBInterface) {
				mockDB.EXPECT().WithContext(gomock.Any()).Return(&gorm.DB{})
				mockDB.EXPECT().First(gomock.Any(), uint(1)).Return(&gorm.DB{})
				mockDB.EXPECT().GetError().Return(errors.New("user not found"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:    "Invalid Anime ID",
			userID:  1,
			animeID: "invalid",
			setupMock: func(mockDB *mocks.MockDBInterface) {
				// No mock setup needed for invalid anime ID
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create a mock DB
			mockDB := mocks.NewMockDBInterface(ctrl)
			tt.setupMock(mockDB)

			// Create a new handler with the mock DB
			handler := NewFavoriteHandler(mockDB)

			// Create a new request
			req := httptest.NewRequest(http.MethodDelete, "/favorites/"+tt.animeID, nil)
			ctx := context.WithValue(req.Context(), middleware.UserContextKey, tt.userID)
			req = req.WithContext(ctx)

			// Set up the URL variables for the request
			vars := map[string]string{
				"id": tt.animeID,
			}
			req = mux.SetURLVars(req, vars)

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handler.RemoveFavoriteHandler(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
