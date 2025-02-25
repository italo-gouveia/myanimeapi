package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"myanimeapi/internal/mocks"
	"myanimeapi/pkg/handlers"
	"myanimeapi/pkg/models"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
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
		DoAndReturn(func(dest interface{}, query string, id uint64) error {
			*dest.(*models.Anime) = testAnime
			return nil
		})

	req := httptest.NewRequest("GET", "/anime/1", nil)
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

func TestGetAnimeHandler_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBInterface(ctrl)
	handler := handlers.AnimeHandler{DB: mockDB}

	// Mock DB returning an error
	mockDB.EXPECT().
		First(gomock.Any(), gomock.Any(), uint64(1)).
		Return(errors.New("record not found"))

	req := httptest.NewRequest("GET", "/anime/1", nil)
	rec := httptest.NewRecorder()
	handler.GetAnimeHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected HTTP status 404 Not Found, got %d", resp.StatusCode)
	}
}
