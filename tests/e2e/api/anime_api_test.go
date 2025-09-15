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

func TestAnimeAPI_E2E(t *testing.T) {
	suite := suites.NewBaseSuite(t)

	// Ensure JWT secret and obtain a token for protected routes
	if os.Getenv("JWT_SECRET_KEY") == "" {
		os.Setenv("JWT_SECRET_KEY", "test-secret")
	}
	token, err := middleware.GenerateToken("1", true)
	require.NoError(t, err)
	authHeader := "Bearer " + token

	t.Run("CreateAnime", func(t *testing.T) {
		// Arrange
		reqBody := map[string]interface{}{
			"title":       "Test Anime",
			"description": "A test anime for testing purposes",
			"start_date":  "2024-01-01T00:00:00Z",
			"status":      "Ongoing",
			"episodes":    12,
			"rating":      8.5,
		}

		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/v1/animes", bytes.NewBuffer(reqJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader)

		// Act
		rec := httptest.NewRecorder()
		suite.Router.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusCreated, rec.Code)

		// Parse response
		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Test Anime", response["title"])
		assert.Equal(t, "A test anime for testing purposes", response["description"])
		assert.Equal(t, "Ongoing", response["status"])
		assert.Equal(t, float64(12), response["episodes"])
		assert.Equal(t, 8.5, response["rating"])
	})

	t.Run("GetAnimeByID", func(t *testing.T) {
		// Arrange
		// First create an anime
		createReqBody := map[string]interface{}{
			"title":       "Get Test Anime",
			"description": "Anime for get test",
			"status":      "Ongoing",
			"episodes":    24,
			"rating":      9.0,
		}

		createReqJSON, _ := json.Marshal(createReqBody)
		createReq := httptest.NewRequest("POST", "/v1/animes", bytes.NewBuffer(createReqJSON))
		createReq.Header.Set("Content-Type", "application/json")
		createReq.Header.Set("Authorization", authHeader)

		createRec := httptest.NewRecorder()
		suite.Router.ServeHTTP(createRec, createReq)

		require.Equal(t, http.StatusCreated, createRec.Code)

		var createResponse map[string]interface{}
		err := json.Unmarshal(createRec.Body.Bytes(), &createResponse)
		require.NoError(t, err)

		animeID := createResponse["id"].(float64)

		// Act - Get the created anime
		getReq := httptest.NewRequest("GET", "/v1/animes/"+fmt.Sprintf("%.0f", animeID), nil)
		getRec := httptest.NewRecorder()
		suite.Router.ServeHTTP(getRec, getReq)

		// Assert
		assert.Equal(t, http.StatusOK, getRec.Code)

		var getResponse map[string]interface{}
		err = json.Unmarshal(getRec.Body.Bytes(), &getResponse)
		require.NoError(t, err)

		assert.Equal(t, "Get Test Anime", getResponse["title"])
		assert.Equal(t, "Anime for get test", getResponse["description"])
		assert.Equal(t, "Ongoing", getResponse["status"])
		assert.Equal(t, float64(24), getResponse["episodes"])
		assert.Equal(t, 9.0, getResponse["rating"])
	})

	t.Run("UpdateAnime", func(t *testing.T) {
		// Arrange
		// First create an anime
		createReqBody := map[string]interface{}{
			"title":       "Update Test Anime",
			"description": "Anime for update test",
			"status":      "Ongoing",
			"episodes":    12,
			"rating":      7.5,
		}

		createReqJSON, _ := json.Marshal(createReqBody)
		createReq := httptest.NewRequest("POST", "/v1/animes", bytes.NewBuffer(createReqJSON))
		createReq.Header.Set("Content-Type", "application/json")
		createReq.Header.Set("Authorization", authHeader)

		createRec := httptest.NewRecorder()
		suite.Router.ServeHTTP(createRec, createReq)

		require.Equal(t, http.StatusCreated, createRec.Code)

		var createResponse map[string]interface{}
		err := json.Unmarshal(createRec.Body.Bytes(), &createResponse)
		require.NoError(t, err)

		animeID := createResponse["id"].(float64)

		// Act - Update the anime
		updateReqBody := map[string]interface{}{
			"title":       "Updated Anime Title",
			"description": "Updated description",
			"rating":      8.5,
		}

		updateReqJSON, _ := json.Marshal(updateReqBody)
		updateReq := httptest.NewRequest("PUT", "/v1/animes/"+fmt.Sprintf("%.0f", animeID), bytes.NewBuffer(updateReqJSON))
		updateReq.Header.Set("Content-Type", "application/json")
		updateReq.Header.Set("Authorization", authHeader)

		updateRec := httptest.NewRecorder()
		suite.Router.ServeHTTP(updateRec, updateReq)

		// Assert
		assert.Equal(t, http.StatusOK, updateRec.Code)

		var updateResponse map[string]interface{}
		err = json.Unmarshal(updateRec.Body.Bytes(), &updateResponse)
		require.NoError(t, err)

		assert.Equal(t, "Updated Anime Title", updateResponse["title"])
		assert.Equal(t, "Updated description", updateResponse["description"])
		assert.Equal(t, 8.5, updateResponse["rating"])
	})

	t.Run("DeleteAnime", func(t *testing.T) {
		// Arrange
		// First create an anime
		createReqBody := map[string]interface{}{
			"title":       "Delete Test Anime",
			"description": "Anime for delete test",
			"status":      "Ongoing",
			"episodes":    12,
			"rating":      7.0,
		}

		createReqJSON, _ := json.Marshal(createReqBody)
		createReq := httptest.NewRequest("POST", "/v1/animes", bytes.NewBuffer(createReqJSON))
		createReq.Header.Set("Content-Type", "application/json")
		createReq.Header.Set("Authorization", authHeader)

		createRec := httptest.NewRecorder()
		suite.Router.ServeHTTP(createRec, createReq)

		require.Equal(t, http.StatusCreated, createRec.Code)

		var createResponse map[string]interface{}
		err := json.Unmarshal(createRec.Body.Bytes(), &createResponse)
		require.NoError(t, err)

		animeID := createResponse["id"].(float64)

		// Act - Delete the anime
		deleteReq := httptest.NewRequest("DELETE", "/v1/animes/"+fmt.Sprintf("%.0f", animeID), nil)
		deleteReq.Header.Set("Authorization", authHeader)
		deleteRec := httptest.NewRecorder()
		suite.Router.ServeHTTP(deleteRec, deleteReq)

		// Assert
		assert.Equal(t, http.StatusNoContent, deleteRec.Code)

		// Verify deletion
		getReq := httptest.NewRequest("GET", "/v1/animes/"+fmt.Sprintf("%.0f", animeID), nil)
		getRec := httptest.NewRecorder()
		suite.Router.ServeHTTP(getRec, getReq)

		assert.Equal(t, http.StatusNotFound, getRec.Code)
	})
}
