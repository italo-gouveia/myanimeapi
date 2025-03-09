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

	"myanimeapi/internal/db"
	"myanimeapi/pkg/middleware"
	"myanimeapi/pkg/models"
	"myanimeapi/pkg/validation"

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

	// Validate ID
	id, err := validation.ValidateID(idStr)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var review models.Review
	result := h.DB.Preload("User", r.Context()).Preload("Anime", r.Context()).First(&review, uint(id))
	if result.Error != nil {
		log.Printf("Error fetching review: %v", result.Error)
		http.Error(w, "Review not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(review); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
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

	log.Printf("Review created successfully: ID %d", payload.ID)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
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

	// Validate ID
	id, err := validation.ValidateID(idStr)
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
	result := h.DB.Preload("User", r.Context()).Preload("Anime", r.Context()).First(&review, uint(id))
	if result.Error != nil {
		log.Printf("Review not found: %v", result.Error)
		http.Error(w, "Review not found", http.StatusNotFound)
		return
	}

	// Ensure the user and anime still exist
	var user models.User
	result = h.DB.First(r.Context(), &user, review.UserID)
	if result.Error != nil {
		log.Printf("User not found: %v", result.Error)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var anime models.Anime
	result = h.DB.First(r.Context(), &anime, review.AnimeID)
	if result.Error != nil {
		log.Printf("Anime not found: %v", result.Error)
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	// Update the review fields
	review.Content = payload.Content
	review.Rating = payload.Rating

	// Save the updated review
	result = h.DB.Save(r.Context(), &review)
	if result.Error != nil {
		log.Printf("Failed to update review: %v", result.Error)
		http.Error(w, "Failed to update review", http.StatusInternalServerError)
		return
	}

	log.Printf("Review updated successfully: ID %d", review.ID)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(review); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
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

	// Validate ID
	id, err := validation.ValidateID(idStr)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var review models.Review
	// Use Unscoped to include soft-deleted records
	result := h.DB.Unscoped(r.Context()).Preload("User", r.Context()).Preload("Anime", r.Context()).First(&review, uint(id))
	if result.Error != nil {
		log.Printf("Review not found: %v", result.Error)
		http.Error(w, "Review not found", http.StatusNotFound)
		return
	}

	// Delete the review
	result = h.DB.Delete(r.Context(), &models.Review{}, id)
	if result.Error != nil {
		log.Printf("Error deleting review: %v", result.Error)
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
	protectedRouter.Handle("", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.CreateReviewHandler), models.Review{})).Methods("POST")
	protectedRouter.Handle("/{id:[0-9]+}", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.UpdateReviewHandler), models.Review{})).Methods("PUT")
	protectedRouter.HandleFunc("/{id:[0-9]+}", h.DeleteReviewHandler).Methods("DELETE")
}
