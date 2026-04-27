package httphandler

import (
	"encoding/json"
	"net/http"

	"myanimeapi/api/middleware"
	"myanimeapi/api/services"
	"myanimeapi/api/utils"
	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"

	"github.com/gorilla/mux"
)

// FavoriteHandler handles favorite-related HTTP requests
type FavoriteHandler struct {
	favoriteService services.FavoriteServiceInterface
	logger          *logger.Logger
}

// NewFavoriteHandler creates a new FavoriteHandler instance
func NewFavoriteHandler(favoriteService services.FavoriteServiceInterface) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteService: favoriteService,
		logger:          logger.New(),
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
	h.logger.Info("Starting to add favorite")

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value(middleware.UserContextKey).(uint)
	if !ok {
		h.logger.Error("Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "User not authenticated", "Authentication required to access this resource", nil)
		return
	}
	h.logger.WithField("user_id", userID).Info("Retrieved user ID from context")

	// Get anime ID from URL
	vars := mux.Vars(r)
	animeIDStr := vars["anime_id"]
	animeID, err := utils.ValidateID(animeIDStr)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"anime_id": animeIDStr,
			"error":    err.Error(),
		}).Error("Invalid anime ID format")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid anime ID format", "The provided anime ID must be a valid unsigned integer", nil)
		return
	}
	h.logger.WithField("anime_id", animeID).Info("Retrieved anime ID from URL")

	// Add favorite using service
	favorite, err := h.favoriteService.AddFavorite(r.Context(), userID, animeID)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"user_id":  userID,
			"anime_id": animeID,
			"error":    err.Error(),
		}).Error("Failed to add favorite")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Database error", "Failed to add favorite to database", nil)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(favorite); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"user_id":  userID,
			"anime_id": animeID,
			"error":    err.Error(),
		}).Error("Failed to encode response")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response", nil)
		return
	}
	h.logger.WithFields(map[string]interface{}{
		"user_id":  userID,
		"anime_id": animeID,
	}).Info("Successfully added favorite")
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
	h.logger.Info("Starting to remove favorite")

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value(middleware.UserContextKey).(uint)
	if !ok {
		h.logger.Error("Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "User not authenticated", "Authentication required to access this resource", nil)
		return
	}
	h.logger.WithField("user_id", userID).Info("Retrieved user ID from context")

	// Get anime ID from URL
	vars := mux.Vars(r)
	animeIDStr := vars["anime_id"]
	animeID, err := utils.ValidateID(animeIDStr)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"anime_id": animeIDStr,
			"error":    err.Error(),
		}).Error("Invalid anime ID format")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid anime ID format", "The provided anime ID must be a valid unsigned integer", nil)
		return
	}
	h.logger.WithField("anime_id", animeID).Info("Retrieved anime ID from URL")

	// Remove favorite using service
	if err := h.favoriteService.RemoveFavorite(r.Context(), userID, animeID); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"user_id":  userID,
			"anime_id": animeID,
			"error":    err.Error(),
		}).Error("Failed to remove favorite")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Database error", "Failed to remove favorite from database", nil)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusNoContent)
	h.logger.WithFields(map[string]interface{}{
		"user_id":  userID,
		"anime_id": animeID,
	}).Info("Successfully removed favorite")
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
	h.logger.Info("Starting to retrieve favorites")

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value(middleware.UserContextKey).(uint)
	if !ok {
		h.logger.Error("Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "User not authenticated", "Authentication required to access this resource", nil)
		return
	}
	h.logger.WithField("user_id", userID).Info("Retrieved user ID from context")

	// Get favorites using service
	favorites, err := h.favoriteService.GetFavorites(r.Context(), userID)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to retrieve favorites")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Database error", "Failed to retrieve favorites from database", nil)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(favorites); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to encode response")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response", nil)
		return
	}
	h.logger.WithFields(map[string]interface{}{
		"user_id": userID,
		"count":   len(favorites),
	}).Info("Successfully retrieved favorites")
}

// RegisterFavoriteRoutes registers all favorite-related routes
func (h *FavoriteHandler) RegisterFavoriteRoutes(router *mux.Router) {
	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/favorites").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware) // Apply authentication middleware

	// Protected routes
	protectedRouter.HandleFunc("", h.GetFavoritesHandler).Methods("GET")
	protectedRouter.HandleFunc("/{anime_id}", h.AddFavoriteHandler).Methods("POST")
	protectedRouter.HandleFunc("/{anime_id}", h.RemoveFavoriteHandler).Methods("DELETE")
}
