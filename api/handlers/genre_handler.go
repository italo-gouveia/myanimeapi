package handlers

import (
	"encoding/json"
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

// RegisterGenreRoutes registers all genre-related routes
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
// @Failure 400 {object} errors.ErrorResponse "Invalid genre ID or request body"
// @Failure 404 {object} errors.ErrorResponse "Genre not found"
// @Failure 409 {object} errors.ErrorResponse "Genre name already exists"
// @Failure 500 {object} errors.ErrorResponse "Failed to update genre"
// @Router /genres/{id} [put]
// @Security BearerAuth
// @Example
//
//	{
//	  "name": "Action-Adventure"
//	}
//
// @ExampleResponse
//
//	{
//	  "id": 1,
//	  "name": "Action-Adventure",
//	  "created_at": "2025-02-20T19:27:00Z",
//	  "updated_at": "2025-02-20T19:27:00Z"
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

// DeleteGenreHandler handles deleting a genre.
// It validates the ID, uses the service to delete the genre,
// and returns a 204 No Content response.
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
