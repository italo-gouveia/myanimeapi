// api/models/filter.go
// AnimeFilter is a value object that carries all optional filter and sort
// parameters accepted by the anime list endpoint.
package models

import (
	"fmt"
	"strings"
)

// AllowedSortFields maps public-facing sort_by values to the fully-qualified
// column name used in ORDER BY. Only whitelisted values are accepted to prevent
// SQL injection.
var AllowedSortFields = map[string]string{
	"title":      "animes.title",
	"rating":     "animes.rating",
	"episodes":   "animes.episodes",
	"created_at": "animes.created_at",
	"start_date": "animes.start_date",
}

// AllowedStatuses is the set of valid status values for filtering.
var AllowedStatuses = map[string]struct{}{
	"Airing":    {},
	"Completed": {},
	"Upcoming":  {},
}

// AnimeFilter carries all optional filter and sort parameters for the anime
// list endpoint. Zero values mean "no constraint / use default".
type AnimeFilter struct {
	// Title performs a case-insensitive LIKE search on animes.title.
	Title string
	// Status filters by exact status: "Airing", "Completed", or "Upcoming" (single, legacy).
	Status string
	// Statuses filters to animes whose status is in the given list (OR semantics).
	Statuses []string
	// Genre filters to animes that belong to a genre with this exact name (single, legacy).
	Genre string
	// Genres filters to animes that have ALL of the named genres (AND semantics).
	Genres []string
	// Tag filters to animes that have a tag with this exact name (single, legacy).
	Tag string
	// Tags filters to animes that have ALL of the named tags (AND semantics).
	Tags []string
	// RatingMin includes only animes with rating >= RatingMin. 0 = no lower bound.
	RatingMin float64
	// RatingMax includes only animes with rating <= RatingMax. 0 = no upper bound.
	RatingMax float64
	// EpisodesMin includes only animes with episodes >= EpisodesMin. 0 = no lower bound.
	EpisodesMin int
	// EpisodesMax includes only animes with episodes <= EpisodesMax. 0 = no upper bound.
	EpisodesMax int
	// YearFrom includes only animes whose start_date year >= YearFrom. 0 = no lower bound.
	YearFrom int
	// YearTo includes only animes whose start_date year <= YearTo. 0 = no upper bound.
	YearTo int
	// SortBy is the field to sort by. Must be a key in AllowedSortFields.
	SortBy string
	// SortOrder is "asc" (default) or "desc".
	SortOrder string
}

// Validate returns a descriptive error if any field contains an invalid value.
func (f AnimeFilter) Validate() error {
	if f.Status != "" {
		if _, ok := AllowedStatuses[f.Status]; !ok {
			return fmt.Errorf("invalid status %q: allowed values are Airing, Completed, Upcoming", f.Status)
		}
	}
	for _, s := range f.Statuses {
		if _, ok := AllowedStatuses[s]; !ok {
			return fmt.Errorf("invalid status %q in statuses: allowed values are Airing, Completed, Upcoming", s)
		}
	}
	if f.SortBy != "" {
		if _, ok := AllowedSortFields[f.SortBy]; !ok {
			keys := make([]string, 0, len(AllowedSortFields))
			for k := range AllowedSortFields {
				keys = append(keys, k)
			}
			return fmt.Errorf("invalid sort_by %q: allowed values are %s", f.SortBy, strings.Join(keys, ", "))
		}
	}
	if f.SortOrder != "" && f.SortOrder != "asc" && f.SortOrder != "desc" {
		return fmt.Errorf("invalid order %q: allowed values are asc, desc", f.SortOrder)
	}
	if f.RatingMin < 0 || f.RatingMin > 10 {
		return fmt.Errorf("rating_min must be between 0 and 10")
	}
	if f.RatingMax < 0 || f.RatingMax > 10 {
		return fmt.Errorf("rating_max must be between 0 and 10")
	}
	if f.RatingMin > 0 && f.RatingMax > 0 && f.RatingMin > f.RatingMax {
		return fmt.Errorf("rating_min cannot be greater than rating_max")
	}
	if f.EpisodesMin < 0 {
		return fmt.Errorf("episodes_min must be >= 0")
	}
	if f.EpisodesMax < 0 {
		return fmt.Errorf("episodes_max must be >= 0")
	}
	if f.EpisodesMin > 0 && f.EpisodesMax > 0 && f.EpisodesMin > f.EpisodesMax {
		return fmt.Errorf("episodes_min cannot be greater than episodes_max")
	}
	if f.YearFrom < 0 {
		return fmt.Errorf("year_from must be >= 0")
	}
	if f.YearTo < 0 {
		return fmt.Errorf("year_to must be >= 0")
	}
	if f.YearFrom > 0 && f.YearTo > 0 && f.YearFrom > f.YearTo {
		return fmt.Errorf("year_from cannot be greater than year_to")
	}
	return nil
}

// OrderClause returns a safe SQL ORDER BY clause derived from the filter.
// Defaults to "animes.created_at DESC" when no sort field is specified.
func (f AnimeFilter) OrderClause() string {
	col, ok := AllowedSortFields[f.SortBy]
	if !ok {
		return "animes.created_at DESC"
	}
	dir := "ASC"
	if strings.ToLower(f.SortOrder) == "desc" {
		dir = "DESC"
	}
	return fmt.Sprintf("%s %s", col, dir)
}

// CacheKeySuffix returns a deterministic string fragment encoding all filter
// fields. Two equal AnimeFilter values always produce the same suffix, making
// it safe to embed in cache keys.
func (f AnimeFilter) CacheKeySuffix() string {
	return fmt.Sprintf(
		"st=%s:sts=%s:g=%s:gs=%s:t=%s:ts=%s:rmin=%.2f:rmax=%.2f:emin=%d:emax=%d:yf=%d:yt=%d:sort=%s:%s:ti=%s",
		f.Status, strings.Join(f.Statuses, ","), f.Genre, strings.Join(f.Genres, ","),
		f.Tag, strings.Join(f.Tags, ","),
		f.RatingMin, f.RatingMax,
		f.EpisodesMin, f.EpisodesMax,
		f.YearFrom, f.YearTo,
		f.SortBy, f.SortOrder, f.Title,
	)
}
