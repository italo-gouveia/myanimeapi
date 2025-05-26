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

// FavoriteRepositoryImpl implements the FavoriteRepository interface
type FavoriteRepositoryImpl struct {
	db     db.DBInterface
	logger *logger.Logger
}

// NewFavoriteRepository creates a new FavoriteRepositoryImpl instance
func NewFavoriteRepository(db db.DBInterface) FavoriteRepository {
	return &FavoriteRepositoryImpl{
		db:     db,
		logger: logger.New(),
	}
}

// GetByID retrieves a favorite by its ID
func (r *FavoriteRepositoryImpl) GetByID(ctx context.Context, id uint) (interface{}, error) {
	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Retrieving favorite by ID")

	var favorite models.Favorite
	if err := r.db.WithContext(ctx).Preload("User").Preload("Anime").First(&favorite, id).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to retrieve favorite")
		return nil, errors.NewError(errors.ErrResourceNotFound, "Favorite not found", fmt.Sprintf("Favorite with ID %d not found", id), http.StatusNotFound, map[string]interface{}{"id": id}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Successfully retrieved favorite")
	return &favorite, nil
}

// GetAll retrieves all favorites with optional pagination
func (r *FavoriteRepositoryImpl) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	r.logger.WithFields(map[string]interface{}{
		"page":  page,
		"limit": limit,
	}).Info("Retrieving all favorites")

	var favorites []models.Favorite
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Favorite{}).Count(&total).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to count favorites")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count favorites", err.Error(), http.StatusInternalServerError, nil, err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve favorites with pagination
	if err := r.db.WithContext(ctx).Preload("User").Preload("Anime").Offset(offset).Limit(limit).Find(&favorites).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to retrieve favorites")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve favorites", err.Error(), http.StatusInternalServerError, nil, err)
	}

	// Convert to interface slice
	result := make([]interface{}, len(favorites))
	for i, favorite := range favorites {
		result[i] = &favorite
	}

	r.logger.WithFields(map[string]interface{}{
		"count": len(favorites),
		"total": total,
	}).Info("Successfully retrieved favorites")
	return result, total, nil
}

// Create creates a new favorite
func (r *FavoriteRepositoryImpl) Create(ctx context.Context, entity interface{}) error {
	r.logger.Info("Creating new favorite")

	favorite, ok := entity.(*models.Favorite)
	if !ok {
		r.logger.Error("Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.Favorite", http.StatusBadRequest, nil, nil)
	}

	if err := r.db.WithContext(ctx).Create(favorite).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to create favorite")
		return errors.NewError(errors.ErrInternalServer, "Failed to create favorite", err.Error(), http.StatusInternalServerError, nil, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": favorite.ID,
	}).Info("Successfully created favorite")
	return nil
}

// Update updates an existing favorite
func (r *FavoriteRepositoryImpl) Update(ctx context.Context, entity interface{}) error {
	r.logger.Info("Updating favorite")

	favorite, ok := entity.(*models.Favorite)
	if !ok {
		r.logger.Error("Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.Favorite", http.StatusBadRequest, nil, nil)
	}

	if err := r.db.WithContext(ctx).Save(favorite).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to update favorite")
		return errors.NewError(errors.ErrInternalServer, "Failed to update favorite", err.Error(), http.StatusInternalServerError, nil, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": favorite.ID,
	}).Info("Successfully updated favorite")
	return nil
}

// Delete deletes a favorite by its ID
func (r *FavoriteRepositoryImpl) Delete(ctx context.Context, id uint) error {
	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Deleting favorite")

	if err := r.db.WithContext(ctx).Delete(&models.Favorite{}, id).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to delete favorite")
		return errors.NewError(errors.ErrInternalServer, "Failed to delete favorite", err.Error(), http.StatusInternalServerError, map[string]interface{}{"id": id}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Successfully deleted favorite")
	return nil
}

// GetByUserID retrieves favorites by user ID
func (r *FavoriteRepositoryImpl) GetByUserID(ctx context.Context, userID uint) ([]models.Favorite, error) {
	r.logger.WithFields(map[string]interface{}{
		"userID": userID,
	}).Info("Retrieving favorites by user ID")

	var favorites []models.Favorite
	if err := r.db.WithContext(ctx).Preload("Anime").Where("user_id = ?", userID).Find(&favorites).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"userID": userID,
			"error":  err.Error(),
		}).Error("Failed to retrieve favorites")
		return nil, errors.NewError(errors.ErrInternalServer, "Failed to retrieve favorites", err.Error(), http.StatusInternalServerError, map[string]interface{}{"userID": userID}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"userID": userID,
		"count":  len(favorites),
	}).Info("Successfully retrieved favorites")
	return favorites, nil
}

// GetByAnimeID retrieves favorites by anime ID
func (r *FavoriteRepositoryImpl) GetByAnimeID(ctx context.Context, animeID uint, page, limit int) ([]models.Favorite, int64, error) {
	r.logger.WithFields(map[string]interface{}{
		"animeID": animeID,
		"page":    page,
		"limit":   limit,
	}).Info("Retrieving favorites by anime ID")

	var favorites []models.Favorite
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Favorite{}).Where("anime_id = ?", animeID).Count(&total).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"animeID": animeID,
			"error":   err.Error(),
		}).Error("Failed to count favorites")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count favorites", err.Error(), http.StatusInternalServerError, map[string]interface{}{"animeID": animeID}, err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve favorites with pagination
	if err := r.db.WithContext(ctx).Preload("User").Where("anime_id = ?", animeID).Offset(offset).Limit(limit).Find(&favorites).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"animeID": animeID,
			"error":   err.Error(),
		}).Error("Failed to retrieve favorites")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve favorites", err.Error(), http.StatusInternalServerError, map[string]interface{}{"animeID": animeID}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"animeID": animeID,
		"count":   len(favorites),
		"total":   total,
	}).Info("Successfully retrieved favorites")
	return favorites, total, nil
}

// GetByUserAndAnime retrieves a favorite by user ID and anime ID
func (r *FavoriteRepositoryImpl) GetByUserAndAnime(ctx context.Context, userID, animeID uint) (*models.Favorite, error) {
	r.logger.WithFields(map[string]interface{}{
		"userID":  userID,
		"animeID": animeID,
	}).Info("Retrieving favorite by user and anime IDs")

	var favorite models.Favorite
	if err := r.db.WithContext(ctx).Where("user_id = ? AND anime_id = ?", userID, animeID).First(&favorite).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"userID":  userID,
			"animeID": animeID,
			"error":   err.Error(),
		}).Error("Failed to retrieve favorite")
		return nil, errors.NewError(errors.ErrResourceNotFound, "Favorite not found", fmt.Sprintf("Favorite for user %d and anime %d not found", userID, animeID), http.StatusNotFound, map[string]interface{}{"userID": userID, "animeID": animeID}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": favorite.ID,
	}).Info("Successfully retrieved favorite")
	return &favorite, nil
}

// DeleteByUserAndAnime deletes a favorite by user ID and anime ID
func (r *FavoriteRepositoryImpl) DeleteByUserAndAnime(ctx context.Context, userID, animeID uint) error {
	r.logger.WithFields(map[string]interface{}{
		"userID":  userID,
		"animeID": animeID,
	}).Info("Deleting favorite by user and anime IDs")

	if err := r.db.WithContext(ctx).Where("user_id = ? AND anime_id = ?", userID, animeID).Delete(&models.Favorite{}).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"userID":  userID,
			"animeID": animeID,
			"error":   err.Error(),
		}).Error("Failed to delete favorite")
		return errors.NewError(errors.ErrInternalServer, "Failed to delete favorite", err.Error(), http.StatusInternalServerError, map[string]interface{}{"userID": userID, "animeID": animeID}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"userID":  userID,
		"animeID": animeID,
	}).Info("Successfully deleted favorite")
	return nil
}
