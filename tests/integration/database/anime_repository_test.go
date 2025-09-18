package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"myanimeapi/api/models"
	modelsfixtures "myanimeapi/tests/fixtures/models"
	"myanimeapi/tests/suites"
)

func TestAnimeRepository_Integration(t *testing.T) {
	suite := suites.NewBaseSuite(t)

	t.Run("CreateAnime", func(t *testing.T) {
		t.Parallel()

		// Arrange
		fixture := modelsfixtures.TestAnime()

		// Act
		err := suite.DB.Create(fixture).Error

		// Assert
		require.NoError(t, err)
		assert.NotEqual(t, uint(0), fixture.ID)

		// Cleanup
		suite.DB.Delete(fixture)
	})

	t.Run("GetAnimeByID", func(t *testing.T) {
		t.Parallel()

		// Arrange
		fixture := modelsfixtures.TestAnime()
		suite.DB.Create(fixture)

		// Act
		var foundAnime models.Anime
		err := suite.DB.First(&foundAnime, fixture.ID).Error

		// Assert
		require.NoError(t, err)
		assert.Equal(t, fixture.Title, foundAnime.Title)
		assert.Equal(t, fixture.Description, foundAnime.Description)
		assert.Equal(t, fixture.Status, foundAnime.Status)

		// Cleanup
		suite.DB.Delete(fixture)
	})

	t.Run("UpdateAnime", func(t *testing.T) {
		t.Parallel()

		// Arrange
		fixture := modelsfixtures.TestAnime()
		suite.DB.Create(fixture)

		// Act
		updateData := map[string]interface{}{
			"title":       "Updated Anime Title",
			"description": "Updated description",
			"rating":      9.0,
		}
		err := suite.DB.Model(&models.Anime{}).Where("id = ?", fixture.ID).Updates(updateData).Error

		// Assert
		require.NoError(t, err)

		var updatedAnime models.Anime
		err = suite.DB.First(&updatedAnime, fixture.ID).Error
		require.NoError(t, err)
		assert.Equal(t, "Updated Anime Title", updatedAnime.Title)
		assert.Equal(t, "Updated description", updatedAnime.Description)
		assert.Equal(t, 9.0, updatedAnime.Rating)

		// Cleanup
		suite.DB.Delete(fixture)
	})

	t.Run("DeleteAnime", func(t *testing.T) {
		t.Parallel()

		// Arrange
		fixture := modelsfixtures.TestAnime()
		suite.DB.Create(fixture)

		// Act
		err := suite.DB.Delete(fixture).Error

		// Assert
		require.NoError(t, err)

		// Verify deletion
		var foundAnime models.Anime
		err = suite.DB.First(&foundAnime, fixture.ID).Error
		assert.Error(t, err) // Should not find the deleted anime
	})

	t.Run("ListAnimes", func(t *testing.T) {
		t.Parallel()

		// Arrange
		// Create multiple animes
		anime1 := modelsfixtures.TestAnime()
		anime2 := modelsfixtures.OngoingAnime()
		anime3 := modelsfixtures.CompletedAnime()

		suite.DB.Create(anime1)
		suite.DB.Create(anime2)
		suite.DB.Create(anime3)

		// Act
		var animes []models.Anime
		err := suite.DB.Find(&animes).Error

		// Assert
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(animes), 3)

		// Cleanup
		suite.DB.Delete(anime1)
		suite.DB.Delete(anime2)
		suite.DB.Delete(anime3)
	})
}
