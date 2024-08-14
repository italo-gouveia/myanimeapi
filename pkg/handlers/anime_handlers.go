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
func GetAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var anime models.Anime
	if err := database.First(&anime, id).Error; err != nil {
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(anime)
}

// GetAllAnimesHandler retrieves all anime entries
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
func CreateAnimeHandler(w http.ResponseWriter, r *http.Request) {
	var anime models.Anime
	if err := json.NewDecoder(r.Body).Decode(&anime); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
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
func DeleteAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := database.Delete(&models.Anime{}, id).Error; err != nil {
		http.Error(w, "Failed to delete anime", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RegisterAnimeRoutes registers all anime-related routes
func RegisterAnimeRoutes(router *mux.Router) {
	router.Use(middleware.Authenticate)
	router.HandleFunc("/anime", GetAllAnimesHandler).Methods("GET")
	router.HandleFunc("/anime/{id:[0-9]+}", GetAnimeHandler).Methods("GET")
	router.HandleFunc("/anime", CreateAnimeHandler).Methods("POST")
	router.HandleFunc("/anime/{id:[0-9]+}", UpdateAnimeHandler).Methods("PUT")
	router.HandleFunc("/anime/{id:[0-9]+}", DeleteAnimeHandler).Methods("DELETE")
	router.HandleFunc("/anime/{id:[0-9]+}/reviews", GetPaginatedReviewsForAnimeHandler).Methods("GET")
}
