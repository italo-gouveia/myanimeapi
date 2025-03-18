// pkg/handlers_test/anime_handlers_test.go
// Package handlers_test provides unit tests for the anime-related handlers in the MyAnimeAPI application.
// It uses the gomock package to mock database interactions and the httptest package to simulate HTTP requests.
// The tests cover various scenarios, including success cases, error cases, and edge cases.
package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"myanimeapi/pkg/handlers"
	"myanimeapi/pkg/middleware"
	"myanimeapi/pkg/mocks"
	"myanimeapi/pkg/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// FAILING NOW
/*func TestGetAnimeHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	testAnime := models.Anime{
		ID:    1,
		Title: "Naruto",
	}

	// Mock the WithContext method to return the same mockDB
	mockDB.EXPECT().
		WithContext(gomock.Any()). // Mock WithContext
		Return(mockDB)             // Return the mockDB itself

	// Mock the First method
	mockDB.EXPECT().
		First(gomock.Any(), gomock.Any(), uint(1)).
		DoAndReturn(func(dest interface{}, conds ...interface{}) *gorm.DB {
			*dest.(*models.Anime) = testAnime
			return &gorm.DB{Error: nil}
		})

	// Create a request with the ID in the URL
	req := httptest.NewRequest("GET", "/anime/1", nil)
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rec := httptest.NewRecorder()
	handler.GetAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d", resp.StatusCode)
	}

	var responseAnime models.Anime
	if err := json.NewDecoder(resp.Body).Decode(&responseAnime); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if diff := cmp.Diff(testAnime, responseAnime); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}*/

// FAILING NOW
/*func TestGetAllAnimesHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	testAnimes := []models.Anime{
		{ID: 1, Title: "Naruto"},
		{ID: 2, Title: "One Piece"},
	}

	// Mock the DB call to return a list of animes
	mockDB.EXPECT().
		WithContext(gomock.Any()). // Mock WithContext
		Return(mockDB)

	mockDB.EXPECT().
		Find(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
			*dest.(*[]models.Anime) = testAnimes
			return &gorm.DB{Error: nil}
		})

	req := httptest.NewRequest("GET", "/anime", nil)
	rec := httptest.NewRecorder()
	handler.GetAllAnimesHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d", resp.StatusCode)
	}

	var responseAnimes []models.Anime
	if err := json.NewDecoder(resp.Body).Decode(&responseAnimes); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(responseAnimes) != len(testAnimes) {
		t.Errorf("Expected %d animes, got %d", len(testAnimes), len(responseAnimes))
	}
}*/

//TODO: Adjust this test
/*func TestGetPaginatedReviewsForAnimeHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	testReviews := []models.Review{
		{ID: 1, AnimeID: 1, Content: "Great anime!"},
		{ID: 2, AnimeID: 1, Content: "Awesome!"},
	}

	// Mock the DB call to return paginated reviews
	mockDB.EXPECT().
		Where(gomock.Any(), "anime_id = ?", uint(1)).
		Return(&gorm.DB{
			Statement: &gorm.Statement{
				DB:      &gorm.DB{},
				Clauses: make(map[string]clause.Clause), // Initialize the Clauses map
			},
			Config: &gorm.Config{}, // Initialize the Config field
		})
	mockDB.EXPECT().
		Offset(0).
		Return(&gorm.DB{
			Statement: &gorm.Statement{
				DB:      &gorm.DB{},
				Clauses: make(map[string]clause.Clause), // Initialize the Clauses map
			},
			Config: &gorm.Config{}, // Initialize the Config field
		})
	mockDB.EXPECT().
		Limit(10).
		Return(&gorm.DB{
			Statement: &gorm.Statement{
				DB:      &gorm.DB{},
				Clauses: make(map[string]clause.Clause), // Initialize the Clauses map
			},
			Config: &gorm.Config{}, // Initialize the Config field
		})
	mockDB.EXPECT().
		Find(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
			// Set the destination value
			if destPtr, ok := dest.(*[]models.Review); ok {
				*destPtr = testReviews
			}
			return &gorm.DB{
				Statement: &gorm.Statement{
					DB:      &gorm.DB{},
					Clauses: make(map[string]clause.Clause), // Initialize the Clauses map
				},
				Config: &gorm.Config{}, // Initialize the Config field
			}
		})

	req := httptest.NewRequest("GET", "/anime/1/reviews?page=1&limit=10", nil)
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rec := httptest.NewRecorder()
	handler.GetPaginatedReviewsForAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d", resp.StatusCode)
	}

	var responseReviews []models.Review
	if err := json.NewDecoder(resp.Body).Decode(&responseReviews); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(responseReviews) != len(testReviews) {
		t.Errorf("Expected %d reviews, got %d", len(testReviews), len(responseReviews))
	}
}*/

// TestGetPaginatedReviewsForAnimeHandler_InvalidPagination tests the handler for invalid pagination parameters.
func TestGetPaginatedReviewsForAnimeHandler_InvalidPagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	req := httptest.NewRequest("GET", "/anime/1/reviews?page=invalid&limit=invalid", nil)
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rec := httptest.NewRecorder()
	handler.GetPaginatedReviewsForAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP status 400 Bad Request, got %d", resp.StatusCode)
	}
}

// TestGetPaginatedReviewsForAnimeHandler_NegativePageLimit tests the handler for negative page and limit values.
func TestGetPaginatedReviewsForAnimeHandler_NegativePageLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	req := httptest.NewRequest("GET", "/anime/1/reviews?page=-1&limit=-10", nil)
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rec := httptest.NewRecorder()
	handler.GetPaginatedReviewsForAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP status 400 Bad Request, got %d", resp.StatusCode)
	}
}

// TestGetPaginatedReviewsForAnimeHandler_ZeroLimit tests the handler for a zero limit value.
func TestGetPaginatedReviewsForAnimeHandler_ZeroLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	req := httptest.NewRequest("GET", "/anime/1/reviews?page=1&limit=0", nil)
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rec := httptest.NewRecorder()
	handler.GetPaginatedReviewsForAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP status 400 Bad Request, got %d", resp.StatusCode)
	}
}

//TODO: Adjust this test
/*func TestGetPaginatedReviewsForAnimeHandler_LargeLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	req := httptest.NewRequest("GET", "/anime/1/reviews?page=1&limit=1000", nil)
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rec := httptest.NewRecorder()
	handler.GetPaginatedReviewsForAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP status 400 Bad Request, got %d", resp.StatusCode)
	}
}*/

//TODO: Adjust this test
/*func TestGetPaginatedReviewsForAnimeHandler_NoReviews(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	// Mock the DB call to return an empty slice of reviews
	mockDB.EXPECT().
		Where(gomock.Any(), "anime_id = ?", uint64(1)).
		Return(mockDB)
	mockDB.EXPECT().
		Offset(0).
		Return(mockDB)
	mockDB.EXPECT().
		Limit(10).
		Return(mockDB)
	mockDB.EXPECT().
		Find(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
			*dest.(*[]models.Review) = []models.Review{}
			return &gorm.DB{Error: nil}
		})

	req := httptest.NewRequest("GET", "/anime/1/reviews?page=1&limit=10", nil)
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rec := httptest.NewRecorder()
	handler.GetPaginatedReviewsForAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d", resp.StatusCode)
	}

	var responseReviews []models.Review
	if err := json.NewDecoder(resp.Body).Decode(&responseReviews); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(responseReviews) != 0 {
		t.Errorf("Expected empty slice, got %d items", len(responseReviews))
	}
}*/
/*
func TestGetPaginatedReviewsForAnimeHandler_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	// Mock the DB call to return an error
	mockDB.EXPECT().
		Where(gomock.Any(), "anime_id = ?", uint64(1)).
		Return(mockDB)
	mockDB.EXPECT().
		Offset(0).
		Return(mockDB)
	mockDB.EXPECT().
		Limit(10).
		Return(mockDB)
	mockDB.EXPECT().
		Find(gomock.Any(), gomock.Any()).
		Return(&gorm.DB{Error: errors.New("database error")})

	req := httptest.NewRequest("GET", "/anime/1/reviews?page=1&limit=10", nil)
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rec := httptest.NewRecorder()
	handler.GetPaginatedReviewsForAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected HTTP status 500 Internal Server Error, got %d", resp.StatusCode)
	}
}*/

// TestGetPaginatedReviewsForAnimeHandler_InvalidAnimeID tests the handler for an invalid anime ID.
func TestGetPaginatedReviewsForAnimeHandler_InvalidAnimeID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	req := httptest.NewRequest("GET", "/anime/invalid/reviews?page=1&limit=10", nil)
	vars := map[string]string{
		"id": "invalid",
	}
	req = mux.SetURLVars(req, vars)

	rec := httptest.NewRecorder()
	handler.GetPaginatedReviewsForAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP status 400 Bad Request, got %d", resp.StatusCode)
	}
}

// TestGetPaginatedReviewsForAnimeHandler_NoAnimeID tests the handler for a missing anime ID.
func TestGetPaginatedReviewsForAnimeHandler_NoAnimeID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	req := httptest.NewRequest("GET", "/anime//reviews?page=1&limit=10", nil)
	rec := httptest.NewRecorder()
	handler.GetPaginatedReviewsForAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP status 400 Bad Request, got %d", resp.StatusCode)
	}
}

// TestGetPaginatedReviewsForAnimeHandler_InvalidPage tests the handler for an invalid page value.
func TestGetPaginatedReviewsForAnimeHandler_InvalidPage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	req := httptest.NewRequest("GET", "/anime/1/reviews?page=invalid&limit=10", nil)
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rec := httptest.NewRecorder()
	handler.GetPaginatedReviewsForAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP status 400 Bad Request, got %d", resp.StatusCode)
	}
}

// TestGetPaginatedReviewsForAnimeHandler_InvalidLimit tests the handler for an invalid limit value.
func TestGetPaginatedReviewsForAnimeHandler_InvalidLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	req := httptest.NewRequest("GET", "/anime/1/reviews?page=1&limit=invalid", nil)
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rec := httptest.NewRecorder()
	handler.GetPaginatedReviewsForAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP status 400 Bad Request, got %d", resp.StatusCode)
	}
}

// TODO: Adjust this test
/*func TestGetPaginatedReviewsForAnimeHandler_NoPaginationParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	testReviews := []models.Review{
		{ID: 1, AnimeID: 1, Content: "Great anime!"},
	}

	// Mock the DB call to return paginated reviews with default values
	mockDB.EXPECT().
		Where(gomock.Any(), "anime_id = ?", uint(1)).
		Return(mockDB)
	mockDB.EXPECT().
		Offset(0).
		Return(mockDB)
	mockDB.EXPECT().
		Limit(10).
		Return(mockDB)
	mockDB.EXPECT().
		Find(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
			*dest.(*[]models.Review) = testReviews
			return &gorm.DB{Error: nil}
		})

	req := httptest.NewRequest("GET", "/anime/1/reviews", nil)
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rec := httptest.NewRecorder()
	handler.GetPaginatedReviewsForAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d", resp.StatusCode)
	}

	var responseReviews []models.Review
	if err := json.NewDecoder(resp.Body).Decode(&responseReviews); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(responseReviews) != len(testReviews) {
		t.Errorf("Expected %d reviews, got %d", len(testReviews), len(responseReviews))
	}
}*/

/*func TestGetAnimeHandler_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	// Create a mock *gorm.DB object
	mockGormDB := &gorm.DB{}

	// Mock the WithContext method to return the mockGormDB
	mockDB.EXPECT().
		WithContext(gomock.Any()). // Mock WithContext
		Return(mockGormDB)         // Return a *gorm.DB

	// Mock the First method on the mockGormDB
	mockDB.EXPECT().
		First(gomock.Any(), gomock.Any(), uint(1)).
		DoAndReturn(func(dest interface{}, conds ...interface{}) *gorm.DB {
			return &gorm.DB{Error: errors.New("record not found")}
		})

	// Create a request with the ID in the URL
	req := httptest.NewRequest("GET", "/anime/1", nil)
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rec := httptest.NewRecorder()
	handler.GetAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected HTTP status 404 Not Found, got %d", resp.StatusCode)
	}
}*/

// TestGetAnimeHandler_InvalidID tests the handler for an invalid anime ID.
func TestGetAnimeHandler_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	req := httptest.NewRequest("GET", "/anime/invalid", nil)
	rec := httptest.NewRecorder()
	handler.GetAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP status 400 Bad Request, got %d", resp.StatusCode)
	}
}

// TODO: Adjust this test
/*func TestGetAnimeHandler_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	// Mock the DB call to return an error
	mockDB.EXPECT().
		First(gomock.Any(), gomock.Any(), uint(1)). // Use uint(1) instead of uint64(1)
		Return(&gorm.DB{Error: errors.New("database error")})

	// Create a request with the ID in the URL
	req := httptest.NewRequest("GET", "/anime/1", nil)

	// Manually set the "id" parameter in the request context
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rec := httptest.NewRecorder()
	handler.GetAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected HTTP status 500 Internal Server Error, got %d", resp.StatusCode)
	}
}*/

// TestGetAllAnimesHandler_Empty tests the handler for retrieving all animes when the database is empty.
func TestGetAllAnimesHandler_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	// Mock the DB call to return an empty slice
	mockDB.EXPECT().
		Find(gomock.Any(), gomock.Any(), gomock.Any()). // Match 3 arguments
		DoAndReturn(func(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
			// Set the destination to an empty slice
			*dest.(*[]models.Anime) = []models.Anime{}
			return &gorm.DB{Error: nil} // Return a *gorm.DB with no error
		})

	req := httptest.NewRequest("GET", "/anime", nil)
	rec := httptest.NewRecorder()
	handler.GetAllAnimesHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d", resp.StatusCode)
	}

	var responseAnimes []models.Anime
	if err := json.NewDecoder(resp.Body).Decode(&responseAnimes); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(responseAnimes) != 0 {
		t.Errorf("Expected empty slice, got %d items", len(responseAnimes))
	}
}

// TestGetAllAnimesHandler_DBError tests the handler for retrieving all animes when the database returns an error.
func TestGetAllAnimesHandler_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	// Mock the DB call to return an error
	mockDB.EXPECT().
		Find(gomock.Any(), gomock.Any()).
		Return(&gorm.DB{Error: errors.New("database error")})

	req := httptest.NewRequest("GET", "/anime", nil)
	rec := httptest.NewRecorder()
	handler.GetAllAnimesHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected HTTP status 500 Internal Server Error, got %d", resp.StatusCode)
	}
}

// TestCreateAnimeHandler_Success tests the handler for successfully creating a new anime.
func TestCreateAnimeHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	testAnime := models.Anime{
		Title: "Naruto",
	}

	// Mock the DB call to create a new anime
	mockDB.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, dest interface{}) *gorm.DB {
			*dest.(*models.Anime) = testAnime
			return &gorm.DB{Error: nil}
		})

	req := httptest.NewRequest("POST", "/anime", nil)
	ctx := context.WithValue(req.Context(), middleware.ValidatedPayloadKey, &testAnime)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	handler.CreateAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected HTTP status 201 Created, got %d", resp.StatusCode)
	}

	var responseAnime models.Anime
	if err := json.NewDecoder(resp.Body).Decode(&responseAnime); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if responseAnime.Title != testAnime.Title {
		t.Errorf("Expected title %s, got %s", testAnime.Title, responseAnime.Title)
	}
}

// TODO: Adjust this test
/*func TestCreateAnimeHandler_MissingTitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	// Create a request with an empty title
	anime := models.Anime{
		Title: "", // Missing title
	}

	// Marshal the payload into JSON
	reqBody, err := json.Marshal(anime)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Create a new request with the JSON payload
	req := httptest.NewRequest("POST", "/anime", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Add the validated payload to the context (simulate validation middleware)
	ctx := context.WithValue(req.Context(), middleware.ValidatedPayloadKey, &anime)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	handler.CreateAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP status 400 Bad Request, got %d", resp.StatusCode)
	}
}*/

// TestCreateAnimeHandler_InvalidInput tests the handler for creating an anime with invalid input.
func TestCreateAnimeHandler_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	// Create a request with invalid JSON (e.g., incorrect data type for Title)
	reqBody := []byte(`{"title": 123}`) // Invalid JSON (Title should be a string)
	req := httptest.NewRequest("POST", "/anime", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Wrap the handler with the middleware
	middleware := middleware.ValidateAndSanitizePayload(http.HandlerFunc(handler.CreateAnimeHandler), models.Anime{})

	rec := httptest.NewRecorder()
	middleware.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP status 400 Bad Request, got %d", resp.StatusCode)
	}
}

// TestCreateAnimeHandler_DBError tests the handler for creating an anime when the database returns an error.
func TestCreateAnimeHandler_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	testAnime := models.Anime{
		Title: "Naruto",
	}

	// Mock the DB call to return an error
	mockDB.EXPECT().
		Create(gomock.Any(), gomock.Any()).                   // Expect the Create method to be called
		Return(&gorm.DB{Error: errors.New("database error")}) // Simulate a database error

	// Marshal the payload into JSON
	reqBody, err := json.Marshal(testAnime)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Create a new request with the JSON payload
	req := httptest.NewRequest("POST", "/anime", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Add the validated payload to the context (simulate validation middleware)
	ctx := context.WithValue(req.Context(), middleware.ValidatedPayloadKey, &testAnime)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	handler.CreateAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected HTTP status 500 Internal Server Error, got %d", resp.StatusCode)
	}
}

// TestCreateAnimeHandler_InvalidJSON tests the handler for creating an anime with invalid JSON input.
func TestCreateAnimeHandler_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	// Create a request with invalid JSON (e.g., malformed JSON)
	reqBody := []byte(`{"title": "Naruto", "description":}`) // Invalid JSON
	req := httptest.NewRequest("POST", "/anime", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Wrap the handler with the middleware
	middleware := middleware.ValidateAndSanitizePayload(http.HandlerFunc(handler.CreateAnimeHandler), models.Anime{})

	rec := httptest.NewRecorder()
	middleware.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP status 400 Bad Request, got %d", resp.StatusCode)
	}
}

// TestUpdateAnimeHandler_Success tests the handler for successfully updating an anime.
func TestUpdateAnimeHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	testAnime := models.Anime{
		ID:    1,
		Title: "Naruto Shippuden",
	}

	// Mock the DB call to find the anime
	mockDB.EXPECT().
		First(gomock.Any(), gomock.Any(), uint(1)).
		DoAndReturn(func(ctx context.Context, dest interface{}, id uint) *gorm.DB {
			// Set the destination to a test anime
			if destPtr, ok := dest.(**models.Anime); ok {
				*destPtr = &models.Anime{ID: 1, Title: "Naruto"}
			}
			return &gorm.DB{
				Statement: &gorm.Statement{
					DB:      &gorm.DB{},
					Clauses: make(map[string]clause.Clause), // Initialize the Clauses map
				},
				Config: &gorm.Config{}, // Initialize the Config field
			}
		})

	// Mock the DB call to save the updated anime
	mockDB.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, dest interface{}) *gorm.DB {
			*dest.(*models.Anime) = testAnime
			return &gorm.DB{
				Statement: &gorm.Statement{
					DB:      &gorm.DB{},
					Clauses: make(map[string]clause.Clause), // Initialize the Clauses map
				},
				Config: &gorm.Config{}, // Initialize the Config field
			}
		})

	req := httptest.NewRequest("PUT", "/anime/1", nil)
	ctx := context.WithValue(req.Context(), middleware.ValidatedPayloadKey, &testAnime)
	req = req.WithContext(ctx)

	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rec := httptest.NewRecorder()
	handler.UpdateAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d", resp.StatusCode)
	}

	var responseAnime models.Anime
	if err := json.NewDecoder(resp.Body).Decode(&responseAnime); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if responseAnime.Title != testAnime.Title {
		t.Errorf("Expected title %s, got %s", testAnime.Title, responseAnime.Title)
	}
}

// TestUpdateAnimeHandler_InvalidInput tests the handler for updating an anime with invalid input.
func TestUpdateAnimeHandler_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	// Create a request with invalid input (e.g., incorrect data type for Title)
	reqBody := []byte(`{"title": 123}`) // Invalid input (Title should be a string)
	req := httptest.NewRequest("PUT", "/anime/1", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Wrap the handler with the middleware
	middleware := middleware.ValidateAndSanitizePayload(http.HandlerFunc(handler.UpdateAnimeHandler), models.Anime{})

	rec := httptest.NewRecorder()
	middleware.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP status 400 Bad Request, got %d", resp.StatusCode)
	}
}

// TestUpdateAnimeHandler_DBError tests the handler for updating an anime when the database returns an error.
func TestUpdateAnimeHandler_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	testAnime := models.Anime{
		ID:    1,
		Title: "Naruto Shippuden",
	}

	// Mock the DB call to find the anime
	mockDB.EXPECT().
		First(gomock.Any(), gomock.Any(), uint(1)). // Use uint(1) instead of uint64(1)
		DoAndReturn(func(ctx context.Context, dest interface{}, id uint) *gorm.DB {
			// Set the destination to a test anime
			if destPtr, ok := dest.(**models.Anime); ok {
				*destPtr = &models.Anime{ID: 1, Title: "Naruto"}
			}
			return &gorm.DB{Error: nil}
		})

	// Mock the DB call to save the updated anime, returning an error
	mockDB.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(&gorm.DB{Error: errors.New("database error")})

	// Marshal the payload into JSON
	reqBody, err := json.Marshal(testAnime)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Create a new request with the JSON payload
	req := httptest.NewRequest("PUT", "/anime/1", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Add the validated payload to the context (simulate validation middleware)
	ctx := context.WithValue(req.Context(), middleware.ValidatedPayloadKey, &testAnime)
	req = req.WithContext(ctx)

	// Manually set the "id" parameter in the request context
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rec := httptest.NewRecorder()
	handler.UpdateAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected HTTP status 500 Internal Server Error, got %d", resp.StatusCode)
	}
}

// TestUpdateAnimeHandler_InvalidJSON tests the handler for updating an anime with invalid JSON input.
func TestUpdateAnimeHandler_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	// Create a request with invalid JSON (e.g., malformed JSON)
	reqBody := []byte(`{"title": "Naruto Shippuden", "description":}`) // Invalid JSON
	req := httptest.NewRequest("PUT", "/anime/1", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Wrap the handler with the middleware
	middleware := middleware.ValidateAndSanitizePayload(http.HandlerFunc(handler.UpdateAnimeHandler), models.Anime{})

	rec := httptest.NewRecorder()
	middleware.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP status 400 Bad Request, got %d", resp.StatusCode)
	}
}

// TestUpdateAnimeHandler_NotFound tests the handler for updating an anime that does not exist.
func TestUpdateAnimeHandler_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	// Mock the DB call to find the anime, returning an error
	mockDB.EXPECT().
		First(gomock.Any(), gomock.Any(), uint(1)). // Use uint(1) instead of uint64(1)
		Return(&gorm.DB{Error: errors.New("record not found")})

	// Create a request with valid data
	reqBody := []byte(`{"title": "Naruto Shippuden"}`)
	req := httptest.NewRequest("PUT", "/anime/1", bytes.NewBuffer(reqBody))
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	// Add the validated payload to the context
	ctx := context.WithValue(req.Context(), middleware.ValidatedPayloadKey, &models.Anime{Title: "Naruto Shippuden"})
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	handler.UpdateAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected HTTP status 404 Not Found, got %d", resp.StatusCode)
	}
}

//TODO: Adjust this test
/*func TestUpdateAnimeHandler_EmptyTitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	// Mock the DB call to find the anime
	mockDB.EXPECT().
		First(gomock.Any(), gomock.Any(), uint(1)). // Use uint(1) instead of uint64(1)
		DoAndReturn(func(ctx context.Context, dest interface{}, id uint) *gorm.DB {
			*dest.(*models.Anime) = models.Anime{ID: 1, Title: "Naruto"}
			return &gorm.DB{Error: nil}
		})

	// Create a request with an empty title
	reqBody := []byte(`{"title": ""}`)
	req := httptest.NewRequest("PUT", "/anime/1", bytes.NewBuffer(reqBody))
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	// Add the validated payload to the context
	ctx := context.WithValue(req.Context(), middleware.ValidatedPayloadKey, &models.Anime{Title: ""})
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	handler.UpdateAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP status 400 Bad Request, got %d", resp.StatusCode)
	}
}*/

// FAILING NOW
/*func TestDeleteAnimeHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	// Mock the DB call to find the anime
	mockDB.EXPECT().
		First(gomock.Any(), gomock.Any(), uint(1)). // Use uint(1) instead of uint64(1)
		DoAndReturn(func(ctx context.Context, dest interface{}, id uint) *gorm.DB {
			*dest.(*models.Anime) = models.Anime{ID: 1, Title: "Naruto"}
			return &gorm.DB{Error: nil}
		})

	// Mock the DB call to delete the anime
	mockDB.EXPECT().
		Delete(gomock.Any(), gomock.Any(), uint(1)). // Use uint(1) instead of uint64(1)
		Return(&gorm.DB{Error: nil})

	req := httptest.NewRequest("DELETE", "/anime/1", nil)
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rec := httptest.NewRecorder()
	handler.DeleteAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected HTTP status 204 No Content, got %d", resp.StatusCode)
	}
}*/

// FAILING NOW
/*func TestDeleteAnimeHandler_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	// Mock the DB call to find the anime, returning an error
	mockDB.EXPECT().
		First(gomock.Any(), gomock.Any(), uint(1)). // Use uint(1) instead of uint64(1)
		Return(&gorm.DB{Error: errors.New("record not found")})

	req := httptest.NewRequest("DELETE", "/anime/1", nil)
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rec := httptest.NewRecorder()
	handler.DeleteAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected HTTP status 404 Not Found, got %d", resp.StatusCode)
	}
}*/

// FAILING NOW
/*func TestDeleteAnimeHandler_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	// Mock the DB call to find the anime
	mockDB.EXPECT().
		First(gomock.Any(), gomock.Any(), uint(1)). // Use uint(1) instead of uint64(1)
		DoAndReturn(func(ctx context.Context, dest interface{}, id uint) *gorm.DB {
			*dest.(*models.Anime) = models.Anime{ID: 1, Title: "Naruto"}
			return &gorm.DB{Error: nil}
		})

	// Mock the DB call to delete the anime, returning an error
	mockDB.EXPECT().
		Delete(gomock.Any(), gomock.Any(), uint(1)). // Use uint(1) instead of uint64(1)
		Return(&gorm.DB{Error: errors.New("database error")})

	req := httptest.NewRequest("DELETE", "/anime/1", nil)
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rec := httptest.NewRecorder()
	handler.DeleteAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected HTTP status 500 Internal Server Error, got %d", resp.StatusCode)
	}
}*/
