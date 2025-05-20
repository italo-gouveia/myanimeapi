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

// ReviewRepositoryImpl implements the ReviewRepository interface
type ReviewRepositoryImpl struct {
	db     db.DBInterface
	logger *logger.Logger
}

// NewReviewRepository creates a new ReviewRepositoryImpl instance
func NewReviewRepository(db db.DBInterface) ReviewRepository {
	return &ReviewRepositoryImpl{
		db:     db,
		logger: logger.New(),
	}
}

// GetByID retrieves a review by its ID
func (r *ReviewRepositoryImpl) GetByID(ctx context.Context, id uint) (interface{}, error) {
	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Retrieving review by ID")

	var review models.Review
	if err := r.db.WithContext(ctx).Preload("User").Preload("Anime").First(&review, id).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to retrieve review")
		return nil, errors.NewError(errors.ErrResourceNotFound, "Review not found", fmt.Sprintf("Review with ID %d not found", id), http.StatusNotFound, map[string]interface{}{"id": id}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Successfully retrieved review")
	return &review, nil
}

// GetAll retrieves all reviews with optional pagination
func (r *ReviewRepositoryImpl) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	r.logger.WithFields(map[string]interface{}{
		"page":  page,
		"limit": limit,
	}).Info("Retrieving all reviews")

	var reviews []models.Review
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Review{}).Count(&total).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to count reviews")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count reviews", err.Error(), http.StatusInternalServerError, nil, err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve reviews with pagination
	if err := r.db.WithContext(ctx).Preload("User").Preload("Anime").Offset(offset).Limit(limit).Find(&reviews).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to retrieve reviews")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve reviews", err.Error(), http.StatusInternalServerError, nil, err)
	}

	// Convert to interface slice
	result := make([]interface{}, len(reviews))
	for i, review := range reviews {
		result[i] = &review
	}

	r.logger.WithFields(map[string]interface{}{
		"count": len(reviews),
		"total": total,
	}).Info("Successfully retrieved reviews")
	return result, total, nil
}

// Create creates a new review
func (r *ReviewRepositoryImpl) Create(ctx context.Context, entity interface{}) error {
	r.logger.Info("Creating new review")

	review, ok := entity.(*models.Review)
	if !ok {
		r.logger.Error("Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.Review", http.StatusBadRequest, nil, nil)
	}

	if err := r.db.WithContext(ctx).Create(review).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to create review")
		return errors.NewError(errors.ErrInternalServer, "Failed to create review", err.Error(), http.StatusInternalServerError, nil, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": review.ID,
	}).Info("Successfully created review")
	return nil
}

// Update updates an existing review
func (r *ReviewRepositoryImpl) Update(ctx context.Context, entity interface{}) error {
	r.logger.Info("Updating review")

	review, ok := entity.(*models.Review)
	if !ok {
		r.logger.Error("Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.Review", http.StatusBadRequest, nil, nil)
	}

	if err := r.db.WithContext(ctx).Save(review).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to update review")
		return errors.NewError(errors.ErrInternalServer, "Failed to update review", err.Error(), http.StatusInternalServerError, nil, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": review.ID,
	}).Info("Successfully updated review")
	return nil
}

// Delete deletes a review by its ID
func (r *ReviewRepositoryImpl) Delete(ctx context.Context, id uint) error {
	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Deleting review")

	if err := r.db.WithContext(ctx).Delete(&models.Review{}, id).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to delete review")
		return errors.NewError(errors.ErrInternalServer, "Failed to delete review", err.Error(), http.StatusInternalServerError, map[string]interface{}{"id": id}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Successfully deleted review")
	return nil
}

// GetByUserID retrieves reviews by user ID
func (r *ReviewRepositoryImpl) GetByUserID(ctx context.Context, userID uint, page, limit int) ([]models.Review, int64, error) {
	r.logger.WithFields(map[string]interface{}{
		"user_id": userID,
		"page":    page,
		"limit":   limit,
	}).Info("Retrieving reviews by user ID")

	var reviews []models.Review
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Review{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to count reviews")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count reviews", err.Error(), http.StatusInternalServerError, map[string]interface{}{"user_id": userID}, err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve reviews with pagination
	if err := r.db.WithContext(ctx).Preload("Anime").Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&reviews).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to retrieve reviews")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve reviews", err.Error(), http.StatusInternalServerError, map[string]interface{}{"user_id": userID}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"count":   len(reviews),
		"total":   total,
		"user_id": userID,
	}).Info("Successfully retrieved reviews by user ID")
	return reviews, total, nil
}

// GetByAnimeID retrieves reviews by anime ID
func (r *ReviewRepositoryImpl) GetByAnimeID(ctx context.Context, animeID uint, page, limit int) ([]models.Review, int64, error) {
	r.logger.WithFields(map[string]interface{}{
		"anime_id": animeID,
		"page":     page,
		"limit":    limit,
	}).Info("Retrieving reviews by anime ID")

	var reviews []models.Review
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Review{}).Where("anime_id = ?", animeID).Count(&total).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to count reviews")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count reviews", err.Error(), http.StatusInternalServerError, map[string]interface{}{"anime_id": animeID}, err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve reviews with pagination
	if err := r.db.WithContext(ctx).Preload("User").Where("anime_id = ?", animeID).Offset(offset).Limit(limit).Find(&reviews).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to retrieve reviews")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve reviews", err.Error(), http.StatusInternalServerError, map[string]interface{}{"anime_id": animeID}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"count":    len(reviews),
		"total":    total,
		"anime_id": animeID,
	}).Info("Successfully retrieved reviews by anime ID")
	return reviews, total, nil
}

// GetByUserAndAnime retrieves a review by user ID and anime ID
func (r *ReviewRepositoryImpl) GetByUserAndAnime(ctx context.Context, userID, animeID uint) (*models.Review, error) {
	r.logger.WithFields(map[string]interface{}{
		"user_id":  userID,
		"anime_id": animeID,
	}).Info("Retrieving review by user and anime")

	var review models.Review
	if err := r.db.WithContext(ctx).Where("user_id = ? AND anime_id = ?", userID, animeID).First(&review).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to retrieve review")
		return nil, errors.NewError(errors.ErrResourceNotFound, "Review not found", fmt.Sprintf("Review not found for user %d and anime %d", userID, animeID), http.StatusNotFound, map[string]interface{}{
			"user_id":  userID,
			"anime_id": animeID,
		}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id":       review.ID,
		"user_id":  userID,
		"anime_id": animeID,
	}).Info("Successfully retrieved review")
	return &review, nil
}
