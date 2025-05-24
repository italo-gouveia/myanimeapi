// api/handlers/review_handlers.go
// Package handlers provides HTTP handlers for review-related routes in the MyAnimeAPI application.
// It defines methods to handle requests for retrieving, creating, updating, and deleting reviews.
// The package uses the Gorilla Mux router for routing, GORM for database interactions, and middleware for request validation and authentication.
//
// Example usage:
//
//	reviewService := services.NewReviewService(repository)
//	reviewHandler := handlers.NewReviewHandler(reviewService)
//	router := mux.NewRouter()
//	reviewHandler.RegisterReviewRoutes(router)
//
//	http.ListenAndServe(":8080", router)
package handlers

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
// It contains a review service for handling business logic.
type ReviewHandler struct {
	reviewService services.ReviewServiceInterface
	storageSvc    services.StorageServiceInterface
	logger        *logger.Logger
}

// NewReviewHandler creates a new instance of ReviewHandler.
// It accepts a review service and returns a pointer to a ReviewHandler.
//
// Example:
//
//	reviewService := services.NewReviewService(repository)
//	reviewHandler := NewReviewHandler(reviewService)
func NewReviewHandler(reviewService services.ReviewServiceInterface, storageSvc services.StorageServiceInterface) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
		storageSvc:    storageSvc,
		logger:        logger.New(),
	}
}

// RegisterReviewRoutes registers all review-related routes with a *mux.Router.
// It sets up the routes for review management, including public and protected endpoints.
//
// Routes registered:
// - GET /reviews/{id} - Get a specific review (public)
// - POST /reviews - Create a new review (protected)
// - PUT /reviews/{id} - Update a review (protected)
// - DELETE /reviews/{id} - Delete a review (protected)
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
// It validates the ID, queries the service, and returns the review as a JSON response.
// If the ID is invalid or the review is not found, it returns an appropriate error response.
//
// @Summary Get a review by ID
// @Description Retrieve a review by its ID
// @Tags reviews
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} models.ReviewResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid ID format"
// @Failure 404 {object} errors.ErrorResponse "Review not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve review"
// @Router /reviews/{id} [get]
// @ExampleResponse
//
//	{
//	  "id": 1,
//	  "user_id": 1,
//	  "anime_id": 1,
//	  "content": "Great anime!",
//	  "rating": 9,
//	  "created_at": "2023-10-01T12:00:00Z",
//	  "updated_at": "2023-10-01T12:00:00Z",
//	  "user": {
//	    "id": 1,
//	    "username": "johndoe",
//	    "is_admin": false
//	  },
//	  "anime": {
//	    "id": 1,
//	    "title": "Naruto",
//	    "description": "A story about ninjas.",
//	    "rating": 8.5
//	  }
//	}
func (h *ReviewHandler) GetReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Debug: Log the raw ID string
	h.logger.WithField("raw_id", idStr).Debug("Raw ID string")

	// Validate ID
	id, err := utils.ValidateID(idStr)
	if err != nil {
		h.logger.WithField("err", err).Warning("Invalid ID format")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", "The provided ID is not a valid unsigned integer.", map[string]interface{}{
			"id": idStr,
		})
		return
	}

	// Debug: Log the parsed ID
	h.logger.WithField("parsed_id", id).Debug("Parsed ID")

	// Get review from service
	review, err := h.reviewService.GetReviewByID(r.Context(), id)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{"review_id": id, "err": err}).Warning("Failed to get review")
		errors.WriteErrorResponse(w, http.StatusNotFound, errors.ErrResourceNotFound, "Review not found", fmt.Sprintf("Review with ID '%d' not found.", id), map[string]interface{}{
			"id": id,
		})
		return
	}

	// Create a custom response to exclude sensitive data
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
// It validates the input payload and uses the service to create the review.
// If successful, it returns the created review as a JSON response.
//
// @Summary Create a new review
// @Description Create a new review for an anime
// @Tags reviews
// @Accept multipart/form-data
// @Produce json
// @Param review body models.ReviewCreateRequest true "Review data"
// @Param media formData file false "Media files (images or videos)"
// @Success 201 {object} models.ReviewResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid input or missing required fields"
// @Failure 404 {object} errors.ErrorResponse "User or anime not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to create review"
// @Router /reviews [post]
// @Security BearerAuth
// @Example
//
//	{
//	  "anime_id": 1,
//	  "content": "Great anime!",
//	  "rating": 9
//	}
//
// @ExampleResponse
//
//	{
//	  "id": 1,
//	  "user_id": 1,
//	  "anime_id": 1,
//	  "content": "Great anime!",
//	  "rating": 9,
//	  "created_at": "2023-10-01T12:00:00Z",
//	  "updated_at": "2023-10-01T12:00:00Z",
//	  "user": {
//	    "id": 1,
//	    "username": "johndoe",
//	    "is_admin": false
//	  },
//	  "anime": {
//	    "id": 1,
//	    "title": "Naruto",
//	    "description": "A story about ninjas.",
//	    "rating": 8.5
//	  }
//	}
func (h *ReviewHandler) CreateReviewHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	if err := r.ParseMultipartForm(100 << 20); err != nil { // 100MB max
		h.logger.WithField("err", err).Warning("Failed to parse multipart form")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Failed to parse form data", "The request form data could not be parsed.", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

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
				// Determine media type from file extension
				ext := strings.ToLower(filepath.Ext(file.Filename))
				mediaType := "image"
				if ext == ".mp4" || ext == ".webm" || ext == ".mov" {
					mediaType = "video"
				}

				// Upload file
				fileURL, err := h.storageSvc.SaveFile(file, mediaType)
				if err != nil {
					h.logger.WithField("err", err).Warning("Failed to upload file")
					errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to upload file", "An internal server error occurred while uploading the file.", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				// Add media attachment to review
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
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to create review", "An internal server error occurred while creating the review.", map[string]interface{}{
			"error": err.Error(),
		})
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

	// Create response
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

	// Set response headers and encode response
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
// It validates the ID and input payload, uses the service to update the review,
// and returns the updated review as a JSON response.
//
// @Summary Update a review
// @Description Update an existing review's details
// @Tags reviews
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Review ID"
// @Param review body models.ReviewUpdateRequest true "Updated review data"
// @Param media formData file false "Media files (images or videos)"
// @Success 200 {object} models.ReviewResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid input or missing required fields"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized to update this review"
// @Failure 404 {object} errors.ErrorResponse "Review not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to update review"
// @Router /reviews/{id} [put]
// @Security BearerAuth
// @Example
//
//	{
//	  "content": "Updated review content",
//	  "rating": 8
//	}
//
// @ExampleResponse
//
//	{
//	  "id": 1,
//	  "user_id": 1,
//	  "anime_id": 1,
//	  "content": "Updated review content",
//	  "rating": 8,
//	  "created_at": "2023-10-01T12:00:00Z",
//	  "updated_at": "2023-10-01T13:00:00Z",
//	  "user": {
//	    "id": 1,
//	    "username": "johndoe",
//	    "is_admin": false
//	  },
//	  "anime": {
//	    "id": 1,
//	    "title": "Naruto",
//	    "description": "A story about ninjas.",
//	    "rating": 8.5
//	  }
//	}
func (h *ReviewHandler) UpdateReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Validate ID
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

	// Update only the fields that were provided in the request
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
				// Determine media type from file extension
				ext := strings.ToLower(filepath.Ext(file.Filename))
				mediaType := "image"
				if ext == ".mp4" || ext == ".webm" || ext == ".mov" {
					mediaType = "video"
				}

				// Upload file
				fileURL, err := h.storageSvc.SaveFile(file, mediaType)
				if err != nil {
					h.logger.WithField("err", err).Warning("Failed to upload file")
					errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to upload file", "An internal server error occurred while uploading the file.", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				// Add media attachment to review
				review.MediaAttachments = append(review.MediaAttachments, models.MediaAttachment{
					Type: mediaType,
					URL:  fileURL,
				})
			}
		}
	}

	// Update review using service
	if err := h.reviewService.UpdateReview(r.Context(), review); err != nil {
		h.logger.WithField("err", err).Warning("Failed to update review")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to update review", "An internal server error occurred while updating the review.", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Get the updated review with user and anime details
	updatedReview, err := h.reviewService.GetReviewByID(r.Context(), review.ID)
	if err != nil {
		h.logger.WithField("err", err).Warning("Failed to get updated review")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to retrieve updated review", "An internal server error occurred while retrieving the updated review.", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Create a custom response to exclude sensitive data
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
// It validates the ID and uses the service to delete the review.
// If successful, it returns a success message as a JSON response.
//
// @Summary Delete a review
// @Description Delete an existing review by its ID
// @Tags reviews
// @Produce json
// @Param id path int true "Review ID"
// @Success 204 "No Content"
// @Failure 400 {object} errors.ErrorResponse "Invalid ID format"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized to delete this review"
// @Failure 404 {object} errors.ErrorResponse "Review not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete review"
// @Router /reviews/{id} [delete]
// @Security BearerAuth
// @ExampleResponse
//
//	{
//	  "status": "success",
//	  "message": "Review deleted successfully"
//	}
func (h *ReviewHandler) DeleteReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Validate ID
	id, err := utils.ValidateID(idStr)
	if err != nil {
		h.logger.WithField("err", err).Warning("Invalid ID format")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", "The provided ID is not a valid unsigned integer.", map[string]interface{}{
			"id": idStr,
		})
		return
	}

	// Get user ID from context
	userID := middleware.GetUserFromContext(r.Context())
	if userID == 0 {
		h.logger.Error("Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "User not authenticated", "You must be logged in to delete a review.", nil)
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

	// Check if user is authorized to delete the review
	if review.UserID != userID {
		h.logger.WithFields(map[string]interface{}{
			"review_user_id":  review.UserID,
			"request_user_id": userID,
		}).Warning("Unauthorized review deletion attempt")
		errors.WriteErrorResponse(w, http.StatusForbidden, errors.ErrForbidden, "Unauthorized", "You are not authorized to delete this review.", nil)
		return
	}

	// Delete review using service
	if err := h.reviewService.DeleteReview(r.Context(), id); err != nil {
		h.logger.WithField("err", err).Warning("Failed to delete review")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to delete review", "An internal server error occurred while deleting the review.", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
