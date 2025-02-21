// internal/handlers/anime_handlers.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"myanimeapi/pkg/models"

	"myanimeapi/pkg/middleware"

	"github.com/gorilla/mux"
)

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
func GetAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert the ID from string to uint
	id, err := strconv.ParseUint(idStr, 10, 32) // Convert to uint32
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var anime models.Anime
	if err := database.First(&anime, id).Error; err != nil {
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(anime)
}

// GetAllAnimesHandler retrieves all anime entries
// @Summary Get all anime entries
// @Description Retrieve a list of all anime entries
// @Tags anime
// @Produce json
// @Success 200 {array} models.Anime
// @Failure 500 {string} string "Failed to retrieve animes"
// @Router /animes [get]
func GetAllAnimesHandler(w http.ResponseWriter, r *http.Request) {
	var animes []models.Anime
	if err := database.Find(&animes).Error; err != nil {
		http.Error(w, "Failed to retrieve animes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animes)
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
func GetPaginatedReviewsForAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert the ID from string to uint
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Set default values if not provided
	page := 1
	limit := 10

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit number", http.StatusBadRequest)
			return
		}
	}

	var reviews []models.Review
	offset := (page - 1) * limit

	if err := database.Where("anime_id = ?", id).Offset(offset).Limit(limit).Find(&reviews).Error; err != nil {
		http.Error(w, "Failed to retrieve reviews", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
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
func CreateAnimeHandler(w http.ResponseWriter, r *http.Request) {
	var anime models.Anime
	if err := json.NewDecoder(r.Body).Decode(&anime); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if title is provided
	if anime.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if err := database.Create(&anime).Error; err != nil {
		http.Error(w, "Failed to create anime", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(anime)
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
func UpdateAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert the ID from string to uint
	id, err := strconv.ParseUint(idStr, 10, 32) // Convert to uint32
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var anime models.Anime
	if err := database.First(&anime, uint(id)).Error; err != nil {
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&anime); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	anime.ID = uint(id) // Ensure the ID is preserved during update
	if err := database.Save(&anime).Error; err != nil {
		http.Error(w, "Failed to update anime", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(anime)
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
func DeleteAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert the ID from string to uint
	id, err := strconv.ParseUint(idStr, 10, 32) // Convert to uint32
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var anime models.Anime
	if err := database.First(&anime, uint(id)).Error; err != nil {
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	if err := database.Delete(&models.Anime{}, id).Error; err != nil {
		http.Error(w, "Failed to delete anime", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RegisterAnimeRoutes registers all anime-related routes
func RegisterAnimeRoutes(router *mux.Router) {
	// Public routes (no authentication required)
	router.HandleFunc("/anime", GetAllAnimesHandler).Methods("GET")
	router.HandleFunc("/anime/{id:[0-9]+}", GetAnimeHandler).Methods("GET")
	router.HandleFunc("/anime/{id:[0-9]+}/reviews", GetPaginatedReviewsForAnimeHandler).Methods("GET")

	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/anime").Subrouter()
	protectedRouter.Use(middleware.Authenticate) // Apply authentication middleware

	// Protected routes (require authentication)
	protectedRouter.HandleFunc("", CreateAnimeHandler).Methods("POST")
	protectedRouter.HandleFunc("/{id:[0-9]+}", UpdateAnimeHandler).Methods("PUT")
	protectedRouter.HandleFunc("/{id:[0-9]+}", DeleteAnimeHandler).Methods("DELETE")
}
