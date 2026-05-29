// api/repositories/anime_repository_filter_test.go
// Unit tests for AnimeRepositoryImpl.GetWithFilters using an in-memory SQLite database.
// These tests exercise the filter/sort/pagination logic without a live Postgres instance.
package repositories_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"myanimeapi/api/models"
	"myanimeapi/api/repositories"
	appdb "myanimeapi/internal/db"
)

// newTestRepo spins up an isolated in-memory SQLite database per test, runs
// auto-migrations, and returns a ready-to-use AnimeRepository.
// MaxOpenConns(1) ensures all GORM pool connections share the same :memory: DB.
func newTestRepo(t *testing.T) (repositories.AnimeRepository, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// Prevent GORM from spawning multiple connections to different :memory: DBs.
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)

	err = db.AutoMigrate(
		&models.Anime{},
		&models.Genre{},
		&models.Tag{},
		&models.Review{},
	)
	require.NoError(t, err)

	repo := repositories.NewAnimeRepository(appdb.NewGormDB(db))
	return repo, db
}

// seed inserts an anime (with optional genres/tags) into the test database.
func seed(t *testing.T, db *gorm.DB, anime *models.Anime) *models.Anime {
	t.Helper()
	require.NoError(t, db.Create(anime).Error)
	return anime
}

// ---- Tests ------------------------------------------------------------------

func TestGetWithFilters_NoFilter_ReturnsAll(t *testing.T) {
	repo, db := newTestRepo(t)
	ctx := context.Background()

	seed(t, db, &models.Anime{Title: "Naruto", Status: "Completed", Rating: 8.0})
	seed(t, db, &models.Anime{Title: "One Piece", Status: "Airing", Rating: 9.0})

	result, total, err := repo.GetWithFilters(ctx, models.AnimeFilter{}, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
}

func TestGetWithFilters_FilterByStatus(t *testing.T) {
	repo, db := newTestRepo(t)
	ctx := context.Background()

	seed(t, db, &models.Anime{Title: "Naruto", Status: "Completed"})
	seed(t, db, &models.Anime{Title: "Bleach", Status: "Completed"})
	seed(t, db, &models.Anime{Title: "One Piece", Status: "Airing"})

	result, total, err := repo.GetWithFilters(ctx, models.AnimeFilter{Status: "Completed"}, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
	for _, r := range result {
		assert.Equal(t, "Completed", r.(*models.Anime).Status)
	}
}

func TestGetWithFilters_FilterByTitleLike(t *testing.T) {
	repo, db := newTestRepo(t)
	ctx := context.Background()

	seed(t, db, &models.Anime{Title: "Fullmetal Alchemist"})
	seed(t, db, &models.Anime{Title: "Fullmetal Alchemist: Brotherhood"})
	seed(t, db, &models.Anime{Title: "Naruto"})

	result, total, err := repo.GetWithFilters(ctx, models.AnimeFilter{Title: "Fullmetal"}, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
}

func TestGetWithFilters_FilterByRatingRange(t *testing.T) {
	repo, db := newTestRepo(t)
	ctx := context.Background()

	seed(t, db, &models.Anime{Title: "Low", Rating: 4.0})
	seed(t, db, &models.Anime{Title: "Mid", Rating: 7.5})
	seed(t, db, &models.Anime{Title: "High", Rating: 9.5})

	result, total, err := repo.GetWithFilters(ctx, models.AnimeFilter{RatingMin: 7.0, RatingMax: 9.0}, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, result, 1)
	assert.Equal(t, "Mid", result[0].(*models.Anime).Title)
}

func TestGetWithFilters_FilterByGenre(t *testing.T) {
	repo, db := newTestRepo(t)
	ctx := context.Background()

	action := models.Genre{Name: "Action"}
	romance := models.Genre{Name: "Romance"}
	require.NoError(t, db.Create(&action).Error)
	require.NoError(t, db.Create(&romance).Error)

	a1 := &models.Anime{Title: "Naruto", Genres: []models.Genre{action}}
	a2 := &models.Anime{Title: "Clannad", Genres: []models.Genre{romance}}
	a3 := &models.Anime{Title: "Bleach", Genres: []models.Genre{action}}
	seed(t, db, a1)
	seed(t, db, a2)
	seed(t, db, a3)

	result, total, err := repo.GetWithFilters(ctx, models.AnimeFilter{Genre: "Action"}, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
}

func TestGetWithFilters_FilterByTag(t *testing.T) {
	repo, db := newTestRepo(t)
	ctx := context.Background()

	ninja := models.Tag{Name: "Ninja"}
	mech := models.Tag{Name: "Mech"}
	require.NoError(t, db.Create(&ninja).Error)
	require.NoError(t, db.Create(&mech).Error)

	seed(t, db, &models.Anime{Title: "Naruto", Tags: []models.Tag{ninja}})
	seed(t, db, &models.Anime{Title: "Evangelion", Tags: []models.Tag{mech}})

	result, total, err := repo.GetWithFilters(ctx, models.AnimeFilter{Tag: "Ninja"}, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, result, 1)
	assert.Equal(t, "Naruto", result[0].(*models.Anime).Title)
}

func TestGetWithFilters_SortByRatingDesc(t *testing.T) {
	repo, db := newTestRepo(t)
	ctx := context.Background()

	seed(t, db, &models.Anime{Title: "B", Rating: 7.0})
	seed(t, db, &models.Anime{Title: "A", Rating: 9.5})
	seed(t, db, &models.Anime{Title: "C", Rating: 5.0})

	result, _, err := repo.GetWithFilters(ctx, models.AnimeFilter{SortBy: "rating", SortOrder: "desc"}, 1, 10)
	require.NoError(t, err)
	require.Len(t, result, 3)
	assert.Equal(t, "A", result[0].(*models.Anime).Title)
	assert.Equal(t, "B", result[1].(*models.Anime).Title)
	assert.Equal(t, "C", result[2].(*models.Anime).Title)
}

func TestGetWithFilters_SortByTitleAsc(t *testing.T) {
	repo, db := newTestRepo(t)
	ctx := context.Background()

	seed(t, db, &models.Anime{Title: "Zoro"})
	seed(t, db, &models.Anime{Title: "Araragi"})
	seed(t, db, &models.Anime{Title: "Monkey"})

	result, _, err := repo.GetWithFilters(ctx, models.AnimeFilter{SortBy: "title", SortOrder: "asc"}, 1, 10)
	require.NoError(t, err)
	require.Len(t, result, 3)
	assert.Equal(t, "Araragi", result[0].(*models.Anime).Title)
	assert.Equal(t, "Monkey", result[1].(*models.Anime).Title)
	assert.Equal(t, "Zoro", result[2].(*models.Anime).Title)
}

func TestGetWithFilters_Pagination(t *testing.T) {
	repo, db := newTestRepo(t)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		seed(t, db, &models.Anime{Title: "Anime"})
	}

	// Page 1 — 2 items
	result, total, err := repo.GetWithFilters(ctx, models.AnimeFilter{}, 1, 2)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, result, 2)

	// Page 3 — 1 item
	result, total, err = repo.GetWithFilters(ctx, models.AnimeFilter{}, 3, 2)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, result, 1)
}

func TestGetWithFilters_EmptyResult(t *testing.T) {
	repo, db := newTestRepo(t)
	ctx := context.Background()

	seed(t, db, &models.Anime{Title: "Naruto", Status: "Airing"})

	result, total, err := repo.GetWithFilters(ctx, models.AnimeFilter{Status: "Upcoming"}, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Len(t, result, 0)
}

func TestGetWithFilters_CombinedFilters(t *testing.T) {
	repo, db := newTestRepo(t)
	ctx := context.Background()

	action := models.Genre{Name: "Action"}
	require.NoError(t, db.Create(&action).Error)

	seed(t, db, &models.Anime{Title: "Naruto", Status: "Completed", Rating: 8.5, Genres: []models.Genre{action}})
	seed(t, db, &models.Anime{Title: "Bleach", Status: "Completed", Rating: 7.0, Genres: []models.Genre{action}})
	seed(t, db, &models.Anime{Title: "Clannad", Status: "Completed", Rating: 9.0})

	// Status=Completed + Genre=Action + RatingMin=8.0
	result, total, err := repo.GetWithFilters(ctx, models.AnimeFilter{
		Status:    "Completed",
		Genre:     "Action",
		RatingMin: 8.0,
	}, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, result, 1)
	assert.Equal(t, "Naruto", result[0].(*models.Anime).Title)
}

// ---- Advanced search tests --------------------------------------------------

func TestGetWithFilters_MultipleGenres_ANDSemantics(t *testing.T) {
	repo, db := newTestRepo(t)
	ctx := context.Background()

	action := models.Genre{Name: "Action"}
	adventure := models.Genre{Name: "Adventure"}
	romance := models.Genre{Name: "Romance"}
	require.NoError(t, db.Create(&action).Error)
	require.NoError(t, db.Create(&adventure).Error)
	require.NoError(t, db.Create(&romance).Error)

	// "Naruto" has both Action and Adventure.
	seed(t, db, &models.Anime{Title: "Naruto", Genres: []models.Genre{action, adventure}})
	// "Bleach" has only Action.
	seed(t, db, &models.Anime{Title: "Bleach", Genres: []models.Genre{action}})
	// "Clannad" has only Romance.
	seed(t, db, &models.Anime{Title: "Clannad", Genres: []models.Genre{romance}})

	// Filter by Action AND Adventure — only Naruto qualifies.
	result, total, err := repo.GetWithFilters(ctx, models.AnimeFilter{
		Genres: []string{"Action", "Adventure"},
	}, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, result, 1)
	assert.Equal(t, "Naruto", result[0].(*models.Anime).Title)
}

func TestGetWithFilters_MultipleTags_ANDSemantics(t *testing.T) {
	repo, db := newTestRepo(t)
	ctx := context.Background()

	ninja := models.Tag{Name: "Ninja"}
	magic := models.Tag{Name: "Magic"}
	require.NoError(t, db.Create(&ninja).Error)
	require.NoError(t, db.Create(&magic).Error)

	// "NinjaMage" has both Ninja and Magic.
	seed(t, db, &models.Anime{Title: "NinjaMage", Tags: []models.Tag{ninja, magic}})
	// "NinjaOnly" has only Ninja.
	seed(t, db, &models.Anime{Title: "NinjaOnly", Tags: []models.Tag{ninja}})

	result, total, err := repo.GetWithFilters(ctx, models.AnimeFilter{
		Tags: []string{"Ninja", "Magic"},
	}, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, result, 1)
	assert.Equal(t, "NinjaMage", result[0].(*models.Anime).Title)
}

func TestGetWithFilters_EpisodesRange(t *testing.T) {
	repo, db := newTestRepo(t)
	ctx := context.Background()

	seed(t, db, &models.Anime{Title: "Short", Episodes: 12})
	seed(t, db, &models.Anime{Title: "Medium", Episodes: 26})
	seed(t, db, &models.Anime{Title: "Long", Episodes: 100})

	// Episodes between 20 and 50.
	result, total, err := repo.GetWithFilters(ctx, models.AnimeFilter{
		EpisodesMin: 20,
		EpisodesMax: 50,
	}, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, result, 1)
	assert.Equal(t, "Medium", result[0].(*models.Anime).Title)
}

func TestGetWithFilters_EpisodesMin_Only(t *testing.T) {
	repo, db := newTestRepo(t)
	ctx := context.Background()

	seed(t, db, &models.Anime{Title: "Short", Episodes: 12})
	seed(t, db, &models.Anime{Title: "Long", Episodes: 100})

	result, total, err := repo.GetWithFilters(ctx, models.AnimeFilter{EpisodesMin: 50}, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "Long", result[0].(*models.Anime).Title)
}
