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
	"log"
	"net/http"
	"strconv"

	"myanimeapi/api/middleware"
	"myanimeapi/api/models"
	"myanimeapi/api/services"
	"myanimeapi/api/utils"
	"myanimeapi/internal/errors"

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
	var request models.AnimeCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create the anime
	anime := &models.Anime{
		Title:       request.Title,
		Description: request.Description,
		Rating:      request.Rating,
	}

	// Create the anime first
	if err := h.service.CreateAnime(r.Context(), anime); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create anime")
		return
	}

	// Add genres if provided
	if len(request.GenreIDs) > 0 {
		if err := h.service.AddGenresToAnime(r.Context(), anime.ID, request.GenreIDs); err != nil {
			// Log the error but don't fail the request
			log.Printf("Failed to add genres to anime: %v", err)
		}
	}

	// Add tags if provided
	if len(request.TagIDs) > 0 {
		if err := h.service.AddTagsToAnime(r.Context(), anime.ID, request.TagIDs); err != nil {
			// Log the error but don't fail the request
			log.Printf("Failed to add tags to anime: %v", err)
		}
	}

	// Get the complete anime with genres and tags
	createdAnime, err := h.service.GetAnimeByID(r.Context(), anime.ID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve created anime")
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, createdAnime)
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
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid anime ID")
		return
	}

	anime, err := h.service.GetAnimeByID(r.Context(), uint(id))
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve anime")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, anime)
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
	page, limit := utils.GetPaginationParams(r)
	animes, total, err := h.service.GetAllAnimes(r.Context(), page, limit)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve animes")
		return
	}

	response := map[string]interface{}{
		"animes": animes,
		"total":  total,
		"page":   page,
		"limit":  limit,
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
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
// @Param anime body models.Anime true "Updated anime details"
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
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid anime ID")
		return
	}

	var anime models.Anime
	if err := json.NewDecoder(r.Body).Decode(&anime); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	anime.ID = uint(id)
	if err := h.service.UpdateAnime(r.Context(), &anime); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update anime")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, anime)
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
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid anime ID")
		return
	}

	if err := h.service.DeleteAnime(r.Context(), uint(id)); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete anime")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Anime deleted successfully"})
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
	animeID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid anime ID")
		return
	}

	var request struct {
		GenreIDs []uint `json:"genre_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.AddGenresToAnime(r.Context(), uint(animeID), request.GenreIDs); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to add genres to anime")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Genres added successfully"})
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
	animeID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid anime ID")
		return
	}

	var request struct {
		GenreIDs []uint `json:"genre_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.RemoveGenresFromAnime(r.Context(), uint(animeID), request.GenreIDs); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to remove genres from anime")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Genres removed successfully"})
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
	animeID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid anime ID")
		return
	}

	var request struct {
		TagIDs []uint `json:"tag_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.AddTagsToAnime(r.Context(), uint(animeID), request.TagIDs); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to add tags to anime")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Tags added successfully"})
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
	animeID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid anime ID")
		return
	}

	var request struct {
		TagIDs []uint `json:"tag_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.RemoveTagsFromAnime(r.Context(), uint(animeID), request.TagIDs); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to remove tags from anime")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Tags removed successfully"})
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
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Title parameter is required")
		return
	}

	page, limit := utils.GetPaginationParams(r)
	animes, total, err := h.service.GetAnimesByTitle(r.Context(), title, page, limit)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to search animes by title")
		return
	}

	response := map[string]interface{}{
		"animes": animes,
		"total":  total,
		"page":   page,
		"limit":  limit,
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
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
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Genre parameter is required")
		return
	}

	page, limit := utils.GetPaginationParams(r)
	animes, total, err := h.service.GetAnimesByGenre(r.Context(), genre, page, limit)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to search animes by genre")
		return
	}

	response := map[string]interface{}{
		"animes": animes,
		"total":  total,
		"page":   page,
		"limit":  limit,
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// RegisterAnimeRoutes registers all anime-related routes
func (h *AnimeHandler) RegisterAnimeRoutes(router *mux.Router) {
	// Public routes (no authentication required)
	router.HandleFunc("/animes", h.GetAllAnimesHandler).Methods("GET")
	router.HandleFunc("/animes/search", h.GetAnimesByTitleHandler).Methods("GET")
	router.HandleFunc("/animes/genre/{genre}", h.GetAnimesByGenreHandler).Methods("GET")
	router.HandleFunc("/animes/{id}", h.GetAnimeHandler).Methods("GET")

	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/animes").Subrouter()
	protectedRouter.Use(middleware.Authenticate) // Apply authentication middleware

	// Protected routes with payload validation
	protectedRouter.Handle("", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.CreateAnimeHandler), models.AnimeCreateRequest{})).Methods("POST")
	protectedRouter.Handle("/{id}", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.UpdateAnimeHandler), models.Anime{})).Methods("PUT")
	protectedRouter.HandleFunc("/{id}", h.DeleteAnimeHandler).Methods("DELETE")

	// Protected routes for genres with payload validation
	protectedRouter.Handle("/{id}/genres", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.AddGenresToAnimeHandler), struct {
		GenreIDs []uint `json:"genre_ids" validate:"required,min=1"`
	}{})).Methods("POST")
	protectedRouter.Handle("/{id}/genres", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.RemoveGenresFromAnimeHandler), struct {
		GenreIDs []uint `json:"genre_ids" validate:"required,min=1"`
	}{})).Methods("DELETE")

	// Protected routes for tags with payload validation
	protectedRouter.Handle("/{id}/tags", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.AddTagsToAnimeHandler), struct {
		TagIDs []uint `json:"tag_ids" validate:"required,min=1"`
	}{})).Methods("POST")
	protectedRouter.Handle("/{id}/tags", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.RemoveTagsFromAnimeHandler), struct {
		TagIDs []uint `json:"tag_ids" validate:"required,min=1"`
	}{})).Methods("DELETE")
}
