package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"myanimeapi/api/models"
	"myanimeapi/api/repositories"
	modelsfixtures "myanimeapi/tests/fixtures/models"
	"myanimeapi/tests/suites"
	appdb "myanimeapi/internal/db"
)

func TestAnimeRepository_GetWithFilters(t *testing.T) {
	suite := suites.NewBaseSuite(t)
	repo := repositories.NewAnimeRepository(appdb.NewGormDB(suite.DB))
	ctx := context.Background()

	t.Run("FilterByStatuses", func(t *testing.T) {
		airing := modelsfixtures.NewAnimeFixture(
			modelsfixtures.WithTitle("Statuses Filter Test Airing"),
			modelsfixtures.WithStatus("Airing"),
		)
		completed := modelsfixtures.NewAnimeFixture(
			modelsfixtures.WithTitle("Statuses Filter Test Completed"),
			modelsfixtures.WithStatus("Completed"),
		)
		upcoming := modelsfixtures.NewAnimeFixture(
			modelsfixtures.WithTitle("Statuses Filter Test Upcoming"),
			modelsfixtures.WithStatus("Upcoming"),
		)
		require.NoError(t, suite.DB.Create(airing).Error)
		require.NoError(t, suite.DB.Create(completed).Error)
		require.NoError(t, suite.DB.Create(upcoming).Error)
		t.Cleanup(func() {
			suite.DB.Delete(airing)
			suite.DB.Delete(completed)
			suite.DB.Delete(upcoming)
		})

		results, _, err := repo.GetWithFilters(ctx, models.AnimeFilter{
			Statuses: []string{"Airing", "Completed"},
		}, 1, 200)
		require.NoError(t, err)

		seen := map[string]bool{}
		for _, r := range results {
			a := r.(*models.Anime)
			seen[a.Title] = true
			assert.Contains(t, []string{"Airing", "Completed"}, a.Status)
		}
		assert.True(t, seen["Statuses Filter Test Airing"])
		assert.True(t, seen["Statuses Filter Test Completed"])
		assert.False(t, seen["Statuses Filter Test Upcoming"])
	})

	t.Run("FilterByYearRange", func(t *testing.T) {
		anime2020 := modelsfixtures.NewAnimeFixture(
			modelsfixtures.WithTitle("YearRange Filter Test 2020"),
			modelsfixtures.WithStartDate(time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC)),
		)
		anime2022 := modelsfixtures.NewAnimeFixture(
			modelsfixtures.WithTitle("YearRange Filter Test 2022"),
			modelsfixtures.WithStartDate(time.Date(2022, 3, 15, 0, 0, 0, 0, time.UTC)),
		)
		anime2024 := modelsfixtures.NewAnimeFixture(
			modelsfixtures.WithTitle("YearRange Filter Test 2024"),
			modelsfixtures.WithStartDate(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		)
		require.NoError(t, suite.DB.Create(anime2020).Error)
		require.NoError(t, suite.DB.Create(anime2022).Error)
		require.NoError(t, suite.DB.Create(anime2024).Error)
		t.Cleanup(func() {
			suite.DB.Delete(anime2020)
			suite.DB.Delete(anime2022)
			suite.DB.Delete(anime2024)
		})

		results, _, err := repo.GetWithFilters(ctx, models.AnimeFilter{
			YearFrom: 2020,
			YearTo:   2022,
		}, 1, 200)
		require.NoError(t, err)

		seen := map[string]bool{}
		for _, r := range results {
			a := r.(*models.Anime)
			seen[a.Title] = true
		}
		assert.True(t, seen["YearRange Filter Test 2020"])
		assert.True(t, seen["YearRange Filter Test 2022"])
		assert.False(t, seen["YearRange Filter Test 2024"])
	})

	t.Run("FilterByYearFrom", func(t *testing.T) {
		anime2018 := modelsfixtures.NewAnimeFixture(
			modelsfixtures.WithTitle("YearFrom Filter Test 2018"),
			modelsfixtures.WithStartDate(time.Date(2018, 12, 31, 0, 0, 0, 0, time.UTC)),
		)
		anime2021 := modelsfixtures.NewAnimeFixture(
			modelsfixtures.WithTitle("YearFrom Filter Test 2021"),
			modelsfixtures.WithStartDate(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
		)
		require.NoError(t, suite.DB.Create(anime2018).Error)
		require.NoError(t, suite.DB.Create(anime2021).Error)
		t.Cleanup(func() {
			suite.DB.Delete(anime2018)
			suite.DB.Delete(anime2021)
		})

		results, _, err := repo.GetWithFilters(ctx, models.AnimeFilter{
			YearFrom: 2019,
		}, 1, 200)
		require.NoError(t, err)

		seen := map[string]bool{}
		for _, r := range results {
			a := r.(*models.Anime)
			seen[a.Title] = true
		}
		assert.False(t, seen["YearFrom Filter Test 2018"])
		assert.True(t, seen["YearFrom Filter Test 2021"])
	})
}
