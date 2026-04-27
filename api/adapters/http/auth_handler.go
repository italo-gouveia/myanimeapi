// api/adapters/http/auth_handler.go
// Package httphandler provides HTTP handlers for authentication-related routes in the MyAnimeAPI application.
package httphandler

import (
	"encoding/json"
	"net/http"

	"myanimeapi/api/middleware"
	"myanimeapi/api/models"
	"myanimeapi/api/services"
	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"

	"github.com/gorilla/mux"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authService services.AuthServiceInterface
	logger      *logger.Logger
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAuthHandler(authService services.AuthServiceInterface, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// writeJSONResponse writes a JSON response to the http.ResponseWriter
func (h *AuthHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}, log *logger.Logger, logFields map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logFields["error"] = err.Error()
		log.WithFields(logFields).Error("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// RegisterUserHandler handles user registration requests.
//
// @Summary Register a new user
// @Description Register a new user with the provided credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User registration information"
// @Success 201 {object} models.Response
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	log := h.logger
	logFields := map[string]interface{}{
		"method":     r.Method,
		"path":       r.URL.Path,
		"request_id": r.Context().Value(middleware.RequestIDContextKey),
	}

	// Decode request body
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to decode request body")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid request body", "The request body could not be decoded.", logFields)
		return
	}

	log.WithFields(logFields).WithField("username", user.Username).WithField("email", user.Email).Info("Processing user registration request")

	// Register user
	if err := h.authService.RegisterUser(r.Context(), &user); err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to register user")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, logFields)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to register user", err.Error(), logFields)
		return
	}

	userResponse := models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	response := models.Response{
		Status:  "success",
		Message: "User registered successfully",
		Data:    userResponse,
	}

	h.writeJSONResponse(w, http.StatusCreated, response, log, logFields)
	log.WithFields(logFields).WithField("user_id", user.ID).Info("User registered successfully")
}

// AuthenticateUserHandler handles user authentication requests.
//
// @Summary Authenticate a user
// @Description Authenticate a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.UserCredentials true "User credentials"
// @Success 200 {object} models.Response
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /auth/authenticate [post]
func (h *AuthHandler) AuthenticateUserHandler(w http.ResponseWriter, r *http.Request) {
	log := h.logger
	logFields := map[string]interface{}{
		"method":     r.Method,
		"path":       r.URL.Path,
		"request_id": r.Context().Value(middleware.RequestIDContextKey),
	}

	// Decode request body
	var credentials models.UserCredentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to decode request body")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid request body", "The request body could not be decoded.", logFields)
		return
	}

	log.WithFields(logFields).WithField("username", credentials.Username).Info("Processing authentication request")

	// Authenticate user
	token, err := h.authService.AuthenticateUser(r.Context(), &credentials)
	if err != nil {
		log.WithFields(logFields).WithField("error", err).Warning("Authentication failed")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, logFields)
			return
		}
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "Authentication failed", "Invalid credentials.", logFields)
		return
	}

	authResponse := models.AuthResponse{
		Token: token,
	}

	response := models.Response{
		Status:  "success",
		Message: "Authentication successful",
		Data:    authResponse,
	}

	h.writeJSONResponse(w, http.StatusOK, response, log, logFields)
	log.WithFields(logFields).Info("Authentication successful")
}

// RegisterAuthRoutes registers the authentication routes with the router.
func (h *AuthHandler) RegisterAuthRoutes(router *mux.Router) {
	router.HandleFunc("/auth/register", h.RegisterUserHandler).Methods("POST")
	router.HandleFunc("/auth/authenticate", h.AuthenticateUserHandler).Methods("POST")
}
