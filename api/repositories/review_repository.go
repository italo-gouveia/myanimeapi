package repositories

import (
	"context"
	"log"
	"myanimeapi/api/models"
	"myanimeapi/internal/db"
	"myanimeapi/internal/errors"
	"net/http"
)

// ReviewRepositoryImpl implements the ReviewRepository interface
type ReviewRepositoryImpl struct {
	db db.DBInterface
}

// NewReviewRepository creates a new ReviewRepositoryImpl instance
func NewReviewRepository(db db.DBInterface) ReviewRepository {
	return &ReviewRepositoryImpl{db: db}
}

// GetByID retrieves a review by its ID
func (r *ReviewRepositoryImpl) GetByID(ctx context.Context, id uint) (interface{}, error) {
	log.Printf("ReviewRepository.GetByID: Retrieving review with ID %d", id)

	var review models.Review
	if err := r.db.WithContext(ctx).Preload("User").Preload("Anime").First(&review, id).Error; err != nil {
		log.Printf("ReviewRepository.GetByID: Failed to retrieve review: %v", err)
		return nil, errors.NewError(errors.ErrResourceNotFound, "Review not found", err.Error(), http.StatusNotFound)
	}

	log.Printf("ReviewRepository.GetByID: Successfully retrieved review with ID %d", id)
	return &review, nil
}

// GetAll retrieves all reviews with optional pagination
func (r *ReviewRepositoryImpl) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	log.Printf("ReviewRepository.GetAll: Retrieving all reviews with page %d and limit %d", page, limit)

	var reviews []models.Review
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Review{}).Count(&total).Error; err != nil {
		log.Printf("ReviewRepository.GetAll: Failed to count reviews: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count reviews", err.Error(), http.StatusInternalServerError)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve reviews with pagination
	if err := r.db.WithContext(ctx).Preload("User").Preload("Anime").Offset(offset).Limit(limit).Find(&reviews).Error; err != nil {
		log.Printf("ReviewRepository.GetAll: Failed to retrieve reviews: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve reviews", err.Error(), http.StatusInternalServerError)
	}

	// Convert to interface slice
	result := make([]interface{}, len(reviews))
	for i, review := range reviews {
		result[i] = review
	}

	log.Printf("ReviewRepository.GetAll: Successfully retrieved %d reviews", len(reviews))
	return result, total, nil
}

// Create creates a new review
func (r *ReviewRepositoryImpl) Create(ctx context.Context, entity interface{}) error {
	log.Printf("ReviewRepository.Create: Creating new review")

	review, ok := entity.(*models.Review)
	if !ok {
		log.Printf("ReviewRepository.Create: Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.Review", http.StatusBadRequest)
	}

	if err := r.db.WithContext(ctx).Create(review).Error; err != nil {
		log.Printf("ReviewRepository.Create: Failed to create review: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to create review", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("ReviewRepository.Create: Successfully created review with ID %d", review.ID)
	return nil
}

// Update updates an existing review
func (r *ReviewRepositoryImpl) Update(ctx context.Context, entity interface{}) error {
	log.Printf("ReviewRepository.Update: Updating review")

	review, ok := entity.(*models.Review)
	if !ok {
		log.Printf("ReviewRepository.Update: Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.Review", http.StatusBadRequest)
	}

	if err := r.db.WithContext(ctx).Save(review).Error; err != nil {
		log.Printf("ReviewRepository.Update: Failed to update review: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to update review", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("ReviewRepository.Update: Successfully updated review with ID %d", review.ID)
	return nil
}

// Delete deletes a review by its ID
func (r *ReviewRepositoryImpl) Delete(ctx context.Context, id uint) error {
	log.Printf("ReviewRepository.Delete: Deleting review with ID %d", id)

	if err := r.db.WithContext(ctx).Delete(&models.Review{}, id).Error; err != nil {
		log.Printf("ReviewRepository.Delete: Failed to delete review: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to delete review", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("ReviewRepository.Delete: Successfully deleted review with ID %d", id)
	return nil
}

// GetByUserID retrieves reviews by user ID
func (r *ReviewRepositoryImpl) GetByUserID(ctx context.Context, userID uint, page, limit int) ([]models.Review, int64, error) {
	log.Printf("ReviewRepository.GetByUserID: Retrieving reviews for user ID %d, page %d, limit %d", userID, page, limit)

	var reviews []models.Review
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Review{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		log.Printf("ReviewRepository.GetByUserID: Failed to count reviews: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count reviews", err.Error(), http.StatusInternalServerError)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve reviews with pagination
	if err := r.db.WithContext(ctx).Preload("Anime").Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&reviews).Error; err != nil {
		log.Printf("ReviewRepository.GetByUserID: Failed to retrieve reviews: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve reviews", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("ReviewRepository.GetByUserID: Successfully retrieved %d reviews for user ID %d", len(reviews), userID)
	return reviews, total, nil
}

// GetByAnimeID retrieves reviews by anime ID
func (r *ReviewRepositoryImpl) GetByAnimeID(ctx context.Context, animeID uint, page, limit int) ([]models.Review, int64, error) {
	log.Printf("ReviewRepository.GetByAnimeID: Retrieving reviews for anime ID %d, page %d, limit %d", animeID, page, limit)

	var reviews []models.Review
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Review{}).Where("anime_id = ?", animeID).Count(&total).Error; err != nil {
		log.Printf("ReviewRepository.GetByAnimeID: Failed to count reviews: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count reviews", err.Error(), http.StatusInternalServerError)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve reviews with pagination
	if err := r.db.WithContext(ctx).Preload("User").Where("anime_id = ?", animeID).Offset(offset).Limit(limit).Find(&reviews).Error; err != nil {
		log.Printf("ReviewRepository.GetByAnimeID: Failed to retrieve reviews: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve reviews", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("ReviewRepository.GetByAnimeID: Successfully retrieved %d reviews for anime ID %d", len(reviews), animeID)
	return reviews, total, nil
}

// GetByUserAndAnime retrieves a review by user ID and anime ID
func (r *ReviewRepositoryImpl) GetByUserAndAnime(ctx context.Context, userID, animeID uint) (*models.Review, error) {
	log.Printf("ReviewRepository.GetByUserAndAnime: Retrieving review for user ID %d and anime ID %d", userID, animeID)

	var review models.Review
	if err := r.db.WithContext(ctx).Where("user_id = ? AND anime_id = ?", userID, animeID).First(&review).Error; err != nil {
		log.Printf("ReviewRepository.GetByUserAndAnime: Failed to retrieve review: %v", err)
		return nil, errors.NewError(errors.ErrResourceNotFound, "Review not found", err.Error(), http.StatusNotFound)
	}

	log.Printf("ReviewRepository.GetByUserAndAnime: Successfully retrieved review with ID %d", review.ID)
	return &review, nil
}
