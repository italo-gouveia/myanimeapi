// pkg/handlers/review_handlers.go
// Package handlers provides HTTP handlers for review-related routes in the MyAnimeAPI application.
// It defines methods to handle requests for retrieving, creating, updating, and deleting reviews.
// The package uses the Gorilla Mux router for routing, GORM for database interactions, and middleware for request validation and authentication.
//
// Example usage:
//
//	db := // initialize your database connection
//	reviewHandler := handlers.NewReviewHandler(db)
//	router := mux.NewRouter()
//	reviewHandler.RegisterReviewRoutes(router)
//
//	http.ListenAndServe(":8080", router)
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"myanimeapi/internal/db"
	"myanimeapi/pkg/middleware"
	"myanimeapi/pkg/models"
	"myanimeapi/pkg/utils"

	"github.com/gorilla/mux"
)

// ReviewHandler defines the handlers for review-related routes.
// It contains a database interface for interacting with the database.
type ReviewHandler struct {
	DB db.DBInterface
}

// NewReviewHandler creates a new instance of ReviewHandler.
// It accepts a database interface and returns a pointer to a ReviewHandler.
//
// Example:
//
//	db := // initialize your database connection
//	reviewHandler := NewReviewHandler(db)
func NewReviewHandler(db db.DBInterface) *ReviewHandler {
	return &ReviewHandler{DB: db}
}

// GetReviewHandler retrieves a review by its ID.
// It validates the ID, queries the database, and returns the review as a JSON response.
// If the ID is invalid or the review is not found, it returns an appropriate error response.
//
// @Summary Get a review by ID
// @Description Retrieve a review by its ID
// @Tags reviews
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} models.ReviewResponse
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 404 {object} map[string]string "Review not found"
// @Failure 500 {object} map[string]string "Failed to retrieve review"
// @Router /v1/reviews/{id} [get]
// @ExampleResponse
//
//	{
//	  "id": 1,
//	  "user_id": 1,
//	  "anime_id": 1,
//	  "content": "Great anime!",
//	  "rating": 9,
//	  "created_at": "2023-10-01T12:00:00Z",
//	  "updated_at": "2023-10-01T12:00:00Z"
//	}
//
// @Security []
func (h *ReviewHandler) GetReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Debug: Log the raw ID string
	log.Printf("Raw ID string: %s", idStr)

	// Validate ID
	id, err := utils.ValidateID(idStr)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Debug: Log the parsed ID
	log.Printf("Parsed ID: %d", id)

	var review models.Review
	// Use the correct context for the database query
	result := h.DB.WithContext(r.Context()).Preload("User").Preload("Anime").First(&review, id)
	if result.Error != nil {
		log.Printf("Error fetching review: %v", result.Error)
		http.Error(w, "Review not found", http.StatusNotFound)
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
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// CreateReviewHandler creates a new review in the database.
// It validates the input payload, checks for the existence of the associated user and anime,
// and creates the review. If successful, it returns the created review as a JSON response.
//
// @Summary Create a new review
// @Description Create a new review for an anime
// @Tags reviews
// @Accept json
// @Produce json
// @Param review body models.ReviewCreateRequest true "Review data"
// @Success 201 {object} models.Review
// @Failure 400 {object} map[string]string "Invalid input or missing required fields"
// @Failure 404 {object} map[string]string "Anime or user not found"
// @Failure 500 {object} map[string]string "Failed to create review"
// @Router /v1/reviews [post]
// @Example
//
//	{
//	  "user_id": 1,
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
//	  "updated_at": "2023-10-01T12:00:00Z"
//	}
//
// @Security ApiKeyAuth
func (h *ReviewHandler) CreateReviewHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the validated and sanitized payload from the context
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.Review)
	if !ok {
		http.Error(w, "Invalid payload", http.StatusInternalServerError)
		return
	}

	// Check if the user exists
	var user models.User
	result := h.DB.First(r.Context(), &user, payload.UserID)
	if result.Error != nil {
		log.Printf("User not found: %v", result.Error)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if the anime exists
	var anime models.Anime
	result = h.DB.First(r.Context(), &anime, payload.AnimeID)
	if result.Error != nil {
		log.Printf("Anime not found: %v", result.Error)
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	// Create the review
	result = h.DB.Create(r.Context(), payload)
	if result.Error != nil {
		log.Printf("Failed to create review: %v", result.Error)
		http.Error(w, "Failed to create review", http.StatusInternalServerError)
		return
	}

	// Remove sensitive data before returning the response
	payload.User = models.User{}   // Clear User field
	payload.Anime = models.Anime{} // Clear Anime field

	log.Printf("Review created successfully: ID %d", payload.ID)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// UpdateReviewHandler updates an existing review in the database.
// It validates the ID and input payload, updates the review, and returns the updated review as a JSON response.
// If the ID or input is invalid, or the update fails, it returns an error response.
//
// @Summary Update a review by ID
// @Description Update a review with the input payload
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Param review body models.ReviewCreateRequest true "Updated review data"
// @Success 200 {object} models.ReviewResponse
// @Failure 400 {object} map[string]string "Invalid input or ID format"
// @Failure 404 {object} map[string]string "Review not found"
// @Failure 500 {object} map[string]string "Failed to update review"
// @Router /v1/reviews/{id} [put]
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
//	  "updated_at": "2023-10-01T12:00:00Z"
//	}
//
// @Security ApiKeyAuth
func (h *ReviewHandler) UpdateReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Validate ID
	id, err := utils.ValidateID(idStr)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the validated and sanitized payload from the context
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.Review)
	if !ok {
		http.Error(w, "Invalid payload", http.StatusInternalServerError)
		return
	}

	// Fetch the existing review
	var review models.Review
	result := h.DB.WithContext(r.Context()).First(&review, id)
	if result.Error != nil {
		log.Printf("Review not found: %v", result.Error)
		http.Error(w, "Review not found", http.StatusNotFound)
		return
	}

	// Update only the allowed fields
	review.Content = payload.Content
	review.Rating = payload.Rating

	// Save the updated review
	result = h.DB.WithContext(r.Context()).Save(&review)
	if result.Error != nil {
		log.Printf("Failed to update review: %v", result.Error)
		http.Error(w, "Failed to update review", http.StatusInternalServerError)
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
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// DeleteReviewHandler deletes a review from the database.
// It validates the ID, deletes the review, and returns a 204 No Content response.
// If the ID is invalid or the deletion fails, it returns an error response.
//
// @Summary Delete a review by ID
// @Description Delete a review by its ID
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 404 {string} string "Review not found"
// @Failure 500 {string} string "Failed to delete review"
// @Router /reviews/{id} [delete]
// @Security ApiKeyAuth
func (h *ReviewHandler) DeleteReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Validate ID
	id, err := utils.ValidateID(idStr)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch the existing review (including soft-deleted records)
	var review models.Review
	result := h.DB.WithContext(r.Context()).Unscoped().First(&review, id)
	if result.Error != nil {
		log.Printf("Review not found: %v", result.Error)
		http.Error(w, "Review not found", http.StatusNotFound)
		return
	}

	// Delete the review
	result = h.DB.WithContext(r.Context()).Delete(&models.Review{}, id)
	if result.Error != nil {
		log.Printf("Error deleting review: %v", result.Error)
		http.Error(w, "Failed to delete review", http.StatusInternalServerError)
		return
	}

	log.Printf("Review deleted successfully: ID %d", id)
	w.WriteHeader(http.StatusNoContent)
}

// RegisterReviewRoutes registers all review-related routes with the provided router.
// It defines public routes (GET) and protected routes (POST, PUT, DELETE) that require authentication.
//
// Example:
//
//	router := mux.NewRouter()
//	reviewHandler.RegisterReviewRoutes(router)
func (h *ReviewHandler) RegisterReviewRoutes(router *mux.Router) {
	// Public routes (no authentication required)
	router.HandleFunc("/reviews/{id:[0-9]+}", h.GetReviewHandler).Methods("GET")

	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/reviews").Subrouter()
	protectedRouter.Use(middleware.Authenticate) // Apply authentication middleware

	// Protected routes (require authentication)
	protectedRouter.Handle("", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.CreateReviewHandler), models.Review{})).Methods("POST")
	protectedRouter.Handle("/{id:[0-9]+}", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.UpdateReviewHandler), models.Review{})).Methods("PUT")
	protectedRouter.HandleFunc("/{id:[0-9]+}", h.DeleteReviewHandler).Methods("DELETE")
}
