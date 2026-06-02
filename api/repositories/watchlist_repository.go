package repositories

import (
	"context"
	"fmt"
	"myanimeapi/api/models"
	"myanimeapi/internal/db"
	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"
	"net/http"

	gormerrors "gorm.io/gorm"
)

// WatchlistRepositoryImpl implements the WatchlistRepository interface
type WatchlistRepositoryImpl struct {
	db     db.DBInterface
	logger *logger.Logger
}

// NewWatchlistRepository creates a new WatchlistRepositoryImpl instance
func NewWatchlistRepository(db db.DBInterface) WatchlistRepository {
	return &WatchlistRepositoryImpl{
		db:     db,
		logger: logger.New(),
	}
}

// GetByID retrieves a watchlist entry by its ID
func (r *WatchlistRepositoryImpl) GetByID(ctx context.Context, id uint) (interface{}, error) {
	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Retrieving watchlist entry by ID")

	var entry models.WatchlistEntry
	if err := r.db.WithContext(ctx).Preload("Anime").First(&entry, id).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to retrieve watchlist entry")
		return nil, errors.NewError(errors.ErrResourceNotFound, "Watchlist entry not found", fmt.Sprintf("Watchlist entry with ID %d not found", id), http.StatusNotFound, map[string]interface{}{"id": id}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Successfully retrieved watchlist entry")
	return &entry, nil
}

// GetAll retrieves all watchlist entries with optional pagination
func (r *WatchlistRepositoryImpl) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	r.logger.WithFields(map[string]interface{}{
		"page":  page,
		"limit": limit,
	}).Info("Retrieving all watchlist entries")

	var entries []models.WatchlistEntry
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.WatchlistEntry{}).Count(&total).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to count watchlist entries")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count watchlist entries", err.Error(), http.StatusInternalServerError, nil, err)
	}

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).Preload("Anime").Offset(offset).Limit(limit).Find(&entries).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to retrieve watchlist entries")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve watchlist entries", err.Error(), http.StatusInternalServerError, nil, err)
	}

	result := make([]interface{}, len(entries))
	for i, entry := range entries {
		e := entry
		result[i] = &e
	}

	r.logger.WithFields(map[string]interface{}{
		"count": len(entries),
		"total": total,
	}).Info("Successfully retrieved watchlist entries")
	return result, total, nil
}

// Create creates a new watchlist entry
func (r *WatchlistRepositoryImpl) Create(ctx context.Context, entity interface{}) error {
	r.logger.Info("Creating new watchlist entry")

	entry, ok := entity.(*models.WatchlistEntry)
	if !ok {
		r.logger.Error("Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.WatchlistEntry", http.StatusBadRequest, nil, nil)
	}

	if err := r.db.WithContext(ctx).Create(entry).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to create watchlist entry")
		return errors.NewError(errors.ErrInternalServer, "Failed to create watchlist entry", err.Error(), http.StatusInternalServerError, nil, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": entry.ID,
	}).Info("Successfully created watchlist entry")
	return nil
}

// Update updates an existing watchlist entry
func (r *WatchlistRepositoryImpl) Update(ctx context.Context, entity interface{}) error {
	r.logger.Info("Updating watchlist entry")

	entry, ok := entity.(*models.WatchlistEntry)
	if !ok {
		r.logger.Error("Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.WatchlistEntry", http.StatusBadRequest, nil, nil)
	}

	if err := r.db.WithContext(ctx).Save(entry).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to update watchlist entry")
		return errors.NewError(errors.ErrInternalServer, "Failed to update watchlist entry", err.Error(), http.StatusInternalServerError, nil, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": entry.ID,
	}).Info("Successfully updated watchlist entry")
	return nil
}

// Delete deletes a watchlist entry by its ID
func (r *WatchlistRepositoryImpl) Delete(ctx context.Context, id uint) error {
	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Deleting watchlist entry")

	if err := r.db.WithContext(ctx).Delete(&models.WatchlistEntry{}, id).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to delete watchlist entry")
		return errors.NewError(errors.ErrInternalServer, "Failed to delete watchlist entry", err.Error(), http.StatusInternalServerError, map[string]interface{}{"id": id}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Successfully deleted watchlist entry")
	return nil
}

// GetByUserID retrieves watchlist entries by user ID
func (r *WatchlistRepositoryImpl) GetByUserID(ctx context.Context, userID uint) ([]models.WatchlistEntry, error) {
	r.logger.WithFields(map[string]interface{}{
		"userID": userID,
	}).Info("Retrieving watchlist entries by user ID")

	var entries []models.WatchlistEntry
	if err := r.db.WithContext(ctx).Preload("Anime").Where("user_id = ?", userID).Find(&entries).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"userID": userID,
			"error":  err.Error(),
		}).Error("Failed to retrieve watchlist entries")
		return nil, errors.NewError(errors.ErrInternalServer, "Failed to retrieve watchlist entries", err.Error(), http.StatusInternalServerError, map[string]interface{}{"userID": userID}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"userID": userID,
		"count":  len(entries),
	}).Info("Successfully retrieved watchlist entries")
	return entries, nil
}

// GetByUserAndAnime retrieves a watchlist entry by user ID and anime ID
func (r *WatchlistRepositoryImpl) GetByUserAndAnime(ctx context.Context, userID, animeID uint) (*models.WatchlistEntry, error) {
	r.logger.WithFields(map[string]interface{}{
		"userID":  userID,
		"animeID": animeID,
	}).Info("Retrieving watchlist entry by user and anime IDs")

	var entry models.WatchlistEntry
	if err := r.db.WithContext(ctx).Where("user_id = ? AND anime_id = ?", userID, animeID).First(&entry).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"userID":  userID,
			"animeID": animeID,
			"error":   err.Error(),
		}).Error("Failed to retrieve watchlist entry")
		return nil, errors.NewError(errors.ErrResourceNotFound, "Watchlist entry not found", fmt.Sprintf("Watchlist entry for user %d and anime %d not found", userID, animeID), http.StatusNotFound, map[string]interface{}{"userID": userID, "animeID": animeID}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": entry.ID,
	}).Info("Successfully retrieved watchlist entry")
	return &entry, nil
}

// UpsertByUserAndAnime creates or updates a watchlist entry for a user/anime pair
func (r *WatchlistRepositoryImpl) UpsertByUserAndAnime(ctx context.Context, userID, animeID uint, status models.WatchlistStatus) (*models.WatchlistEntry, error) {
	r.logger.WithFields(map[string]interface{}{
		"userID":  userID,
		"animeID": animeID,
		"status":  status,
	}).Info("Upserting watchlist entry")

	var entry models.WatchlistEntry
	err := r.db.WithContext(ctx).Where("user_id = ? AND anime_id = ?", userID, animeID).First(&entry).Error

	if err != nil {
		if err == gormerrors.ErrRecordNotFound {
			// Create new entry
			entry = models.WatchlistEntry{
				UserID:  userID,
				AnimeID: animeID,
				Status:  status,
			}
			if createErr := r.db.WithContext(ctx).Create(&entry).Error; createErr != nil {
				r.logger.WithFields(map[string]interface{}{
					"userID":  userID,
					"animeID": animeID,
					"error":   createErr.Error(),
				}).Error("Failed to create watchlist entry")
				return nil, errors.NewError(errors.ErrInternalServer, "Failed to create watchlist entry", createErr.Error(), http.StatusInternalServerError, map[string]interface{}{"userID": userID, "animeID": animeID}, createErr)
			}
		} else {
			r.logger.WithFields(map[string]interface{}{
				"userID":  userID,
				"animeID": animeID,
				"error":   err.Error(),
			}).Error("Failed to query watchlist entry")
			return nil, errors.NewError(errors.ErrInternalServer, "Failed to query watchlist entry", err.Error(), http.StatusInternalServerError, map[string]interface{}{"userID": userID, "animeID": animeID}, err)
		}
	} else {
		// Update existing entry
		entry.Status = status
		if saveErr := r.db.WithContext(ctx).Save(&entry).Error; saveErr != nil {
			r.logger.WithFields(map[string]interface{}{
				"id":    entry.ID,
				"error": saveErr.Error(),
			}).Error("Failed to update watchlist entry")
			return nil, errors.NewError(errors.ErrInternalServer, "Failed to update watchlist entry", saveErr.Error(), http.StatusInternalServerError, map[string]interface{}{"id": entry.ID}, saveErr)
		}
	}

	// Reload with preloads
	if err := r.db.WithContext(ctx).Preload("Anime").First(&entry, entry.ID).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"id":    entry.ID,
			"error": err.Error(),
		}).Error("Failed to reload watchlist entry")
		return nil, errors.NewError(errors.ErrInternalServer, "Failed to reload watchlist entry", err.Error(), http.StatusInternalServerError, map[string]interface{}{"id": entry.ID}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id":      entry.ID,
		"userID":  userID,
		"animeID": animeID,
		"status":  status,
	}).Info("Successfully upserted watchlist entry")
	return &entry, nil
}

// DeleteByUserAndAnime deletes a watchlist entry by user ID and anime ID
func (r *WatchlistRepositoryImpl) DeleteByUserAndAnime(ctx context.Context, userID, animeID uint) error {
	r.logger.WithFields(map[string]interface{}{
		"userID":  userID,
		"animeID": animeID,
	}).Info("Deleting watchlist entry by user and anime IDs")

	if err := r.db.WithContext(ctx).Where("user_id = ? AND anime_id = ?", userID, animeID).Delete(&models.WatchlistEntry{}).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"userID":  userID,
			"animeID": animeID,
			"error":   err.Error(),
		}).Error("Failed to delete watchlist entry")
		return errors.NewError(errors.ErrInternalServer, "Failed to delete watchlist entry", err.Error(), http.StatusInternalServerError, map[string]interface{}{"userID": userID, "animeID": animeID}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"userID":  userID,
		"animeID": animeID,
	}).Info("Successfully deleted watchlist entry")
	return nil
}
