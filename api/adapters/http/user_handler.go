package httphandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"myanimeapi/api/middleware"
	"myanimeapi/api/models"
	"myanimeapi/api/services"
	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"

	"github.com/gorilla/mux"
)

// UserLoginRequest represents the request payload for user login
type UserLoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50" example:"john_doe"`
	Password string `json:"password" validate:"required,min=5,max=100" example:"password123"`
}

// UserHandler defines the handlers for user-related routes.
type UserHandler struct {
	userService      services.UserServiceInterface
	genreService     services.GenreServiceInterface
	passwordResetSvc services.PasswordResetServiceInterface
}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler(userService services.UserServiceInterface, genreService services.GenreServiceInterface, passwordResetSvc services.PasswordResetServiceInterface) *UserHandler {
	return &UserHandler{
		userService:      userService,
		genreService:     genreService,
		passwordResetSvc: passwordResetSvc,
	}
}

// writeJSONResponse writes a JSON response to the http.ResponseWriter
func (h *UserHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}, log *logger.Logger, logFields map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to encode response")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An error occurred while encoding the response.", logFields)
		return
	}
}

// Register handles user registration requests.
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	logFields := map[string]interface{}{
		"method":     r.Method,
		"path":       r.URL.Path,
		"request_id": r.Context().Value(middleware.RequestIDContextKey),
	}

	var req models.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to decode request body")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid request body", "The request body is not a valid JSON object.", logFields)
		return
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		IsActive: true,
	}

	if err := h.userService.CreateUser(r.Context(), user); err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to create user")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, logFields)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to create user", "An internal server error occurred while creating the user.", logFields)
		return
	}

	response := models.Response{
		Status:  "success",
		Message: "User registered successfully",
		Data:    user,
	}

	h.writeJSONResponse(w, http.StatusCreated, response, log, logFields)
	log.WithFields(logFields).WithField("user_id", user.ID).Info("User registered successfully")
}

// Login handles user authentication requests.
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	logFields := map[string]interface{}{
		"method":     r.Method,
		"path":       r.URL.Path,
		"request_id": r.Context().Value(middleware.RequestIDContextKey),
	}

	var req UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to decode request body")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid request body", "The request body is not a valid JSON object.", logFields)
		return
	}

	user, err := h.userService.ValidateUser(r.Context(), req.Username, req.Password)
	if err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Authentication failed")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, logFields)
			return
		}
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "Authentication failed", "Invalid username or password.", logFields)
		return
	}

	token, err := middleware.GenerateToken(fmt.Sprintf("%d", user.ID), user.IsAdmin)
	if err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to generate token")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to generate token", "An internal server error occurred while generating the authentication token.", logFields)
		return
	}

	response := models.Response{
		Status:  "success",
		Message: "Login successful",
		Data: models.AuthResponse{
			Token: token,
		},
	}

	h.writeJSONResponse(w, http.StatusOK, response, log, logFields)
	log.WithFields(logFields).WithField("user_id", user.ID).Info("User logged in successfully")
}

// GetProfile retrieves the authenticated user's profile information.
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	logFields := map[string]interface{}{
		"method":     r.Method,
		"path":       r.URL.Path,
		"request_id": r.Context().Value(middleware.RequestIDContextKey),
	}

	userID, ok := r.Context().Value(middleware.UserContextKey).(uint)
	if !ok {
		log.WithFields(logFields).Error("Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "User not authenticated", "The user is not authenticated.", logFields)
		return
	}

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to get user profile")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, logFields)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to get user profile", "An internal server error occurred while retrieving the user profile.", logFields)
		return
	}

	response := models.Response{
		Status:  "success",
		Message: "Profile retrieved successfully",
		Data:    user,
	}

	h.writeJSONResponse(w, http.StatusOK, response, log, logFields)
	log.WithFields(logFields).WithField("user_id", user.ID).Info("Profile retrieved successfully")
}

// UpdateProfile updates the authenticated user's profile information.
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	logFields := map[string]interface{}{
		"method":     r.Method,
		"path":       r.URL.Path,
		"request_id": r.Context().Value(middleware.RequestIDContextKey),
	}

	var req models.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to decode request body")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid request body", "The request body is not a valid JSON object.", logFields)
		return
	}

	userID, ok := r.Context().Value(middleware.UserContextKey).(uint)
	if !ok {
		log.WithFields(logFields).Error("Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "User not authenticated", "The user is not authenticated.", logFields)
		return
	}

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to get user")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, logFields)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to get user", "An internal server error occurred while retrieving the user.", logFields)
		return
	}

	// Update user fields if provided
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.ProfilePic != "" {
		user.ProfilePic = req.ProfilePic
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.SocialLinks != nil {
		user.SocialLinks = req.SocialLinks
	}

	// Update genres if provided
	if len(req.GenreIDs) > 0 {
		genres, err := h.genreService.GetGenresByIDs(r.Context(), req.GenreIDs)
		if err != nil {
			log.WithFields(logFields).WithField("error", err).Error("Failed to get genres")
			if appErr, ok := err.(*errors.AppError); ok {
				errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, logFields)
				return
			}
			errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to get genres", "An internal server error occurred while retrieving the genres.", logFields)
			return
		}
		user.Genres = genres
	}

	if err := h.userService.UpdateUser(r.Context(), user); err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to update user")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, logFields)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to update user", "An internal server error occurred while updating the user.", logFields)
		return
	}

	response := models.Response{
		Status:  "success",
		Message: "Profile updated successfully",
		Data:    user,
	}

	h.writeJSONResponse(w, http.StatusOK, response, log, logFields)
	log.WithFields(logFields).WithField("user_id", user.ID).Info("Profile updated successfully")
}

// ChangePassword handles password change requests for authenticated users.
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	logFields := map[string]interface{}{
		"method":     r.Method,
		"path":       r.URL.Path,
		"request_id": r.Context().Value(middleware.RequestIDContextKey),
	}

	var req models.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to decode request body")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid request body", "The request body is not a valid JSON object.", logFields)
		return
	}

	userID, ok := r.Context().Value(middleware.UserContextKey).(uint)
	if !ok {
		log.WithFields(logFields).Error("Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "User not authenticated", "The user is not authenticated.", logFields)
		return
	}

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to get user")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, logFields)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to get user", "An internal server error occurred while retrieving the user.", logFields)
		return
	}

	if err := h.userService.ChangePassword(r.Context(), user.ID, req.CurrentPassword, req.NewPassword); err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to change password")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, logFields)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to change password", "An internal server error occurred while changing the password.", logFields)
		return
	}

	response := models.Response{
		Status:  "success",
		Message: "Password changed successfully",
	}

	h.writeJSONResponse(w, http.StatusOK, response, log, logFields)
	log.WithFields(logFields).WithField("user_id", user.ID).Info("Password changed successfully")
}

// DeactivateAccount handles account deactivation requests for authenticated users.
func (h *UserHandler) DeactivateAccount(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	logFields := map[string]interface{}{
		"method":     r.Method,
		"path":       r.URL.Path,
		"request_id": r.Context().Value(middleware.RequestIDContextKey),
	}

	var req models.DeactivateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to decode request body")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid request body", "The request body is not a valid JSON object.", logFields)
		return
	}

	userID, ok := r.Context().Value(middleware.UserContextKey).(uint)
	if !ok {
		log.WithFields(logFields).Error("Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "User not authenticated", "The user is not authenticated.", logFields)
		return
	}

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to get user")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, logFields)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to get user", "An internal server error occurred while retrieving the user.", logFields)
		return
	}

	if err := h.userService.DeactivateAccount(r.Context(), user.ID, req.Password); err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to deactivate account")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, logFields)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to deactivate account", "An internal server error occurred while deactivating the account.", logFields)
		return
	}

	response := models.Response{
		Status:  "success",
		Message: "Account deactivated successfully",
	}

	h.writeJSONResponse(w, http.StatusOK, response, log, logFields)
	log.WithFields(logFields).WithField("user_id", user.ID).Info("Account deactivated successfully")
}

// RequestPasswordReset handles requests to reset a password
func (h *UserHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	logFields := map[string]interface{}{
		"method":     r.Method,
		"path":       r.URL.Path,
		"request_id": r.Context().Value(middleware.RequestIDContextKey),
	}

	var req models.PasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to decode request body")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid request body", "The request body is not a valid JSON object.", logFields)
		return
	}

	if err := h.passwordResetSvc.RequestPasswordReset(r.Context(), req.Email); err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to process password reset request")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, logFields)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to process password reset request", "An internal server error occurred while processing the password reset request.", logFields)
		return
	}

	response := models.Response{
		Status:  "success",
		Message: "Password reset instructions sent to your email",
	}

	h.writeJSONResponse(w, http.StatusOK, response, log, logFields)
	log.WithFields(logFields).WithField("email", req.Email).Info("Password reset request processed successfully")
}

// ResetPassword handles requests to set a new password
func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	logFields := map[string]interface{}{
		"method":     r.Method,
		"path":       r.URL.Path,
		"request_id": r.Context().Value(middleware.RequestIDContextKey),
	}

	var req models.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to decode request body")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid request body", "The request body is not a valid JSON object.", logFields)
		return
	}

	if err := h.passwordResetSvc.ResetPassword(r.Context(), req.Token, req.NewPassword); err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to reset password")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, logFields)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to reset password", "An internal server error occurred while resetting the password.", logFields)
		return
	}

	response := models.Response{
		Status:  "success",
		Message: "Password reset successfully",
	}

	h.writeJSONResponse(w, http.StatusOK, response, log, logFields)
	log.WithFields(logFields).Info("Password reset successfully")
}

// PasswordResetRateLimiter is the interface for the rate limiter used on the password reset endpoint.
type PasswordResetRateLimiter interface {
	StrictRateLimitMiddleware(maxRequests int, window time.Duration) func(http.Handler) http.Handler
}

// passwordResetRateLimit reads FORGOT_PASSWORD_RATE_LIMIT (max requests) and
// FORGOT_PASSWORD_RATE_WINDOW_HOURS (window in hours) from env, defaulting to 5/1h.
func passwordResetRateLimit() (int, time.Duration) {
	maxReqs := 5
	if v := os.Getenv("FORGOT_PASSWORD_RATE_LIMIT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			maxReqs = n
		}
	}
	window := time.Hour
	if v := os.Getenv("FORGOT_PASSWORD_RATE_WINDOW_HOURS"); v != "" {
		if h, err := strconv.Atoi(v); err == nil && h > 0 {
			window = time.Duration(h) * time.Hour
		}
	}
	return maxReqs, window
}

// RegisterUserRoutes registers the user-related routes with the router.
func (h *UserHandler) RegisterUserRoutes(router *mux.Router, rateLimiter PasswordResetRateLimiter) {
	maxReqs, window := passwordResetRateLimit()
	resetLimiter := rateLimiter.StrictRateLimitMiddleware(maxReqs, window)

	// Public routes
	router.HandleFunc("/users/register", h.Register).Methods(http.MethodPost)
	router.HandleFunc("/users/login", h.Login).Methods(http.MethodPost)
	router.Handle("/users/reset-password/request", resetLimiter(http.HandlerFunc(h.RequestPasswordReset))).Methods(http.MethodPost)
	router.HandleFunc("/users/reset-password", h.ResetPassword).Methods(http.MethodPost)

	// Protected routes
	protected := router.PathPrefix("/users").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/profile", h.GetProfile).Methods(http.MethodGet)
	protected.HandleFunc("/profile", h.UpdateProfile).Methods(http.MethodPut)
	protected.HandleFunc("/change-password", h.ChangePassword).Methods(http.MethodPost)
	protected.HandleFunc("/deactivate", h.DeactivateAccount).Methods(http.MethodPost)
	protected.HandleFunc("/delete", h.DeleteAccount).Methods(http.MethodDelete)
}

// DeleteAccount handles account deletion for authenticated users.
func (h *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	logFields := map[string]interface{}{
		"method":     r.Method,
		"path":       r.URL.Path,
		"request_id": r.Context().Value(middleware.RequestIDContextKey),
	}

	userID, ok := r.Context().Value(middleware.UserContextKey).(uint)
	if !ok {
		log.WithFields(logFields).Error("Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "User not authenticated", "The user is not authenticated.", logFields)
		return
	}

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to get user")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, logFields)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to get user", "An internal server error occurred while retrieving the user.", logFields)
		return
	}

	if err := h.userService.DeleteUser(r.Context(), user.ID); err != nil {
		log.WithFields(logFields).WithField("error", err).Error("Failed to delete user")
		if appErr, ok := err.(*errors.AppError); ok {
			errors.WriteErrorResponse(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details, logFields)
			return
		}
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to delete user", "An internal server error occurred while deleting the user.", logFields)
		return
	}

	response := models.Response{
		Status:  "success",
		Message: "Account deleted successfully",
	}

	h.writeJSONResponse(w, http.StatusOK, response, log, logFields)
	log.WithFields(logFields).WithField("user_id", user.ID).Info("Account deleted successfully")
}
