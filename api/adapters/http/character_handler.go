package httphandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"myanimeapi/api/middleware"
	"myanimeapi/api/models"
	"myanimeapi/api/services"
	"myanimeapi/api/utils"
	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"
)

// CharacterHandler handles character-related HTTP requests.
type CharacterHandler struct {
	characterService services.CharacterServiceInterface
	logger           *logger.Logger
}

// NewCharacterHandler creates a new CharacterHandler.
func NewCharacterHandler(characterService services.CharacterServiceInterface) *CharacterHandler {
	return &CharacterHandler{characterService: characterService, logger: logger.New()}
}

// GetCharactersHandler returns a paginated list of all characters.
// @Summary Get all characters
// @Description Get a paginated list of all characters
// @Tags characters
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.CharacterListResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /characters [get]
func (h *CharacterHandler) GetCharactersHandler(w http.ResponseWriter, r *http.Request) {
	page, limit := paginationParams(r)

	characters, total, err := h.characterService.GetAllCharacters(r.Context(), page, limit)
	if err != nil {
		h.writeServiceError(w, err, "Failed to retrieve characters")
		return
	}

	responses := make([]models.CharacterResponse, len(characters))
	for i, c := range characters {
		responses[i] = c.ToResponse()
	}
	writeJSON(w, http.StatusOK, models.CharacterListResponse{
		Characters: responses, Total: total, Page: page, Limit: limit,
	})
}

// GetCharacterHandler returns a single character by ID.
// @Summary Get character by ID
// @Description Retrieve a character and its anime appearances
// @Tags characters
// @Produce json
// @Param id path int true "Character ID"
// @Success 200 {object} models.CharacterResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /characters/{id} [get]
func (h *CharacterHandler) GetCharacterHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}

	character, err := h.characterService.GetCharacterByID(r.Context(), id)
	if err != nil {
		h.writeServiceError(w, err, "Failed to retrieve character")
		return
	}

	writeJSON(w, http.StatusOK, character.ToResponse())
}

// SearchCharactersHandler searches characters by name.
// @Summary Search characters by name
// @Description Find characters whose name contains the given query
// @Tags characters
// @Produce json
// @Param name query string true "Name substring to search"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.CharacterListResponse
// @Failure 400 {object} errors.ErrorResponse
// @Router /characters/search [get]
func (h *CharacterHandler) SearchCharactersHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput,
			"Missing query parameter", "The 'name' query parameter is required", nil)
		return
	}

	page, limit := paginationParams(r)
	characters, total, err := h.characterService.GetCharactersByName(r.Context(), name, page, limit)
	if err != nil {
		h.writeServiceError(w, err, "Failed to search characters")
		return
	}

	responses := make([]models.CharacterResponse, len(characters))
	for i, c := range characters {
		responses[i] = c.ToResponse()
	}
	writeJSON(w, http.StatusOK, models.CharacterListResponse{
		Characters: responses, Total: total, Page: page, Limit: limit,
	})
}

// GetAnimeCharactersHandler returns characters for a specific anime.
// @Summary Get characters for an anime
// @Description Returns the cast of a given anime
// @Tags characters
// @Produce json
// @Param animeId path int true "Anime ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.CharacterListResponse
// @Failure 400 {object} errors.ErrorResponse
// @Router /animes/{animeId}/characters [get]
func (h *CharacterHandler) GetAnimeCharactersHandler(w http.ResponseWriter, r *http.Request) {
	animeID, ok := pathID(w, r, "animeId")
	if !ok {
		return
	}

	page, limit := paginationParams(r)
	characters, total, err := h.characterService.GetCharactersByAnime(r.Context(), animeID, page, limit)
	if err != nil {
		h.writeServiceError(w, err, "Failed to retrieve anime characters")
		return
	}

	responses := make([]models.CharacterResponse, len(characters))
	for i, c := range characters {
		responses[i] = c.ToResponse()
	}
	writeJSON(w, http.StatusOK, models.CharacterListResponse{
		Characters: responses, Total: total, Page: page, Limit: limit,
	})
}

// CreateCharacterHandler creates a new character (admin only).
// @Summary Create character
// @Description Create a new character entry
// @Tags characters
// @Accept json
// @Produce json
// @Param body body models.CharacterCreateRequest true "Character data"
// @Success 201 {object} models.CharacterResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /characters [post]
// @Security BearerAuth
func (h *CharacterHandler) CreateCharacterHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CharacterCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput,
			"Invalid request body", "Failed to parse JSON", nil)
		return
	}

	character := &models.Character{
		Name:        req.Name,
		Description: req.Description,
		VoiceActor:  req.VoiceActor,
		ImageURL:    req.ImageURL,
	}

	if err := h.characterService.CreateCharacter(r.Context(), character); err != nil {
		h.writeServiceError(w, err, "Failed to create character")
		return
	}

	writeJSON(w, http.StatusCreated, character.ToResponse())
}

// UpdateCharacterHandler updates an existing character (admin only).
// @Summary Update character
// @Description Update character fields
// @Tags characters
// @Accept json
// @Produce json
// @Param id path int true "Character ID"
// @Param body body models.CharacterUpdateRequest true "Fields to update"
// @Success 200 {object} models.CharacterResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /characters/{id} [put]
// @Security BearerAuth
func (h *CharacterHandler) UpdateCharacterHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}

	character, err := h.characterService.GetCharacterByID(r.Context(), id)
	if err != nil {
		h.writeServiceError(w, err, "Character not found")
		return
	}

	var req models.CharacterUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput,
			"Invalid request body", "Failed to parse JSON", nil)
		return
	}

	if req.Name != "" {
		character.Name = req.Name
	}
	if req.Description != "" {
		character.Description = req.Description
	}
	if req.VoiceActor != "" {
		character.VoiceActor = req.VoiceActor
	}
	if req.ImageURL != nil {
		character.ImageURL = req.ImageURL
	}

	if err := h.characterService.UpdateCharacter(r.Context(), character); err != nil {
		h.writeServiceError(w, err, "Failed to update character")
		return
	}

	writeJSON(w, http.StatusOK, character.ToResponse())
}

// DeleteCharacterHandler deletes a character (admin only).
// @Summary Delete character
// @Tags characters
// @Param id path int true "Character ID"
// @Success 204 "No Content"
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /characters/{id} [delete]
// @Security BearerAuth
func (h *CharacterHandler) DeleteCharacterHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}

	if err := h.characterService.DeleteCharacter(r.Context(), id); err != nil {
		h.writeServiceError(w, err, "Failed to delete character")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddAnimeCharactersHandler associates characters with an anime (admin only).
// @Summary Add characters to anime
// @Tags characters
// @Accept json
// @Param animeId path int true "Anime ID"
// @Param body body models.AnimeCharacterRequest true "Character IDs"
// @Success 204 "No Content"
// @Failure 400 {object} errors.ErrorResponse
// @Router /animes/{animeId}/characters [post]
// @Security BearerAuth
func (h *CharacterHandler) AddAnimeCharactersHandler(w http.ResponseWriter, r *http.Request) {
	animeID, ok := pathID(w, r, "animeId")
	if !ok {
		return
	}

	var req models.AnimeCharacterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput,
			"Invalid request body", "Failed to parse JSON", nil)
		return
	}

	if err := h.characterService.AddCharactersToAnime(r.Context(), animeID, req.CharacterIDs); err != nil {
		h.writeServiceError(w, err, "Failed to add characters to anime")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveAnimeCharacterHandler removes a character from an anime (admin only).
// @Summary Remove character from anime
// @Tags characters
// @Param animeId path int true "Anime ID"
// @Param characterId path int true "Character ID"
// @Success 204 "No Content"
// @Failure 400 {object} errors.ErrorResponse
// @Router /animes/{animeId}/characters/{characterId} [delete]
// @Security BearerAuth
func (h *CharacterHandler) RemoveAnimeCharacterHandler(w http.ResponseWriter, r *http.Request) {
	animeID, ok := pathID(w, r, "animeId")
	if !ok {
		return
	}
	characterID, ok := pathID(w, r, "characterId")
	if !ok {
		return
	}

	if err := h.characterService.RemoveCharacterFromAnime(r.Context(), animeID, characterID); err != nil {
		h.writeServiceError(w, err, "Failed to remove character from anime")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RegisterCharacterRoutes wires up all character routes.
func (h *CharacterHandler) RegisterCharacterRoutes(router *mux.Router) {
	// Public read routes
	charRouter := router.PathPrefix("/characters").Subrouter()
	charRouter.HandleFunc("", h.GetCharactersHandler).Methods(http.MethodGet)
	charRouter.HandleFunc("/search", h.SearchCharactersHandler).Methods(http.MethodGet)
	charRouter.HandleFunc("/{id}", h.GetCharacterHandler).Methods(http.MethodGet)

	// Admin-only write routes for /characters
	adminChar := charRouter.NewRoute().Subrouter()
	adminChar.Use(middleware.AuthMiddleware)
	adminChar.HandleFunc("", h.CreateCharacterHandler).Methods(http.MethodPost)
	adminChar.HandleFunc("/{id}", h.UpdateCharacterHandler).Methods(http.MethodPut)
	adminChar.HandleFunc("/{id}", h.DeleteCharacterHandler).Methods(http.MethodDelete)

	// Anime-scoped character routes
	animeRouter := router.PathPrefix("/animes/{animeId}/characters").Subrouter()
	animeRouter.HandleFunc("", h.GetAnimeCharactersHandler).Methods(http.MethodGet)

	adminAnimeChar := animeRouter.NewRoute().Subrouter()
	adminAnimeChar.Use(middleware.AuthMiddleware)
	adminAnimeChar.HandleFunc("", h.AddAnimeCharactersHandler).Methods(http.MethodPost)
	adminAnimeChar.HandleFunc("/{characterId}", h.RemoveAnimeCharacterHandler).Methods(http.MethodDelete)
}

// ── helpers ──────────────────────────────────────────────────────────────────

func paginationParams(r *http.Request) (page, limit int) {
	page, limit = 1, 10
	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
		limit = l
	}
	return
}

func pathID(w http.ResponseWriter, r *http.Request, key string) (uint, bool) {
	idStr := mux.Vars(r)[key]
	id, err := utils.ValidateID(idStr)
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput,
			"Invalid ID", "The provided ID must be a positive integer", nil)
		return 0, false
	}
	return id, true
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *CharacterHandler) writeServiceError(w http.ResponseWriter, err error, fallbackMsg string) {
	h.logger.WithField("error", err.Error()).Error(fallbackMsg)
	if appErr, ok := err.(*errors.AppError); ok {
		errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, appErr.Context)
		return
	}
	errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, fallbackMsg, "", nil)
}
