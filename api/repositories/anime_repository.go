package repositories

import (
	"context"
	"fmt"
	"myanimeapi/api/models"
	"myanimeapi/internal/db"
	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"
	"net/http"
)

// AnimeRepositoryImpl implements the AnimeRepository interface
type AnimeRepositoryImpl struct {
	db     db.DBInterface
	logger *logger.Logger
}

// NewAnimeRepository creates a new AnimeRepositoryImpl instance
func NewAnimeRepository(db db.DBInterface) AnimeRepository {
	return &AnimeRepositoryImpl{
		db:     db,
		logger: logger.New(),
	}
}

// GetByID retrieves an anime by its ID
func (r *AnimeRepositoryImpl) GetByID(ctx context.Context, id uint) (interface{}, error) {
	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Retrieving anime by ID")

	var anime models.Anime
	if err := r.db.WithContext(ctx).Preload("Reviews").Preload("Favorites").First(&anime, id).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to retrieve anime")
		return nil, errors.NewError(errors.ErrResourceNotFound, "Anime not found", fmt.Sprintf("Anime with ID %d not found", id), http.StatusNotFound, map[string]interface{}{"id": id}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Successfully retrieved anime")
	return &anime, nil
}

// GetAll retrieves all animes with optional pagination
func (r *AnimeRepositoryImpl) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	r.logger.WithFields(map[string]interface{}{
		"page":  page,
		"limit": limit,
	}).Info("Retrieving all animes")

	var animes []models.Anime
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Anime{}).Count(&total).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to count animes")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count animes", err.Error(), http.StatusInternalServerError, nil, err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve animes with pagination
	if err := r.db.WithContext(ctx).Preload("Reviews").Preload("Favorites").Offset(offset).Limit(limit).Find(&animes).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to retrieve animes")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve animes", err.Error(), http.StatusInternalServerError, nil, err)
	}

	// Convert to interface slice
	result := make([]interface{}, len(animes))
	for i, anime := range animes {
		result[i] = &anime
	}

	r.logger.WithFields(map[string]interface{}{
		"count": len(animes),
		"total": total,
	}).Info("Successfully retrieved animes")
	return result, total, nil
}

// Create creates a new anime
func (r *AnimeRepositoryImpl) Create(ctx context.Context, entity interface{}) error {
	r.logger.Info("Creating new anime")

	anime, ok := entity.(*models.Anime)
	if !ok {
		r.logger.Error("Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.Anime", http.StatusBadRequest, nil, nil)
	}

	if err := r.db.WithContext(ctx).Create(anime).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to create anime")
		return errors.NewError(errors.ErrInternalServer, "Failed to create anime", err.Error(), http.StatusInternalServerError, nil, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": anime.ID,
	}).Info("Successfully created anime")
	return nil
}

// Update updates an existing anime
func (r *AnimeRepositoryImpl) Update(ctx context.Context, entity interface{}) error {
	r.logger.Info("Updating anime")

	anime, ok := entity.(*models.Anime)
	if !ok {
		r.logger.Error("Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.Anime", http.StatusBadRequest, nil, nil)
	}

	if err := r.db.WithContext(ctx).Save(anime).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to update anime")
		return errors.NewError(errors.ErrInternalServer, "Failed to update anime", err.Error(), http.StatusInternalServerError, nil, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": anime.ID,
	}).Info("Successfully updated anime")
	return nil
}

// Delete deletes an anime by its ID
func (r *AnimeRepositoryImpl) Delete(ctx context.Context, id uint) error {
	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Deleting anime")

	if err := r.db.WithContext(ctx).Delete(&models.Anime{}, id).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to delete anime")
		return errors.NewError(errors.ErrInternalServer, "Failed to delete anime", err.Error(), http.StatusInternalServerError, map[string]interface{}{"id": id}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Successfully deleted anime")
	return nil
}

// GetByTitle retrieves animes by title
func (r *AnimeRepositoryImpl) GetByTitle(ctx context.Context, title string, page, limit int) ([]models.Anime, int64, error) {
	r.logger.WithFields(map[string]interface{}{
		"title": title,
		"page":  page,
		"limit": limit,
	}).Info("Retrieving animes by title")

	var animes []models.Anime
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Anime{}).Where("title LIKE ?", "%"+title+"%").Count(&total).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to count animes")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count animes", err.Error(), http.StatusInternalServerError, map[string]interface{}{"title": title}, err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve animes with pagination
	if err := r.db.WithContext(ctx).Preload("Reviews").Preload("Favorites").Where("title LIKE ?", "%"+title+"%").Offset(offset).Limit(limit).Find(&animes).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to retrieve animes")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve animes", err.Error(), http.StatusInternalServerError, map[string]interface{}{"title": title}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"count": len(animes),
		"total": total,
	}).Info("Successfully retrieved animes by title")
	return animes, total, nil
}

// GetByGenre retrieves animes by genre
func (r *AnimeRepositoryImpl) GetByGenre(ctx context.Context, genre string, page, limit int) ([]models.Anime, int64, error) {
	r.logger.WithFields(map[string]interface{}{
		"genre": genre,
		"page":  page,
		"limit": limit,
	}).Info("Retrieving animes by genre")

	var animes []models.Anime
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Anime{}).Where("genre = ?", genre).Count(&total).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to count animes")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count animes", err.Error(), http.StatusInternalServerError, map[string]interface{}{"genre": genre}, err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve animes with pagination
	if err := r.db.WithContext(ctx).Preload("Reviews").Preload("Favorites").Where("genre = ?", genre).Offset(offset).Limit(limit).Find(&animes).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to retrieve animes")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve animes", err.Error(), http.StatusInternalServerError, map[string]interface{}{"genre": genre}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"count": len(animes),
		"total": total,
	}).Info("Successfully retrieved animes by genre")
	return animes, total, nil
}

// GetReviewsForAnime retrieves reviews for a specific anime
func (r *AnimeRepositoryImpl) GetReviewsForAnime(ctx context.Context, animeID uint, page, limit int) ([]models.Review, int64, error) {
	r.logger.WithFields(map[string]interface{}{
		"anime_id": animeID,
		"page":     page,
		"limit":    limit,
	}).Info("Retrieving reviews for anime")

	var reviews []models.Review
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&models.Review{}).Where("anime_id = ?", animeID).Count(&total).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to count reviews")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count reviews", err.Error(), http.StatusInternalServerError, map[string]interface{}{"anime_id": animeID}, err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated reviews
	if err := r.db.WithContext(ctx).Where("anime_id = ?", animeID).Offset(offset).Limit(limit).Find(&reviews).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to retrieve reviews")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve reviews", err.Error(), http.StatusInternalServerError, map[string]interface{}{"anime_id": animeID}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"count": len(reviews),
		"total": total,
	}).Info("Successfully retrieved reviews for anime")
	return reviews, total, nil
}
