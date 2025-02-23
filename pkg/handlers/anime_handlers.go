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
	"strconv"

	"myanimeapi/internal/db"
	"myanimeapi/pkg/middleware"
	"myanimeapi/pkg/models"

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

	// Convert the ID from string to uint
	id, err := strconv.ParseUint(idStr, 10, 32) // Convert to uint32
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var anime models.Anime
	if err := h.DB.First(r.Context(), &anime, id).Error; err != nil {
		log.Printf("Anime not found: %v", err)
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
func (h *AnimeHandler) GetAllAnimesHandler(w http.ResponseWriter, r *http.Request) {
	var animes []models.Anime
	if err := h.DB.Find(r.Context(), &animes).Error; err != nil {
		log.Printf("Failed to retrieve animes: %v", err)
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
func (h *AnimeHandler) GetPaginatedReviewsForAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert the ID from string to uint
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
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
			log.Printf("Invalid page number: %v", err)
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			log.Printf("Invalid limit number: %v", err)
			http.Error(w, "Invalid limit number", http.StatusBadRequest)
			return
		}
	}

	var reviews []models.Review
	offset := (page - 1) * limit

	if err := h.DB.Where(r.Context(), "anime_id = ?", id).Offset(offset).Limit(limit).Find(&reviews).Error; err != nil {
		log.Printf("Failed to retrieve reviews: %v", err)
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
func (h *AnimeHandler) CreateAnimeHandler(w http.ResponseWriter, r *http.Request) {
	var anime models.Anime
	if err := json.NewDecoder(r.Body).Decode(&anime); err != nil {
		log.Printf("Invalid input: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if title is provided
	if anime.Title == "" {
		log.Println("Title is required")
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if err := h.DB.Create(r.Context(), &anime).Error; err != nil {
		log.Printf("Failed to create anime: %v", err)
		http.Error(w, "Failed to create anime", http.StatusInternalServerError)
		return
	}

	log.Printf("Anime %s created successfully", anime.Title)
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
func (h *AnimeHandler) UpdateAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert the ID from string to uint
	id, err := strconv.ParseUint(idStr, 10, 32) // Convert to uint32
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var anime models.Anime
	if err := h.DB.First(r.Context(), &anime, uint(id)).Error; err != nil {
		log.Printf("Anime not found: %v", err)
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&anime); err != nil {
		log.Printf("Invalid input: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	anime.ID = uint(id) // Ensure the ID is preserved during update
	if err := h.DB.Save(r.Context(), &anime).Error; err != nil {
		log.Printf("Failed to update anime: %v", err)
		http.Error(w, "Failed to update anime", http.StatusInternalServerError)
		return
	}

	log.Printf("Anime %s updated successfully", anime.Title)
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
func (h *AnimeHandler) DeleteAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert the ID from string to uint
	id, err := strconv.ParseUint(idStr, 10, 32) // Convert to uint32
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var anime models.Anime
	if err := h.DB.First(r.Context(), &anime, uint(id)).Error; err != nil {
		log.Printf("Anime not found: %v", err)
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	if err := h.DB.Delete(r.Context(), &models.Anime{}, id).Error; err != nil {
		log.Printf("Failed to delete anime: %v", err)
		http.Error(w, "Failed to delete anime", http.StatusInternalServerError)
		return
	}

	log.Printf("Anime %s deleted successfully", anime.Title)
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
	protectedRouter.HandleFunc("", h.CreateAnimeHandler).Methods("POST")
	protectedRouter.HandleFunc("/{id:[0-9]+}", h.UpdateAnimeHandler).Methods("PUT")
	protectedRouter.HandleFunc("/{id:[0-9]+}", h.DeleteAnimeHandler).Methods("DELETE")
}
