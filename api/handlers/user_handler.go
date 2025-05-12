package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"myanimeapi/api/auth"
	"myanimeapi/api/middleware"
	"myanimeapi/api/models"
	"myanimeapi/api/services"
	"myanimeapi/internal/errors"

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

// UserHandler defines the handlers for user-related routes.
// It contains a user service for handling user-related business logic and a genre service for managing user genre preferences.
type UserHandler struct {
	userService  *services.UserService
	genreService *services.GenreService
}

// NewUserHandler creates a new instance of UserHandler.
// It accepts a user service and genre service and returns a pointer to a UserHandler.
//
// Example:
//
//	userService := services.NewUserService(userRepo)
//	genreService := services.NewGenreService(genreRepo)
//	userHandler := NewUserHandler(userService, genreService)
func NewUserHandler(userService *services.UserService, genreService *services.GenreService) *UserHandler {
	return &UserHandler{
		userService:  userService,
		genreService: genreService,
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
	var req models.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid request body", "The request body is not a valid JSON object.")
		return
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		IsActive: true,
	}

	if err := h.userService.CreateUser(r.Context(), user); err != nil {
		if err.(*errors.AppError).Code == errors.ErrConflict {
			log.Printf("Conflict: %v", err)
			errors.WriteErrorResponse(w, http.StatusConflict, errors.ErrConflict, "Username or email already exists", err.Error())
			return
		}
		log.Printf("Failed to create user: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to create user", "An internal server error occurred while creating the user.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.Response{
		Status:  "success",
		Message: "User registered successfully",
		Data:    user,
	})
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
	var credentials models.UserCredentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid request body", "The request body is not a valid JSON object.")
		return
	}

	user, err := h.userService.ValidateUser(r.Context(), credentials.Username, credentials.Password)
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "Invalid credentials", "The provided username or password is incorrect.")
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.IsAdmin)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to generate token", "An internal server error occurred while generating the authentication token.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{
		Status:  "success",
		Message: "Login successful",
		Data: map[string]interface{}{
			"token": token,
			"user":  user,
		},
	})
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
	userID, ok := r.Context().Value(middleware.UserContextKey).(uint)
	if !ok {
		log.Println("Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "User not authenticated", "The user is not authenticated.")
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to get user", "An internal server error occurred while retrieving the user.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{
		Status:  "success",
		Message: "Profile retrieved successfully",
		Data:    user,
	})
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
	var req models.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid request body", "The request body is not a valid JSON object.")
		return
	}

	userID, ok := r.Context().Value(middleware.UserContextKey).(uint)
	if !ok {
		log.Println("Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "User not authenticated", "The user is not authenticated.")
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to get user", "An internal server error occurred while retrieving the user.")
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
			log.Printf("Failed to get genres: %v", err)
			errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to get genres", "An internal server error occurred while retrieving the genres.")
			return
		}
		user.Genres = genres
	}

	if err := h.userService.UpdateUser(r.Context(), user); err != nil {
		if err.(*errors.AppError).Code == errors.ErrConflict {
			log.Printf("Conflict: %v", err)
			errors.WriteErrorResponse(w, http.StatusConflict, errors.ErrConflict, "Username or email already exists", err.Error())
			return
		}
		log.Printf("Failed to update user: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to update user", "An internal server error occurred while updating the user.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{
		Status:  "success",
		Message: "Profile updated successfully",
		Data:    user,
	})
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
// @Param request body models.PasswordChangeRequest true "Password change request"
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
	var req models.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid request body", "The request body is not a valid JSON object.")
		return
	}

	userID, ok := r.Context().Value(middleware.UserContextKey).(uint)
	if !ok {
		log.Println("Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "User not authenticated", "The user is not authenticated.")
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to get user", "An internal server error occurred while retrieving the user.")
		return
	}

	valid, err := auth.CheckPasswordHash(req.CurrentPassword, user.Password)
	if err != nil || !valid {
		log.Println("Current password is incorrect")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "Current password is incorrect", "The provided current password does not match the user's password.")
		return
	}

	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to hash password", "An internal server error occurred while hashing the password.")
		return
	}

	user.Password = hashedPassword
	if err := h.userService.UpdateUser(r.Context(), user); err != nil {
		log.Printf("Failed to update password: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to update password", "An internal server error occurred while updating the password.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{
		Status:  "success",
		Message: "Password changed successfully",
	})
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
// @Param request body models.AccountDeactivationRequest true "Account deactivation request"
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
	var req models.DeactivateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid request body", "The request body is not a valid JSON object.")
		return
	}

	userID, ok := r.Context().Value(middleware.UserContextKey).(uint)
	if !ok {
		log.Println("Failed to get user ID from context")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "User not authenticated", "The user is not authenticated.")
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to get user", "An internal server error occurred while retrieving the user.")
		return
	}

	valid, err := auth.CheckPasswordHash(req.Password, user.Password)
	if err != nil || !valid {
		log.Println("Password is incorrect")
		errors.WriteErrorResponse(w, http.StatusUnauthorized, errors.ErrUnauthorized, "Password is incorrect", "The provided password does not match the user's password.")
		return
	}

	user.IsActive = false
	if err := h.userService.UpdateUser(r.Context(), user); err != nil {
		log.Printf("Failed to deactivate account: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to deactivate account", "An internal server error occurred while deactivating the account.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{
		Status:  "success",
		Message: "Account deactivated successfully",
	})
}

// RegisterUserRoutes registers all user-related routes with a *mux.Router.
// It sets up the routes for user registration, authentication, profile management, and account operations.
//
// Routes registered:
// - POST /users/register - Register a new user
// - POST /users/login - Authenticate a user
// - GET /users/profile - Get user profile
// - PUT /users/profile - Update user profile
// - POST /users/change-password - Change user password
// - POST /users/deactivate - Deactivate user account
func (h *UserHandler) RegisterUserRoutes(router *mux.Router) {
	// Public routes (no authentication required)
	router.HandleFunc("/users/register", h.Register).Methods("POST")
	router.HandleFunc("/users/login", h.Login).Methods("POST")

	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/users").Subrouter()
	protectedRouter.Use(middleware.Authenticate) // Apply authentication middleware

	// Protected routes (require authentication)
	protectedRouter.HandleFunc("/profile", h.GetProfile).Methods("GET")
	protectedRouter.HandleFunc("/profile", h.UpdateProfile).Methods("PUT")
	protectedRouter.HandleFunc("/change-password", h.ChangePassword).Methods("POST")
	protectedRouter.HandleFunc("/deactivate", h.DeactivateAccount).Methods("POST")
}
