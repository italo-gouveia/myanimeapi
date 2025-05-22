package services

import (
	"context"
	"fmt"
	"myanimeapi/api/models"
	"myanimeapi/api/repositories"
	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"
	"net/http"
	"time"
)

// ReviewServiceInterface defines the interface for review service operations
type ReviewServiceInterface interface {
	// GetReviewByID retrieves a review by ID
	GetReviewByID(ctx context.Context, id uint) (*models.Review, error)
	// GetAllReviews retrieves all reviews with pagination
	GetAllReviews(ctx context.Context, page, limit int) ([]models.Review, int64, error)
	// CreateReview creates a new review
	CreateReview(ctx context.Context, review *models.Review) error
	// UpdateReview updates an existing review
	UpdateReview(ctx context.Context, review *models.Review) error
	// DeleteReview deletes a review by ID
	DeleteReview(ctx context.Context, id uint) error
	// GetReviewsByUserID retrieves reviews by user ID with pagination
	GetReviewsByUserID(ctx context.Context, userID uint, page, limit int) ([]models.Review, int64, error)
	// GetReviewsByAnimeID retrieves reviews by anime ID with pagination
	GetReviewsByAnimeID(ctx context.Context, animeID uint, page, limit int) ([]models.Review, int64, error)
}

// ReviewService handles business logic for review operations
// It implements the ReviewServiceInterface.
type ReviewService struct {
	reviewRepo repositories.ReviewRepository
	userRepo   repositories.UserRepository
	animeRepo  repositories.AnimeRepository
	storageSvc *StorageService
	logger     *logger.Logger
}

// NewReviewService creates a new ReviewService instance
// It returns a ReviewServiceInterface implementation.
func NewReviewService(reviewRepo repositories.ReviewRepository, userRepo repositories.UserRepository, animeRepo repositories.AnimeRepository, storageSvc *StorageService) ReviewServiceInterface {
	return &ReviewService{
		reviewRepo: reviewRepo,
		userRepo:   userRepo,
		animeRepo:  animeRepo,
		storageSvc: storageSvc,
		logger:     logger.New(),
	}
}

// GetReviewByID retrieves a review by ID
func (s *ReviewService) GetReviewByID(ctx context.Context, id uint) (*models.Review, error) {
	s.logger.WithField("review_id", id).Info("Retrieving review")

	reviewInterface, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"review_id": id,
			"error":     err.Error(),
		}).Error("Failed to retrieve review")
		return nil, errors.NewError(errors.ErrResourceNotFound, "Review not found", err.Error(), http.StatusNotFound,
			map[string]interface{}{
				"review_id": id,
			},
			err)
	}

	review, ok := reviewInterface.(*models.Review)
	if !ok {
		s.logger.WithField("review_id", id).Error("Invalid review type returned from repository")
		return nil, errors.NewError(errors.ErrInternalServer, "Invalid review type returned from repository", "Type assertion failed", http.StatusInternalServerError,
			map[string]interface{}{
				"review_id": id,
			},
			nil)
	}

	s.logger.WithField("review_id", id).Info("Successfully retrieved review")
	return review, nil
}

// GetAllReviews retrieves all reviews with pagination
func (s *ReviewService) GetAllReviews(ctx context.Context, page, limit int) ([]models.Review, int64, error) {
	s.logger.WithFields(map[string]interface{}{
		"page":  page,
		"limit": limit,
	}).Info("Retrieving all reviews")

	reviewsInterface, total, err := s.reviewRepo.GetAll(ctx, page, limit)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"page":  page,
			"limit": limit,
			"error": err.Error(),
		}).Error("Failed to retrieve reviews")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve reviews", err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"page":  page,
				"limit": limit,
			},
			err)
	}

	// Convert interface slice to Review slice
	reviews := make([]models.Review, len(reviewsInterface))
	for i, reviewInterface := range reviewsInterface {
		review, ok := reviewInterface.(models.Review)
		if !ok {
			s.logger.WithFields(map[string]interface{}{
				"index": i,
				"error": "Type assertion failed",
			}).Error("Invalid review type in slice")
			return nil, 0, errors.NewError(errors.ErrInternalServer, "Invalid review type in repository response", fmt.Sprintf("Type assertion failed at index %d", i), http.StatusInternalServerError,
				map[string]interface{}{
					"index": i,
				},
				nil)
		}
		reviews[i] = review
	}

	s.logger.WithField("count", len(reviews)).Info("Successfully retrieved reviews")
	return reviews, total, nil
}

// CreateReview creates a new review
func (s *ReviewService) CreateReview(ctx context.Context, review *models.Review) error {
	s.logger.WithFields(map[string]interface{}{
		"user_id":  review.UserID,
		"anime_id": review.AnimeID,
	}).Info("Creating new review")

	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, review.UserID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id": review.UserID,
			"error":   err.Error(),
		}).Error("Failed to retrieve user")
		return errors.NewError(errors.ErrResourceNotFound, "User not found", err.Error(), http.StatusNotFound,
			map[string]interface{}{
				"user_id": review.UserID,
			},
			err)
	}

	// Check if anime exists
	_, err = s.animeRepo.GetByID(ctx, review.AnimeID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"anime_id": review.AnimeID,
			"error":    err.Error(),
		}).Error("Failed to retrieve anime")
		return errors.NewError(errors.ErrResourceNotFound, "Anime not found", err.Error(), http.StatusNotFound,
			map[string]interface{}{
				"anime_id": review.AnimeID,
			},
			err)
	}

	// Check if user has already reviewed this anime
	existingReview, err := s.reviewRepo.GetByUserAndAnime(ctx, review.UserID, review.AnimeID)
	if err == nil && existingReview != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id":  review.UserID,
			"anime_id": review.AnimeID,
		}).Warning("User has already reviewed this anime")
		return errors.NewError(errors.ErrConflict, "User has already reviewed this anime", fmt.Sprintf("User %d has an existing review for anime %d", review.UserID, review.AnimeID), http.StatusConflict,
			map[string]interface{}{
				"user_id":  review.UserID,
				"anime_id": review.AnimeID,
			},
			nil)
	}

	// Validate rating
	if review.Rating < 1 || review.Rating > 10 {
		s.logger.WithFields(map[string]interface{}{
			"user_id":  review.UserID,
			"anime_id": review.AnimeID,
			"rating":   review.Rating,
		}).Error("Invalid rating")
		return errors.NewError(errors.ErrInvalidInput, "Rating must be between 1 and 10", fmt.Sprintf("Provided rating: %d", review.Rating), http.StatusBadRequest,
			map[string]interface{}{
				"user_id":  review.UserID,
				"anime_id": review.AnimeID,
				"rating":   review.Rating,
			},
			nil)
	}

	// Set timestamps
	review.CreatedAt = time.Now()
	review.UpdatedAt = time.Now()

	// Set timestamps for media attachments
	for i := range review.MediaAttachments {
		review.MediaAttachments[i].CreatedAt = time.Now()
		review.MediaAttachments[i].UpdatedAt = time.Now()
	}

	// Create review
	if err := s.reviewRepo.Create(ctx, review); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id":  review.UserID,
			"anime_id": review.AnimeID,
			"error":    err.Error(),
		}).Error("Failed to create review")
		return errors.NewError(errors.ErrInternalServer, "Failed to create review", err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"user_id":  review.UserID,
				"anime_id": review.AnimeID,
			},
			err)
	}

	s.logger.WithFields(map[string]interface{}{
		"review_id": review.ID,
		"user_id":   review.UserID,
		"anime_id":  review.AnimeID,
	}).Info("Successfully created review")
	return nil
}

// UpdateReview updates an existing review
func (s *ReviewService) UpdateReview(ctx context.Context, review *models.Review) error {
	s.logger.WithField("review_id", review.ID).Info("Updating review")

	// Check if review exists
	existingReview, err := s.reviewRepo.GetByID(ctx, review.ID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"review_id": review.ID,
			"error":     err.Error(),
		}).Error("Failed to retrieve review")
		return errors.NewError(errors.ErrResourceNotFound, "Review not found", err.Error(), http.StatusNotFound,
			map[string]interface{}{
				"review_id": review.ID,
			},
			err)
	}

	// Update the review
	existingReviewModel, ok := existingReview.(*models.Review)
	if !ok {
		s.logger.WithField("review_id", review.ID).Error("Invalid review type returned from repository")
		return errors.NewError(errors.ErrInternalServer, "Invalid review type returned from repository", "Type assertion failed", http.StatusInternalServerError,
			map[string]interface{}{
				"review_id": review.ID,
			},
			nil)
	}

	review.CreatedAt = existingReviewModel.CreatedAt
	review.UpdatedAt = time.Now()

	// Validate rating
	if review.Rating < 1 || review.Rating > 10 {
		s.logger.WithFields(map[string]interface{}{
			"review_id": review.ID,
			"rating":    review.Rating,
		}).Error("Invalid rating")
		return errors.NewError(errors.ErrInvalidInput, "Invalid rating", "Rating must be between 1 and 10", http.StatusBadRequest,
			map[string]interface{}{
				"review_id": review.ID,
				"rating":    review.Rating,
			},
			nil)
	}

	// Set timestamps for new media attachments
	for i := range review.MediaAttachments {
		if review.MediaAttachments[i].ID == 0 { // New attachment
			review.MediaAttachments[i].CreatedAt = time.Now()
			review.MediaAttachments[i].UpdatedAt = time.Now()
		}
	}

	// Update review
	if err := s.reviewRepo.Update(ctx, review); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"review_id": review.ID,
			"error":     err.Error(),
		}).Error("Failed to update review")
		return errors.NewError(errors.ErrInternalServer, "Failed to update review", err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"review_id": review.ID,
			},
			err)
	}

	s.logger.WithField("review_id", review.ID).Info("Successfully updated review")
	return nil
}

// DeleteReview deletes a review by ID
func (s *ReviewService) DeleteReview(ctx context.Context, id uint) error {
	s.logger.WithField("review_id", id).Info("Deleting review")

	// Get the review first to handle media attachments
	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"review_id": id,
			"error":     err.Error(),
		}).Error("Failed to retrieve review")
		return errors.NewError(errors.ErrResourceNotFound, "Review not found", err.Error(), http.StatusNotFound,
			map[string]interface{}{
				"review_id": id,
			},
			err)
	}

	reviewModel, ok := review.(*models.Review)
	if !ok {
		s.logger.WithField("review_id", id).Error("Invalid review type returned from repository")
		return errors.NewError(errors.ErrInternalServer, "Invalid review type returned from repository", "Type assertion failed", http.StatusInternalServerError,
			map[string]interface{}{
				"review_id": id,
			},
			nil)
	}

	// Delete associated media files
	for _, attachment := range reviewModel.MediaAttachments {
		if err := s.storageSvc.DeleteFile(ctx, attachment.URL); err != nil {
			s.logger.WithFields(map[string]interface{}{
				"review_id": id,
				"url":       attachment.URL,
				"error":     err.Error(),
			}).Error("Failed to delete media file")
			// Continue with deletion even if file deletion fails
		}
	}

	// Delete review
	if err := s.reviewRepo.Delete(ctx, id); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"review_id": id,
			"error":     err.Error(),
		}).Error("Failed to delete review")
		return errors.NewError(errors.ErrInternalServer, "Failed to delete review", err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"review_id": id,
			},
			err)
	}

	s.logger.WithField("review_id", id).Info("Successfully deleted review")
	return nil
}

// GetReviewsByUserID retrieves reviews by user ID with pagination
func (s *ReviewService) GetReviewsByUserID(ctx context.Context, userID uint, page, limit int) ([]models.Review, int64, error) {
	s.logger.WithFields(map[string]interface{}{
		"user_id": userID,
		"page":    page,
		"limit":   limit,
	}).Info("Retrieving reviews for user")

	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to retrieve user")
		return nil, 0, errors.NewError(errors.ErrResourceNotFound, "User not found", err.Error(), http.StatusNotFound,
			map[string]interface{}{
				"user_id": userID,
			},
			err)
	}

	reviews, total, err := s.reviewRepo.GetByUserID(ctx, userID, page, limit)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to retrieve reviews for user")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve reviews for user", err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"user_id": userID,
			},
			err)
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id": userID,
		"count":   len(reviews),
	}).Info("Successfully retrieved reviews for user")
	return reviews, total, nil
}

// GetReviewsByAnimeID retrieves reviews by anime ID with pagination
func (s *ReviewService) GetReviewsByAnimeID(ctx context.Context, animeID uint, page, limit int) ([]models.Review, int64, error) {
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
		}).Error("Failed to retrieve anime")
		return nil, 0, errors.NewError(errors.ErrResourceNotFound, "Anime not found", err.Error(), http.StatusNotFound,
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
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve reviews for anime", err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"anime_id": animeID,
			},
			err)
	}

	s.logger.WithFields(map[string]interface{}{
		"anime_id": animeID,
		"count":    len(reviews),
	}).Info("Successfully retrieved reviews for anime")
	return reviews, total, nil
}
