// internal/handlers/review_handlers.go
// Package handlers provides the handler functions for the different routes.
// It defines the handlers for review-related routes.
// It uses the db package to interact with the database.
// It uses the models package to interact with the data models.
// It uses the middleware package to authenticate users.
// It uses the gorilla/mux package to handle HTTP requests.
// It uses the encoding/json package to encode and decode JSON data.
// It uses the log package to log messages.
// It uses the net/http package to write HTTP responses.
// It uses the strconv package to convert strings to integers.
// It uses the fmt package to format strings.
// It uses the myanimeapi/internal/db package to interact with the database.
// It uses the myanimeapi/pkg/models package to interact with the data models.
// It uses the myanimeapi/pkg/middleware package to authenticate users.
// It uses the myanimeapi/pkg/db package to interact with the database.
// It uses the myanimeapi/pkg/handlers package to handle the requests.
// It uses the myanimeapi/pkg/models package to interact with the data models.
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"myanimeapi/internal/db"
	"myanimeapi/pkg/middleware"
	"myanimeapi/pkg/models"

	"github.com/gorilla/mux"
)

// ReviewHandler defines the handlers for review-related routes
type ReviewHandler struct {
	DB db.DBInterface
}

// NewReviewHandler creates a new ReviewHandler instance
func NewReviewHandler(db db.DBInterface) *ReviewHandler {
	return &ReviewHandler{DB: db}
}

// GetReviewHandler godoc
// @Summary Get a review by ID
// @Description Get a review by its ID
// @Tags reviews
// @Accept  json
// @Produce  json
// @Param id path int true "Review ID"
// @Success 200 {object} models.Review
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /reviews/{id} [get]
func (h *ReviewHandler) GetReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert the ID from string to uint
	id, err := strconv.ParseUint(idStr, 10, 32) // Convert to uint32
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var review models.Review
	if err := h.DB.Preload("User", r.Context()).Preload("Anime", r.Context()).First(&review, uint(id)).Error; err != nil {
		log.Printf("Error fetching review: %v", err)
		http.Error(w, "Review not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

// CreateReviewHandler godoc
// @Summary Create a new review
// @Description Create a new review with the input payload
// @Tags reviews
// @Accept  json
// @Produce  json
// @Param review body models.Review true "Review object"
// @Success 201 {object} models.Review
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reviews [post]
func (h *ReviewHandler) CreateReviewHandler(w http.ResponseWriter, r *http.Request) {
	var review models.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		log.Printf("Invalid input: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if the user exists
	var user models.User
	if err := h.DB.First(r.Context(), &user, review.UserID).Error; err != nil {
		log.Printf("User not found: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if the anime exists
	var anime models.Anime
	if err := h.DB.First(r.Context(), &anime, review.AnimeID).Error; err != nil {
		log.Printf("Anime not found: %v", err)
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	// Validate if the content is provided
	if review.Content == "" {
		log.Println("Content is required")
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	// Validate if the rating is between 0 to 10
	if review.Rating < 0 || review.Rating > 10 {
		log.Println("Rating should be between 0 to 10")
		http.Error(w, "Rating should be between 0 to 10", http.StatusBadRequest)
		return
	}

	// Create the review
	if err := h.DB.Create(r.Context(), &review).Error; err != nil {
		log.Printf("Failed to create review: %v", err)
		http.Error(w, "Failed to create review", http.StatusInternalServerError)
		return
	}

	log.Printf("Review created successfully: ID %d", review.ID)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

// UpdateReviewHandler godoc
// @Summary Update a review by ID
// @Description Update a review with the input payload
// @Tags reviews
// @Accept  json
// @Produce  json
// @Param id path int true "Review ID"
// @Param review body models.Review true "Review object"
// @Success 200 {object} models.Review
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reviews/{id} [put]
func (h *ReviewHandler) UpdateReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert the ID from string to uint
	id, err := strconv.ParseUint(idStr, 10, 32) // Convert to uint32
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var review models.Review
	if err := h.DB.Preload("User", r.Context()).Preload("Anime", r.Context()).First(&review, uint(id)).Error; err != nil {
		log.Printf("Review not found: %v", err)
		http.Error(w, "Review not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		log.Printf("Invalid input: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Ensure the user and anime still exist
	var user models.User
	if err := h.DB.First(r.Context(), &user, review.UserID).Error; err != nil {
		log.Printf("User not found: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var anime models.Anime
	if err := h.DB.First(r.Context(), &anime, review.AnimeID).Error; err != nil {
		log.Printf("Anime not found: %v", err)
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	if err := h.DB.Save(r.Context(), &review).Error; err != nil {
		log.Printf("Failed to update review: %v", err)
		http.Error(w, "Failed to update review", http.StatusInternalServerError)
		return
	}

	log.Printf("Review updated successfully: ID %d", review.ID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

// DeleteReviewHandler godoc
// @Summary Delete a review by ID
// @Description Delete a review by its ID
// @Tags reviews
// @Accept  json
// @Produce  json
// @Param id path int true "Review ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reviews/{id} [delete]
func (h *ReviewHandler) DeleteReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert the ID from string to uint
	id, err := strconv.ParseUint(idStr, 10, 32) // Convert to uint32
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var review models.Review
	// Use Unscoped to include soft-deleted records
	if err := h.DB.Unscoped(r.Context()).Preload("User", r.Context()).Preload("Anime", r.Context()).First(&review, uint(id)).Error; err != nil {
		log.Printf("Review not found: %v", err)
		http.Error(w, "Review not found", http.StatusNotFound)
		return
	}

	// Delete the review
	if err := h.DB.Delete(r.Context(), &models.Review{}, id).Error; err != nil {
		log.Printf("Error deleting review: %v", err)
		http.Error(w, "Failed to delete review", http.StatusInternalServerError)
		return
	}

	log.Printf("Review deleted successfully: ID %d", id)
	w.WriteHeader(http.StatusNoContent)
}

// RegisterReviewRoutes registers all review-related routes
func (h *ReviewHandler) RegisterReviewRoutes(router *mux.Router) {
	// Public routes (no authentication required)
	router.HandleFunc("/reviews/{id:[0-9]+}", h.GetReviewHandler).Methods("GET")

	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/reviews").Subrouter()
	protectedRouter.Use(middleware.Authenticate) // Apply authentication middleware

	// Protected routes (require authentication)
	protectedRouter.HandleFunc("", h.CreateReviewHandler).Methods("POST")
	protectedRouter.HandleFunc("/{id:[0-9]+}", h.UpdateReviewHandler).Methods("PUT")
	protectedRouter.HandleFunc("/{id:[0-9]+}", h.DeleteReviewHandler).Methods("DELETE")
}
