// api/models/filter_test.go
package models_test

import (
	"strings"
	"testing"

	"myanimeapi/api/models"

	"github.com/stretchr/testify/assert"
)

// ---- Validate -----------------------------------------------------------

func TestAnimeFilter_Validate_Empty(t *testing.T) {
	assert.NoError(t, models.AnimeFilter{}.Validate())
}

func TestAnimeFilter_Validate_ValidStatus(t *testing.T) {
	for _, s := range []string{"Airing", "Completed", "Upcoming"} {
		assert.NoError(t, models.AnimeFilter{Status: s}.Validate(), "status=%s", s)
	}
}

func TestAnimeFilter_Validate_InvalidStatus(t *testing.T) {
	err := models.AnimeFilter{Status: "Unknown"}.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid status")
}

func TestAnimeFilter_Validate_ValidSortBy(t *testing.T) {
	for key := range models.AllowedSortFields {
		assert.NoError(t, models.AnimeFilter{SortBy: key}.Validate(), "sort_by=%s", key)
	}
}

func TestAnimeFilter_Validate_InvalidSortBy(t *testing.T) {
	err := models.AnimeFilter{SortBy: "unknown_col"}.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sort_by")
}

func TestAnimeFilter_Validate_SortOrder(t *testing.T) {
	assert.NoError(t, models.AnimeFilter{SortOrder: "asc"}.Validate())
	assert.NoError(t, models.AnimeFilter{SortOrder: "desc"}.Validate())
	assert.Error(t, models.AnimeFilter{SortOrder: "random"}.Validate())
}

func TestAnimeFilter_Validate_RatingRange(t *testing.T) {
	assert.NoError(t, models.AnimeFilter{RatingMin: 5, RatingMax: 9}.Validate())
	assert.Error(t, models.AnimeFilter{RatingMin: -1}.Validate())
	assert.Error(t, models.AnimeFilter{RatingMax: 11}.Validate())
	assert.Error(t, models.AnimeFilter{RatingMin: 9, RatingMax: 5}.Validate())
}

// ---- OrderClause --------------------------------------------------------

func TestAnimeFilter_OrderClause_Default(t *testing.T) {
	clause := models.AnimeFilter{}.OrderClause()
	assert.Equal(t, "animes.created_at DESC", clause)
}

func TestAnimeFilter_OrderClause_RatingDesc(t *testing.T) {
	clause := models.AnimeFilter{SortBy: "rating", SortOrder: "desc"}.OrderClause()
	assert.Equal(t, "animes.rating DESC", clause)
}

func TestAnimeFilter_OrderClause_TitleAsc(t *testing.T) {
	clause := models.AnimeFilter{SortBy: "title", SortOrder: "asc"}.OrderClause()
	assert.Equal(t, "animes.title ASC", clause)
}

func TestAnimeFilter_OrderClause_DefaultDirectionIsAsc(t *testing.T) {
	// When SortOrder is empty but SortBy is set, direction defaults to ASC.
	clause := models.AnimeFilter{SortBy: "episodes"}.OrderClause()
	assert.Equal(t, "animes.episodes ASC", clause)
}

// ---- CacheKeySuffix -----------------------------------------------------

func TestAnimeFilter_CacheKeySuffix_Empty(t *testing.T) {
	suffix := models.AnimeFilter{}.CacheKeySuffix()
	assert.True(t, strings.HasPrefix(suffix, "st=:g=:t=:"), "suffix=%s", suffix)
}

func TestAnimeFilter_CacheKeySuffix_Deterministic(t *testing.T) {
	f := models.AnimeFilter{Status: "Completed", Genre: "Action", RatingMin: 7.5, SortBy: "rating", SortOrder: "desc"}
	assert.Equal(t, f.CacheKeySuffix(), f.CacheKeySuffix())
}

func TestAnimeFilter_CacheKeySuffix_DifferentFiltersProduceDifferentKeys(t *testing.T) {
	a := models.AnimeFilter{Status: "Airing"}.CacheKeySuffix()
	b := models.AnimeFilter{Status: "Completed"}.CacheKeySuffix()
	assert.NotEqual(t, a, b)
}

func TestAnimeFilter_CacheKeySuffix_TitleIncluded(t *testing.T) {
	a := models.AnimeFilter{Title: "Naruto"}.CacheKeySuffix()
	b := models.AnimeFilter{Title: "Bleach"}.CacheKeySuffix()
	assert.Contains(t, a, "ti=Naruto")
	assert.NotEqual(t, a, b)
}
