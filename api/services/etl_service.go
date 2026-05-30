package services

import (
	"context"
	"errors"
	"fmt"
	"myanimeapi/api/models"
	"myanimeapi/internal/db"
	"myanimeapi/internal/etl/jikan"
	"myanimeapi/internal/logger"
	"time"

	gormerrors "gorm.io/gorm"
)

// ETLServiceInterface defines the contract for data-ingestion operations.
type ETLServiceInterface interface {
	SyncFromJikan(ctx context.Context, pages int) (SyncResult, error)
}

// SyncResult summarises a completed sync run.
type SyncResult struct {
	Imported int `json:"imported"`
	Updated  int `json:"updated"`
	Skipped  int `json:"skipped"`
	Errors   int `json:"errors"`
	Pages    int `json:"pages"`
}

// ETLService implements ETLServiceInterface.
// It bypasses the high-level service/repository interfaces so it can use GORM
// bulk operations (FirstOrCreate, Association.Replace) that the generic interface
// does not expose.
type ETLService struct {
	db     db.DBInterface
	jikan  *jikan.Client
	logger *logger.Logger
}

// NewETLService creates an ETLService backed by the given database.
func NewETLService(database db.DBInterface) ETLServiceInterface {
	return &ETLService{
		db:     database,
		jikan:  jikan.NewClient(),
		logger: logger.New(),
	}
}

// SyncFromJikan fetches `pages` pages from the Jikan top-anime list
// and upserts each entry into the local database.
// Pages are 1-indexed; Jikan returns 25 entries per page.
func (s *ETLService) SyncFromJikan(ctx context.Context, pages int) (SyncResult, error) {
	result := SyncResult{Pages: pages}

	for page := 1; page <= pages; page++ {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}

		s.logger.WithFields(map[string]interface{}{
			"page":  page,
			"total": pages,
		}).Info("ETL: fetching Jikan page")

		resp, err := s.jikan.GetTopAnime(ctx, page)
		if err != nil {
			s.logger.WithFields(map[string]interface{}{
				"page":  page,
				"error": err.Error(),
			}).Error("ETL: failed to fetch from Jikan")
			result.Errors++
			// Wait before retrying the next page on transient errors
			time.Sleep(2 * time.Second)
			continue
		}

		for _, entry := range resp.Data {
			imported, err := s.upsertAnime(ctx, entry)
			if err != nil {
				s.logger.WithFields(map[string]interface{}{
					"mal_id": entry.MALID,
					"title":  entry.Title,
					"error":  err.Error(),
				}).Error("ETL: failed to upsert anime")
				result.Errors++
				continue
			}
			if imported {
				result.Imported++
			} else {
				result.Updated++
			}
		}

		if !resp.Pagination.HasNextPage {
			result.Pages = page // actual pages consumed
			break
		}

		// Respect Jikan's rate limit (~3 req/s).  500 ms gives comfortable margin.
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}
	}

	s.logger.WithFields(map[string]interface{}{
		"imported": result.Imported,
		"updated":  result.Updated,
		"errors":   result.Errors,
	}).Info("ETL: sync complete")

	return result, nil
}

// upsertAnime inserts a new anime or updates the existing one matched by mal_id.
// Returns (true, nil) on insert, (false, nil) on update.
func (s *ETLService) upsertAnime(ctx context.Context, entry jikan.AnimeEntry) (bool, error) {
	gdb := s.db.WithContext(ctx)

	genres, err := s.findOrCreateGenres(ctx, entry.Genres)
	if err != nil {
		return false, fmt.Errorf("genres: %w", err)
	}
	tags, err := s.findOrCreateTags(ctx, entry.Themes)
	if err != nil {
		return false, fmt.Errorf("tags: %w", err)
	}

	startDate := parseJikanDate(entry.Aired.From)
	endDate := parseJikanDate(entry.Aired.To)
	status := normalizeStatus(entry.Status)

	coverURL := entry.Images.JPG.LargeImageURL
	if coverURL == "" {
		coverURL = entry.Images.JPG.ImageURL
	}
	malID := entry.MALID

	// Look up by mal_id first
	var existing models.Anime
	lookupErr := gdb.Where("mal_id = ?", malID).First(&existing).Error

	if lookupErr != nil && !errors.Is(lookupErr, gormerrors.ErrRecordNotFound) {
		return false, fmt.Errorf("lookup anime mal_id=%d: %w", malID, lookupErr)
	}

	isNew := errors.Is(lookupErr, gormerrors.ErrRecordNotFound)

	anime := existing // copy (zero-value if new)
	anime.Title = entry.Title
	anime.Description = entry.Synopsis
	anime.Rating = entry.Score
	anime.Episodes = entry.Episodes
	anime.Status = status
	anime.StartDate = startDate
	anime.EndDate = endDate
	anime.MALId = &malID
	anime.CoverURL = &coverURL

	if isNew {
		anime.Genres = genres
		anime.Tags = tags
		if err := gdb.Create(&anime).Error; err != nil {
			return false, fmt.Errorf("create anime: %w", err)
		}
		return true, nil
	}

	// Update scalar fields
	if err := gdb.Save(&anime).Error; err != nil {
		return false, fmt.Errorf("update anime: %w", err)
	}
	// Replace many2many associations atomically
	if err := gdb.Model(&anime).Association("Genres").Replace(genres); err != nil {
		return false, fmt.Errorf("replace genres: %w", err)
	}
	if err := gdb.Model(&anime).Association("Tags").Replace(tags); err != nil {
		return false, fmt.Errorf("replace tags: %w", err)
	}
	return false, nil
}

// findOrCreateGenres resolves a slice of Jikan genre names into Genre model rows,
// creating any that do not yet exist.
func (s *ETLService) findOrCreateGenres(ctx context.Context, jikanTags []jikan.Tag) ([]models.Genre, error) {
	gdb := s.db.WithContext(ctx)
	result := make([]models.Genre, 0, len(jikanTags))
	for _, jt := range jikanTags {
		g := models.Genre{Name: jt.Name}
		if err := gdb.Where(models.Genre{Name: jt.Name}).FirstOrCreate(&g).Error; err != nil {
			return nil, fmt.Errorf("genre %q: %w", jt.Name, err)
		}
		result = append(result, g)
	}
	return result, nil
}

// findOrCreateTags resolves Jikan "themes" into Tag model rows.
func (s *ETLService) findOrCreateTags(ctx context.Context, themes []jikan.Tag) ([]models.Tag, error) {
	gdb := s.db.WithContext(ctx)
	result := make([]models.Tag, 0, len(themes))
	for _, t := range themes {
		tag := models.Tag{Name: t.Name}
		if err := gdb.Where(models.Tag{Name: t.Name}).FirstOrCreate(&tag).Error; err != nil {
			return nil, fmt.Errorf("tag %q: %w", t.Name, err)
		}
		result = append(result, tag)
	}
	return result, nil
}

// parseJikanDate converts a nullable Jikan RFC3339 date string to time.Time.
// Returns zero value on nil or parse failure.
func parseJikanDate(s *string) time.Time {
	if s == nil || *s == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return time.Time{}
	}
	return t
}

// normalizeStatus maps Jikan status strings to the vocabulary used in the DB.
func normalizeStatus(s string) string {
	switch s {
	case "Finished Airing":
		return "Completed"
	case "Currently Airing":
		return "Airing"
	case "Not yet aired":
		return "Upcoming"
	default:
		return s
	}
}
