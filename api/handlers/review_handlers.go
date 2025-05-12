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
	"log"
	"net/http"

	"myanimeapi/api/middleware"
	"myanimeapi/api/models"
	"myanimeapi/api/services"
	"myanimeapi/api/utils"
	"myanimeapi/internal/errors"

	"github.com/gorilla/mux"
)

// ReviewHandler defines the handlers for review-related routes.
// It contains a review service for handling business logic.
type ReviewHandler struct {
	reviewService *services.ReviewService
}

// NewReviewHandler creates a new instance of ReviewHandler.
// It accepts a review service and returns a pointer to a ReviewHandler.
//
// Example:
//
//	reviewService := services.NewReviewService(repository)
//	reviewHandler := NewReviewHandler(reviewService)
func NewReviewHandler(reviewService *services.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
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
	log.Printf("Raw ID string: %s", idStr)

	// Validate ID
	id, err := utils.ValidateID(idStr)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", "The provided ID is not a valid unsigned integer.")
		return
	}

	// Debug: Log the parsed ID
	log.Printf("Parsed ID: %d", id)

	// Get review from service
	review, err := h.reviewService.GetReviewByID(r.Context(), id)
	if err != nil {
		log.Printf("Failed to get review: %v", err)
		errors.WriteErrorResponse(w, http.StatusNotFound, errors.ErrResourceNotFound, "Review not found", fmt.Sprintf("Review with ID '%d' not found.", id))
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
		log.Printf("Failed to encode response: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.")
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
// @Accept json
// @Produce json
// @Param review body models.ReviewCreateRequest true "Review data"
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
	// Retrieve the validated and sanitized payload from the context
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.ReviewCreateRequest)
	if !ok {
		log.Printf("Invalid payload: %v", payload)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Invalid payload", "The request payload could not be retrieved.")
		return
	}

	// Create a new Review from the request
	review := &models.Review{
		UserID:  payload.UserID,
		AnimeID: payload.AnimeID,
		Content: payload.Content,
		Rating:  payload.Rating,
	}

	// Create review using service
	if err := h.reviewService.CreateReview(r.Context(), review); err != nil {
		log.Printf("Failed to create review: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to create review", "An internal server error occurred while creating the review.")
		return
	}

	// Get the created review with user and anime details
	createdReview, err := h.reviewService.GetReviewByID(r.Context(), review.ID)
	if err != nil {
		log.Printf("Failed to get created review: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to retrieve created review", "An internal server error occurred while retrieving the created review.")
		return
	}

	// Create a custom response to exclude sensitive data
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
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.")
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
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Param review body models.ReviewUpdateRequest true "Updated review data"
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
		log.Printf("Invalid ID format: %v", err)
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", "The provided ID is not a valid unsigned integer.")
		return
	}

	// Retrieve the validated and sanitized payload from the context
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.ReviewUpdateRequest)
	if !ok {
		log.Printf("Invalid payload: %v", payload)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Invalid payload", "The request payload could not be retrieved.")
		return
	}

	// Get the existing review
	review, err := h.reviewService.GetReviewByID(r.Context(), id)
	if err != nil {
		log.Printf("Failed to get review: %v", err)
		errors.WriteErrorResponse(w, http.StatusNotFound, errors.ErrResourceNotFound, "Review not found", fmt.Sprintf("Review with ID '%d' not found.", id))
		return
	}

	// Update only the fields that were provided in the request
	if payload.Content != "" {
		review.Content = payload.Content
	}
	if payload.Rating != 0 {
		review.Rating = payload.Rating
	}

	// Update review using service
	if err := h.reviewService.UpdateReview(r.Context(), review); err != nil {
		log.Printf("Failed to update review: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to update review", "An internal server error occurred while updating the review.")
		return
	}

	// Get the updated review
	updatedReview, err := h.reviewService.GetReviewByID(r.Context(), id)
	if err != nil {
		log.Printf("Failed to get updated review: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to retrieve updated review", "An internal server error occurred while retrieving the updated review.")
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
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.")
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
// @Success 200 {object} models.Response
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
		log.Printf("Invalid ID format: %v", err)
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", "The provided ID is not a valid unsigned integer.")
		return
	}

	// Delete review using service
	if err := h.reviewService.DeleteReview(r.Context(), id); err != nil {
		log.Printf("Failed to delete review: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to delete review", "An internal server error occurred while deleting the review.")
		return
	}

	log.Printf("Review deleted successfully: ID %d", id)
	w.WriteHeader(http.StatusNoContent)
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
	router.HandleFunc("/reviews/{id:[0-9]+}", h.GetReviewHandler).Methods("GET")

	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/reviews").Subrouter()
	protectedRouter.Use(middleware.Authenticate) // Apply authentication middleware

	// Protected routes (require authentication)
	protectedRouter.Handle("", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.CreateReviewHandler), models.ReviewCreateRequest{})).Methods("POST")
	protectedRouter.Handle("/{id:[0-9]+}", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.UpdateReviewHandler), models.ReviewUpdateRequest{})).Methods("PUT")
	protectedRouter.HandleFunc("/{id:[0-9]+}", h.DeleteReviewHandler).Methods("DELETE")
}
