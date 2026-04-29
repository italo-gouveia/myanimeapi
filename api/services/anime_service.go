package services

import (
	"context"
	"encoding/json"
	"fmt"
	"myanimeapi/api/adapters/cache"
	"myanimeapi/api/models"
	"myanimeapi/api/repositories"
	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"
	"net/http"
	"time"
)

// animeListCache is used to serialise paginated anime lists into/out of cache.
type animeListCache struct {
	Items []*models.Anime `json:"items"`
	Total int64           `json:"total"`
}

// AnimeServiceInterface defines the interface for anime service operations
type AnimeServiceInterface interface {
	GetAnimeByID(ctx context.Context, id uint) (*models.Anime, error)
	GetAllAnimes(ctx context.Context, page, limit int) ([]*models.Anime, int64, error)
	CreateAnime(ctx context.Context, anime *models.Anime) error
	UpdateAnime(ctx context.Context, anime *models.Anime) error
	DeleteAnime(ctx context.Context, id uint) error
	GetAnimesByTitle(ctx context.Context, title string, page, limit int) ([]*models.Anime, int64, error)
	GetAnimesByGenre(ctx context.Context, genre string, page, limit int) ([]*models.Anime, int64, error)
	GetReviewsForAnime(ctx context.Context, animeID uint, page, limit int) ([]*models.Review, int64, error)
	AddGenresToAnime(ctx context.Context, animeID uint, genreIDs []uint) error
	RemoveGenresFromAnime(ctx context.Context, animeID uint, genreIDs []uint) error
	AddTagsToAnime(ctx context.Context, animeID uint, tagIDs []uint) error
	RemoveTagsFromAnime(ctx context.Context, animeID uint, tagIDs []uint) error
}

// AnimeService handles business logic for anime operations
type AnimeService struct {
	animeRepo  repositories.AnimeRepository
	genreRepo  repositories.GenreRepository
	tagRepo    repositories.TagRepository
	reviewRepo repositories.ReviewRepository
	logger     *logger.Logger
	cache      cache.CacheInterface
}

// NewAnimeService creates a new AnimeService instance
func NewAnimeService(
	animeRepo repositories.AnimeRepository,
	genreRepo repositories.GenreRepository,
	tagRepo repositories.TagRepository,
	reviewRepo repositories.ReviewRepository,
	cacheImpl cache.CacheInterface,
) *AnimeService {
	return &AnimeService{
		animeRepo:  animeRepo,
		genreRepo:  genreRepo,
		tagRepo:    tagRepo,
		reviewRepo: reviewRepo,
		logger:     logger.New(),
		cache:      cacheImpl,
	}
}

// GetAnimeByID retrieves an anime by ID
func (s *AnimeService) GetAnimeByID(ctx context.Context, id uint) (*models.Anime, error) {
	s.logger.WithField("anime_id", id).Info("Retrieving anime")

	// Cache-aside: check cache first
	cacheKey := cache.KeyAnime(id)
	if data, err := s.cache.Get(ctx, cacheKey); err == nil {
		var anime models.Anime
		if jsonErr := json.Unmarshal(data, &anime); jsonErr == nil {
			s.logger.WithField("anime_id", id).Info("Anime cache hit")
			return &anime, nil
		}
	}

	animeInterface, err := s.animeRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"anime_id": id,
			"error":    err.Error(),
		}).Error("Failed to retrieve anime")
		return nil, errors.NewError(errors.ErrResourceNotFound, "Anime not found",
			fmt.Sprintf("Anime with ID %d not found", id), http.StatusNotFound,
			map[string]interface{}{
				"anime_id": id,
			},
			err)
	}

	anime, ok := animeInterface.(*models.Anime)
	if !ok {
		s.logger.WithField("anime_id", id).Error("Invalid anime type")
		return nil, errors.NewError(errors.ErrInternalServer, "Invalid anime type", "Type assertion failed", http.StatusInternalServerError,
			map[string]interface{}{
				"anime_id": id,
			},
			nil)
	}

	// Populate cache
	if data, jsonErr := json.Marshal(anime); jsonErr == nil {
		_ = s.cache.Set(ctx, cacheKey, data, cache.TTLAnime)
	}

	s.logger.WithFields(map[string]interface{}{
		"anime_id":   id,
		"anime_name": anime.Title,
	}).Info("Successfully retrieved anime")
	return anime, nil
}

// GetAllAnimes retrieves all animes with pagination
func (s *AnimeService) GetAllAnimes(ctx context.Context, page, limit int) ([]*models.Anime, int64, error) {
	s.logger.WithFields(map[string]interface{}{
		"page":  page,
		"limit": limit,
	}).Info("Retrieving all animes")

	// Cache-aside: check cache first
	cacheKey := cache.KeyAnimeList(page, limit)
	if data, err := s.cache.Get(ctx, cacheKey); err == nil {
		var cached animeListCache
		if jsonErr := json.Unmarshal(data, &cached); jsonErr == nil {
			s.logger.WithFields(map[string]interface{}{"page": page, "limit": limit}).Info("Anime list cache hit")
			return cached.Items, cached.Total, nil
		}
	}

	animesInterface, total, err := s.animeRepo.GetAll(ctx, page, limit)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"page":  page,
			"limit": limit,
			"error": err.Error(),
		}).Error("Failed to retrieve animes")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve animes",
			err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"page":  page,
				"limit": limit,
			},
			err)
	}

	animes := make([]*models.Anime, len(animesInterface))
	for i, animeInterface := range animesInterface {
		anime, ok := animeInterface.(*models.Anime)
		if !ok {
			s.logger.WithFields(map[string]interface{}{
				"index": i,
			}).Error("Invalid anime type in slice")
			return nil, 0, errors.NewError(errors.ErrInternalServer, "Invalid anime type", fmt.Sprintf("Type assertion failed at index %d", i), http.StatusInternalServerError,
				map[string]interface{}{
					"index": i,
				},
				nil)
		}
		animes[i] = anime
	}

	// Populate cache
	if data, jsonErr := json.Marshal(animeListCache{Items: animes, Total: total}); jsonErr == nil {
		_ = s.cache.Set(ctx, cacheKey, data, cache.TTLList)
	}

	s.logger.WithFields(map[string]interface{}{
		"count": len(animes),
		"total": total,
	}).Info("Successfully retrieved animes")
	return animes, total, nil
}

// CreateAnime creates a new anime
func (s *AnimeService) CreateAnime(ctx context.Context, anime *models.Anime) error {
	s.logger.WithField("anime_title", anime.Title).Info("Creating new anime")

	// Check if title already exists
	existingAnimes, _, err := s.animeRepo.GetByTitle(ctx, anime.Title, 1, 1)
	if err == nil && len(existingAnimes) > 0 {
		s.logger.WithField("anime_title", anime.Title).Warning("Anime with same title already exists")
		return errors.NewError(errors.ErrConflict, "Title already exists", fmt.Sprintf("Anime with title '%s' already exists", anime.Title), http.StatusConflict,
			map[string]interface{}{
				"title": anime.Title,
			},
			nil)
	}

	// Set timestamps
	anime.CreatedAt = time.Now()
	anime.UpdatedAt = time.Now()

	// Create anime
	if err := s.animeRepo.Create(ctx, anime); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"anime_title": anime.Title,
			"error":       err.Error(),
		}).Error("Failed to create anime")
		return errors.NewError(errors.ErrInternalServer, "Failed to create anime",
			err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"title": anime.Title,
			},
			err)
	}

	// Invalidate list caches — new item may appear in any page
	_ = s.cache.DeleteByPattern(ctx, cache.PatternAllAnimes)

	s.logger.WithFields(map[string]interface{}{
		"anime_id":   anime.ID,
		"anime_name": anime.Title,
	}).Info("Successfully created anime")
	return nil
}

// UpdateAnime updates an existing anime
func (s *AnimeService) UpdateAnime(ctx context.Context, anime *models.Anime) error {
	s.logger.WithFields(map[string]interface{}{
		"anime_id":   anime.ID,
		"anime_name": anime.Title,
	}).Info("Updating anime")

	// Check if anime exists
	existingAnime, err := s.animeRepo.GetByID(ctx, anime.ID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"anime_id": anime.ID,
			"error":    err.Error(),
		}).Error("Failed to retrieve anime for update")
		return errors.NewError(errors.ErrResourceNotFound, "Anime not found",
			fmt.Sprintf("Anime with ID %d not found", anime.ID), http.StatusNotFound,
			map[string]interface{}{
				"anime_id": anime.ID,
			},
			err)
	}

	// Check if title is being changed and if it already exists
	if anime.Title != existingAnime.(*models.Anime).Title {
		titleAnimes, _, err := s.animeRepo.GetByTitle(ctx, anime.Title, 1, 1)
		if err == nil && len(titleAnimes) > 0 && titleAnimes[0].ID != anime.ID {
			s.logger.WithField("anime_title", anime.Title).Warning("Anime with same title already exists")
			return errors.NewError(errors.ErrConflict, "Title already exists", fmt.Sprintf("Anime with title '%s' is already taken by another anime", anime.Title), http.StatusConflict,
				map[string]interface{}{
					"title": anime.Title,
				},
				nil)
		}
	}

	// Set updated timestamp
	anime.UpdatedAt = time.Now()

	// Update anime
	if err := s.animeRepo.Update(ctx, anime); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"anime_id":   anime.ID,
			"anime_name": anime.Title,
			"error":      err.Error(),
		}).Error("Failed to update anime")
		return errors.NewError(errors.ErrInternalServer, "Failed to update anime",
			err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"anime_id": anime.ID,
				"title":    anime.Title,
			},
			err)
	}

	// Invalidate the single-item entry and all list pages
	_ = s.cache.Delete(ctx, cache.KeyAnime(anime.ID))
	_ = s.cache.DeleteByPattern(ctx, cache.PatternAllAnimes)

	s.logger.WithFields(map[string]interface{}{
		"anime_id":   anime.ID,
		"anime_name": anime.Title,
	}).Info("Successfully updated anime")
	return nil
}

// DeleteAnime deletes an anime by ID
func (s *AnimeService) DeleteAnime(ctx context.Context, id uint) error {
	s.logger.WithField("anime_id", id).Info("Deleting anime")

	// Check if anime exists
	existingAnime, err := s.animeRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"anime_id": id,
			"error":    err.Error(),
		}).Error("Failed to retrieve anime for deletion")
		return errors.NewError(errors.ErrResourceNotFound, "Anime not found",
			fmt.Sprintf("Anime with ID %d not found", id), http.StatusNotFound,
			map[string]interface{}{
				"anime_id": id,
			},
			err)
	}

	animeModel, ok := existingAnime.(*models.Anime)
	if !ok {
		s.logger.WithField("anime_id", id).Error("Invalid anime type")
		return errors.NewError(errors.ErrInternalServer, "Invalid anime type", "Type assertion failed", http.StatusInternalServerError,
			map[string]interface{}{
				"anime_id": id,
			},
			nil)
	}

	// Delete anime
	if err := s.animeRepo.Delete(ctx, id); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"anime_id": id,
			"error":    err.Error(),
		}).Error("Failed to delete anime")
		return errors.NewError(errors.ErrInternalServer, "Failed to delete anime",
			err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"anime_id": id,
			},
			err)
	}

	// Invalidate the single-item entry and all list pages
	_ = s.cache.Delete(ctx, cache.KeyAnime(id))
	_ = s.cache.DeleteByPattern(ctx, cache.PatternAllAnimes)

	s.logger.WithFields(map[string]interface{}{
		"anime_id":   id,
		"anime_name": animeModel.Title,
	}).Info("Successfully deleted anime")
	return nil
}

// GetAnimesByTitle retrieves animes by title with pagination
func (s *AnimeService) GetAnimesByTitle(ctx context.Context, title string, page, limit int) ([]*models.Anime, int64, error) {
	s.logger.WithFields(map[string]interface{}{
		"title": title,
		"page":  page,
		"limit": limit,
	}).Info("Retrieving animes by title")

	// Cache-aside: check cache first
	cacheKey := cache.KeyAnimeSearch(title, page, limit)
	if data, err := s.cache.Get(ctx, cacheKey); err == nil {
		var cached animeListCache
		if jsonErr := json.Unmarshal(data, &cached); jsonErr == nil {
			s.logger.WithField("title", title).Info("Anime search cache hit")
			return cached.Items, cached.Total, nil
		}
	}

	animes, total, err := s.animeRepo.GetByTitle(ctx, title, page, limit)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"title": title,
			"error": err.Error(),
		}).Error("Failed to retrieve animes by title")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve animes by title",
			err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"title": title,
				"page":  page,
				"limit": limit,
			},
			err)
	}

	// Convert to pointer slice
	result := make([]*models.Anime, len(animes))
	for i := range animes {
		result[i] = &animes[i]
	}

	// Populate cache
	if data, jsonErr := json.Marshal(animeListCache{Items: result, Total: total}); jsonErr == nil {
		_ = s.cache.Set(ctx, cacheKey, data, cache.TTLList)
	}

	s.logger.WithFields(map[string]interface{}{
		"title": title,
		"count": len(result),
		"total": total,
	}).Info("Successfully retrieved animes by title")
	return result, total, nil
}

// GetAnimesByGenre retrieves animes by genre with pagination
func (s *AnimeService) GetAnimesByGenre(ctx context.Context, genre string, page, limit int) ([]*models.Anime, int64, error) {
	s.logger.WithFields(map[string]interface{}{
		"genre": genre,
		"page":  page,
		"limit": limit,
	}).Info("Retrieving animes by genre")

	animes, total, err := s.animeRepo.GetByGenre(ctx, genre, page, limit)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"genre": genre,
			"error": err.Error(),
		}).Error("Failed to retrieve animes by genre")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve animes by genre",
			err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"genre": genre,
				"page":  page,
				"limit": limit,
			},
			err)
	}

	result := make([]*models.Anime, len(animes))
	for i := range animes {
		result[i] = &animes[i]
	}

	s.logger.WithFields(map[string]interface{}{
		"genre": genre,
		"count": len(result),
		"total": total,
	}).Info("Successfully retrieved animes by genre")
	return result, total, nil
}

// GetReviewsForAnime retrieves paginated reviews for an anime
func (s *AnimeService) GetReviewsForAnime(ctx context.Context, animeID uint, page, limit int) ([]*models.Review, int64, error) {
	s.logger.WithFields(map[string]interface{}{
		"anime_id": animeID,
		"page":     page,
		"limit":    limit,
	}).Info("Retrieving reviews for anime")

	// Check if anime exists
	_, err := s.animeRepo.GetByID(ctx, animeID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"anime_id": animeID,
			"error":    err.Error(),
		}).Error("Failed to retrieve anime for reviews")
		return nil, 0, errors.NewError(errors.ErrResourceNotFound, "Anime not found",
			fmt.Sprintf("Anime with ID %d not found", animeID), http.StatusNotFound,
			map[string]interface{}{
				"anime_id": animeID,
			},
			err)
	}

	reviews, total, err := s.reviewRepo.GetByAnimeID(ctx, animeID, page, limit)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"anime_id": animeID,
			"error":    err.Error(),
		}).Error("Failed to retrieve reviews for anime")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve reviews for anime",
			err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"anime_id": animeID,
				"page":     page,
				"limit":    limit,
			},
			err)
	}

	result := make([]*models.Review, len(reviews))
	for i := range reviews {
		result[i] = &reviews[i]
	}

	s.logger.WithFields(map[string]interface{}{
		"anime_id": animeID,
		"count":    len(result),
		"total":    total,
	}).Info("Successfully retrieved reviews for anime")
	return result, total, nil
}

// AddGenresToAnime adds genres to an anime
func (s *AnimeService) AddGenresToAnime(ctx context.Context, animeID uint, genreIDs []uint) error {
	s.logger.WithFields(map[string]interface{}{
		"anime_id":  animeID,
		"genre_ids": genreIDs,
	}).Info("Adding genres to anime")

	// Check if anime exists
	animeInterface, err := s.animeRepo.GetByID(ctx, animeID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"anime_id": animeID,
			"error":    err.Error(),
		}).Error("Failed to retrieve anime for adding genres")
		return errors.NewError(errors.ErrResourceNotFound, "Anime not found",
			fmt.Sprintf("Anime with ID %d not found", animeID), http.StatusNotFound,
			map[string]interface{}{
				"anime_id":  animeID,
				"genre_ids": genreIDs,
			},
			err)
	}

	anime, ok := animeInterface.(*models.Anime)
	if !ok {
		s.logger.WithField("anime_id", animeID).Error("Invalid anime type")
		return errors.NewError(errors.ErrInternalServer, "Invalid anime type", "Type assertion failed", http.StatusInternalServerError,
			map[string]interface{}{
				"anime_id":  animeID,
				"genre_ids": genreIDs,
			},
			nil)
	}

	// Get genres by IDs
	genres, err := s.genreRepo.GetByIDs(ctx, genreIDs)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"genre_ids": genreIDs,
			"error":     err.Error(),
		}).Error("Failed to retrieve genres")
		return errors.NewError(errors.ErrInternalServer, "Failed to retrieve genres",
			err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"genre_ids": genreIDs,
			},
			err)
	}

	// Add genres to anime
	anime.Genres = append(anime.Genres, genres...)
	if err := s.animeRepo.Update(ctx, anime); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"anime_id":  animeID,
			"genre_ids": genreIDs,
			"error":     err.Error(),
		}).Error("Failed to update anime with new genres")
		return errors.NewError(errors.ErrInternalServer, "Failed to update anime with new genres",
			err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"anime_id":  animeID,
				"genre_ids": genreIDs,
			},
			err)
	}

	// Invalidate single-item and list caches
	_ = s.cache.Delete(ctx, cache.KeyAnime(animeID))
	_ = s.cache.DeleteByPattern(ctx, cache.PatternAllAnimes)

	s.logger.WithFields(map[string]interface{}{
		"anime_id":  animeID,
		"genre_ids": genreIDs,
	}).Info("Successfully added genres to anime")
	return nil
}

// RemoveGenresFromAnime removes genres from an anime
func (s *AnimeService) RemoveGenresFromAnime(ctx context.Context, animeID uint, genreIDs []uint) error {
	s.logger.WithFields(map[string]interface{}{
		"anime_id":  animeID,
		"genre_ids": genreIDs,
	}).Info("Removing genres from anime")

	// Check if anime exists
	anime, err := s.animeRepo.GetByID(ctx, animeID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"anime_id": animeID,
			"error":    err.Error(),
		}).Error("Failed to retrieve anime for removing genres")
		return errors.NewError(errors.ErrResourceNotFound, "Anime not found",
			fmt.Sprintf("Anime with ID %d not found", animeID), http.StatusNotFound,
			map[string]interface{}{
				"anime_id":  animeID,
				"genre_ids": genreIDs,
			},
			err)
	}

	animeObj, ok := anime.(*models.Anime)
	if !ok {
		s.logger.WithField("anime_id", animeID).Error("Invalid anime type")
		return errors.NewError(errors.ErrInternalServer, "Invalid anime type", "Type assertion failed", http.StatusInternalServerError,
			map[string]interface{}{
				"anime_id":  animeID,
				"genre_ids": genreIDs,
			},
			nil)
	}

	// Create a map of genre IDs to remove
	genreMap := make(map[uint]bool)
	for _, id := range genreIDs {
		genreMap[id] = true
	}

	// Filter out genres to remove
	var filteredGenres []models.Genre
	for _, genre := range animeObj.Genres {
		if !genreMap[genre.ID] {
			filteredGenres = append(filteredGenres, genre)
		}
	}

	// Update anime with filtered genres
	animeObj.Genres = filteredGenres
	if err := s.animeRepo.Update(ctx, animeObj); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"anime_id":  animeID,
			"genre_ids": genreIDs,
			"error":     err.Error(),
		}).Error("Failed to update anime after removing genres")
		return errors.NewError(errors.ErrInternalServer, "Failed to update anime after removing genres",
			err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"anime_id":  animeID,
				"genre_ids": genreIDs,
			},
			err)
	}

	// Invalidate single-item and list caches
	_ = s.cache.Delete(ctx, cache.KeyAnime(animeID))
	_ = s.cache.DeleteByPattern(ctx, cache.PatternAllAnimes)

	s.logger.WithFields(map[string]interface{}{
		"anime_id":  animeID,
		"genre_ids": genreIDs,
	}).Info("Successfully removed genres from anime")
	return nil
}

// AddTagsToAnime adds tags to an anime
func (s *AnimeService) AddTagsToAnime(ctx context.Context, animeID uint, tagIDs []uint) error {
	s.logger.WithFields(map[string]interface{}{
		"anime_id": animeID,
		"tag_ids":  tagIDs,
	}).Info("Adding tags to anime")

	// Check if anime exists
	animeInterface, err := s.animeRepo.GetByID(ctx, animeID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"anime_id": animeID,
			"error":    err.Error(),
		}).Error("Failed to retrieve anime for adding tags")
		return errors.NewError(errors.ErrResourceNotFound, "Anime not found",
			fmt.Sprintf("Anime with ID %d not found", animeID), http.StatusNotFound,
			map[string]interface{}{
				"anime_id": animeID,
				"tag_ids":  tagIDs,
			},
			err)
	}

	anime, ok := animeInterface.(*models.Anime)
	if !ok {
		s.logger.WithField("anime_id", animeID).Error("Invalid anime type")
		return errors.NewError(errors.ErrInternalServer, "Invalid anime type", "Type assertion failed", http.StatusInternalServerError,
			map[string]interface{}{
				"anime_id": animeID,
				"tag_ids":  tagIDs,
			},
			nil)
	}

	// Get tags by IDs
	tags, err := s.tagRepo.GetByIDs(ctx, tagIDs)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"tag_ids": tagIDs,
			"error":   err.Error(),
		}).Error("Failed to retrieve tags")
		return errors.NewError(errors.ErrInternalServer, "Failed to retrieve tags",
			err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"tag_ids": tagIDs,
			},
			err)
	}

	// Add tags to anime
	anime.Tags = append(anime.Tags, tags...)
	if err := s.animeRepo.Update(ctx, anime); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"anime_id": animeID,
			"tag_ids":  tagIDs,
			"error":    err.Error(),
		}).Error("Failed to update anime with new tags")
		return errors.NewError(errors.ErrInternalServer, "Failed to update anime with new tags",
			err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"anime_id": animeID,
				"tag_ids":  tagIDs,
			},
			err)
	}

	// Invalidate single-item and list caches
	_ = s.cache.Delete(ctx, cache.KeyAnime(animeID))
	_ = s.cache.DeleteByPattern(ctx, cache.PatternAllAnimes)

	s.logger.WithFields(map[string]interface{}{
		"anime_id": animeID,
		"tag_ids":  tagIDs,
	}).Info("Successfully added tags to anime")
	return nil
}

// RemoveTagsFromAnime removes tags from an anime
func (s *AnimeService) RemoveTagsFromAnime(ctx context.Context, animeID uint, tagIDs []uint) error {
	s.logger.WithFields(map[string]interface{}{
		"anime_id": animeID,
		"tag_ids":  tagIDs,
	}).Info("Removing tags from anime")

	// Check if anime exists
	anime, err := s.animeRepo.GetByID(ctx, animeID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"anime_id": animeID,
			"error":    err.Error(),
		}).Error("Failed to retrieve anime for removing tags")
		return errors.NewError(errors.ErrResourceNotFound, "Anime not found",
			fmt.Sprintf("Anime with ID %d not found", animeID), http.StatusNotFound,
			map[string]interface{}{
				"anime_id": animeID,
				"tag_ids":  tagIDs,
			},
			err)
	}

	animeObj, ok := anime.(*models.Anime)
	if !ok {
		s.logger.WithField("anime_id", animeID).Error("Invalid anime type")
		return errors.NewError(errors.ErrInternalServer, "Invalid anime type", "Type assertion failed", http.StatusInternalServerError,
			map[string]interface{}{
				"anime_id": animeID,
				"tag_ids":  tagIDs,
			},
			nil)
	}

	// Create a map of tag IDs to remove
	tagMap := make(map[uint]bool)
	for _, id := range tagIDs {
		tagMap[id] = true
	}

	// Filter out tags to remove
	var filteredTags []models.Tag
	for _, tag := range animeObj.Tags {
		if !tagMap[tag.ID] {
			filteredTags = append(filteredTags, tag)
		}
	}

	// Update anime with filtered tags
	animeObj.Tags = filteredTags
	if err := s.animeRepo.Update(ctx, animeObj); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"anime_id": animeID,
			"tag_ids":  tagIDs,
			"error":    err.Error(),
		}).Error("Failed to update anime after removing tags")
		return errors.NewError(errors.ErrInternalServer, "Failed to update anime after removing tags",
			err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"anime_id": animeID,
				"tag_ids":  tagIDs,
			},
			err)
	}

	// Invalidate single-item and list caches
	_ = s.cache.Delete(ctx, cache.KeyAnime(animeID))
	_ = s.cache.DeleteByPattern(ctx, cache.PatternAllAnimes)

	s.logger.WithFields(map[string]interface{}{
		"anime_id": animeID,
		"tag_ids":  tagIDs,
	}).Info("Successfully removed tags from anime")
	return nil
}
