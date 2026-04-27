package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"myanimeapi/api/middleware"
	"myanimeapi/tests/suites"
)

// TestFavoriteFlow_E2E verifies the complete favorite flow:
// create an anime → add to favorites → list favorites → remove from favorites.
func TestFavoriteFlow_E2E(t *testing.T) {
	suite := suites.NewBaseSuite(t)

	if os.Getenv("JWT_SECRET_KEY") == "" {
		_ = os.Setenv("JWT_SECRET_KEY", "test-secret")
	}

	token, err := middleware.GenerateToken("1", true)
	require.NoError(t, err)
	authHeader := "Bearer " + token

	// -----------------------------------------------------------------------
	// Step 1: Create an anime to use as the favorite target
	// -----------------------------------------------------------------------
	var animeID float64

	t.Run("CreateAnimeForFavorite", func(t *testing.T) {
		body := map[string]interface{}{
			"title":       "Favorite Flow Anime",
			"description": "Used for the favorite E2E flow test",
			"status":      "Ongoing",
			"episodes":    12,
			"rating":      8.0,
			"start_date":  "2024-01-01T00:00:00Z",
		}
		b, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/v1/animes", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader)
		rec := httptest.NewRecorder()
		suite.Router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusCreated, rec.Code, "anime creation must succeed")

		var resp map[string]interface{}
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		animeID = resp["id"].(float64)
		assert.NotEqual(t, float64(0), animeID)
	})

	// -----------------------------------------------------------------------
	// Step 2: Add anime to favorites
	// -----------------------------------------------------------------------
	t.Run("AddFavorite", func(t *testing.T) {
		path := fmt.Sprintf("/v1/favorites/%.0f", animeID)
		req := httptest.NewRequest(http.MethodPost, path, nil)
		req.Header.Set("Authorization", authHeader)
		rec := httptest.NewRecorder()
		suite.Router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		assert.Equal(t, animeID, resp["anime_id"])
	})

	// -----------------------------------------------------------------------
	// Step 3: List favorites — the anime should appear
	// -----------------------------------------------------------------------
	t.Run("GetFavorites", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/favorites", nil)
		req.Header.Set("Authorization", authHeader)
		rec := httptest.NewRecorder()
		suite.Router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var favorites []interface{}
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&favorites))

		// Verify our anime is in the favorites list
		found := false
		for _, fav := range favorites {
			if favMap, ok := fav.(map[string]interface{}); ok {
				if favMap["anime_id"] == animeID {
					found = true
					break
				}
			}
		}
		assert.True(t, found, "created anime should appear in favorites list")
	})

	// -----------------------------------------------------------------------
	// Step 4: Remove favorite
	// -----------------------------------------------------------------------
	t.Run("RemoveFavorite", func(t *testing.T) {
		path := fmt.Sprintf("/v1/favorites/%.0f", animeID)
		req := httptest.NewRequest(http.MethodDelete, path, nil)
		req.Header.Set("Authorization", authHeader)
		rec := httptest.NewRecorder()
		suite.Router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	// -----------------------------------------------------------------------
	// Step 5: List favorites again — the anime should no longer appear
	// -----------------------------------------------------------------------
	t.Run("GetFavoritesAfterRemoval", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/favorites", nil)
		req.Header.Set("Authorization", authHeader)
		rec := httptest.NewRecorder()
		suite.Router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var favorites []interface{}
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&favorites))

		for _, fav := range favorites {
			if favMap, ok := fav.(map[string]interface{}); ok {
				assert.NotEqual(t, animeID, favMap["anime_id"],
					"removed anime should not appear in favorites list")
			}
		}
	})

	// -----------------------------------------------------------------------
	// Step 6: Verify unauthorized access is rejected
	// -----------------------------------------------------------------------
	t.Run("GetFavorites_Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/favorites", nil)
		rec := httptest.NewRecorder()
		suite.Router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
