// internal/handlers/review_handlers.go
package handlers

import (
	"encoding/json"
	"log"
	"myanimeapi/pkg/middleware"
	"myanimeapi/pkg/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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
func GetReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert the ID from string to uint
	id, err := strconv.ParseUint(idStr, 10, 32) // Convert to uint32
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var review models.Review
	if err := database.Preload("User").Preload("Anime").First(&review, uint(id)).Error; err != nil {
		log.Printf("Error fetching review: %v\n", err)
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
func CreateReviewHandler(w http.ResponseWriter, r *http.Request) {
	var review models.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	// Check if the user exists
	var user models.User
	if err := database.First(&user, review.UserID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if the anime exists
	var anime models.Anime
	if err := database.First(&anime, review.AnimeID).Error; err != nil {
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	// Validate if the content is provided
	if review.Content == "" {
		http.Error(w, "Failed to create review: Content is required", http.StatusBadRequest)
		return
	}

	// Validate if the rating is between 0 to 10
	if review.Rating < 0 || review.Rating > 10 {
		http.Error(w, "Failed to create review: Rating should be between 0 to 10", http.StatusBadRequest)
		return
	}

	// Create the review
	if err := database.Create(&review).Error; err != nil {
		http.Error(w, "Failed to create review", http.StatusInternalServerError)
		return
	}

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
func UpdateReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert the ID from string to uint
	id, err := strconv.ParseUint(idStr, 10, 32) // Convert to uint32
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var review models.Review
	if err := database.Preload("User").Preload("Anime").First(&review, uint(id)).Error; err != nil {
		http.Error(w, "Review not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Ensure the user and anime still exist
	var user models.User
	if err := database.First(&user, review.UserID).Error; err != nil {
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	var anime models.Anime
	if err := database.First(&anime, review.AnimeID).Error; err != nil {
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	if err := database.Save(&review).Error; err != nil {
		http.Error(w, "Failed to update review", http.StatusInternalServerError)
		return
	}

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
func DeleteReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert the ID from string to uint
	id, err := strconv.ParseUint(idStr, 10, 32) // Convert to uint32
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var review models.Review
	// Use Unscoped to include soft-deleted records
	if err := database.Unscoped().Preload("User").Preload("Anime").First(&review, uint(id)).Error; err != nil {
		log.Printf("Review not found: %v\n", err)
		http.Error(w, "Review not found", http.StatusNotFound)
		return
	}

	// Delete the review
	if err := database.Delete(&models.Review{}, id).Error; err != nil {
		log.Printf("Error deleting review: %v\n", err)
		http.Error(w, "Failed to delete review", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func RegisterReviewRoutes(router *mux.Router) {
	// Public routes (no authentication required)
	router.HandleFunc("/reviews/{id:[0-9]+}", GetReviewHandler).Methods("GET")

	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/reviews").Subrouter()
	protectedRouter.Use(middleware.Authenticate) // Apply authentication middleware

	// Protected routes (require authentication)
	protectedRouter.HandleFunc("", CreateReviewHandler).Methods("POST")
	protectedRouter.HandleFunc("/{id:[0-9]+}", UpdateReviewHandler).Methods("PUT")
	protectedRouter.HandleFunc("/{id:[0-9]+}", DeleteReviewHandler).Methods("DELETE")
}
