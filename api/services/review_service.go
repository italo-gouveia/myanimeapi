package services

import (
	"context"
	"fmt"
	"log"
	"myanimeapi/api/models"
	"myanimeapi/api/repositories"
	"myanimeapi/internal/errors"
	"net/http"
	"time"
)

// ReviewService handles business logic for review operations
type ReviewService struct {
	reviewRepo repositories.ReviewRepository
	userRepo   repositories.UserRepository
	animeRepo  repositories.AnimeRepository
	storageSvc *StorageService
}

// NewReviewService creates a new ReviewService instance
func NewReviewService(reviewRepo repositories.ReviewRepository, userRepo repositories.UserRepository, animeRepo repositories.AnimeRepository, storageSvc *StorageService) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		userRepo:   userRepo,
		animeRepo:  animeRepo,
		storageSvc: storageSvc,
	}
}

// GetReviewByID retrieves a review by ID
func (s *ReviewService) GetReviewByID(ctx context.Context, id uint) (*models.Review, error) {
	log.Printf("ReviewService.GetReviewByID: Retrieving review with ID %d", id)

	reviewInterface, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		log.Printf("ReviewService.GetReviewByID: Failed to retrieve review: %v", err)
		return nil, errors.NewError(errors.ErrResourceNotFound, "Review not found", err.Error(), http.StatusNotFound)
	}

	review, ok := reviewInterface.(*models.Review)
	if !ok {
		log.Printf("ReviewService.GetReviewByID: Invalid review type returned from repository")
		return nil, errors.NewError(errors.ErrInternalServer, "Invalid review type returned from repository", "Type assertion failed", http.StatusInternalServerError)
	}

	log.Printf("ReviewService.GetReviewByID: Successfully retrieved review with ID %d", id)
	return review, nil
}

// GetAllReviews retrieves all reviews with pagination
func (s *ReviewService) GetAllReviews(ctx context.Context, page, limit int) ([]models.Review, int64, error) {
	log.Printf("ReviewService.GetAllReviews: Retrieving all reviews with page %d and limit %d", page, limit)

	reviewsInterface, total, err := s.reviewRepo.GetAll(ctx, page, limit)
	if err != nil {
		log.Printf("ReviewService.GetAllReviews: Failed to retrieve reviews: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve reviews", err.Error(), http.StatusInternalServerError)
	}

	// Convert interface slice to Review slice
	reviews := make([]models.Review, len(reviewsInterface))
	for i, reviewInterface := range reviewsInterface {
		review, ok := reviewInterface.(models.Review)
		if !ok {
			log.Printf("ReviewService.GetAllReviews: Invalid review type in slice at index %d", i)
			return nil, 0, errors.NewError(errors.ErrInternalServer, "Invalid review type in repository response", fmt.Sprintf("Type assertion failed at index %d", i), http.StatusInternalServerError)
		}
		reviews[i] = review
	}

	log.Printf("ReviewService.GetAllReviews: Successfully retrieved %d reviews", len(reviews))
	return reviews, total, nil
}

// CreateReview creates a new review
func (s *ReviewService) CreateReview(ctx context.Context, review *models.Review) error {
	log.Printf("ReviewService.CreateReview: Creating new review for user %d and anime %d", review.UserID, review.AnimeID)

	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, review.UserID)
	if err != nil {
		log.Printf("ReviewService.CreateReview: Failed to retrieve user: %v", err)
		return errors.NewError(errors.ErrResourceNotFound, "User not found", err.Error(), http.StatusNotFound)
	}

	// Check if anime exists
	_, err = s.animeRepo.GetByID(ctx, review.AnimeID)
	if err != nil {
		log.Printf("ReviewService.CreateReview: Failed to retrieve anime: %v", err)
		return errors.NewError(errors.ErrResourceNotFound, "Anime not found", err.Error(), http.StatusNotFound)
	}

	// Check if user has already reviewed this anime
	existingReview, err := s.reviewRepo.GetByUserAndAnime(ctx, review.UserID, review.AnimeID)
	if err == nil && existingReview != nil {
		log.Printf("ReviewService.CreateReview: User %d has already reviewed anime %d", review.UserID, review.AnimeID)
		return errors.NewError(errors.ErrConflict, "User has already reviewed this anime", fmt.Sprintf("User %d has an existing review for anime %d", review.UserID, review.AnimeID), http.StatusConflict)
	}

	// Validate rating
	if review.Rating < 1 || review.Rating > 10 {
		log.Printf("ReviewService.CreateReview: Invalid rating %d", review.Rating)
		return errors.NewError(errors.ErrInvalidInput, "Rating must be between 1 and 10", fmt.Sprintf("Provided rating: %d", review.Rating), http.StatusBadRequest)
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
		log.Printf("ReviewService.CreateReview: Failed to create review: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to create review", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("ReviewService.CreateReview: Successfully created review with ID %d", review.ID)
	return nil
}

// UpdateReview updates an existing review
func (s *ReviewService) UpdateReview(ctx context.Context, review *models.Review) error {
	log.Printf("ReviewService.UpdateReview: Updating review with ID %d", review.ID)

	// Check if review exists
	existingReview, err := s.reviewRepo.GetByID(ctx, review.ID)
	if err != nil {
		log.Printf("ReviewService.UpdateReview: Failed to retrieve review: %v", err)
		return errors.NewError(errors.ErrResourceNotFound, "Review not found", err.Error(), http.StatusNotFound)
	}

	// Update the review
	existingReviewModel, ok := existingReview.(*models.Review)
	if !ok {
		log.Printf("ReviewService.UpdateReview: Invalid review type returned from repository")
		return errors.NewError(errors.ErrInternalServer, "Invalid review type returned from repository", "Type assertion failed", http.StatusInternalServerError)
	}

	review.CreatedAt = existingReviewModel.CreatedAt
	review.UpdatedAt = time.Now()

	// Validate rating
	if review.Rating < 1 || review.Rating > 10 {
		return errors.NewError(errors.ErrInvalidInput, "Invalid rating", "Rating must be between 1 and 10", http.StatusBadRequest)
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
		log.Printf("ReviewService.UpdateReview: Failed to update review: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to update review", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("ReviewService.UpdateReview: Successfully updated review with ID %d", review.ID)
	return nil
}

// DeleteReview deletes a review by ID
func (s *ReviewService) DeleteReview(ctx context.Context, id uint) error {
	log.Printf("ReviewService.DeleteReview: Deleting review with ID %d", id)

	// Get the review first to handle media attachments
	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		log.Printf("ReviewService.DeleteReview: Failed to retrieve review: %v", err)
		return errors.NewError(errors.ErrResourceNotFound, "Review not found", err.Error(), http.StatusNotFound)
	}

	reviewModel, ok := review.(*models.Review)
	if !ok {
		log.Printf("ReviewService.DeleteReview: Invalid review type returned from repository")
		return errors.NewError(errors.ErrInternalServer, "Invalid review type returned from repository", "Type assertion failed", http.StatusInternalServerError)
	}

	// Delete associated media files
	for _, attachment := range reviewModel.MediaAttachments {
		if err := s.storageSvc.DeleteFile(ctx, attachment.URL); err != nil {
			log.Printf("ReviewService.DeleteReview: Failed to delete media file %s: %v", attachment.URL, err)
			// Continue with deletion even if file deletion fails
		}
	}

	// Delete review
	if err := s.reviewRepo.Delete(ctx, id); err != nil {
		log.Printf("ReviewService.DeleteReview: Failed to delete review: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to delete review", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("ReviewService.DeleteReview: Successfully deleted review with ID %d", id)
	return nil
}

// GetReviewsByUserID retrieves reviews by user ID with pagination
func (s *ReviewService) GetReviewsByUserID(ctx context.Context, userID uint, page, limit int) ([]models.Review, int64, error) {
	log.Printf("ReviewService.GetReviewsByUserID: Retrieving reviews for user %d", userID)

	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		log.Printf("ReviewService.GetReviewsByUserID: Failed to retrieve user: %v", err)
		return nil, 0, errors.NewError(errors.ErrResourceNotFound, "User not found", err.Error(), http.StatusNotFound)
	}

	reviews, total, err := s.reviewRepo.GetByUserID(ctx, userID, page, limit)
	if err != nil {
		log.Printf("ReviewService.GetReviewsByUserID: Failed to retrieve reviews: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve reviews for user", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("ReviewService.GetReviewsByUserID: Successfully retrieved %d reviews", len(reviews))
	return reviews, total, nil
}

// GetReviewsByAnimeID retrieves reviews by anime ID with pagination
func (s *ReviewService) GetReviewsByAnimeID(ctx context.Context, animeID uint, page, limit int) ([]models.Review, int64, error) {
	log.Printf("ReviewService.GetReviewsByAnimeID: Retrieving reviews for anime %d", animeID)

	// Check if anime exists
	_, err := s.animeRepo.GetByID(ctx, animeID)
	if err != nil {
		log.Printf("ReviewService.GetReviewsByAnimeID: Failed to retrieve anime: %v", err)
		return nil, 0, errors.NewError(errors.ErrResourceNotFound, "Anime not found", err.Error(), http.StatusNotFound)
	}

	reviews, total, err := s.reviewRepo.GetByAnimeID(ctx, animeID, page, limit)
	if err != nil {
		log.Printf("ReviewService.GetReviewsByAnimeID: Failed to retrieve reviews: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve reviews for anime", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("ReviewService.GetReviewsByAnimeID: Successfully retrieved %d reviews", len(reviews))
	return reviews, total, nil
}
