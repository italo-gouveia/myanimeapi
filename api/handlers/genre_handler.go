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

// GenreHandler handles HTTP requests for genre operations.
// It contains a genre service for handling business logic.
type GenreHandler struct {
	genreService services.GenreServiceInterface
	logger       *logger.Logger
}

// NewGenreHandler creates a new GenreHandler instance.
// It accepts a genre service interface and returns a pointer to a GenreHandler.
//
// Example:
//
//	genreService := services.NewGenreService(repository)
//	genreHandler := NewGenreHandler(genreService)
func NewGenreHandler(genreService services.GenreServiceInterface) *GenreHandler {
	return &GenreHandler{
		genreService: genreService,
		logger:       logger.New(),
	}
}

// RegisterGenreRoutes registers all genre-related routes with a *mux.Router.
// It sets up the routes for genre management, including public and protected endpoints.
//
// Routes registered:
// - GET /genres - Get all genres (public)
// - GET /genres/{id} - Get a specific genre (public)
// - GET /genres/search - Search genres by name (public)
// - POST /genres - Create a new genre (protected)
// - PUT /genres/{id} - Update a genre (protected)
// - DELETE /genres/{id} - Delete a genre (protected)
// - POST /genres/bulk - Create multiple genres (protected)
// - DELETE /genres/bulk - Delete multiple genres (protected)
func (h *GenreHandler) RegisterGenreRoutes(router *mux.Router) {
	// Public routes (no authentication required)
	router.HandleFunc("/genres", h.GetAllGenresHandler).Methods("GET")
	router.HandleFunc("/genres/{id}", h.GetGenreHandler).Methods("GET")
	router.HandleFunc("/genres/search", h.SearchGenresHandler).Methods("GET")

	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/genres").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware) // Apply authentication middleware

	// Protected routes with payload validation
	protectedRouter.Handle("", middleware.ValidateAndSanitizePayload(models.Genre{})(http.HandlerFunc(h.CreateGenreHandler))).Methods("POST")
	protectedRouter.Handle("/{id}", middleware.ValidateAndSanitizePayload(models.Genre{})(http.HandlerFunc(h.UpdateGenreHandler))).Methods("PUT")
	protectedRouter.HandleFunc("/{id}", h.DeleteGenreHandler).Methods("DELETE")
	protectedRouter.Handle("/bulk", middleware.ValidateAndSanitizePayload(BulkCreateGenresRequest{})(http.HandlerFunc(h.BulkCreateGenresHandler))).Methods("POST")
	protectedRouter.Handle("/bulk", middleware.ValidateAndSanitizePayload(BulkDeleteGenresRequest{})(http.HandlerFunc(h.BulkDeleteGenresHandler))).Methods("DELETE")
}

// CreateGenreHandler handles the creation of a new genre.
// It retrieves the validated and sanitized payload from the context,
// uses the service to create the genre, and returns the created genre as a JSON response.
//
// @Summary Create a new genre
// @Description Create a new genre with the provided details
// @Tags genres
// @Accept json
// @Produce json
// @Param genre body models.Genre true "Genre details"
// @Success 201 {object} models.GenreResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 409 {object} errors.ErrorResponse "Genre name already exists"
// @Failure 500 {object} errors.ErrorResponse "Failed to create genre"
// @Router /genres [post]
// @Security BearerAuth
// @Example
//
//	{
//	  "name": "Action"
//	}
//
// @ExampleResponse
//
//	{
//	  "id": 1,
//	  "name": "Action",
//	  "created_at": "2025-02-20T19:27:00Z",
//	  "updated_at": "2025-02-20T19:27:00Z"
//	}
func (h *GenreHandler) CreateGenreHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the validated and sanitized payload from the context
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.Genre)
	if !ok {
		h.logger.Error("Invalid payload")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Invalid payload", "The request payload could not be retrieved.", nil)
		return
	}

	if err := h.genreService.CreateGenre(r.Context(), payload); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to create genre")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to create genre", "An internal server error occurred while creating the genre.", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(payload.ToResponse()); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to encode response")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.", nil)
		return
	}
}

// GetGenreHandler handles retrieving a genre by ID.
// It validates the ID, queries the service, and returns the genre as a JSON response.
// If the ID is invalid or the genre is not found, it returns an appropriate error response.
//
// @Summary Get a genre by ID
// @Description Get a genre's details by its ID
// @Tags genres
// @Produce json
// @Param id path int true "Genre ID"
// @Success 200 {object} models.GenreResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid genre ID"
// @Failure 404 {object} errors.ErrorResponse "Genre not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve genre"
// @Router /genres/{id} [get]
// @ExampleResponse
//
//	{
//	  "id": 1,
//	  "name": "Action",
//	  "created_at": "2025-02-20T19:27:00Z",
//	  "updated_at": "2025-02-20T19:27:00Z"
//	}
func (h *GenreHandler) GetGenreHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := utils.ValidateID(vars["id"])
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"id":    vars["id"],
		}).Error("Invalid genre ID")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid genre ID", "The provided ID is not a valid unsigned integer.", map[string]interface{}{
			"id": vars["id"],
		})
		return
	}

	genre, err := h.genreService.GetGenreByID(r.Context(), id)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		}).Error("Failed to retrieve genre")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to retrieve genre", "An internal server error occurred while retrieving the genre.", map[string]interface{}{
			"id": id,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(genre.ToResponse()); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to encode response")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.", nil)
		return
	}
}

// GetAllGenresHandler handles retrieving all genres.
// It queries the service and returns the genres as a JSON response.
// If an error occurs, it returns an appropriate error response.
//
// @Summary Get all genres
// @Description Get a list of all genres
// @Tags genres
// @Produce json
// @Success 200 {array} models.GenreResponse
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve genres"
// @Router /genres [get]
// @ExampleResponse
//
//	[
//	  {
//	    "id": 1,
//	    "name": "Action",
//	    "created_at": "2025-02-20T19:27:00Z",
//	    "updated_at": "2025-02-20T19:27:00Z"
//	  },
//	  {
//	    "id": 2,
//	    "name": "Comedy",
//	    "created_at": "2025-02-20T19:27:00Z",
//	    "updated_at": "2025-02-20T19:27:00Z"
//	  }
//	]
func (h *GenreHandler) GetAllGenresHandler(w http.ResponseWriter, r *http.Request) {
	genres, _, err := h.genreService.GetAllGenres(r.Context(), 1, 10)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to retrieve genres")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to retrieve genres", "An internal server error occurred while retrieving genres.", nil)
		return
	}

	responses := make([]models.GenreResponse, len(genres))
	for i, genre := range genres {
		responses[i] = genre.ToResponse()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(responses); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to encode response")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.", nil)
		return
	}
}

// UpdateGenreHandler handles updating an existing genre.
// It retrieves the validated and sanitized payload from the context,
// uses the service to update the genre, and returns the updated genre as a JSON response.
//
// @Summary Update a genre
// @Description Update an existing genre's details
// @Tags genres
// @Accept json
// @Produce json
// @Param id path int true "Genre ID"
// @Param genre body models.Genre true "Updated genre details"
// @Success 200 {object} models.GenreResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 404 {object} errors.ErrorResponse "Genre not found"
// @Failure 409 {object} errors.ErrorResponse "Genre name already exists"
// @Failure 500 {object} errors.ErrorResponse "Failed to update genre"
// @Router /genres/{id} [put]
// @Security BearerAuth
// @Example
//
//	{
//	  "name": "Updated Action"
//	}
//
// @ExampleResponse
//
//	{
//	  "id": 1,
//	  "name": "Updated Action",
//	  "created_at": "2025-02-20T19:27:00Z",
//	  "updated_at": "2025-02-20T19:28:00Z"
//	}
func (h *GenreHandler) UpdateGenreHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := utils.ValidateID(vars["id"])
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"id":    vars["id"],
		}).Error("Invalid genre ID")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid genre ID", "The provided ID is not a valid unsigned integer.", map[string]interface{}{
			"id": vars["id"],
		})
		return
	}

	// Retrieve the validated and sanitized payload from the context
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.Genre)
	if !ok {
		h.logger.Error("Invalid payload")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Invalid payload", "The request payload could not be retrieved.", nil)
		return
	}

	payload.ID = id
	if err := h.genreService.UpdateGenre(r.Context(), payload); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		}).Error("Failed to update genre")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to update genre", "An internal server error occurred while updating the genre.", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(payload.ToResponse()); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to encode response")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.", nil)
		return
	}
}

// DeleteGenreHandler handles deleting a genre by ID.
// It validates the ID, uses the service to delete the genre,
// and returns a success response if the deletion is successful.
//
// @Summary Delete a genre
// @Description Delete a genre by its ID
// @Tags genres
// @Produce json
// @Param id path int true "Genre ID"
// @Success 204 "No Content"
// @Failure 400 {object} errors.ErrorResponse "Invalid genre ID"
// @Failure 404 {object} errors.ErrorResponse "Genre not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete genre"
// @Router /genres/{id} [delete]
// @Security BearerAuth
func (h *GenreHandler) DeleteGenreHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := utils.ValidateID(vars["id"])
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"id":    vars["id"],
		}).Error("Invalid genre ID")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid genre ID", "The provided ID is not a valid unsigned integer.", map[string]interface{}{
			"id": vars["id"],
		})
		return
	}

	if err := h.genreService.DeleteGenre(r.Context(), id); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		}).Error("Failed to delete genre")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to delete genre", "An internal server error occurred while deleting the genre.", map[string]interface{}{
			"id": id,
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SearchGenresHandler handles searching genres by name.
// It validates the query parameter, uses the service to search for genres,
// and returns the matching genres as a JSON response.
//
// @Summary Search genres
// @Description Search genres by name
// @Tags genres
// @Produce json
// @Param query query string true "Search query"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Success 200 {array} models.GenreResponse
// @Failure 400 {object} errors.ErrorResponse "Missing search query"
// @Failure 500 {object} errors.ErrorResponse "Failed to search genres"
// @Router /genres/search [get]
func (h *GenreHandler) SearchGenresHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		h.logger.Error("Missing search query")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Missing search query", "", map[string]interface{}{
			"query": query,
		})
		return
	}

	page, limit := utils.GetPaginationParams(r)
	genres, total, err := h.genreService.SearchGenres(r.Context(), query, page, limit)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"query": query,
		}).Error("Failed to search genres")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to search genres", "An internal server error occurred while searching genres.", map[string]interface{}{
			"query": query,
		})
		return
	}

	responses := make([]models.GenreResponse, len(genres))
	for i, genre := range genres {
		responses[i] = genre.ToResponse()
	}

	response := map[string]interface{}{
		"genres": responses,
		"total":  total,
		"page":   page,
		"limit":  limit,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to encode response")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.", nil)
		return
	}
}

// BulkCreateGenresRequest represents a request to create multiple genres.
type BulkCreateGenresRequest struct {
	Genres []models.Genre `json:"genres" validate:"required,dive"` // List of genres to create
}

// BulkCreateGenresHandler handles creating multiple genres in bulk.
// It validates the input payload, uses the service to create the genres,
// and returns the created genres as a JSON response.
//
// @Summary Create multiple genres
// @Description Create multiple genres in bulk
// @Tags genres
// @Accept json
// @Produce json
// @Param request body BulkCreateGenresRequest true "List of genres to create"
// @Success 201 {array} models.GenreResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 409 {object} errors.ErrorResponse "Genre name already exists"
// @Failure 500 {object} errors.ErrorResponse "Failed to create genres"
// @Router /genres/bulk [post]
// @Security BearerAuth
func (h *GenreHandler) BulkCreateGenresHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the validated and sanitized payload from the context
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*BulkCreateGenresRequest)
	if !ok {
		h.logger.Error("Invalid payload")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Invalid payload", "The request payload could not be retrieved.", nil)
		return
	}

	if err := h.genreService.BulkCreateGenres(r.Context(), payload.Genres); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"count": len(payload.Genres),
		}).Error("Failed to bulk create genres")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to bulk create genres", "An internal server error occurred while creating genres in bulk.", nil)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// BulkDeleteGenresRequest represents a request to delete multiple genres.
type BulkDeleteGenresRequest struct {
	IDs []uint `json:"ids" validate:"required,dive"` // List of genre IDs to delete
}

// BulkDeleteGenresHandler handles deleting multiple genres in bulk.
// It validates the input payload, uses the service to delete the genres,
// and returns a success response if the deletion is successful.
//
// @Summary Delete multiple genres
// @Description Delete multiple genres in bulk
// @Tags genres
// @Accept json
// @Produce json
// @Param request body BulkDeleteGenresRequest true "List of genre IDs to delete"
// @Success 204 "No Content"
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 404 {object} errors.ErrorResponse "Genre not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete genres"
// @Router /genres/bulk [delete]
// @Security BearerAuth
func (h *GenreHandler) BulkDeleteGenresHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the validated and sanitized payload from the context
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*BulkDeleteGenresRequest)
	if !ok {
		h.logger.Error("Invalid payload")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Invalid payload", "The request payload could not be retrieved.", nil)
		return
	}

	if err := h.genreService.BulkDeleteGenres(r.Context(), payload.IDs); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"ids":   payload.IDs,
		}).Error("Failed to bulk delete genres")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to bulk delete genres", "An internal server error occurred while deleting genres in bulk.", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
