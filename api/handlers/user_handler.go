package handlers

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

// api/handlers/user_handler.go
// Package handlers provides HTTP handlers for user-related routes in the MyAnimeAPI application.
// It defines methods to handle user registration, authentication, profile management, and account operations.
// The package uses the Gorilla Mux router for routing, GORM for database interactions, and middleware for request validation and authentication.
//
// Example usage:
//
//	userService := services.NewUserService(userRepo)
//	genreService := services.NewGenreService(genreRepo)
//	userHandler := handlers.NewUserHandler(userService, genreService)
//	router := mux.NewRouter()
//	userHandler.RegisterUserRoutes(router)
//
//	http.ListenAndServe(":8080", router)

// UserLoginRequest represents the request payload for user login
type UserLoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50" example:"john_doe"`     // Username for authentication
	Password string `json:"password" validate:"required,min=5,max=100" example:"password123"` // Password for authentication
}

// UserHandler defines the handlers for user-related routes.
// It contains a user service for handling user-related business logic and a genre service for managing user genre preferences.
type UserHandler struct {
	userService      services.UserServiceInterface
	genreService     services.GenreServiceInterface
	passwordResetSvc services.PasswordResetServiceInterface
}

// NewUserHandler creates a new instance of UserHandler.
// It accepts a user service and genre service and returns a pointer to a UserHandler.
//
// Example:
//
//	userService := services.NewUserService(userRepo)
//	genreService := services.NewGenreService(genreRepo)
//	userHandler := NewUserHandler(userService, genreService)
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
// It validates the input payload and uses the service to create the user.
// If successful, it returns the created user as a JSON response.
//
// @Summary Register a new user
// @Description Register a new user with the provided credentials
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.UserCreateRequest true "User registration data"
// @Success 201 {object} models.Response
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 409 {object} errors.ErrorResponse "Username or email already exists"
// @Failure 500 {object} errors.ErrorResponse "Failed to create user"
// @Router /users/register [post]
// @Example
//
//	{
//	  "username": "johndoe",
//	  "email": "john@example.com",
//	  "password": "securepassword123"
//	}
//
// @ExampleResponse
//
//	{
//	  "status": "success",
//	  "message": "User registered successfully",
//	  "data": {
//	    "id": 1,
//	    "username": "johndoe",
//	    "email": "john@example.com",
//	    "created_at": "2024-02-20T19:27:00Z"
//	  }
//	}
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
// It validates the credentials and generates a JWT token upon successful authentication.
// If successful, it returns the token and user information as a JSON response.
//
// @Summary Login user
// @Description Authenticate a user with username and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.UserCredentials true "User credentials"
// @Success 200 {object} models.Response
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 401 {object} errors.ErrorResponse "Invalid credentials"
// @Failure 500 {object} errors.ErrorResponse "Failed to authenticate user"
// @Router /users/login [post]
// @Example
//
//	{
//	  "username": "johndoe",
//	  "password": "securepassword123"
//	}
//
// @ExampleResponse
//
//	{
//	  "status": "success",
//	  "message": "Login successful",
//	  "data": {
//	    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
//	    "user": {
//	      "id": 1,
//	      "username": "johndoe",
//	      "email": "john@example.com"
//	    }
//	  }
//	}
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
// It uses the user ID from the request context to fetch the user's details.
// If successful, it returns the user's profile as a JSON response.
//
// @Summary Get user profile
// @Description Get the authenticated user's profile information
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 401 {object} errors.ErrorResponse "User not authenticated"
// @Failure 500 {object} errors.ErrorResponse "Failed to get user profile"
// @Router /users/profile [get]
// @ExampleResponse
//
//	{
//	  "status": "success",
//	  "message": "Profile retrieved successfully",
//	  "data": {
//	    "id": 1,
//	    "username": "johndoe",
//	    "email": "john@example.com",
//	    "profile_pic": "https://example.com/profile.jpg",
//	    "bio": "Anime enthusiast",
//	    "social_links": {
//	      "twitter": "https://twitter.com/johndoe",
//	      "instagram": "https://instagram.com/johndoe"
//	    },
//	    "genres": [
//	      {
//	        "id": 1,
//	        "name": "Action"
//	      }
//	    ]
//	  }
//	}
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
// It validates the input payload and uses the service to update the user's details.
// If successful, it returns the updated user profile as a JSON response.
//
// @Summary Update user profile
// @Description Update the authenticated user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.UserUpdateRequest true "Profile update request"
// @Success 200 {object} models.Response
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 401 {object} errors.ErrorResponse "User not authenticated"
// @Failure 500 {object} errors.ErrorResponse "Failed to update profile"
// @Router /users/profile [put]
// @Example
//
//	{
//	  "username": "newusername",
//	  "email": "newemail@example.com",
//	  "profile_pic": "https://example.com/new-profile.jpg",
//	  "bio": "Updated bio",
//	  "social_links": {
//	    "twitter": "https://twitter.com/newusername"
//	  },
//	  "genre_ids": [1, 2, 3]
//	}
//
// @ExampleResponse
//
//	{
//	  "status": "success",
//	  "message": "Profile updated successfully",
//	  "data": {
//	    "id": 1,
//	    "username": "newusername",
//	    "email": "newemail@example.com",
//	    "profile_pic": "https://example.com/new-profile.jpg",
//	    "bio": "Updated bio",
//	    "social_links": {
//	      "twitter": "https://twitter.com/newusername"
//	    },
//	    "genres": [
//	      {
//	        "id": 1,
//	        "name": "Action"
//	      }
//	    ]
//	  }
//	}
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
// It validates the current password and updates to the new password if valid.
// If successful, it returns a success message as a JSON response.
//
// @Summary Change user password
// @Description Change the authenticated user's password
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ChangePasswordRequest true "Password change request"
// @Success 200 {object} models.Response
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 401 {object} errors.ErrorResponse "Invalid current password"
// @Failure 500 {object} errors.ErrorResponse "Failed to change password"
// @Router /users/change-password [post]
// @Example
//
//	{
//	  "current_password": "oldpassword123",
//	  "new_password": "newpassword456"
//	}
//
// @ExampleResponse
//
//	{
//	  "status": "success",
//	  "message": "Password changed successfully"
//	}
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
// It validates the user's password and deactivates the account if valid.
// If successful, it returns a success message as a JSON response.
//
// @Summary Deactivate user account
// @Description Deactivate the authenticated user's account
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.DeactivateAccountRequest true "Account deactivation request"
// @Success 200 {object} models.Response
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 401 {object} errors.ErrorResponse "Invalid password"
// @Failure 500 {object} errors.ErrorResponse "Failed to deactivate account"
// @Router /users/deactivate [post]
// @Example
//
//	{
//	  "password": "currentpassword123"
//	}
//
// @ExampleResponse
//
//	{
//	  "status": "success",
//	  "message": "Account deactivated successfully"
//	}
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
// @Summary Request password reset
// @Description Send a password reset email to the user's email address
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.PasswordResetRequest true "Password reset request"
// @Success 200 {object} models.Response
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 500 {object} errors.ErrorResponse "Failed to process password reset request"
// @Router /users/forgot-password [post]
// @Example
//
//	{
//	  "email": "user@example.com"
//	}
//
// @ExampleResponse
//
//	{
//	  "status": "success",
//	  "message": "If an account exists with this email, you will receive password reset instructions."
//	}
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
// @Summary Reset password
// @Description Reset user's password using a valid reset token
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.ResetPasswordRequest true "Password reset request"
// @Success 200 {object} models.Response
// @Failure 400 {object} errors.ErrorResponse "Invalid request body or token"
// @Failure 500 {object} errors.ErrorResponse "Failed to reset password"
// @Router /users/reset-password [post]
// @Example
//
//	{
//	  "token": "valid-reset-token",
//	  "new_password": "newSecurePassword123"
//	}
//
// @ExampleResponse
//
//	{
//	  "status": "success",
//	  "message": "Password has been reset successfully."
//	}
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
// rateLimiter is applied with a configurable strict limit on the password-reset endpoint.
// Env vars: FORGOT_PASSWORD_RATE_LIMIT (default 5), FORGOT_PASSWORD_RATE_WINDOW_HOURS (default 1).
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
