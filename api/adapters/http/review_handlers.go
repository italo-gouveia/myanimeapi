// api/adapters/http/review_handlers.go
// Package httphandler provides HTTP handlers for review-related routes in the MyAnimeAPI application.
package httphandler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"myanimeapi/api/middleware"
	"myanimeapi/api/models"
	"myanimeapi/api/services"
	"myanimeapi/api/utils"
	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"

	"github.com/gorilla/mux"
)

// ReviewHandler defines the handlers for review-related routes.
type ReviewHandler struct {
	reviewService services.ReviewServiceInterface
	storageSvc    services.StorageServiceInterface
	logger        *logger.Logger
}

// NewReviewHandler creates a new instance of ReviewHandler.
func NewReviewHandler(reviewService services.ReviewServiceInterface, storageSvc services.StorageServiceInterface) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
		storageSvc:    storageSvc,
		logger:        logger.New(),
	}
}

// mediaTypeFromExt maps a lowercase file extension to a media type string.
func mediaTypeFromExt(ext string) string {
	switch ext {
	case ".mp4", ".webm", ".mov":
		return "video"
	default:
		return "image"
	}
}

// processMediaFiles reads all files from the multipart form and uploads them
// via the storage service. It returns the list of created attachments, or an
// error if any upload fails.
func (h *ReviewHandler) processMediaFiles(r *http.Request) ([]models.MediaAttachment, error) {
	form := r.MultipartForm
	if form == nil || form.File == nil {
		return nil, nil
	}
	var attachments []models.MediaAttachment
	for _, files := range form.File {
		for _, file := range files {
			mediaType := mediaTypeFromExt(strings.ToLower(filepath.Ext(file.Filename)))
			fileURL, err := h.storageSvc.SaveFile(file, mediaType)
			if err != nil {
				return nil, err
			}
			attachments = append(attachments, models.MediaAttachment{Type: mediaType, URL: fileURL})
		}
	}
	return attachments, nil
}

// toReviewResponse converts a Review model to a ReviewResponse DTO.
func toReviewResponse(review *models.Review) models.ReviewResponse {
	return models.ReviewResponse{
		ID:        review.ID,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
		UserID:    review.UserID,
		AnimeID:   review.AnimeID,
		Content:   review.Content,
		Rating:    review.Rating,
		User: models.UserResponse{
			ID:       review.User.ID,
			Username: review.User.Username,
			IsAdmin:  review.User.IsAdmin,
		},
		Anime: models.AnimeResponse{
			ID:          review.Anime.ID,
			Title:       review.Anime.Title,
			Description: review.Anime.Description,
			Rating:      review.Anime.Rating,
		},
		MediaAttachments: review.MediaAttachments,
	}
}

// RegisterReviewRoutes registers all review-related routes with a *mux.Router.
func (h *ReviewHandler) RegisterReviewRoutes(router *mux.Router) {
	// Public routes (no authentication required)
	router.HandleFunc("/reviews/{id}", h.GetReviewHandler).Methods("GET")

	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/reviews").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware) // Apply authentication middleware

	// Protected routes with payload validation
	protectedRouter.Handle("", middleware.ValidateAndSanitizePayload(models.ReviewCreateRequest{})(http.HandlerFunc(h.CreateReviewHandler))).Methods("POST")
	protectedRouter.Handle("/{id}", middleware.ValidateAndSanitizePayload(models.ReviewUpdateRequest{})(http.HandlerFunc(h.UpdateReviewHandler))).Methods("PUT")
	protectedRouter.HandleFunc("/{id}", h.DeleteReviewHandler).Methods("DELETE")
}

// fetchReview loads a review by id, writing a 404 response and returning nil on failure.
func (h *ReviewHandler) fetchReview(w http.ResponseWriter, r *http.Request, id uint) *models.Review {
	review, err := h.reviewService.GetReviewByID(r.Context(), id)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{"review_id": id, "err": err}).Warning("Failed to get review")
		errors.WriteErrorResponse(w, http.StatusNotFound, errors.ErrResourceNotFound, "Review not found", fmt.Sprintf("Review with ID '%d' not found.", id), map[string]interface{}{
			"id": id,
		})
		return nil
	}
	return review
}

// GetReviewHandler retrieves a review by its ID.
func (h *ReviewHandler) GetReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	h.logger.WithField("raw_id", idStr).Debug("Raw ID string")

	id, err := utils.ValidateID(idStr)
	if err != nil {
		h.logger.WithField("err", err).Warning("Invalid ID format")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", "The provided ID is not a valid unsigned integer.", map[string]interface{}{
			"id": idStr,
		})
		return
	}

	h.logger.WithField("parsed_id", id).Debug("Parsed ID")

	review := h.fetchReview(w, r, id)
	if review == nil {
		return
	}

	writeJSON(w, http.StatusOK, toReviewResponse(review), h.logger)
}

// CreateReviewHandler creates a new review in the database.
func (h *ReviewHandler) CreateReviewHandler(w http.ResponseWriter, r *http.Request) {
	// ParseMultipartForm only when the client sends files — JSON-only requests
	// will fail this call gracefully; we continue without multipart in that case.
	_ = r.ParseMultipartForm(100 << 20)

	// Retrieve the validated and sanitized payload from the context
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.ReviewCreateRequest)
	if !ok {
		respondInvalidPayload(w, h.logger)
		return
	}

	// Get user ID from context
	userID := middleware.GetUserFromContext(r.Context())
	if userID == 0 {
		h.logger.Error("Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, msgNotAuthenticated, msgNotAuthReview, nil)
		return
	}

	// Create a new Review from the request
	review := &models.Review{
		UserID:  userID,
		AnimeID: payload.AnimeID,
		Content: payload.Content,
		Rating:  payload.Rating,
	}

	// Handle media file uploads
	attachments, err := h.processMediaFiles(r)
	if err != nil {
		respondUploadError(w, h.logger, err)
		return
	}
	review.MediaAttachments = attachments

	// Create review using service
	if err := h.reviewService.CreateReview(r.Context(), review); err != nil {
		h.logger.WithField("err", err).Warning("Failed to create review")
		if !writeAppError(w, err) {
			errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to create review", err.Error(), nil)
		}
		return
	}

	// Get the created review with user and anime details
	createdReview, err := h.reviewService.GetReviewByID(r.Context(), review.ID)
	if err != nil {
		h.logger.WithField("err", err).Warning("Failed to get created review")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to retrieve created review", "An internal server error occurred while retrieving the created review.", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, toReviewResponse(createdReview), h.logger)
}

// UpdateReviewHandler updates an existing review in the database.
func (h *ReviewHandler) UpdateReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := utils.ValidateID(idStr)
	if err != nil {
		h.logger.WithField("err", err).Warning("Invalid ID format")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", "The provided ID is not a valid unsigned integer.", map[string]interface{}{
			"id": idStr,
		})
		return
	}

	// ParseMultipartForm only when the client sends files — JSON-only requests
	// will fail this call gracefully; we continue without multipart in that case.
	_ = r.ParseMultipartForm(100 << 20)

	// Retrieve the validated and sanitized payload from the context
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.ReviewUpdateRequest)
	if !ok {
		respondInvalidPayload(w, h.logger)
		return
	}

	// Get user ID from context
	userID := middleware.GetUserFromContext(r.Context())
	if userID == 0 {
		h.logger.Error("Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, msgNotAuthenticated, msgNotAuthUpdateReview, nil)
		return
	}

	// Get the existing review
	review := h.fetchReview(w, r, id)
	if review == nil {
		return
	}

	// Check if user is authorized to update the review
	if review.UserID != userID {
		h.logger.WithFields(map[string]interface{}{
			"review_user_id":  review.UserID,
			"request_user_id": userID,
		}).Warning("Unauthorized review update attempt")
		errors.WriteErrorResponse(w, http.StatusForbidden, errors.ErrForbidden, "Unauthorized", msgUnauthorizedUpdate, nil)
		return
	}

	if payload.Content != "" {
		review.Content = payload.Content
	}
	if payload.Rating != 0 {
		review.Rating = payload.Rating
	}

	// Handle media file uploads
	newAttachments, err := h.processMediaFiles(r)
	if err != nil {
		respondUploadError(w, h.logger, err)
		return
	}
	review.MediaAttachments = append(review.MediaAttachments, newAttachments...)

	if err := h.reviewService.UpdateReview(r.Context(), review); err != nil {
		h.logger.WithField("err", err).Warning("Failed to update review")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to update review", "An internal server error occurred while updating the review.", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	updatedReview, err := h.reviewService.GetReviewByID(r.Context(), review.ID)
	if err != nil {
		h.logger.WithField("err", err).Warning("Failed to get updated review")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to retrieve updated review", "An internal server error occurred while retrieving the updated review.", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, toReviewResponse(updatedReview), h.logger)
}

// DeleteReviewHandler handles the deletion of a review.
func (h *ReviewHandler) DeleteReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := utils.ValidateID(idStr)
	if err != nil {
		h.logger.WithField("err", err).Warning("Invalid ID format")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", "The provided ID is not a valid unsigned integer.", map[string]interface{}{
			"id": idStr,
		})
		return
	}

	userID := middleware.GetUserFromContext(r.Context())
	if userID == 0 {
		h.logger.Error("Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, msgNotAuthenticated, msgNotAuthDeleteReview, nil)
		return
	}

	review := h.fetchReview(w, r, id)
	if review == nil {
		return
	}

	if review.UserID != userID {
		h.logger.WithFields(map[string]interface{}{
			"review_user_id":  review.UserID,
			"request_user_id": userID,
		}).Warning("Unauthorized review deletion attempt")
		errors.WriteErrorResponse(w, http.StatusForbidden, errors.ErrForbidden, "Unauthorized", msgUnauthorizedDelete, nil)
		return
	}

	if err := h.reviewService.DeleteReview(r.Context(), id); err != nil {
		h.logger.WithField("err", err).Warning("Failed to delete review")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to delete review", "An internal server error occurred while deleting the review.", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
