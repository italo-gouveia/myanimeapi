// api/handlers/anime_handlers.go
// Package handlers provides HTTP handlers for anime-related routes in the MyAnimeAPI application.
// It defines methods to handle requests for retrieving, creating, updating, and deleting anime entries.
// The package uses the Gorilla Mux router for routing, GORM for database interactions, and middleware for request validation and authentication.
//
// Example usage:
//
//	animeService := services.NewAnimeService(repository)
//	animeHandler := handlers.NewAnimeHandler(animeService)
//	router := mux.NewRouter()
//	animeHandler.RegisterAnimeRoutes(router)
//
//	http.ListenAndServe(":8080", router)
package handlers

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
// @Example
//
//	{
//	  "title": "Naruto",
//	  "description": "A story about ninjas.",
//	  "episodes": 220,
//	  "status": "Completed",
//	  "start_date": "2002-10-03T00:00:00Z",
//	  "end_date": "2007-02-08T00:00:00Z",
//	  "rating": 8.5,
//	  "genre_ids": [1, 2, 3],
//	  "tag_ids": [1, 2, 3]
//	}
//
// @ExampleResponse
//
//	{
//	  "id": 1,
//	  "title": "Naruto",
//	  "description": "A story about ninjas.",
//	  "rating": 8.5,
//	  "episodes": 220,
//	  "status": "Completed",
//	  "start_date": "2002-10-03T00:00:00Z",
//	  "end_date": "2007-02-08T00:00:00Z",
//	  "created_at": "2025-02-20T19:27:00Z",
//	  "updated_at": "2025-02-20T19:27:00Z",
//	  "genres": [
//	    {
//	      "id": 1,
//	      "name": "Action"
//	    }
//	  ],
//	  "tags": [
//	    {
//	      "id": 1,
//	      "name": "Ninja"
//	    }
//	  ]
//	}
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
// It validates the ID, queries the service, and returns the anime as a JSON response.
// If the ID is invalid or the anime is not found, it returns an appropriate error response.
//
// @Summary Get an anime by ID
// @Description Get an anime's details by its ID
// @Tags anime
// @Produce json
// @Param id path int true "Anime ID"
// @Success 200 {object} models.AnimeResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid anime ID"
// @Failure 404 {object} errors.ErrorResponse "Anime not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve anime"
// @Router /animes/{id} [get]
// @ExampleResponse
//
//	{
//	  "id": 1,
//	  "title": "Naruto",
//	  "description": "A story about ninjas.",
//	  "rating": 8.5,
//	  "episodes": 220,
//	  "status": "Completed",
//	  "start_date": "2002-10-03T00:00:00Z",
//	  "end_date": "2007-02-08T00:00:00Z",
//	  "created_at": "2025-02-20T19:27:00Z",
//	  "updated_at": "2025-02-20T19:27:00Z",
//	  "genres": [
//	    {
//	      "id": 1,
//	      "name": "Action"
//	    }
//	  ],
//	  "tags": [
//	    {
//	      "id": 1,
//	      "name": "Ninja"
//	    }
//	  ]
//	}
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

// GetAllAnimesHandler handles retrieving all animes with pagination.
// It validates the pagination parameters, queries the service, and returns the animes as a JSON response.
// If the pagination parameters are invalid, it returns an appropriate error response.
//
// @Summary Get all animes
// @Description Get a list of all animes with pagination
// @Tags anime
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.AnimeListResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid pagination parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve animes"
// @Router /animes [get]
// @ExampleResponse
//
//	{
//	  "animes": [
//	    {
//	      "id": 1,
//	      "title": "Naruto",
//	      "description": "A story about ninjas.",
//	      "rating": 8.5,
//	      "episodes": 220,
//	      "status": "Completed",
//	      "start_date": "2002-10-03T00:00:00Z",
//	      "end_date": "2007-02-08T00:00:00Z",
//	      "created_at": "2025-02-20T19:27:00Z",
//	      "updated_at": "2025-02-20T19:27:00Z",
//	      "genres": [
//	        {
//	          "id": 1,
//	          "name": "Action"
//	        }
//	      ],
//	      "tags": [
//	        {
//	          "id": 1,
//	          "name": "Ninja"
//	        }
//	      ]
//	    }
//	  ],
//	  "total": 1,
//	  "page": 1,
//	  "limit": 10
//	}
func (h *AnimeHandler) GetAllAnimesHandler(w http.ResponseWriter, r *http.Request) {
	// Get pagination parameters from query
	page, limit, err := utils.ValidatePagination(r.URL.Query().Get("page"), r.URL.Query().Get("limit"), 1, 100)
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid pagination parameters", err.Error(), nil)
		return
	}

	animes, total, err := h.service.GetAllAnimes(r.Context(), page, limit)
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
// It validates the ID and input payload, uses the service to update the anime,
// and returns the updated anime as a JSON response.
//
// @Summary Update an anime
// @Description Update an existing anime's details
// @Tags anime
// @Accept json
// @Produce json
// @Param id path int true "Anime ID"
// @Param anime body models.AnimeUpdateRequest true "Updated anime details"
// @Success 200 {object} models.AnimeResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid anime ID or request body"
// @Failure 404 {object} errors.ErrorResponse "Anime not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to update anime"
// @Router /animes/{id} [put]
// @Security BearerAuth
// @Example
//
//	{
//	  "title": "Naruto Shippuden",
//	  "description": "The continuation of Naruto's story.",
//	  "episodes": 500,
//	  "status": "Completed",
//	  "start_date": "2007-02-15T00:00:00Z",
//	  "end_date": "2017-03-23T00:00:00Z",
//	  "rating": 8.7
//	}
//
// @ExampleResponse
//
//	{
//	  "id": 1,
//	  "title": "Naruto Shippuden",
//	  "description": "The continuation of Naruto's story.",
//	  "rating": 8.7,
//	  "episodes": 500,
//	  "status": "Completed",
//	  "start_date": "2007-02-15T00:00:00Z",
//	  "end_date": "2017-03-23T00:00:00Z",
//	  "created_at": "2025-02-20T19:27:00Z",
//	  "updated_at": "2025-02-20T19:27:00Z",
//	  "genres": [
//	    {
//	      "id": 1,
//	      "name": "Action"
//	    }
//	  ],
//	  "tags": [
//	    {
//	      "id": 1,
//	      "name": "Ninja"
//	    }
//	  ]
//	}
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
// It validates the ID, uses the service to delete the anime,
// and returns a success message as a JSON response.
//
// @Summary Delete an anime
// @Description Delete an anime by its ID
// @Tags anime
// @Produce json
// @Param id path int true "Anime ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid anime ID"
// @Failure 404 {object} errors.ErrorResponse "Anime not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete anime"
// @Router /animes/{id} [delete]
// @Security BearerAuth
// @ExampleResponse
//
//	{
//	  "message": "Anime deleted successfully"
//	}
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
// It validates the ID and input payload, uses the service to add genres,
// and returns a success message as a JSON response.
//
// @Summary Add genres to an anime
// @Description Add one or more genres to an existing anime
// @Tags anime
// @Accept json
// @Produce json
// @Param id path int true "Anime ID"
// @Param genres body []uint true "Genre IDs to add"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid anime ID or request body"
// @Failure 404 {object} errors.ErrorResponse "Anime not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to add genres to anime"
// @Router /animes/{id}/genres [post]
// @Security BearerAuth
// @Example
//
//	{
//	  "genre_ids": [1, 2, 3]
//	}
//
// @ExampleResponse
//
//	{
//	  "message": "Genres added successfully"
//	}
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
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*struct {
		GenreIDs []uint `json:"genre_ids" validate:"required,min=1"`
	})
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
// It validates the ID and input payload, uses the service to remove genres,
// and returns a success message as a JSON response.
//
// @Summary Remove genres from an anime
// @Description Remove one or more genres from an existing anime
// @Tags anime
// @Accept json
// @Produce json
// @Param id path int true "Anime ID"
// @Param genres body []uint true "Genre IDs to remove"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid anime ID or request body"
// @Failure 404 {object} errors.ErrorResponse "Anime not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to remove genres from anime"
// @Router /animes/{id}/genres [delete]
// @Security BearerAuth
// @Example
//
//	{
//	  "genre_ids": [1, 2, 3]
//	}
//
// @ExampleResponse
//
//	{
//	  "message": "Genres removed successfully"
//	}
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
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*struct {
		GenreIDs []uint `json:"genre_ids" validate:"required,min=1"`
	})
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
// It validates the ID and input payload, uses the service to add tags,
// and returns a success message as a JSON response.
//
// @Summary Add tags to an anime
// @Description Add one or more tags to an existing anime
// @Tags anime
// @Accept json
// @Produce json
// @Param id path int true "Anime ID"
// @Param tags body []uint true "Tag IDs to add"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid anime ID or request body"
// @Failure 404 {object} errors.ErrorResponse "Anime not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to add tags to anime"
// @Router /animes/{id}/tags [post]
// @Security BearerAuth
// @Example
//
//	{
//	  "tag_ids": [1, 2, 3]
//	}
//
// @ExampleResponse
//
//	{
//	  "message": "Tags added successfully"
//	}
func (h *AnimeHandler) AddTagsToAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	animeID, err := utils.ValidateID(vars["id"])
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", "The provided ID is not a valid unsigned integer.", map[string]interface{}{
			"id": vars["id"],
		})
		return
	}

	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*struct {
		TagIDs []uint `json:"tag_ids" validate:"required,min=1"`
	})
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
// It validates the ID and input payload, uses the service to remove tags,
// and returns a success message as a JSON response.
//
// @Summary Remove tags from an anime
// @Description Remove one or more tags from an existing anime
// @Tags anime
// @Accept json
// @Produce json
// @Param id path int true "Anime ID"
// @Param tags body []uint true "Tag IDs to remove"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid anime ID or request body"
// @Failure 404 {object} errors.ErrorResponse "Anime not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to remove tags from anime"
// @Router /animes/{id}/tags [delete]
// @Security BearerAuth
// @Example
//
//	{
//	  "tag_ids": [1, 2, 3]
//	}
//
// @ExampleResponse
//
//	{
//	  "message": "Tags removed successfully"
//	}
func (h *AnimeHandler) RemoveTagsFromAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	animeID, err := utils.ValidateID(vars["id"])
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", "The provided ID is not a valid unsigned integer.", map[string]interface{}{
			"id": vars["id"],
		})
		return
	}

	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*struct {
		TagIDs []uint `json:"tag_ids" validate:"required,min=1"`
	})
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

// GetAnimesByTitleHandler handles searching animes by title.
// It validates the title parameter, uses the service to search for animes,
// and returns the matching animes as a JSON response.
//
// @Summary Search animes by title
// @Description Search for animes by title with pagination
// @Tags anime
// @Produce json
// @Param title query string true "Title to search for"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.AnimeListResponse
// @Failure 400 {object} errors.ErrorResponse "Title parameter is required"
// @Failure 500 {object} errors.ErrorResponse "Failed to search animes by title"
// @Router /animes/search [get]
// @ExampleResponse
//
//	{
//	  "animes": [
//	    {
//	      "id": 1,
//	      "title": "Naruto",
//	      "description": "A story about ninjas.",
//	      "rating": 8.5,
//	      "episodes": 220,
//	      "status": "Completed",
//	      "start_date": "2002-10-03T00:00:00Z",
//	      "end_date": "2007-02-08T00:00:00Z",
//	      "created_at": "2025-02-20T19:27:00Z",
//	      "updated_at": "2025-02-20T19:27:00Z",
//	      "genres": [
//	        {
//	          "id": 1,
//	          "name": "Action"
//	        }
//	      ],
//	      "tags": [
//	        {
//	          "id": 1,
//	          "name": "Ninja"
//	        }
//	      ]
//	    }
//	  ],
//	  "total": 1,
//	  "page": 1,
//	  "limit": 10
//	}
func (h *AnimeHandler) GetAnimesByTitleHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	if title == "" {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid input", "Title parameter is required", nil)
		return
	}

	page, limit, err := utils.ValidatePagination(r.URL.Query().Get("page"), r.URL.Query().Get("limit"), 1, 100)
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid pagination parameters", err.Error(), nil)
		return
	}

	animes, total, err := h.service.GetAnimesByTitle(r.Context(), title, page, limit)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Database error", "Failed to search animes", nil)
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

// GetAnimesByGenreHandler handles searching animes by genre.
// It validates the genre parameter, uses the service to search for animes,
// and returns the matching animes as a JSON response.
//
// @Summary Search animes by genre
// @Description Search for animes by genre with pagination
// @Tags anime
// @Produce json
// @Param genre path string true "Genre to search for"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.AnimeListResponse
// @Failure 400 {object} errors.ErrorResponse "Genre parameter is required"
// @Failure 500 {object} errors.ErrorResponse "Failed to search animes by genre"
// @Router /animes/genre/{genre} [get]
// @ExampleResponse
//
//	{
//	  "animes": [
//	    {
//	      "id": 1,
//	      "title": "Naruto",
//	      "description": "A story about ninjas.",
//	      "rating": 8.5,
//	      "episodes": 220,
//	      "status": "Completed",
//	      "start_date": "2002-10-03T00:00:00Z",
//	      "end_date": "2007-02-08T00:00:00Z",
//	      "created_at": "2025-02-20T19:27:00Z",
//	      "updated_at": "2025-02-20T19:27:00Z",
//	      "genres": [
//	        {
//	          "id": 1,
//	          "name": "Action"
//	        }
//	      ],
//	      "tags": [
//	        {
//	          "id": 1,
//	          "name": "Ninja"
//	        }
//	      ]
//	    }
//	  ],
//	  "total": 1,
//	  "page": 1,
//	  "limit": 10
//	}
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

// RegisterAnimeRoutes registers all anime-related routes
func (h *AnimeHandler) RegisterAnimeRoutes(router *mux.Router) {
	// Public routes (no authentication required)
	router.HandleFunc("/animes", h.GetAllAnimesHandler).Methods("GET")
	router.HandleFunc("/animes/{id}", h.GetAnimeHandler).Methods("GET")
	router.HandleFunc("/animes/search", h.GetAnimesByTitleHandler).Methods("GET")
	router.HandleFunc("/animes/genre/{genre}", h.GetAnimesByGenreHandler).Methods("GET")

	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/animes").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware) // Apply authentication middleware

	// Protected routes with payload validation
	protectedRouter.Handle("", middleware.ValidateAndSanitizePayload(models.AnimeCreateRequest{})(http.HandlerFunc(h.CreateAnimeHandler))).Methods("POST")
	protectedRouter.Handle("/{id}", middleware.ValidateAndSanitizePayload(models.AnimeUpdateRequest{})(http.HandlerFunc(h.UpdateAnimeHandler))).Methods("PUT")
	protectedRouter.HandleFunc("/{id}", h.DeleteAnimeHandler).Methods("DELETE")

	// Genre management routes
	protectedRouter.Handle("/{id}/genres", middleware.ValidateAndSanitizePayload(struct {
		GenreIDs []uint `json:"genre_ids" validate:"required,min=1"`
	}{})(http.HandlerFunc(h.AddGenresToAnimeHandler))).Methods("POST")
	protectedRouter.Handle("/{id}/genres", middleware.ValidateAndSanitizePayload(struct {
		GenreIDs []uint `json:"genre_ids" validate:"required,min=1"`
	}{})(http.HandlerFunc(h.RemoveGenresFromAnimeHandler))).Methods("DELETE")

	// Tag management routes
	protectedRouter.Handle("/{id}/tags", middleware.ValidateAndSanitizePayload(struct {
		TagIDs []uint `json:"tag_ids" validate:"required,min=1"`
	}{})(http.HandlerFunc(h.AddTagsToAnimeHandler))).Methods("POST")
	protectedRouter.Handle("/{id}/tags", middleware.ValidateAndSanitizePayload(struct {
		TagIDs []uint `json:"tag_ids" validate:"required,min=1"`
	}{})(http.HandlerFunc(h.RemoveTagsFromAnimeHandler))).Methods("DELETE")
}
