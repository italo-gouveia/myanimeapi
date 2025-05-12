package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"myanimeapi/api/middleware"
	"myanimeapi/api/models"
	"myanimeapi/api/services"
	"myanimeapi/api/utils"
	"myanimeapi/internal/errors"

	"github.com/gorilla/mux"
)

// GenreHandler handles HTTP requests for genre operations.
// It contains a genre service for handling business logic.
type GenreHandler struct {
	genreService services.GenreServiceInterface
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

	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/genres").Subrouter()
	protectedRouter.Use(middleware.Authenticate) // Apply authentication middleware

	// Protected routes with payload validation
	protectedRouter.Handle("", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.CreateGenreHandler), models.Genre{})).Methods("POST")
	protectedRouter.Handle("/{id}", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.UpdateGenreHandler), models.Genre{})).Methods("PUT")
	protectedRouter.HandleFunc("/{id}", h.DeleteGenreHandler).Methods("DELETE")
	router.HandleFunc("/genres/search", h.SearchGenresHandler).Methods("GET")
	protectedRouter.Handle("/bulk", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.BulkCreateGenresHandler), BulkCreateGenresRequest{})).Methods("POST")
	protectedRouter.Handle("/bulk", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.BulkDeleteGenresHandler), BulkDeleteGenresRequest{})).Methods("DELETE")
}

// CreateGenreHandler handles the creation of a new genre.
// It validates the input payload and uses the service to create the genre.
// If successful, it returns the created genre as a JSON response.
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
	var genre models.Genre
	if err := json.NewDecoder(r.Body).Decode(&genre); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.genreService.CreateGenre(r.Context(), &genre); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create genre")
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, genre.ToResponse())
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
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid genre ID")
		return
	}

	genre, err := h.genreService.GetGenreByID(r.Context(), uint(id))
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve genre")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, genre.ToResponse())
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
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve genres")
		return
	}

	responses := make([]models.GenreResponse, len(genres))
	for i, genre := range genres {
		responses[i] = genre.ToResponse()
	}

	utils.WriteJSONResponse(w, http.StatusOK, responses)
}

// UpdateGenreHandler handles updating an existing genre.
// It validates the ID and input payload, uses the service to update the genre,
// and returns the updated genre as a JSON response.
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
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid genre ID")
		return
	}

	var genre models.Genre
	if err := json.NewDecoder(r.Body).Decode(&genre); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	genre.ID = uint(id)
	if err := h.genreService.UpdateGenre(r.Context(), &genre); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update genre")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, genre.ToResponse())
}

// DeleteGenreHandler handles the deletion of a genre.
// It validates the ID and uses the service to delete the genre.
// If successful, it returns a success message as a JSON response.
//
// @Summary Delete a genre
// @Description Delete an existing genre by its ID
// @Tags genres
// @Produce json
// @Param id path int true "Genre ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} errors.ErrorResponse "Invalid genre ID"
// @Failure 404 {object} errors.ErrorResponse "Genre not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete genre"
// @Router /genres/{id} [delete]
// @Security BearerAuth
// @ExampleResponse
//
//	{
//	  "status": "success",
//	  "message": "Genre deleted successfully"
//	}
func (h *GenreHandler) DeleteGenreHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid genre ID")
		return
	}

	if err := h.genreService.DeleteGenre(r.Context(), uint(id)); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete genre")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SearchGenresHandler handles searching for genres by name.
// It queries the service and returns matching genres as a JSON response.
// If an error occurs, it returns an appropriate error response.
//
// @Summary Search genres
// @Description Search for genres by name
// @Tags genres
// @Produce json
// @Param query query string true "Search query"
// @Success 200 {array} models.GenreResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid search query"
// @Failure 500 {object} errors.ErrorResponse "Failed to search genres"
// @Router /genres/search [get]
// @ExampleResponse
//
//	[
//	  {
//	    "id": 1,
//	    "name": "Action",
//	    "created_at": "2025-02-20T19:27:00Z",
//	    "updated_at": "2025-02-20T19:27:00Z"
//	  }
//	]
func (h *GenreHandler) SearchGenresHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		errors.WriteErrorResponse(w, http.StatusBadRequest, "Query parameter is required", "Missing query parameter", "Please provide a search query")
		return
	}

	page, limit, err := utils.ValidatePagination(r.URL.Query().Get("page"), r.URL.Query().Get("limit"), 1, 10)
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), "Invalid pagination parameters", "Please provide valid page and limit values")
		return
	}

	genres, total, err := h.genreService.SearchGenres(r.Context(), query, page, limit)
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to search genres", "Search operation failed", "Please try again later")
		return
	}

	response := make([]models.GenreResponse, len(genres))
	for i, genre := range genres {
		response[i] = genre.ToResponse()
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"data": response,
		"pagination": map[string]interface{}{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": int(math.Ceil(float64(total) / float64(limit))),
		},
	})
}

// BulkCreateGenresRequest represents the request payload for creating multiple genres.
// It contains a list of genres to be created.
type BulkCreateGenresRequest struct {
	Genres []models.Genre `json:"genres" validate:"required,dive"` // List of genres to create
}

// BulkCreateGenresHandler handles the creation of multiple genres.
// It validates the input payload and uses the service to create the genres.
// If successful, it returns the created genres as a JSON response.
//
// @Summary Create multiple genres
// @Description Create multiple genres with the provided details
// @Tags genres
// @Accept json
// @Produce json
// @Param request body BulkCreateGenresRequest true "Bulk genre creation request"
// @Success 201 {array} models.GenreResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 409 {object} errors.ErrorResponse "One or more genre names already exist"
// @Failure 500 {object} errors.ErrorResponse "Failed to create genres"
// @Router /genres/bulk [post]
// @Security BearerAuth
// @Example
//
//	{
//	  "genres": [
//	    {
//	      "name": "Action"
//	    },
//	    {
//	      "name": "Comedy"
//	    }
//	  ]
//	}
//
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
func (h *GenreHandler) BulkCreateGenresHandler(w http.ResponseWriter, r *http.Request) {
	var req BulkCreateGenresRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", "Invalid JSON format", "Please provide a valid JSON payload")
		return
	}

	if err := h.genreService.BulkCreateGenres(r.Context(), req.Genres); err != nil {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create genres", "Bulk creation failed", "Please try again later")
		return
	}

	response := make([]models.GenreResponse, len(req.Genres))
	for i, genre := range req.Genres {
		response[i] = genre.ToResponse()
	}

	utils.WriteJSONResponse(w, http.StatusCreated, response)
}

// BulkDeleteGenresRequest represents the request payload for deleting multiple genres.
// It contains a list of genre IDs to be deleted.
type BulkDeleteGenresRequest struct {
	IDs []uint `json:"ids" validate:"required,dive"` // List of genre IDs to delete
}

// BulkDeleteGenresHandler handles the deletion of multiple genres.
// It validates the input payload and uses the service to delete the genres.
// If successful, it returns a success message as a JSON response.
//
// @Summary Delete multiple genres
// @Description Delete multiple genres by their IDs
// @Tags genres
// @Accept json
// @Produce json
// @Param request body BulkDeleteGenresRequest true "Bulk genre deletion request"
// @Success 200 {object} models.Response
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 404 {object} errors.ErrorResponse "One or more genres not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete genres"
// @Router /genres/bulk [delete]
// @Security BearerAuth
// @Example
//
//	{
//	  "ids": [1, 2, 3]
//	}
//
// @ExampleResponse
//
//	{
//	  "status": "success",
//	  "message": "Genres deleted successfully"
//	}
func (h *GenreHandler) BulkDeleteGenresHandler(w http.ResponseWriter, r *http.Request) {
	var req BulkDeleteGenresRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", "Invalid JSON format", "Please provide a valid JSON payload")
		return
	}

	if err := h.genreService.BulkDeleteGenres(r.Context(), req.IDs); err != nil {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete genres", "Bulk deletion failed", "Please try again later")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
