package handlers

import (
	"encoding/json"
	"log"
	"myanimeapi/api/services"
	"myanimeapi/api/utils"
	"myanimeapi/internal/errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// FavoriteHandler handles favorite-related HTTP requests
type FavoriteHandler struct {
	favoriteService *services.FavoriteService
}

// NewFavoriteHandler creates a new FavoriteHandler instance
func NewFavoriteHandler(favoriteService *services.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteService: favoriteService,
	}
}

// AddFavoriteHandler handles adding an anime to user's favorites
// @Summary Add anime to favorites
// @Description Add an anime to the authenticated user's favorites
// @Tags favorites
// @Accept json
// @Produce json
// @Param anime_id path int true "Anime ID"
// @Success 201 {object} models.Favorite
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /favorites/{anime_id} [post]
// @Security BearerAuth
func (h *FavoriteHandler) AddFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("AddFavoriteHandler: Starting to add favorite")

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		log.Println("AddFavoriteHandler: Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "User not authenticated", "")
		return
	}
	log.Printf("AddFavoriteHandler: Retrieved user ID from context: %d", userID)

	// Get anime ID from URL
	vars := mux.Vars(r)
	animeIDStr := vars["anime_id"]
	animeID, err := strconv.ParseUint(animeIDStr, 10, 32)
	if err != nil {
		log.Printf("AddFavoriteHandler: Invalid anime ID format: %v", err)
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid anime ID format", err.Error())
		return
	}
	log.Printf("AddFavoriteHandler: Retrieved anime ID from URL: %d", animeID)

	// Add favorite using service
	favorite, err := h.favoriteService.AddFavorite(r.Context(), userID, uint(animeID))
	if err != nil {
		log.Printf("AddFavoriteHandler: Failed to add favorite: %v", err)
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to add favorite", err.Error())
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(favorite); err != nil {
		log.Printf("AddFavoriteHandler: Failed to encode response: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", err.Error())
		return
	}
	log.Printf("AddFavoriteHandler: Successfully added favorite for user %d and anime %d", userID, animeID)
}

// RemoveFavoriteHandler handles removing an anime from user's favorites
// @Summary Remove anime from favorites
// @Description Remove an anime from the authenticated user's favorites
// @Tags favorites
// @Produce json
// @Param anime_id path int true "Anime ID"
// @Success 204 "No Content"
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /favorites/{anime_id} [delete]
// @Security BearerAuth
func (h *FavoriteHandler) RemoveFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("RemoveFavoriteHandler: Starting to remove favorite")

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		log.Println("RemoveFavoriteHandler: Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "User not authenticated", "")
		return
	}
	log.Printf("RemoveFavoriteHandler: Retrieved user ID from context: %d", userID)

	// Get anime ID from URL
	vars := mux.Vars(r)
	animeIDStr := vars["anime_id"]
	animeID, err := utils.ValidateID(animeIDStr)
	if err != nil {
		log.Printf("RemoveFavoriteHandler: Invalid anime ID format: %v", err)
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, err.Error(), "")
		return
	}
	log.Printf("RemoveFavoriteHandler: Retrieved anime ID from URL: %d", animeID)

	// Remove favorite using service
	if err := h.favoriteService.RemoveFavorite(r.Context(), userID, animeID); err != nil {
		log.Printf("RemoveFavoriteHandler: Failed to remove favorite: %v", err)
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to remove favorite", err.Error())
		return
	}

	// Return success response
	w.WriteHeader(http.StatusNoContent)
	log.Printf("RemoveFavoriteHandler: Successfully removed favorite for user %d and anime %d", userID, animeID)
}

// GetFavoritesHandler handles retrieving all favorites for a user
// @Summary Get user's favorites
// @Description Get all favorites for the authenticated user
// @Tags favorites
// @Produce json
// @Success 200 {array} models.Favorite
// @Failure 401 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /favorites [get]
// @Security BearerAuth
func (h *FavoriteHandler) GetFavoritesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetFavoritesHandler: Starting to retrieve favorites")

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		log.Println("GetFavoritesHandler: Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "User not authenticated", "")
		return
	}
	log.Printf("GetFavoritesHandler: Retrieved user ID from context: %d", userID)

	// Get favorites using service
	favorites, err := h.favoriteService.GetFavorites(r.Context(), userID)
	if err != nil {
		log.Printf("GetFavoritesHandler: Failed to retrieve favorites: %v", err)
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to retrieve favorites", err.Error())
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(favorites); err != nil {
		log.Printf("GetFavoritesHandler: Failed to encode response: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", err.Error())
		return
	}
	log.Printf("GetFavoritesHandler: Successfully retrieved %d favorites for user %d", len(favorites), userID)
}

// RegisterFavoriteRoutes registers all favorite-related routes
func (h *FavoriteHandler) RegisterFavoriteRoutes(router *mux.Router) {
	router.HandleFunc("/favorites", h.GetFavoritesHandler).Methods("GET")
	router.HandleFunc("/favorites/{anime_id}", h.AddFavoriteHandler).Methods("POST")
	router.HandleFunc("/favorites/{anime_id}", h.RemoveFavoriteHandler).Methods("DELETE")
}
