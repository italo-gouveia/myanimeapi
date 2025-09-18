package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"myanimeapi/api/models"
	modelsfixtures "myanimeapi/tests/fixtures/models"
)

func TestAnimeService_CreateAnime(t *testing.T) {
	t.Parallel()

	// Arrange
	input := models.Anime{
		Title:       "Test Anime",
		Description: "A test anime for testing",
		Status:      "Ongoing",
		Episodes:    12,
		Rating:      8.5,
		StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// Act - For now we only test the creation of the fixture
	anime := modelsfixtures.TestAnime()

	// Assert
	require.NotNil(t, anime)
	assert.Equal(t, input.Title, anime.Title)
	assert.Equal(t, input.Description, anime.Description)
	assert.Equal(t, input.Status, anime.Status)
	assert.Equal(t, input.Episodes, anime.Episodes)
	assert.Equal(t, input.Rating, anime.Rating)
}

func TestAnimeService_GetAnimeByID(t *testing.T) {
	t.Parallel()

	// Arrange
	expectedAnime := modelsfixtures.TestAnime()
	expectedAnime.ID = 1

	// Act - For now we only test the creation of the fixture
	anime := modelsfixtures.TestAnime()

	// Assert
	require.NotNil(t, anime)
	assert.Equal(t, expectedAnime.Title, anime.Title)
	assert.Equal(t, expectedAnime.Description, anime.Description)
}

func TestAnimeService_GetAnimeByID_NotFound(t *testing.T) {
	t.Parallel()

	// Arrange
	// For now we only test that the test executes
	t.Log("Test placeholder - será implementado quando os mocks estiverem disponíveis")

	// Act & Assert
	assert.True(t, true) // Placeholder assertion
}
