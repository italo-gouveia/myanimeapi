package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"myanimeapi/api/middleware"
	"myanimeapi/api/models"

	"gorm.io/gorm"
)

// MockDB is a mock implementation of the DBInterface
type MockDB struct {
	mock.Mock
}

func (m *MockDB) WithContext(ctx context.Context) *gorm.DB {
	args := m.Called(ctx)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) First(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(ctx, dest, conds)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Find(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(ctx, dest, conds)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Create(ctx context.Context, value interface{}) *gorm.DB {
	args := m.Called(ctx, value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Delete(ctx context.Context, value interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(ctx, value, conds)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Preload(query string, ctx context.Context, args ...interface{}) *gorm.DB {
	mockArgs := m.Called(query, ctx, args)
	return mockArgs.Get(0).(*gorm.DB)
}

func (m *MockDB) Where(ctx context.Context, query interface{}, args ...interface{}) *gorm.DB {
	mockArgs := m.Called(ctx, query, args)
	return mockArgs.Get(0).(*gorm.DB)
}

// Add missing methods to fully implement the db.DBInterface

func (m *MockDB) Begin(ctx context.Context) *gorm.DB {
	args := m.Called(ctx)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Commit(ctx context.Context) *gorm.DB {
	args := m.Called(ctx)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Rollback(ctx context.Context) *gorm.DB {
	args := m.Called(ctx)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Save(ctx context.Context, value interface{}) *gorm.DB {
	args := m.Called(ctx, value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Unscoped(ctx context.Context) *gorm.DB {
	args := m.Called(ctx)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) GetError() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDB) IsHealthy() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockDB) Offset(offset int) *gorm.DB {
	args := m.Called(offset)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Limit(limit int) *gorm.DB {
	args := m.Called(limit)
	return args.Get(0).(*gorm.DB)
}

// TestGetFavoritesHandler tests the GetFavoritesHandler function
func TestGetFavoritesHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		setupMock      func(*MockDB)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Success",
			userID: 1,
			setupMock: func(mockDB *MockDB) {
				mockDB.On("WithContext", mock.Anything).Return(mockDB)
				mockDB.On("Preload", "Anime", mock.Anything).Return(mockDB)
				mockDB.On("Preload", "User", mock.Anything).Return(mockDB)
				mockDB.On("Where", "user_id = ?", uint(1)).Return(mockDB)
				mockDB.On("Find", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					favorites := args.Get(0).(*[]models.Favorite)
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
				})
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":1,"user_id":1,"anime_id":1,"anime":{"id":1,"title":"Test Anime","description":"Test Description"},"user":{"id":1,"username":"testuser"}}]`,
		},
		{
			name:   "Database Error",
			userID: 1,
			setupMock: func(mockDB *MockDB) {
				mockDB.On("WithContext", mock.Anything).Return(mockDB)
				mockDB.On("Preload", "Anime", mock.Anything).Return(mockDB)
				mockDB.On("Preload", "User", mock.Anything).Return(mockDB)
				mockDB.On("Where", "user_id = ?", uint(1)).Return(mockDB)
				mockDB.On("Find", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":{"code":"ERR-004","message":"Failed to get favorites","details":"database error"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock DB
			mockDB := new(MockDB)
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

			// Verify that all expected mock calls were made
			mockDB.AssertExpectations(t)
		})
	}
}

// TestAddFavoriteHandler tests the AddFavoriteHandler function
func TestAddFavoriteHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		requestBody    string
		setupMock      func(*MockDB)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Success",
			userID:      1,
			requestBody: `{"anime_id": 1}`,
			setupMock: func(mockDB *MockDB) {
				// Mock user retrieval
				mockDB.On("WithContext", mock.Anything).Return(mockDB)
				mockDB.On("First", mock.Anything, uint(1)).Return(nil).Run(func(args mock.Arguments) {
					user := args.Get(0).(*models.User)
					user.ID = 1
					user.Username = "testuser"
				})

				// Mock anime existence check
				mockDB.On("First", mock.Anything, uint(1)).Return(nil).Run(func(args mock.Arguments) {
					anime := args.Get(0).(*models.Anime)
					anime.ID = 1
					anime.Title = "Test Anime"
				})

				// Mock favorite existence check
				mockDB.On("Where", "user_id = ? AND anime_id = ?", uint(1), uint(1)).Return(mockDB)
				mockDB.On("First", mock.Anything, mock.Anything).Return(errors.New("record not found"))

				// Mock favorite creation
				mockDB.On("Create", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					favorite := args.Get(0).(*models.Favorite)
					favorite.ID = 1
					favorite.UserID = 1
					favorite.AnimeID = 1
				})

				// Mock preloading related data
				mockDB.On("Preload", "Anime", mock.Anything).Return(mockDB)
				mockDB.On("Preload", "User", mock.Anything).Return(mockDB)
				mockDB.On("First", mock.Anything, uint(1)).Return(nil).Run(func(args mock.Arguments) {
					favorite := args.Get(0).(*models.Favorite)
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
				})
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":1,"user_id":1,"anime_id":1,"anime":{"id":1,"title":"Test Anime","description":"Test Description"},"user":{"id":1,"username":"testuser"}}`,
		},
		{
			name:        "User Not Found",
			userID:      1,
			requestBody: `{"anime_id": 1}`,
			setupMock: func(mockDB *MockDB) {
				mockDB.On("WithContext", mock.Anything).Return(mockDB)
				mockDB.On("First", mock.Anything, uint(1)).Return(errors.New("user not found"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":{"code":"ERR-001","message":"Unauthorized","details":"User not found"}}`,
		},
		{
			name:        "Anime Not Found",
			userID:      1,
			requestBody: `{"anime_id": 1}`,
			setupMock: func(mockDB *MockDB) {
				// Mock user retrieval
				mockDB.On("WithContext", mock.Anything).Return(mockDB)
				mockDB.On("First", mock.Anything, uint(1)).Return(nil).Run(func(args mock.Arguments) {
					user := args.Get(0).(*models.User)
					user.ID = 1
				})

				// Mock anime existence check
				mockDB.On("First", mock.Anything, uint(1)).Return(errors.New("record not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":{"code":"ERR-002","message":"Anime not found","details":"record not found"}}`,
		},
		{
			name:        "Favorite Already Exists",
			userID:      1,
			requestBody: `{"anime_id": 1}`,
			setupMock: func(mockDB *MockDB) {
				// Mock user retrieval
				mockDB.On("WithContext", mock.Anything).Return(mockDB)
				mockDB.On("First", mock.Anything, uint(1)).Return(nil).Run(func(args mock.Arguments) {
					user := args.Get(0).(*models.User)
					user.ID = 1
				})

				// Mock anime existence check
				mockDB.On("First", mock.Anything, uint(1)).Return(nil).Run(func(args mock.Arguments) {
					anime := args.Get(0).(*models.Anime)
					anime.ID = 1
				})

				// Mock favorite existence check
				mockDB.On("Where", "user_id = ? AND anime_id = ?", uint(1), uint(1)).Return(mockDB)
				mockDB.On("First", mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"error":{"code":"ERR-008","message":"Anime already in favorites","details":"This anime is already in your favorites list"}}`,
		},
		{
			name:        "Invalid Request Body",
			userID:      1,
			requestBody: `{"anime_id": "invalid"}`,
			setupMock: func(mockDB *MockDB) {
				// No mock setup needed for invalid request body
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":"ERR-003","message":"Invalid request body","details":"json: cannot unmarshal string into Go struct field FavoriteCreateRequest.anime_id of type uint"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock DB
			mockDB := new(MockDB)
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

			// Verify that all expected mock calls were made
			mockDB.AssertExpectations(t)
		})
	}
}

// TestRemoveFavoriteHandler tests the RemoveFavoriteHandler function
func TestRemoveFavoriteHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		animeID        string
		setupMock      func(*MockDB)
		expectedStatus int
	}{
		{
			name:    "Success",
			userID:  1,
			animeID: "1",
			setupMock: func(mockDB *MockDB) {
				// Mock user retrieval
				mockDB.On("WithContext", mock.Anything).Return(mockDB)
				mockDB.On("First", mock.Anything, uint(1)).Return(nil).Run(func(args mock.Arguments) {
					user := args.Get(0).(*models.User)
					user.ID = 1
				})

				// Mock favorite deletion
				mockDB.On("Where", "user_id = ? AND anime_id = ?", uint(1), uint(1)).Return(mockDB)
				mockDB.On("Delete", mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:    "User Not Found",
			userID:  1,
			animeID: "1",
			setupMock: func(mockDB *MockDB) {
				mockDB.On("WithContext", mock.Anything).Return(mockDB)
				mockDB.On("First", mock.Anything, uint(1)).Return(errors.New("user not found"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:    "Invalid Anime ID",
			userID:  1,
			animeID: "invalid",
			setupMock: func(mockDB *MockDB) {
				// No mock setup needed for invalid anime ID
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock DB
			mockDB := new(MockDB)
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

			// Verify that all expected mock calls were made
			mockDB.AssertExpectations(t)
		})
	}
}
