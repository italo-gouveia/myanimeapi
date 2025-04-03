package handlers

import (
	"encoding/json"
	"log"
	"myanimeapi/api/middleware"
	"myanimeapi/api/models"
	"myanimeapi/internal/db"
	"myanimeapi/internal/errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// FavoriteHandler handles favorite-related operations
// @Description Handler for managing user's favorite anime
type FavoriteHandler struct {
	db db.DBInterface
}

// NewFavoriteHandler creates a new FavoriteHandler instance
// @Description Creates a new instance of FavoriteHandler with the provided database interface
func NewFavoriteHandler(db db.DBInterface) *FavoriteHandler {
	return &FavoriteHandler{db: db}
}

// AddFavoriteHandler handles adding an anime to favorites
// @Summary Add an anime to user's favorites
// @Description Adds a specific anime to the authenticated user's favorites list
// @Tags favorites
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param request body models.FavoriteCreateRequest true "Anime ID to add to favorites"
// @Success 201 {object} models.Favorite "Favorite created successfully"
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 404 {object} errors.ErrorResponse "Anime not found"
// @Failure 409 {object} errors.ErrorResponse "Anime already in favorites"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /favorites [post]
func (h *FavoriteHandler) AddFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("AddFavoriteHandler: Starting to process request")

	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserContextKey).(uint)
	if !ok {
		log.Printf("AddFavoriteHandler: Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "Unauthorized", "Invalid or missing user ID in context")
		return
	}
	log.Printf("AddFavoriteHandler: Got user ID from context: %d", userID)

	// Get user from database
	var user models.User
	if err := h.db.WithContext(r.Context()).First(&user, userID).Error; err != nil {
		log.Printf("AddFavoriteHandler: Failed to get user from database: %v", err)
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "Unauthorized", "User not found")
		return
	}
	log.Printf("AddFavoriteHandler: Found user in database: %s", user.Username)

	// Parse request body
	var req models.FavoriteCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("AddFavoriteHandler: Failed to decode request body: %v", err)
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid request body", err.Error())
		return
	}
	log.Printf("AddFavoriteHandler: Successfully decoded request body, anime_id: %d", req.AnimeID)

	// Check if anime exists
	var anime models.Anime
	if err := h.db.WithContext(r.Context()).First(&anime, req.AnimeID).Error; err != nil {
		log.Printf("AddFavoriteHandler: Failed to find anime: %v", err)
		errors.WriteErrorResponse(w, http.StatusNotFound, errors.ErrResourceNotFound, "Anime not found", err.Error())
		return
	}
	log.Printf("AddFavoriteHandler: Found anime in database: %s", anime.Title)

	// Check if favorite already exists
	var existingFavorite models.Favorite
	if err := h.db.WithContext(r.Context()).Where("user_id = ? AND anime_id = ?", user.ID, req.AnimeID).First(&existingFavorite).Error; err == nil {
		log.Printf("AddFavoriteHandler: Anime already in favorites for user %d", user.ID)
		errors.WriteErrorResponse(w, http.StatusConflict, errors.ErrConflict, "Anime already in favorites", "This anime is already in your favorites list")
		return
	}
	log.Printf("AddFavoriteHandler: Anime not already in favorites, proceeding to create")

	// Create new favorite
	favorite := models.Favorite{
		UserID:    user.ID,
		AnimeID:   req.AnimeID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.db.WithContext(r.Context()).Create(&favorite).Error; err != nil {
		log.Printf("AddFavoriteHandler: Failed to create favorite: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to add favorite", err.Error())
		return
	}
	log.Printf("AddFavoriteHandler: Successfully created favorite for user %d and anime %d", user.ID, req.AnimeID)

	// Preload the Anime and User data for the response
	if err := h.db.WithContext(r.Context()).Preload("Anime").Preload("User").First(&favorite, favorite.ID).Error; err != nil {
		log.Printf("AddFavoriteHandler: Failed to preload related data: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to load favorite details", err.Error())
		return
	}
	log.Printf("AddFavoriteHandler: Successfully preloaded related data")

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(favorite); err != nil {
		log.Printf("AddFavoriteHandler: Failed to encode response: %v", err)
		return
	}
	log.Printf("AddFavoriteHandler: Successfully sent response")
}

// RemoveFavoriteHandler handles removing an anime from favorites
// @Summary Remove an anime from user's favorites
// @Description Removes a specific anime from the authenticated user's favorites list
// @Tags favorites
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param id path string true "Anime ID to remove from favorites"
// @Success 204 "Favorite removed successfully"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /favorites/{id} [delete]
func (h *FavoriteHandler) RemoveFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserContextKey).(uint)
	if !ok {
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "Unauthorized", "Invalid or missing user ID in context")
		return
	}

	// Get user from database
	var user models.User
	if err := h.db.WithContext(r.Context()).First(&user, userID).Error; err != nil {
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "Unauthorized", "User not found")
		return
	}

	// Get anime ID from URL
	vars := mux.Vars(r)
	animeID := vars["id"]

	// Delete favorite
	if err := h.db.WithContext(r.Context()).Where("user_id = ? AND anime_id = ?", user.ID, animeID).Delete(&models.Favorite{}).Error; err != nil {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to remove favorite", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetFavoritesHandler handles retrieving a user's favorites
// @Summary Get user's favorite anime list
// @Description Retrieves all anime in the authenticated user's favorites list
// @Tags favorites
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {array} models.Favorite "List of user's favorite anime"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /favorites [get]
func (h *FavoriteHandler) GetFavoritesHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserContextKey).(uint)
	if !ok {
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "Unauthorized", "Invalid or missing user ID in context")
		return
	}

	// Get user from database
	var user models.User
	if err := h.db.First(r.Context(), &user, userID).Error; err != nil {
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "Unauthorized", "User not found")
		return
	}

	var favorites []models.Favorite
	if err := h.db.WithContext(r.Context()).Preload("Anime").Preload("User").Where("user_id = ?", user.ID).Find(&favorites).Error; err != nil {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to get favorites", err.Error())
		return
	}

	json.NewEncoder(w).Encode(favorites)
}

// RegisterFavoriteRoutes registers all favorite-related routes
// @Description Registers all favorite-related HTTP routes with the provided router
func (h *FavoriteHandler) RegisterFavoriteRoutes(router *mux.Router) {
	// Create a subrouter for favorites with authentication middleware
	favoritesRouter := router.PathPrefix("/favorites").Subrouter()
	favoritesRouter.Use(middleware.Authenticate)

	// Register routes on the authenticated subrouter
	favoritesRouter.HandleFunc("", h.GetFavoritesHandler).Methods("GET")
	favoritesRouter.HandleFunc("", h.AddFavoriteHandler).Methods("POST")
	favoritesRouter.HandleFunc("/{id}", h.RemoveFavoriteHandler).Methods("DELETE")
}
