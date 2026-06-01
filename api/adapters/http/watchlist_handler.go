package httphandler

import (
	"encoding/json"
	"net/http"

	"myanimeapi/api/middleware"
	"myanimeapi/api/models"
	"myanimeapi/api/services"
	"myanimeapi/api/utils"
	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"

	"github.com/gorilla/mux"
)

// WatchlistHandler handles watchlist-related HTTP requests
type WatchlistHandler struct {
	watchlistService services.WatchlistServiceInterface
	logger           *logger.Logger
}

// NewWatchlistHandler creates a new WatchlistHandler instance
func NewWatchlistHandler(watchlistService services.WatchlistServiceInterface) *WatchlistHandler {
	return &WatchlistHandler{
		watchlistService: watchlistService,
		logger:           logger.New(),
	}
}

// watchlistGrouped is the grouped-by-status response shape
type watchlistGrouped map[models.WatchlistStatus][]models.WatchlistEntryResponse

// watchlistResponse is the full response for GET /watchlist
type watchlistResponse struct {
	Entries []models.WatchlistEntryResponse `json:"entries"`
	Grouped watchlistGrouped                `json:"grouped"`
}

// GetWatchlistHandler retrieves the authenticated user's watchlist, grouped by status
// @Summary Get user's watchlist
// @Description Get all watchlist entries for the authenticated user, with a grouped view
// @Tags watchlist
// @Produce json
// @Success 200 {object} watchlistResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /watchlist [get]
// @Security BearerAuth
func (h *WatchlistHandler) GetWatchlistHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Starting to retrieve watchlist")

	userID := middleware.GetUserFromContext(r.Context())
	if userID == 0 {
		h.logger.Error("Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "User not authenticated", "Authentication required to access this resource", nil)
		return
	}
	h.logger.WithField("user_id", userID).Info("Retrieved user ID from context")

	entries, err := h.watchlistService.GetWatchlist(r.Context(), userID)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to retrieve watchlist")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Database error", "Failed to retrieve watchlist from database", nil)
		return
	}

	// Build flat entries list and grouped map
	responseEntries := make([]models.WatchlistEntryResponse, 0, len(entries))
	grouped := watchlistGrouped{
		models.WatchlistStatusPlanToWatch: {},
		models.WatchlistStatusWatching:    {},
		models.WatchlistStatusCompleted:   {},
		models.WatchlistStatusDropped:     {},
		models.WatchlistStatusOnHold:      {},
	}

	for _, e := range entries {
		resp := e.ToResponse()
		responseEntries = append(responseEntries, resp)
		grouped[e.Status] = append(grouped[e.Status], resp)
	}

	result := watchlistResponse{
		Entries: responseEntries,
		Grouped: grouped,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to encode watchlist response")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response", nil)
		return
	}
	h.logger.WithFields(map[string]interface{}{
		"user_id": userID,
		"count":   len(entries),
	}).Info("Successfully retrieved watchlist")
}

// UpsertWatchlistHandler adds or updates an anime in the user's watchlist
// @Summary Add or update watchlist entry
// @Description Add an anime to the authenticated user's watchlist or update its status
// @Tags watchlist
// @Accept json
// @Produce json
// @Param animeId path int true "Anime ID"
// @Param body body models.WatchlistUpsertRequest true "Watchlist status"
// @Success 200 {object} models.WatchlistEntryResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /watchlist/{animeId} [put]
// @Security BearerAuth
func (h *WatchlistHandler) UpsertWatchlistHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Starting to upsert watchlist entry")

	userID := middleware.GetUserFromContext(r.Context())
	if userID == 0 {
		h.logger.Error("Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "User not authenticated", "Authentication required to access this resource", nil)
		return
	}
	h.logger.WithField("user_id", userID).Info("Retrieved user ID from context")

	vars := mux.Vars(r)
	animeIDStr := vars["animeId"]
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

	var req models.WatchlistUpsertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to decode request body")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid request body", "Failed to parse the request body as JSON", nil)
		return
	}

	if _, valid := models.ValidWatchlistStatuses[req.Status]; !valid {
		h.logger.WithFields(map[string]interface{}{
			"status": req.Status,
		}).Error("Invalid watchlist status")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid watchlist status", "The provided status is not a valid watchlist status value", map[string]interface{}{
			"status":          req.Status,
			"allowed_statuses": models.ValidWatchlistStatuses,
		})
		return
	}

	entry, err := h.watchlistService.UpsertWatchlistEntry(r.Context(), userID, animeID, req.Status)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"user_id":  userID,
			"anime_id": animeID,
			"error":    err.Error(),
		}).Error("Failed to upsert watchlist entry")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Database error", "Failed to upsert watchlist entry", nil)
		return
	}

	resp := entry.ToResponse()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
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
		"status":   req.Status,
	}).Info("Successfully upserted watchlist entry")
}

// RemoveWatchlistHandler removes an anime from the user's watchlist
// @Summary Remove watchlist entry
// @Description Remove an anime from the authenticated user's watchlist
// @Tags watchlist
// @Produce json
// @Param animeId path int true "Anime ID"
// @Success 204 "No Content"
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /watchlist/{animeId} [delete]
// @Security BearerAuth
func (h *WatchlistHandler) RemoveWatchlistHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Starting to remove watchlist entry")

	userID := middleware.GetUserFromContext(r.Context())
	if userID == 0 {
		h.logger.Error("Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "User not authenticated", "Authentication required to access this resource", nil)
		return
	}
	h.logger.WithField("user_id", userID).Info("Retrieved user ID from context")

	vars := mux.Vars(r)
	animeIDStr := vars["animeId"]
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

	if err := h.watchlistService.RemoveWatchlistEntry(r.Context(), userID, animeID); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"user_id":  userID,
			"anime_id": animeID,
			"error":    err.Error(),
		}).Error("Failed to remove watchlist entry")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Database error", "Failed to remove watchlist entry", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	h.logger.WithFields(map[string]interface{}{
		"user_id":  userID,
		"anime_id": animeID,
	}).Info("Successfully removed watchlist entry")
}

// RegisterWatchlistRoutes registers all watchlist-related routes
func (h *WatchlistHandler) RegisterWatchlistRoutes(router *mux.Router) {
	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/watchlist").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware)

	// Protected routes
	protectedRouter.HandleFunc("", h.GetWatchlistHandler).Methods("GET")
	protectedRouter.HandleFunc("/{animeId}", h.UpsertWatchlistHandler).Methods("PUT")
	protectedRouter.HandleFunc("/{animeId}", h.RemoveWatchlistHandler).Methods("DELETE")
}
