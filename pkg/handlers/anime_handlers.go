// internal/handlers/anime_handlers.go
// This file defines the handlers for anime-related routes.
// It imports the necessary packages and defines the AnimeHandler struct.
// It defines methods to handle requests for retrieving, creating, updating, and deleting anime entries.
// It uses the db package to interact with the database.
// It uses the models package to work with data models.
// It uses the gorilla/mux package to handle HTTP requests.
// It uses the log package to log messages.
// It uses the encoding/json package to encode and decode JSON data.
// It uses the net/http package to write HTTP responses.
// It uses the strconv package to convert strings to other types.
// It uses the middleware package to authenticate requests.
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

// AnimeHandler defines the handlers for anime-related routes
type AnimeHandler struct {
	DB db.DBInterface
}

// NewAnimeHandler creates a new AnimeHandler instance
func NewAnimeHandler(db db.DBInterface) *AnimeHandler {
	return &AnimeHandler{DB: db}
}

// GetAnimeHandler retrieves an anime by ID
// @Summary Get an anime by ID
// @Description Retrieve an anime by its ID
// @Tags anime
// @Produce json
// @Param id path int true "Anime ID"
// @Success 200 {object} models.Anime
// @Failure 400 {string} string "Invalid ID format"
// @Failure 404 {string} string "Anime not found"
// @Router /animes/{id} [get]
func (h *AnimeHandler) GetAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Validate ID
	id, err := validation.ValidateID(idStr)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var anime models.Anime
	result := h.DB.WithContext(r.Context()).First(&anime, id)
	if result.Error != nil {
		log.Printf("Anime not found: %v", result.Error)
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(anime); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetAllAnimesHandler retrieves all anime entries
// @Summary Get all anime entries
// @Description Retrieve a list of all anime entries
// @Tags anime
// @Produce json
// @Success 200 {array} models.Anime
// @Failure 500 {string} string "Failed to retrieve animes"
// @Router /animes [get]
func (h *AnimeHandler) GetAllAnimesHandler(w http.ResponseWriter, r *http.Request) {
	var animes []models.Anime
	result := h.DB.Find(r.Context(), &animes)
	if result.Error != nil {
		log.Printf("Failed to retrieve animes: %v", result.Error)
		http.Error(w, "Failed to retrieve animes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(animes); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetPaginatedReviewsForAnimeHandler retrieves paginated reviews for an anime by ID
// @Summary Get paginated reviews for an anime
// @Description Retrieve paginated reviews for an anime by its ID
// @Tags anime
// @Produce json
// @Param id path int true "Anime ID"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10)"
// @Success 200 {array} models.Review
// @Failure 400 {string} string "Invalid ID format or pagination parameters"
// @Failure 500 {string} string "Failed to retrieve reviews"
// @Router /animes/{id}/reviews [get]
func (h *AnimeHandler) GetPaginatedReviewsForAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Debug: Log the raw ID string
	log.Printf("Raw ID string: %s", idStr)

	// Validate ID
	id, err := validation.ValidateID(idStr)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Debug: Log the parsed ID
	log.Printf("Parsed ID: %d", id)

	// Validate pagination
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	page, limit, err := validation.ValidatePagination(pageStr, limitStr, 1, 10)
	if err != nil {
		log.Printf("Invalid pagination parameters: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Debug: Log pagination parameters
	log.Printf("Page: %d, Limit: %d", page, limit)

	var reviews []models.Review
	offset := (page - 1) * limit

	// Debug: Log the query being executed
	log.Printf("Executing query: SELECT * FROM reviews WHERE anime_id = %d OFFSET %d LIMIT %d", id, offset, limit)

	result := h.DB.WithContext(r.Context()).Where("anime_id = ?", id).Offset(offset).Limit(limit).Find(&reviews)
	if result.Error != nil {
		log.Printf("Failed to retrieve reviews: %v", result.Error)
		http.Error(w, "Failed to retrieve reviews", http.StatusInternalServerError)
		return
	}

	// Debug: Log the number of reviews found
	log.Printf("Found %d reviews for anime ID %d", len(reviews), id)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reviews); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// CreateAnimeHandler creates a new anime entry
// @Summary Create a new anime
// @Description Create a new anime entry with the provided data
// @Tags anime
// @Accept json
// @Produce json
// @Param anime body models.Anime true "Anime data"
// @Success 201 {object} models.Anime
// @Failure 400 {string} string "Invalid input or missing required fields"
// @Failure 500 {string} string "Failed to create anime"
// @Router /animes [post]
func (h *AnimeHandler) CreateAnimeHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the validated and sanitized payload from the context
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.Anime)
	if !ok {
		http.Error(w, "Invalid payload", http.StatusInternalServerError)
		return
	}

	// Proceed with creating the anime
	result := h.DB.Create(r.Context(), payload)
	if result.Error != nil {
		log.Printf("Failed to create anime: %v", result.Error)
		http.Error(w, "Failed to create anime", http.StatusInternalServerError)
		return
	}

	log.Printf("Anime %s created successfully", payload.Title)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// UpdateAnimeHandler updates an existing anime entry
// @Summary Update an anime
// @Description Update an existing anime entry with the provided data
// @Tags anime
// @Accept json
// @Produce json
// @Param id path int true "Anime ID"
// @Param anime body models.Anime true "Updated anime data"
// @Success 200 {object} models.Anime
// @Failure 400 {string} string "Invalid input or ID format"
// @Failure 404 {string} string "Anime not found"
// @Failure 500 {string} string "Failed to update anime"
// @Router /animes/{id} [put]
func (h *AnimeHandler) UpdateAnimeHandler(w http.ResponseWriter, r *http.Request) {
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
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.Anime)
	if !ok {
		http.Error(w, "Invalid payload", http.StatusInternalServerError)
		return
	}

	result := h.DB.First(r.Context(), &payload, uint(id))
	if result.Error != nil {
		log.Printf("Anime not found: %v", result.Error)
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	// Ensure the ID is preserved during update
	payload.ID = uint(id)

	// Save the updated anime
	result = h.DB.Save(r.Context(), payload)
	if result.Error != nil {
		log.Printf("Failed to update anime: %v", result.Error)
		http.Error(w, "Failed to update anime", http.StatusInternalServerError)
		return
	}

	log.Printf("Anime %s updated successfully", payload.Title)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// DeleteAnimeHandler deletes an anime entry
// @Summary Delete an anime
// @Description Delete an anime entry by its ID
// @Tags anime
// @Param id path int true "Anime ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 404 {string} string "Anime not found"
// @Failure 500 {string} string "Failed to delete anime"
// @Router /animes/{id} [delete]
func (h *AnimeHandler) DeleteAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Validate ID
	id, err := validation.ValidateID(idStr)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Delete all reviews associated with the anime
	result := h.DB.WithContext(r.Context()).Where("anime_id = ?", id).Delete(&models.Review{})
	if result.Error != nil {
		log.Printf("Failed to delete reviews: %v", result.Error)
		http.Error(w, "Failed to delete reviews", http.StatusInternalServerError)
		return
	}

	// Delete the anime
	result = h.DB.WithContext(r.Context()).Delete(&models.Anime{}, id)
	if result.Error != nil {
		log.Printf("Failed to delete anime: %v", result.Error)
		http.Error(w, "Failed to delete anime", http.StatusInternalServerError)
		return
	}

	log.Printf("Anime %d deleted successfully", id)
	w.WriteHeader(http.StatusNoContent)
}

// RegisterAnimeRoutes registers all anime-related routes
func (h *AnimeHandler) RegisterAnimeRoutes(router *mux.Router) {
	// Public routes (no authentication required)
	router.HandleFunc("/anime", h.GetAllAnimesHandler).Methods("GET")
	router.HandleFunc("/anime/{id:[0-9]+}", h.GetAnimeHandler).Methods("GET")
	router.HandleFunc("/anime/{id:[0-9]+}/reviews", h.GetPaginatedReviewsForAnimeHandler).Methods("GET")

	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/anime").Subrouter()
	protectedRouter.Use(middleware.Authenticate) // Apply authentication middleware

	// Protected routes (require authentication)
	protectedRouter.Handle("", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.CreateAnimeHandler), models.Anime{})).Methods("POST")
	protectedRouter.Handle("/{id:[0-9]+}", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.UpdateAnimeHandler), models.Anime{})).Methods("PUT")
	protectedRouter.HandleFunc("/{id:[0-9]+}", h.DeleteAnimeHandler).Methods("DELETE")
}
