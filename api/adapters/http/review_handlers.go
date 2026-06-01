// api/adapters/http/review_handlers.go
// Package httphandler provides HTTP handlers for review-related routes in the MyAnimeAPI application.
package httphandler

import (
	"encoding/json"
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

	review, err := h.reviewService.GetReviewByID(r.Context(), id)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{"review_id": id, "err": err}).Warning("Failed to get review")
		errors.WriteErrorResponse(w, http.StatusNotFound, errors.ErrResourceNotFound, "Review not found", fmt.Sprintf("Review with ID '%d' not found.", id), map[string]interface{}{
			"id": id,
		})
		return
	}

	response := models.ReviewResponse{
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
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.WithField("err", err).Error("Failed to encode response")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
}

// CreateReviewHandler creates a new review in the database.
func (h *ReviewHandler) CreateReviewHandler(w http.ResponseWriter, r *http.Request) {
	// ParseMultipartForm only when the client sends files — JSON-only requests
	// will fail this call gracefully; we continue without multipart in that case.
	_ = r.ParseMultipartForm(100 << 20)

	// Retrieve the validated and sanitized payload from the context
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.ReviewCreateRequest)
	if !ok {
		h.logger.Warning("Invalid payload")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Invalid payload", "The request payload could not be retrieved.", nil)
		return
	}

	// Get user ID from context
	userID := middleware.GetUserFromContext(r.Context())
	if userID == 0 {
		h.logger.Error("Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "User not authenticated", "You must be logged in to create a review.", nil)
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
	form := r.MultipartForm
	if form != nil && form.File != nil {
		for _, files := range form.File {
			for _, file := range files {
				ext := strings.ToLower(filepath.Ext(file.Filename))
				mediaType := "image"
				if ext == ".mp4" || ext == ".webm" || ext == ".mov" {
					mediaType = "video"
				}

				fileURL, err := h.storageSvc.SaveFile(file, mediaType)
				if err != nil {
					h.logger.WithField("err", err).Warning("Failed to upload file")
					errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to upload file", "An internal server error occurred while uploading the file.", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				review.MediaAttachments = append(review.MediaAttachments, models.MediaAttachment{
					Type: mediaType,
					URL:  fileURL,
				})
			}
		}
	}

	// Create review using service
	if err := h.reviewService.CreateReview(r.Context(), review); err != nil {
		h.logger.WithField("err", err).Warning("Failed to create review")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
		} else {
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

	response := models.ReviewResponse{
		ID:        createdReview.ID,
		CreatedAt: createdReview.CreatedAt,
		UpdatedAt: createdReview.UpdatedAt,
		UserID:    createdReview.UserID,
		AnimeID:   createdReview.AnimeID,
		Content:   createdReview.Content,
		Rating:    createdReview.Rating,
		User: models.UserResponse{
			ID:       createdReview.User.ID,
			Username: createdReview.User.Username,
			IsAdmin:  createdReview.User.IsAdmin,
		},
		Anime: models.AnimeResponse{
			ID:          createdReview.Anime.ID,
			Title:       createdReview.Anime.Title,
			Description: createdReview.Anime.Description,
			Rating:      createdReview.Anime.Rating,
		},
		MediaAttachments: createdReview.MediaAttachments,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.WithField("err", err).Error("Failed to encode response")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
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

	// Parse multipart form
	if err := r.ParseMultipartForm(100 << 20); err != nil { // 100MB max
		h.logger.WithField("err", err).Warning("Failed to parse multipart form")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Failed to parse form data", "The request form data could not be parsed.", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Retrieve the validated and sanitized payload from the context
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.ReviewUpdateRequest)
	if !ok {
		h.logger.Warning("Invalid payload")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Invalid payload", "The request payload could not be retrieved.", nil)
		return
	}

	// Get user ID from context
	userID := middleware.GetUserFromContext(r.Context())
	if userID == 0 {
		h.logger.Error("Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "User not authenticated", "You must be logged in to update a review.", nil)
		return
	}

	// Get the existing review
	review, err := h.reviewService.GetReviewByID(r.Context(), id)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{"review_id": id, "err": err}).Warning("Failed to get review")
		errors.WriteErrorResponse(w, http.StatusNotFound, errors.ErrResourceNotFound, "Review not found", fmt.Sprintf("Review with ID '%d' not found.", id), map[string]interface{}{
			"id": id,
		})
		return
	}

	// Check if user is authorized to update the review
	if review.UserID != userID {
		h.logger.WithFields(map[string]interface{}{
			"review_user_id":  review.UserID,
			"request_user_id": userID,
		}).Warning("Unauthorized review update attempt")
		errors.WriteErrorResponse(w, http.StatusForbidden, errors.ErrForbidden, "Unauthorized", "You are not authorized to update this review.", nil)
		return
	}

	if payload.Content != "" {
		review.Content = payload.Content
	}
	if payload.Rating != 0 {
		review.Rating = payload.Rating
	}

	// Handle media file uploads
	form := r.MultipartForm
	if form != nil && form.File != nil {
		for _, files := range form.File {
			for _, file := range files {
				ext := strings.ToLower(filepath.Ext(file.Filename))
				mediaType := "image"
				if ext == ".mp4" || ext == ".webm" || ext == ".mov" {
					mediaType = "video"
				}

				fileURL, err := h.storageSvc.SaveFile(file, mediaType)
				if err != nil {
					h.logger.WithField("err", err).Warning("Failed to upload file")
					errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to upload file", "An internal server error occurred while uploading the file.", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				review.MediaAttachments = append(review.MediaAttachments, models.MediaAttachment{
					Type: mediaType,
					URL:  fileURL,
				})
			}
		}
	}

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

	response := models.ReviewResponse{
		ID:        updatedReview.ID,
		CreatedAt: updatedReview.CreatedAt,
		UpdatedAt: updatedReview.UpdatedAt,
		UserID:    updatedReview.UserID,
		AnimeID:   updatedReview.AnimeID,
		Content:   updatedReview.Content,
		Rating:    updatedReview.Rating,
		User: models.UserResponse{
			ID:       updatedReview.User.ID,
			Username: updatedReview.User.Username,
			IsAdmin:  updatedReview.User.IsAdmin,
		},
		Anime: models.AnimeResponse{
			ID:          updatedReview.Anime.ID,
			Title:       updatedReview.Anime.Title,
			Description: updatedReview.Anime.Description,
			Rating:      updatedReview.Anime.Rating,
		},
		MediaAttachments: updatedReview.MediaAttachments,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.WithField("err", err).Error("Failed to encode response")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
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
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "User not authenticated", "You must be logged in to delete a review.", nil)
		return
	}

	review, err := h.reviewService.GetReviewByID(r.Context(), id)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{"review_id": id, "err": err}).Warning("Failed to get review")
		errors.WriteErrorResponse(w, http.StatusNotFound, errors.ErrResourceNotFound, "Review not found", fmt.Sprintf("Review with ID '%d' not found.", id), map[string]interface{}{
			"id": id,
		})
		return
	}

	if review.UserID != userID {
		h.logger.WithFields(map[string]interface{}{
			"review_user_id":  review.UserID,
			"request_user_id": userID,
		}).Warning("Unauthorized review deletion attempt")
		errors.WriteErrorResponse(w, http.StatusForbidden, errors.ErrForbidden, "Unauthorized", "You are not authorized to delete this review.", nil)
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
