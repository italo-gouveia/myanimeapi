// pkg/handlers_test/anime_handlers_test.go
package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"myanimeapi/internal/mocks"
	"myanimeapi/pkg/handlers"
	"myanimeapi/pkg/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func TestGetAnimeHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	testAnime := models.Anime{
		ID:    1,
		Title: "Naruto",
		Reviews: []models.Review{
			{ID: 1, AnimeID: 1, Content: "Great anime!"},
		},
	}

	// Mock the DB call
	mockDB.EXPECT().
		First(gomock.Any(), gomock.Any(), uint64(1)).
		DoAndReturn(func(ctx context.Context, dest interface{}, id uint64) *gorm.DB {
			*dest.(*models.Anime) = testAnime // Set the destination to the test anime
			return &gorm.DB{Error: nil}       // Return a *gorm.DB with no error
		})

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

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d", resp.StatusCode)
	}

	var responseAnime models.Anime
	if err := json.NewDecoder(resp.Body).Decode(&responseAnime); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Use cmp.Diff for deep comparison
	if diff := cmp.Diff(testAnime, responseAnime); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func TestGetAllAnimesHandler_Success(t *testing.T) {
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
}

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
			*dest.(*[]models.Review) = testReviews
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

	if len(responseReviews) != len(testReviews) {
		t.Errorf("Expected %d reviews, got %d", len(testReviews), len(responseReviews))
	}
}*/

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

func TestGetAnimeHandler_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	// Mock DB returning an error
	mockDB.EXPECT().
		First(gomock.Any(), gomock.Any(), uint64(1)).
		DoAndReturn(func(ctx context.Context, dest interface{}, id uint64) *gorm.DB {
			return &gorm.DB{Error: errors.New("record not found")} // Return a *gorm.DB with an error
		})

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

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected HTTP status 404 Not Found, got %d", resp.StatusCode)
	}
}

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

	reqBody, _ := json.Marshal(testAnime)
	req := httptest.NewRequest("POST", "/anime", bytes.NewBuffer(reqBody))
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

// TESTING GOTTING SUCCESS
func TestCreateAnimeHandler_MissingTitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	anime := models.Anime{
		ID:    1,
		Title: "", // Missing title
	}

	reqBody, _ := json.Marshal(anime)
	req := httptest.NewRequest("POST", "/anime", bytes.NewBuffer(reqBody))
	rec := httptest.NewRecorder()
	handler.CreateAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP status 400 Bad Request, got %d", resp.StatusCode)
	}
}
