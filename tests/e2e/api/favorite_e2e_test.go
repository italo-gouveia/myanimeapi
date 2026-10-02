package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"myanimeapi/api/models"
	"myanimeapi/tests/suites"
)

// TestFavoriteFlow_E2E verifies the complete favorite flow:
// register a user → create an anime → add to favorites → list favorites → remove from favorites.
func TestFavoriteFlow_E2E(t *testing.T) {
	suite := suites.NewBaseSuite(t)

	if os.Getenv("JWT_SECRET_KEY") == "" {
		_ = os.Setenv("JWT_SECRET_KEY", "test-secret")
	}

	username := "e2euser_favorite"
	password := "E2Epassword123!"
	email := username + "@example.com"
	var authHeader string

	// -----------------------------------------------------------------------
	// Step 0: Register a user and login to get a real JWT
	// -----------------------------------------------------------------------
	t.Run("RegisterUser", func(t *testing.T) {
		body := map[string]string{
			"username": username,
			"email":    email,
			"password": password,
		}
		b, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		suite.Router.ServeHTTP(rec, req)

		if rec.Code == http.StatusConflict {
			t.Log("User already exists, proceeding to login")
		} else {
			require.Equal(t, http.StatusCreated, rec.Code, "user registration must succeed")
		}
	})

	t.Run("LoginUser", func(t *testing.T) {
		body := map[string]string{
			"username": username,
			"password": password,
		}
		b, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/v1/auth/authenticate", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		suite.Router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code, "user login must succeed")

		var resp map[string]interface{}
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		data, ok := resp["data"].(map[string]interface{})
		require.True(t, ok, "expected data in login response")
		token, _ := data["token"].(string)
		require.NotEmpty(t, token, "expected JWT token")
		authHeader = "Bearer " + token
	})

	// -----------------------------------------------------------------------
	// Step 1: Seed an anime directly via DB (POST /animes requires admin).
	// -----------------------------------------------------------------------
	var animeID float64

	t.Run("CreateAnimeForFavorite", func(t *testing.T) {
		anime := models.Anime{
			Title:       "Favorite Flow Anime",
			Description: "Used for the favorite E2E flow test",
			Status:      "Ongoing",
			Episodes:    12,
			Rating:      8.0,
			StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		require.NoError(t, suite.DB.Create(&anime).Error, "anime seed must succeed")
		animeID = float64(anime.ID)
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
