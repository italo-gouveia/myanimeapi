package handlers

/*import (
	"bytes"
	"context"
	"myanimeapi/api/mocks"
	"myanimeapi/api/models"
	"myanimeapi/internal/errors"
	"net/http"
	"net/http/httptest"
	"testing"

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
					ID:    1,
					Title: "Test Anime",
				}
				mockAnimeService.EXPECT().
					GetAnimeByID(ctx, uint(1)).
					Return(anime, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"title":"Test Anime"}`,
		},
		{
			name:    "Invalid ID",
			animeID: "invalid",
			setupMock: func() {
				// No mock setup needed for invalid ID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid anime ID"}`,
		},
		{
			name:    "Not Found",
			animeID: "999",
			setupMock: func() {
				ctx := context.Background()
				mockAnimeService.EXPECT().
					GetAnimeByID(ctx, uint(999)).
					Return(nil, errors.NewError(errors.ErrResourceNotFound, "Anime not found", "", http.StatusNotFound))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"Anime not found"}`,
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
					{ID: 1, Title: "Anime 1"},
					{ID: 2, Title: "Anime 2"},
				}
				mockAnimeService.EXPECT().
					GetAllAnimes(ctx, 1, 10).
					Return(animes, int64(2), nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"animes":[{"id":1,"title":"Anime 1"},{"id":2,"title":"Anime 2"}],"total":2,"page":1,"limit":10}`,
		},
		{
			name:  "Invalid Pagination",
			query: "?page=invalid&limit=10",
			setupMock: func() {
				// No mock setup needed for invalid pagination
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid pagination parameters"}`,
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
			requestBody: `{"title":"New Anime"}`,
			setupMock: func() {
				ctx := context.Background()
				expectedAnime := &models.Anime{Title: "New Anime"}
				mockAnimeService.EXPECT().
					CreateAnime(ctx, expectedAnime).
					Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"message":"Anime created successfully"}`,
		},
		{
			name:        "Invalid Request Body",
			requestBody: `{"invalid":"json"`,
			setupMock: func() {
				// No mock setup needed for invalid request body
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid request body"}`,
		},
		{
			name:        "Validation Error",
			requestBody: `{"title":""}`,
			setupMock: func() {
				// No mock setup needed for validation error
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Title are required"}`,
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
			requestBody: `{"id":1,"title":"Updated Anime"}`,
			setupMock: func() {
				ctx := context.Background()
				expectedAnime := &models.Anime{ID: 1, Title: "Updated Anime"}
				mockAnimeService.EXPECT().
					UpdateAnime(ctx, expectedAnime).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Anime updated successfully"}`,
		},
		{
			name:        "Invalid ID",
			animeID:     "invalid",
			requestBody: `{"id":1,"title":"Updated Anime"}`,
			setupMock: func() {
				// No mock setup needed for invalid ID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid anime ID"}`,
		},
		{
			name:        "Not Found",
			animeID:     "999",
			requestBody: `{"id":999,"title":"Updated Anime"}`,
			setupMock: func() {
				ctx := context.Background()
				expectedAnime := &models.Anime{ID: 999, Title: "Updated Anime"}
				mockAnimeService.EXPECT().
					UpdateAnime(ctx, expectedAnime).
					Return(errors.NewError(errors.ErrResourceNotFound, "Anime not found", "", http.StatusNotFound))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"Anime not found"}`,
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
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Anime deleted successfully"}`,
		},
		{
			name:    "Invalid ID",
			animeID: "invalid",
			setupMock: func() {
				// No mock setup needed for invalid ID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid anime ID"}`,
		},
		{
			name:    "Not Found",
			animeID: "999",
			setupMock: func() {
				ctx := context.Background()
				mockAnimeService.EXPECT().
					DeleteAnime(ctx, uint(999)).
					Return(errors.NewError(errors.ErrResourceNotFound, "Anime not found", "", http.StatusNotFound))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"Anime not found"}`,
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
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}*/
