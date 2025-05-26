package handlers

import (
	"encoding/json"
	"net/http"

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
	protectedRouter.Use(middleware.AuthMiddleware) // Apply authentication middleware

	// Protected routes with payload validation
	protectedRouter.Handle("", middleware.ValidateAndSanitizePayload(models.TagCreateRequest{})(http.HandlerFunc(h.CreateTagHandler))).Methods("POST")
	protectedRouter.Handle("/{id}", middleware.ValidateAndSanitizePayload(models.TagUpdateRequest{})(http.HandlerFunc(h.UpdateTagHandler))).Methods("PUT")
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
// @Param tag body models.TagCreateRequest true "Tag details"
// @Success 201 {object} models.TagResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 409 {object} errors.ErrorResponse "Tag name already exists"
// @Failure 500 {object} errors.ErrorResponse "Failed to create tag"
// @Router /tags [post]
// @Security BearerAuth
func (h *TagHandler) CreateTagHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the validated and sanitized payload from the context
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.TagCreateRequest)
	if !ok {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Invalid payload", "The request payload could not be retrieved.", nil)
		return
	}

	// Create a new Tag from the request
	tag := &models.Tag{
		Name: payload.Name,
	}

	if err := h.tagService.CreateTag(r.Context(), tag); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to create tag", "An internal server error occurred while creating the tag.", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(tag.ToResponse()); err != nil {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.", nil)
		return
	}
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
func (h *TagHandler) GetTagHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := utils.ValidateID(vars["id"])
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", "The provided ID is not a valid unsigned integer.", map[string]interface{}{
			"id": vars["id"],
		})
		return
	}

	tag, err := h.tagService.GetTagByID(r.Context(), id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to retrieve tag", "An internal server error occurred while retrieving the tag.", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tag.ToResponse()); err != nil {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.", nil)
		return
	}
}

// GetAllTagsHandler handles retrieving all tags.
// It queries the service and returns the tags as a JSON response.
// If an error occurs, it returns an appropriate error response.
//
// @Summary Get all tags
// @Description Get a list of all tags
// @Tags tags
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Success 200 {array} models.TagResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid pagination parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve tags"
// @Router /tags [get]
func (h *TagHandler) GetAllTagsHandler(w http.ResponseWriter, r *http.Request) {
	// Get pagination parameters from query
	page, limit, err := utils.ValidatePagination(r.URL.Query().Get("page"), r.URL.Query().Get("limit"), 1, 100)
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid pagination parameters", err.Error(), nil)
		return
	}

	tags, total, err := h.tagService.GetAllTags(r.Context(), page, limit)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to retrieve tags", "An internal server error occurred while retrieving the tags.", nil)
		return
	}

	responses := make([]models.TagResponse, len(tags))
	for i, tag := range tags {
		responses[i] = tag.ToResponse()
	}

	// Create paginated response
	response := struct {
		Data  []models.TagResponse `json:"data"`
		Total int64                `json:"total"`
		Page  int                  `json:"page"`
		Limit int                  `json:"limit"`
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
// @Param tag body models.TagUpdateRequest true "Updated tag details"
// @Success 200 {object} models.TagResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 404 {object} errors.ErrorResponse "Tag not found"
// @Failure 409 {object} errors.ErrorResponse "Tag name already exists"
// @Failure 500 {object} errors.ErrorResponse "Failed to update tag"
// @Router /tags/{id} [put]
// @Security BearerAuth
func (h *TagHandler) UpdateTagHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := utils.ValidateID(vars["id"])
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", "The provided ID is not a valid unsigned integer.", map[string]interface{}{
			"id": vars["id"],
		})
		return
	}

	// Retrieve the validated and sanitized payload from the context
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.TagUpdateRequest)
	if !ok {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Invalid payload", "The request payload could not be retrieved.", nil)
		return
	}

	// Create a new Tag from the request
	tag := &models.Tag{
		ID:   id,
		Name: payload.Name,
	}

	if err := h.tagService.UpdateTag(r.Context(), tag); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to update tag", "An internal server error occurred while updating the tag.", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tag.ToResponse()); err != nil {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.", nil)
		return
	}
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
// @Success 204 "No Content"
// @Failure 400 {object} errors.ErrorResponse "Invalid tag ID"
// @Failure 404 {object} errors.ErrorResponse "Tag not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete tag"
// @Router /tags/{id} [delete]
// @Security BearerAuth
func (h *TagHandler) DeleteTagHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := utils.ValidateID(vars["id"])
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", "The provided ID is not a valid unsigned integer.", map[string]interface{}{
			"id": vars["id"],
		})
		return
	}

	if err := h.tagService.DeleteTag(r.Context(), id); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to delete tag", "An internal server error occurred while deleting the tag.", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
