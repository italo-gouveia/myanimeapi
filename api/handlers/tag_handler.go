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

// api/handlers/tag_handler.go
// Package handlers provides HTTP handlers for tag-related routes in the MyAnimeAPI application.
// It defines methods to handle tag creation, retrieval, updating, and deletion.
// The package uses the Gorilla Mux router for routing, GORM for database interactions, and middleware for request validation and authentication.
//
// Example usage:
//
//	tagService := services.NewTagService(repository)
//	tagHandler := handlers.NewTagHandler(tagService)
//	router := mux.NewRouter()
//	tagHandler.RegisterTagRoutes(router)
//
//	http.ListenAndServe(":8080", router)

// TagHandler handles HTTP requests for tag operations.
// It contains a tag service for handling business logic.
type TagHandler struct {
	tagService services.TagServiceInterface
}

// NewTagHandler creates a new TagHandler instance.
// It accepts a tag service interface and returns a pointer to a TagHandler.
//
// Example:
//
//	tagService := services.NewTagService(repository)
//	tagHandler := NewTagHandler(tagService)
func NewTagHandler(tagService services.TagServiceInterface) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// RegisterTagRoutes registers all tag-related routes with a *mux.Router.
// It sets up the routes for tag management, including public and protected endpoints.
//
// Routes registered:
// - GET /tags - Get all tags (public)
// - GET /tags/{id} - Get a specific tag (public)
// - POST /tags - Create a new tag (protected)
// - PUT /tags/{id} - Update a tag (protected)
// - DELETE /tags/{id} - Delete a tag (protected)
func (h *TagHandler) RegisterTagRoutes(router *mux.Router) {
	// Public routes (no authentication required)
	router.HandleFunc("/tags", h.GetAllTagsHandler).Methods("GET")
	router.HandleFunc("/tags/{id}", h.GetTagHandler).Methods("GET")

	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/tags").Subrouter()
	protectedRouter.Use(middleware.Authenticate) // Apply authentication middleware

	// Protected routes with payload validation
	protectedRouter.Handle("", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.CreateTagHandler), models.Tag{})).Methods("POST")
	protectedRouter.Handle("/{id}", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.UpdateTagHandler), models.Tag{})).Methods("PUT")
	protectedRouter.HandleFunc("/{id}", h.DeleteTagHandler).Methods("DELETE")
}

// CreateTagHandler handles the creation of a new tag.
// It validates the input payload and uses the service to create the tag.
// If successful, it returns the created tag as a JSON response.
//
// @Summary Create a new tag
// @Description Create a new tag with the provided details
// @Tags tags
// @Accept json
// @Produce json
// @Param tag body models.Tag true "Tag details"
// @Success 201 {object} models.TagResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 409 {object} errors.ErrorResponse "Tag name already exists"
// @Failure 500 {object} errors.ErrorResponse "Failed to create tag"
// @Router /tags [post]
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
func (h *TagHandler) CreateTagHandler(w http.ResponseWriter, r *http.Request) {
	var tag models.Tag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.tagService.CreateTag(r.Context(), &tag); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create tag")
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, tag.ToResponse())
}

// GetTagHandler handles retrieving a tag by ID.
// It validates the ID, queries the service, and returns the tag as a JSON response.
// If the ID is invalid or the tag is not found, it returns an appropriate error response.
//
// @Summary Get a tag by ID
// @Description Get a tag's details by its ID
// @Tags tags
// @Produce json
// @Param id path int true "Tag ID"
// @Success 200 {object} models.TagResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid tag ID"
// @Failure 404 {object} errors.ErrorResponse "Tag not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve tag"
// @Router /tags/{id} [get]
// @ExampleResponse
//
//	{
//	  "id": 1,
//	  "name": "Action",
//	  "created_at": "2025-02-20T19:27:00Z",
//	  "updated_at": "2025-02-20T19:27:00Z"
//	}
func (h *TagHandler) GetTagHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid tag ID")
		return
	}

	tag, err := h.tagService.GetTagByID(r.Context(), uint(id))
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve tag")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, tag.ToResponse())
}

// GetAllTagsHandler handles retrieving all tags.
// It queries the service and returns the tags as a JSON response.
// If an error occurs, it returns an appropriate error response.
//
// @Summary Get all tags
// @Description Get a list of all tags
// @Tags tags
// @Produce json
// @Success 200 {array} models.TagResponse
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve tags"
// @Router /tags [get]
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
func (h *TagHandler) GetAllTagsHandler(w http.ResponseWriter, r *http.Request) {
	tags, _, err := h.tagService.GetAllTags(r.Context(), 1, 10)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve tags")
		return
	}

	responses := make([]models.TagResponse, len(tags))
	for i, tag := range tags {
		responses[i] = tag.ToResponse()
	}

	utils.WriteJSONResponse(w, http.StatusOK, responses)
}

// UpdateTagHandler handles updating an existing tag.
// It validates the ID and input payload, uses the service to update the tag,
// and returns the updated tag as a JSON response.
//
// @Summary Update a tag
// @Description Update an existing tag's details
// @Tags tags
// @Accept json
// @Produce json
// @Param id path int true "Tag ID"
// @Param tag body models.Tag true "Updated tag details"
// @Success 200 {object} models.TagResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 404 {object} errors.ErrorResponse "Tag not found"
// @Failure 409 {object} errors.ErrorResponse "Tag name already exists"
// @Failure 500 {object} errors.ErrorResponse "Failed to update tag"
// @Router /tags/{id} [put]
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
func (h *TagHandler) UpdateTagHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid tag ID")
		return
	}

	var tag models.Tag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	tag.ID = uint(id)
	if err := h.tagService.UpdateTag(r.Context(), &tag); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update tag")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, tag.ToResponse())
}

// DeleteTagHandler handles the deletion of a tag.
// It validates the ID and uses the service to delete the tag.
// If successful, it returns a success message as a JSON response.
//
// @Summary Delete a tag
// @Description Delete an existing tag by its ID
// @Tags tags
// @Produce json
// @Param id path int true "Tag ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} errors.ErrorResponse "Invalid tag ID"
// @Failure 404 {object} errors.ErrorResponse "Tag not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete tag"
// @Router /tags/{id} [delete]
// @Security BearerAuth
// @ExampleResponse
//
//	{
//	  "status": "success",
//	  "message": "Tag deleted successfully"
//	}
func (h *TagHandler) DeleteTagHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid tag ID")
		return
	}

	if err := h.tagService.DeleteTag(r.Context(), uint(id)); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteErrorResponse(w, appErr.StatusCode, appErr.Message)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete tag")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
