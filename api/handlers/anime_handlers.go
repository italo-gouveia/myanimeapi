// api/handlers/anime_handlers.go
// Package handlers provides HTTP handlers for anime-related routes in the MyAnimeAPI application.
// It defines methods to handle requests for retrieving, creating, updating, and deleting anime entries.
// The package uses the Gorilla Mux router for routing, GORM for database interactions, and middleware for request validation and authentication.
//
// Example usage:
//
//	db := // initialize your database connection
//	animeHandler := handlers.NewAnimeHandler(db)
//	router := mux.NewRouter()
//	animeHandler.RegisterAnimeRoutes(router)
//
//	http.ListenAndServe(":8080", router)
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"myanimeapi/api/middleware"
	"myanimeapi/api/models"
	"myanimeapi/api/services"
	"myanimeapi/internal/errors"

	"myanimeapi/api/utils"

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
	return &AnimeHandler{service: service}
}

// GetAnimeHandler retrieves an anime by its ID.
// It validates the ID, queries the database, and returns the anime as a JSON response.
// If the ID is invalid or the anime is not found, it returns an appropriate error response.
//
// @Summary Get an anime by ID
// @Description Retrieve an anime by its ID
// @Tags anime
// @Produce json
// @Param id path int true "Anime ID"
// @Success 200 {object} models.Anime
// @Failure 400 {object} errors.ErrorResponse "Invalid ID format"
// @Failure 404 {object} errors.ErrorResponse "Anime not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve anime"
// @Router /v1/anime/{id} [get]
// @ExampleResponse
//
//	{
//	  "id": 1,
//	  "title": "Naruto",
//	  "description": "A story about ninjas.",
//	  "rating": 8.5,
//	  "created_at": "2023-10-01T12:00:00Z",
//	  "updated_at": "2023-10-01T12:00:00Z"
//	}
//
// @Security []
func (h *AnimeHandler) GetAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := utils.ParseUint(vars["id"])
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

// GetAllAnimesHandler retrieves paginated anime entries from the database.
// It validates the pagination parameters, queries the database, and returns the anime entries as a JSON response.
// If the pagination parameters are invalid or the query fails, it returns an error response.
//
// @Summary Get paginated anime entries
// @Description Retrieve a paginated list of anime entries
// @Tags anime
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10)"
// @Success 200 {array} models.Anime
// @Failure 400 {object} errors.ErrorResponse "Invalid pagination parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve anime"
// @Router /v1/anime [get]
// @ExampleResponse
// [
//
//	{
//	  "id": 1,
//	  "title": "Naruto",
//	  "description": "A story about ninjas.",
//	  "rating": 8.5,
//	  "created_at": "2023-10-01T12:00:00Z",
//	  "updated_at": "2023-10-01T12:00:00Z"
//	}
//
// ]
// @Security []
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

// GetPaginatedReviewsForAnimeHandler retrieves paginated reviews for an anime by its ID.
// It validates the ID and pagination parameters, queries the database, and returns the reviews as a JSON response.
// If the ID or pagination parameters are invalid, or the query fails, it returns an error response.
//
// @Summary Get paginated reviews for an anime
// @Description Retrieve paginated reviews for an anime by its ID
// @Tags anime
// @Produce json
// @Param id path int true "Anime ID"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10)"
// @Success 200 {array} models.Review
// @Failure 400 {object} errors.ErrorResponse "Invalid ID format or pagination parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve reviews"
// @Router /v1/anime/{id}/reviews [get]
// @ExampleResponse
// [
//
//	{
//	  "id": 1,
//	  "user_id": 1,
//	  "anime_id": 1,
//	  "content": "Great anime!",
//	  "rating": 9,
//	  "created_at": "2023-10-01T12:00:00Z",
//	  "updated_at": "2023-10-01T12:00:00Z"
//	}
//
// ]
// @Security []
func (h *AnimeHandler) GetPaginatedReviewsForAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Validate ID
	id, err := utils.ValidateID(idStr)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", "The provided ID is not a valid unsigned integer.")
		return
	}

	// Validate pagination
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	page, limit, err := utils.ValidatePagination(pageStr, limitStr, 1, 10)
	if err != nil {
		log.Printf("Invalid pagination parameters: %v", err)
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid pagination parameters", err.Error())
		return
	}

	// Get reviews for the anime
	reviews, total, err := h.service.GetReviewsForAnime(r.Context(), id, page, limit)
	if err != nil {
		log.Printf("Failed to retrieve reviews: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to retrieve reviews", "An internal server error occurred while retrieving reviews.")
		return
	}

	response := map[string]interface{}{
		"reviews": reviews,
		"total":   total,
		"page":    page,
		"limit":   limit,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.")
		return
	}
}

// CreateAnimeHandler creates a new anime entry in the database.
// It validates the input payload, creates the anime, and returns the created anime as a JSON response.
// If the input is invalid or the creation fails, it returns an error response.
//
// @Summary Create a new anime
// @Description Create a new anime entry with the provided data
// @Tags anime
// @Accept json
// @Produce json
// @Param anime body models.AnimeCreateRequest true "Anime data"
// @Success 201 {object} models.AnimeResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid input or missing required fields"
// @Failure 500 {object} errors.ErrorResponse "Failed to create anime"
// @Router /v1/anime [post]
// @Example
//
//	{
//	  "title": "Naruto",
//	  "description": "A story about ninjas.",
//	  "rating": 8.5
//	}
//
// @ExampleResponse
//
//	{
//	  "id": 1,
//	  "title": "Naruto",
//	  "description": "A story about ninjas.",
//	  "rating": 8.5,
//	  "created_at": "2023-10-01T12:00:00Z",
//	  "updated_at": "2023-10-01T12:00:00Z"
//	}
//
// @Security ApiKeyAuth
func (h *AnimeHandler) CreateAnimeHandler(w http.ResponseWriter, r *http.Request) {
	var payload models.AnimeCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Convert AnimeCreateRequest to Anime
	anime := &models.Anime{
		Title:       payload.Title,
		Description: payload.Description,
		Rating:      payload.Rating,
	}

	err := h.service.CreateAnime(r.Context(), anime)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create anime")
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, anime)
}

// UpdateAnimeHandler updates an existing anime entry in the database.
// It validates the ID and input payload, updates the anime, and returns the updated anime as a JSON response.
// If the ID or input is invalid, or the update fails, it returns an error response.
//
// @Summary Update an anime
// @Description Update an existing anime entry with the provided data
// @Tags anime
// @Accept json
// @Produce json
// @Param id path int true "Anime ID"
// @Param anime body models.Anime true "Updated anime data"
// @Success 200 {object} models.Anime
// @Failure 400 {object} errors.ErrorResponse "Invalid input or ID format"
// @Failure 404 {object} errors.ErrorResponse "Anime not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to update anime"
// @Router /v1/anime/{id} [put]
// @Example
//
//	{
//	  "title": "Naruto Shippuden",
//	  "description": "The continuation of Naruto's journey.",
//	  "rating": 9.0
//	}
//
// @ExampleResponse
//
//	{
//	  "id": 1,
//	  "title": "Naruto Shippuden",
//	  "description": "The continuation of Naruto's journey.",
//	  "rating": 9.0,
//	  "created_at": "2023-10-01T12:00:00Z",
//	  "updated_at": "2023-10-01T12:00:00Z"
//	}
//
// @Security ApiKeyAuth
func (h *AnimeHandler) UpdateAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := utils.ParseUint(vars["id"])
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid anime ID")
		return
	}

	var payload models.Anime
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	payload.ID = uint(id)
	err = h.service.UpdateAnime(r.Context(), &payload)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update anime")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, payload)
}

// DeleteAnimeHandler deletes an anime entry from the database.
// It validates the ID, deletes the anime, and returns a 204 No Content response.
// If the ID is invalid or the deletion fails, it returns an error response.
//
// @Summary Delete an anime
// @Description Delete an anime entry by its ID
// @Tags anime
// @Param id path int true "Anime ID"
// @Success 204 "No Content"
// @Failure 400 {object} errors.ErrorResponse "Invalid ID format"
// @Failure 404 {object} errors.ErrorResponse "Anime not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete anime"
// @Router /v1/anime/{id} [delete]
// @Security ApiKeyAuth
func (h *AnimeHandler) DeleteAnimeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := utils.ParseUint(vars["id"])
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid anime ID")
		return
	}

	err = h.service.DeleteAnime(r.Context(), uint(id))
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete anime")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Anime deleted successfully"})
}

// RegisterAnimeRoutes registers all anime-related routes with the provided router.
// It defines public routes (GET) and protected routes (POST, PUT, DELETE) that require authentication.
//
// Example:
//
//	router := mux.NewRouter()
//	animeHandler.RegisterAnimeRoutes(router)
func (h *AnimeHandler) RegisterAnimeRoutes(router *mux.Router) {
	// Public routes (no authentication required)
	router.HandleFunc("/anime", h.GetAllAnimesHandler).Methods("GET")
	router.HandleFunc("/anime/{id:[0-9]+}", h.GetAnimeHandler).Methods("GET")
	router.HandleFunc("/anime/{id:[0-9]+}/reviews", h.GetPaginatedReviewsForAnimeHandler).Methods("GET")

	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/anime").Subrouter()
	protectedRouter.Use(middleware.Authenticate) // Apply authentication middleware

	// Protected routes (require authentication)
	protectedRouter.Handle("", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.CreateAnimeHandler), models.AnimeCreateRequest{})).Methods("POST")
	protectedRouter.Handle("/{id:[0-9]+}", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.UpdateAnimeHandler), models.Anime{})).Methods("PUT")
	protectedRouter.HandleFunc("/{id:[0-9]+}", h.DeleteAnimeHandler).Methods("DELETE")
}
