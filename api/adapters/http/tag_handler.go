package httphandler

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

// TagHandler handles HTTP requests for tag operations.
type TagHandler struct {
	tagService services.TagServiceInterface
}

// NewTagHandler creates a new TagHandler instance.
func NewTagHandler(tagService services.TagServiceInterface) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// RegisterTagRoutes registers all tag-related routes with a *mux.Router.
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
