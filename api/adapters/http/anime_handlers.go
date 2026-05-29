// api/adapters/http/anime_handlers.go
// Package httphandler provides HTTP handlers for anime-related routes in the MyAnimeAPI application.
// It defines methods to handle requests for retrieving, creating, updating, and deleting anime entries.
// The package uses the Gorilla Mux router for routing, GORM for database interactions, and middleware for request validation and authentication.
//
// Example usage:
//
//	animeService := services.NewAnimeService(repository)
//	animeHandler := httphandler.NewAnimeHandler(animeService)
//	router := mux.NewRouter()
//	animeHandler.RegisterAnimeRoutes(router)
//
//	http.ListenAndServe(":8080", router)
package httphandler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"myanimeapi/api/middleware"
	"myanimeapi/api/models"
	"myanimeapi/api/services"
	"myanimeapi/api/utils"
	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"

	"github.com/gorilla/mux"
)

// AnimeHandler defines the handlers for anime-related routes.
// It contains an anime service for handling business logic.
type AnimeHandler struct {
	service services.AnimeServiceInterface
}

// NewAnimeHandler creates a new instance of AnimeHandler.
// It accepts an anime service interface and returns a pointer to an AnimeHandler.
//
// Example:
//
//	animeService := services.NewAnimeService(repository)
//	animeHandler := NewAnimeHandler(animeService)
func NewAnimeHandler(service services.AnimeServiceInterface) *AnimeHandler {
	return &AnimeHandler{
		service: service,
	}
}

// CreateAnimeHandler handles the creation of a new anime.
// It validates the input payload and uses the service to create the anime.
// If successful, it returns the created anime as a JSON response.
//
// @Summary Create a new anime
// @Description Create a new anime with the provided details
// @Tags anime
// @Accept json
// @Produce json
// @Param anime body models.AnimeCreateRequest true "Anime details"
// @Success 201 {object} models.AnimeResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 500 {object} errors.ErrorResponse "Failed to create anime"
// @Router /animes [post]
// @Security BearerAuth
func (h *AnimeHandler) CreateAnimeHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the validated and sanitized payload from the context
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.AnimeCreateRequest)
	if !ok {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Invalid payload", "The request payload could not be retrieved.", nil)
		return
	}

	// Create the anime
	anime := &models.Anime{
		Title:       payload.Title,
		Description: payload.Description,
		Rating:      payload.Rating,
		Episodes:    payload.Episodes,
		Status:      payload.Status,
		StartDate:   payload.StartDate,
		EndDate:     payload.EndDate,
	}

	// Create the anime first
	if err := h.service.CreateAnime(r.Context(), anime); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to create anime", "An internal server error occurred while creating the anime.", nil)
		return
	}

	// Add genres if provided
	if len(payload.GenreIDs) > 0 {
		if err := h.service.AddGenresToAnime(r.Context(), anime.ID, payload.GenreIDs); err != nil {
			if rollbackErr := h.service.DeleteAnime(r.Context(), anime.ID); rollbackErr != nil {
				logger.Get().WithField("anime_id", anime.ID).WithField("error", rollbackErr.Error()).Error("Failed to rollback anime after genre association error")
			}
			if appErr, ok := err.(*errors.AppError); ok {
				errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
				return
			}
			errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to add genres to anime", "An internal server error occurred while adding genres to the anime.", nil)
			return
		}
	}

	// Add tags if provided
	if len(payload.TagIDs) > 0 {
		if err := h.service.AddTagsToAnime(r.Context(), anime.ID, payload.TagIDs); err != nil {
			if rollbackErr := h.service.DeleteAnime(r.Context(), anime.ID); rollbackErr != nil {
				logger.Get().WithField("anime_id", anime.ID).WithField("error", rollbackErr.Error()).Error("Failed to rollback anime after tag association error")
			}
			if appErr, ok := err.(*errors.AppError); ok {
				errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
				return
			}
			errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to add tags to anime", "An internal server error occurred while adding tags to the anime.", nil)
			return
		}
	}

	// Get the complete anime with genres and tags
	createdAnime, err := h.service.GetAnimeByID(r.Context(), anime.ID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to retrieve created anime", "An internal server error occurred while retrieving the created anime.", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdAnime.ToResponse()); err != nil {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.", nil)
		return
	}
}

// GetAnimeHandler handles retrieving an anime by ID.
func (h *AnimeHandler) GetAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := utils.ValidateID(vars["id"])
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", "The provided ID is not a valid unsigned integer.", map[string]interface{}{
			"id": vars["id"],
		})
		return
	}

	anime, err := h.service.GetAnimeByID(r.Context(), id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to retrieve anime", "An internal server error occurred while retrieving the anime.", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(anime.ToResponse()); err != nil {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.", nil)
		return
	}
}

// GetAllAnimesHandler handles retrieving all animes with pagination, filtering, and sorting.
//
// @Summary      List animes
// @Description  Returns a paginated, filterable, sortable list of animes.
// @Tags         anime
// @Produce      json
// @Param        page          query  int     false  "Page number (default 1)"
// @Param        limit         query  int     false  "Items per page (default 100)"
// @Param        status        query  string  false  "Filter by status"                         Enums(Airing,Completed,Upcoming)
// @Param        genre         query  string  false  "Filter by exact genre name (single)"
// @Param        genres        query  string  false  "Filter by multiple genres, comma-separated (AND logic)"
// @Param        tag           query  string  false  "Filter by exact tag name (single)"
// @Param        tags          query  string  false  "Filter by multiple tags, comma-separated (AND logic)"
// @Param        rating_min    query  number  false  "Minimum rating (0–10)"
// @Param        rating_max    query  number  false  "Maximum rating (0–10)"
// @Param        episodes_min  query  int     false  "Minimum episode count"
// @Param        episodes_max  query  int     false  "Maximum episode count"
// @Param        year_from     query  int     false  "Start-date year lower bound (inclusive)"
// @Param        year_to       query  int     false  "Start-date year upper bound (inclusive)"
// @Param        sort_by       query  string  false  "Sort field"                               Enums(title,rating,episodes,created_at,start_date)
// @Param        order         query  string  false  "Sort direction"                           Enums(asc,desc)
// @Success      200           {object}  object{data=[]models.AnimeResponse,total=integer,page=integer,limit=integer}
// @Failure      400           {object}  errors.ErrorResponse  "Invalid filter or pagination"
// @Failure      500           {object}  errors.ErrorResponse  "Internal server error"
// @Router       /animes [get]
func (h *AnimeHandler) GetAllAnimesHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// Pagination
	page, limit, err := utils.ValidatePagination(q.Get("page"), q.Get("limit"), 1, 100)
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid pagination parameters", err.Error(), nil)
		return
	}

	// Filter & sort
	filter := models.AnimeFilter{
		Status:    q.Get("status"),
		Genre:     q.Get("genre"),
		Tag:       q.Get("tag"),
		SortBy:    q.Get("sort_by"),
		SortOrder: q.Get("order"),
	}
	// Multi-value: genres and tags are comma-separated lists.
	if v := q.Get("genres"); v != "" {
		filter.Genres = splitCSV(v)
	}
	if v := q.Get("tags"); v != "" {
		filter.Tags = splitCSV(v)
	}
	if v := q.Get("rating_min"); v != "" {
		filter.RatingMin, _ = strconv.ParseFloat(v, 64)
	}
	if v := q.Get("rating_max"); v != "" {
		filter.RatingMax, _ = strconv.ParseFloat(v, 64)
	}
	if v := q.Get("episodes_min"); v != "" {
		filter.EpisodesMin, _ = strconv.Atoi(v)
	}
	if v := q.Get("episodes_max"); v != "" {
		filter.EpisodesMax, _ = strconv.Atoi(v)
	}
	if v := q.Get("year_from"); v != "" {
		filter.YearFrom, _ = strconv.Atoi(v)
	}
	if v := q.Get("year_to"); v != "" {
		filter.YearTo, _ = strconv.Atoi(v)
	}
	if err := filter.Validate(); err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid filter parameters", err.Error(), nil)
		return
	}

	animes, total, err := h.service.GetAllAnimes(r.Context(), page, limit, filter)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to retrieve animes", "An internal server error occurred while retrieving the animes.", nil)
		return
	}

	responses := make([]models.AnimeResponse, len(animes))
	for i, anime := range animes {
		responses[i] = anime.ToResponse()
	}

	// Create paginated response
	response := struct {
		Data  []models.AnimeResponse `json:"data"`
		Total int64                  `json:"total"`
		Page  int                    `json:"page"`
		Limit int                    `json:"limit"`
	}{
		Data:  responses,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.", nil)
		return
	}
}

// UpdateAnimeHandler handles updating an existing anime.
func (h *AnimeHandler) UpdateAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := utils.ValidateID(vars["id"])
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", "The provided ID is not a valid unsigned integer.", map[string]interface{}{
			"id": vars["id"],
		})
		return
	}

	// Retrieve the validated and sanitized payload from the context
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.AnimeUpdateRequest)
	if !ok {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Invalid payload", "The request payload could not be retrieved.", nil)
		return
	}

	// Get the existing anime
	anime, err := h.service.GetAnimeByID(r.Context(), id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to retrieve anime", "An internal server error occurred while retrieving the anime.", nil)
		return
	}

	// Update the anime fields if provided in the request
	if payload.Title != "" {
		anime.Title = payload.Title
	}
	if payload.Description != "" {
		anime.Description = payload.Description
	}
	if payload.Episodes != 0 {
		anime.Episodes = payload.Episodes
	}
	if payload.Status != "" {
		anime.Status = payload.Status
	}
	if !payload.StartDate.IsZero() {
		anime.StartDate = payload.StartDate
	}
	if !payload.EndDate.IsZero() {
		anime.EndDate = payload.EndDate
	}
	if payload.Rating != 0 {
		anime.Rating = payload.Rating
	}

	// Update the anime
	if err := h.service.UpdateAnime(r.Context(), anime); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to update anime", "An internal server error occurred while updating the anime.", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(anime.ToResponse()); err != nil {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.", nil)
		return
	}
}

// DeleteAnimeHandler handles deleting an anime.
func (h *AnimeHandler) DeleteAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := utils.ValidateID(vars["id"])
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", "The provided ID is not a valid unsigned integer.", map[string]interface{}{
			"id": vars["id"],
		})
		return
	}

	if err := h.service.DeleteAnime(r.Context(), id); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to delete anime", "An internal server error occurred while deleting the anime.", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddGenresToAnimeHandler handles adding genres to an anime.
func (h *AnimeHandler) AddGenresToAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	animeID, err := utils.ValidateID(vars["id"])
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", "The provided ID is not a valid unsigned integer.", map[string]interface{}{
			"id": vars["id"],
		})
		return
	}

	// Retrieve the validated and sanitized payload from the context
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.AnimeGenresRequest)
	if !ok {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Invalid payload", "The request payload could not be retrieved.", nil)
		return
	}

	if err := h.service.AddGenresToAnime(r.Context(), animeID, payload.GenreIDs); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to add genres to anime", "An internal server error occurred while adding genres to the anime.", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveGenresFromAnimeHandler handles removing genres from an anime.
func (h *AnimeHandler) RemoveGenresFromAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	animeID, err := utils.ValidateID(vars["id"])
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", "The provided ID is not a valid unsigned integer.", map[string]interface{}{
			"id": vars["id"],
		})
		return
	}

	// Retrieve the validated and sanitized payload from the context
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.AnimeGenresRequest)
	if !ok {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Invalid payload", "The request payload could not be retrieved.", nil)
		return
	}

	if err := h.service.RemoveGenresFromAnime(r.Context(), animeID, payload.GenreIDs); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to remove genres from anime", "An internal server error occurred while removing genres from the anime.", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddTagsToAnimeHandler handles adding tags to an anime.
func (h *AnimeHandler) AddTagsToAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	animeID, err := utils.ValidateID(vars["id"])
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", "The provided ID is not a valid unsigned integer.", map[string]interface{}{
			"id": vars["id"],
		})
		return
	}

	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.AnimeTagsRequest)
	if !ok {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Invalid payload", "The request payload could not be retrieved.", nil)
		return
	}

	if err := h.service.AddTagsToAnime(r.Context(), animeID, payload.TagIDs); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to add tags to anime", "An internal server error occurred while adding tags to the anime.", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveTagsFromAnimeHandler handles removing tags from an anime.
func (h *AnimeHandler) RemoveTagsFromAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	animeID, err := utils.ValidateID(vars["id"])
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", "The provided ID is not a valid unsigned integer.", map[string]interface{}{
			"id": vars["id"],
		})
		return
	}

	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.AnimeTagsRequest)
	if !ok {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Invalid payload", "The request payload could not be retrieved.", nil)
		return
	}

	if err := h.service.RemoveTagsFromAnime(r.Context(), animeID, payload.TagIDs); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to remove tags from anime", "An internal server error occurred while removing tags from the anime.", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAnimesByTitleHandler handles searching animes by title with optional sorting.
//
// @Summary      Search animes by title
// @Description  Returns animes whose title contains the given string. Supports sort and pagination.
// @Tags         anime
// @Produce      json
// @Param        title     query  string  true   "Title substring to search for"
// @Param        page      query  int     false  "Page number (default 1)"
// @Param        limit     query  int     false  "Items per page (default 100)"
// @Param        sort_by   query  string  false  "Sort field"     Enums(title,rating,episodes,created_at,start_date)
// @Param        order     query  string  false  "Sort direction" Enums(asc,desc)
// @Success      200       {object}  object{data=[]models.AnimeResponse,total=integer,page=integer,limit=integer}
// @Failure      400       {object}  errors.ErrorResponse
// @Failure      500       {object}  errors.ErrorResponse
// @Router       /animes/search [get]
func (h *AnimeHandler) GetAnimesByTitleHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	title := q.Get("title")
	if title == "" {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid input", "Title parameter is required", nil)
		return
	}

	page, limit, err := utils.ValidatePagination(q.Get("page"), q.Get("limit"), 1, 100)
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid pagination parameters", err.Error(), nil)
		return
	}

	filter := models.AnimeFilter{
		Title:     title,
		SortBy:    q.Get("sort_by"),
		SortOrder: q.Get("order"),
	}
	if err := filter.Validate(); err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid filter parameters", err.Error(), nil)
		return
	}

	animes, total, err := h.service.GetAllAnimes(r.Context(), page, limit, filter)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Database error", "Failed to search animes", nil)
		return
	}

	responses := make([]models.AnimeResponse, len(animes))
	for i, anime := range animes {
		responses[i] = anime.ToResponse()
	}

	response := map[string]interface{}{
		"data":  responses,
		"total": total,
		"page":  page,
		"limit": limit,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.", nil)
		return
	}
}

// GetAnimesByGenreHandler handles searching animes by genre.
func (h *AnimeHandler) GetAnimesByGenreHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	genre := vars["genre"]
	if genre == "" {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid input", "Genre parameter is required", nil)
		return
	}

	page, limit, err := utils.ValidatePagination(r.URL.Query().Get("page"), r.URL.Query().Get("limit"), 1, 100)
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid pagination parameters", err.Error(), nil)
		return
	}

	animes, total, err := h.service.GetAnimesByGenre(r.Context(), genre, page, limit)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Database error", "Failed to search animes by genre", nil)
		return
	}

	response := map[string]interface{}{
		"data":  animes,
		"total": total,
		"page":  page,
		"limit": limit,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.", nil)
		return
	}
}

// splitCSV splits a comma-separated string into a trimmed, non-empty slice.
// Used to parse multi-value query params like ?genres=Action,Adventure.
func splitCSV(s string) []string {
	var result []string
	for _, part := range strings.Split(s, ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// RegisterAnimeRoutes registers all anime-related routes
func (h *AnimeHandler) RegisterAnimeRoutes(router *mux.Router) {
	// Public routes (no authentication required).
	// NOTE: literal segments (/search, /genre/{x}) MUST be registered before the
	// wildcard /{id} so Gorilla Mux does not swallow them as id values.
	router.HandleFunc("/animes", h.GetAllAnimesHandler).Methods("GET")
	router.HandleFunc("/animes/search", h.GetAnimesByTitleHandler).Methods("GET")
	router.HandleFunc("/animes/genre/{genre}", h.GetAnimesByGenreHandler).Methods("GET")
	router.HandleFunc("/animes/{id}", h.GetAnimeHandler).Methods("GET")

	// Create a subrouter for protected routes (requires authentication)
	protectedRouter := router.PathPrefix("/animes").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware)

	// Admin-only routes: create, update, delete anime
	// RequireAdmin is applied per-handler so the middleware chain is explicit.
	protectedRouter.Handle("", middleware.RequireAdmin(
		middleware.ValidateAndSanitizePayload(models.AnimeCreateRequest{})(http.HandlerFunc(h.CreateAnimeHandler)),
	)).Methods("POST")
	protectedRouter.Handle("/{id}", middleware.RequireAdmin(
		middleware.ValidateAndSanitizePayload(models.AnimeUpdateRequest{})(http.HandlerFunc(h.UpdateAnimeHandler)),
	)).Methods("PUT")
	protectedRouter.Handle("/{id}", middleware.RequireAdmin(
		http.HandlerFunc(h.DeleteAnimeHandler),
	)).Methods("DELETE")

	// Genre management routes — admin only
	protectedRouter.Handle("/{id}/genres", middleware.RequireAdmin(
		middleware.ValidateAndSanitizePayload(models.AnimeGenresRequest{})(http.HandlerFunc(h.AddGenresToAnimeHandler)),
	)).Methods("POST")
	protectedRouter.Handle("/{id}/genres", middleware.RequireAdmin(
		middleware.ValidateAndSanitizePayload(models.AnimeGenresRequest{})(http.HandlerFunc(h.RemoveGenresFromAnimeHandler)),
	)).Methods("DELETE")

	// Tag management routes — admin only
	protectedRouter.Handle("/{id}/tags", middleware.RequireAdmin(
		middleware.ValidateAndSanitizePayload(models.AnimeTagsRequest{})(http.HandlerFunc(h.AddTagsToAnimeHandler)),
	)).Methods("POST")
	protectedRouter.Handle("/{id}/tags", middleware.RequireAdmin(
		middleware.ValidateAndSanitizePayload(models.AnimeTagsRequest{})(http.HandlerFunc(h.RemoveTagsFromAnimeHandler)),
	)).Methods("DELETE")
}
