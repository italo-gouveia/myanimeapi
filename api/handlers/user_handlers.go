// api/handlers/user_handlers.go
// Package handlers provides HTTP handlers for user-related routes in the MyAnimeAPI application.
// It defines methods to handle requests for retrieving, creating, updating, and deleting users.
// The package uses the Gorilla Mux router for routing, GORM for database interactions, and middleware for request validation and authentication.
//
// Example usage:
//
//	userService := services.NewUserService(userRepo)
//	userHandler := handlers.NewUserHandler(userService)
//	router := mux.NewRouter()
//	userHandler.RegisterUserRoutes(router)
//
//	http.ListenAndServe(":8080", router)
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"myanimeapi/api/middleware"
	"myanimeapi/api/models"
	"myanimeapi/api/services"
	"myanimeapi/api/utils"
	"myanimeapi/internal/errors"

	"github.com/gorilla/mux"
)

// UserHandler defines the handlers for user-related routes.
// It contains a user service for handling user-related business logic.
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new instance of UserHandler.
// It accepts a user service and returns a pointer to a UserHandler.
//
// Example:
//
//	userService := services.NewUserService(userRepo)
//	userHandler := NewUserHandler(userService)
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetAllUsersHandler retrieves a paginated list of all users.
// This endpoint is restricted to admin users only.
//
// @Summary Get all users
// @Description Retrieve a paginated list of all users. Restricted to admin users only.
// @Tags users
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10, max: 100)"
// @Success 200 {array} models.User
// @Failure 400 {object} errors.ErrorResponse "Invalid pagination parameters"
// @Failure 403 {object} errors.ErrorResponse "Access denied"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve users"
// @Security ApiKeyAuth
// @Router /v1/users [get]
// @ExampleResponse
// [
//
//	{
//	  "id": 1,
//	  "username": "john_doe",
//	  "email": "john@example.com",
//	  "created_at": "2023-10-01T12:00:00Z",
//	  "updated_at": "2023-10-01T12:00:00Z"
//	}
//
// ]
// @Security ApiKeyAuth
func (h *UserHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	isAdmin, ok := r.Context().Value(middleware.IsAdminContextKey).(bool)
	if !ok || !isAdmin {
		log.Println("Access denied: user is not an admin")
		errors.WriteErrorResponse(w, http.StatusForbidden, errors.ErrForbidden, "Access denied", "You do not have permission to access this resource.")
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

	// Get users using the service
	users, total, err := h.userService.GetAllUsers(r.Context(), page, limit)
	if err != nil {
		log.Printf("Failed to retrieve users: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to retrieve users", "An internal server error occurred while retrieving users.")
		return
	}

	// Create response with pagination metadata
	response := map[string]interface{}{
		"users": users,
		"pagination": map[string]interface{}{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	}

	// Set response header and encode the result
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.")
		return
	}
}

// GetUserHandler retrieves a user by their ID.
// This endpoint requires authentication.
//
// @Summary Get a user by ID
// @Description Retrieve a user by their ID. Requires authentication.
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} errors.ErrorResponse "Invalid ID format"
// @Failure 404 {object} errors.ErrorResponse "User not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve user"
// @Security ApiKeyAuth
// @Router /v1/users/{id} [get]
// @ExampleResponse
//
//	{
//	  "id": 1,
//	  "username": "john_doe",
//	  "email": "john@example.com",
//	  "created_at": "2023-10-01T12:00:00Z",
//	  "updated_at": "2023-10-01T12:00:00Z"
//	}
//
// @Security ApiKeyAuth
func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := utils.ValidateID(idStr)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", err.Error())
		return
	}

	// Get user using the service
	user, err := h.userService.GetUserByID(r.Context(), id)
	if err != nil {
		// Check if it's a "not found" error
		if err.(*errors.AppError).Code == errors.ErrResourceNotFound {
			log.Printf("User not found: %v", err)
			errors.WriteErrorResponse(w, http.StatusNotFound, errors.ErrResourceNotFound, "User not found", fmt.Sprintf("No user found with ID %d", id))
			return
		}

		log.Printf("Failed to retrieve user: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to retrieve user", "An internal server error occurred while retrieving the user.")
		return
	}

	// Set response header and encode the result
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Failed to encode response: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.")
		return
	}
}

// CreateUserHandler creates a new user.
// This endpoint is public and does not require authentication.
//
// @Summary Create a new user
// @Description Create a new user account. This endpoint is public and does not require authentication.
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User object"
// @Success 201 {object} models.User
// @Failure 400 {object} errors.ErrorResponse "Invalid input"
// @Failure 409 {object} errors.ErrorResponse "Username or email already exists"
// @Failure 500 {object} errors.ErrorResponse "Failed to create user"
// @Router /v1/users [post]
func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var payload models.User
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid request body", "The request body is not a valid JSON object.")
		return
	}

	// Validate required fields
	if payload.Username == "" || payload.Email == "" || payload.Password == "" {
		log.Println("Missing required fields")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Missing required fields", "Username, email, and password are required.")
		return
	}

	// Create user using the service
	err := h.userService.CreateUser(r.Context(), &payload)
	if err != nil {
		// Check if it's a conflict error (username or email already exists)
		if err.(*errors.AppError).Code == errors.ErrConflict {
			log.Printf("Conflict: %v", err)
			errors.WriteErrorResponse(w, http.StatusConflict, errors.ErrConflict, "Username or email already exists", err.Error())
			return
		}

		log.Printf("Failed to create user: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to create user", "An internal server error occurred while creating the user.")
		return
	}

	// Set response header and encode the result
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Failed to encode response: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.")
		return
	}
}

// UpdateUserHandler updates an existing user.
// This endpoint requires authentication and the user can only update their own profile.
//
// @Summary Update a user
// @Description Update an existing user. Requires authentication and the user can only update their own profile.
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.User true "User object"
// @Success 200 {object} models.User
// @Failure 400 {object} errors.ErrorResponse "Invalid input"
// @Failure 403 {object} errors.ErrorResponse "Access denied"
// @Failure 404 {object} errors.ErrorResponse "User not found"
// @Failure 409 {object} errors.ErrorResponse "Username or email already exists"
// @Failure 500 {object} errors.ErrorResponse "Failed to update user"
// @Security ApiKeyAuth
// @Router /v1/users/{id} [put]
func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := utils.ValidateID(idStr)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", err.Error())
		return
	}

	// Check if the authenticated user is updating their own profile
	userID, ok := r.Context().Value(middleware.UserContextKey).(uint)
	if !ok || userID != id {
		log.Printf("Access denied: user %d attempting to update user %d", userID, id)
		errors.WriteErrorResponse(w, http.StatusForbidden, errors.ErrForbidden, "Access denied", "You do not have permission to update this user.")
		return
	}

	// Parse request body
	var payload models.User
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid request body", "The request body is not a valid JSON object.")
		return
	}

	// Set the ID from the URL
	payload.ID = id

	// Update user using the service
	err = h.userService.UpdateUser(r.Context(), &payload)
	if err != nil {
		// Check if it's a "not found" error
		if err.(*errors.AppError).Code == errors.ErrResourceNotFound {
			log.Printf("User not found: %v", err)
			errors.WriteErrorResponse(w, http.StatusNotFound, errors.ErrResourceNotFound, "User not found", fmt.Sprintf("No user found with ID %d", id))
			return
		}

		// Check if it's a conflict error (username or email already exists)
		if err.(*errors.AppError).Code == errors.ErrConflict {
			log.Printf("Conflict: %v", err)
			errors.WriteErrorResponse(w, http.StatusConflict, errors.ErrConflict, "Username or email already exists", err.Error())
			return
		}

		log.Printf("Failed to update user: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to update user", "An internal server error occurred while updating the user.")
		return
	}

	// Set response header and encode the result
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Failed to encode response: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.")
		return
	}
}

// DeleteUserHandler deletes a user by their ID.
// This endpoint requires authentication and the user can only delete their own profile.
//
// @Summary Delete a user
// @Description Delete a user by their ID. Requires authentication and the user can only delete their own profile.
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} errors.ErrorResponse "Invalid ID format"
// @Failure 403 {object} errors.ErrorResponse "Access denied"
// @Failure 404 {object} errors.ErrorResponse "User not found"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete user"
// @Security ApiKeyAuth
// @Router /v1/users/{id} [delete]
func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := utils.ValidateID(idStr)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid ID format", err.Error())
		return
	}

	// Check if the authenticated user is deleting their own profile
	userID, ok := r.Context().Value(middleware.UserContextKey).(uint)
	if !ok || userID != id {
		log.Printf("Access denied: user %d attempting to delete user %d", userID, id)
		errors.WriteErrorResponse(w, http.StatusForbidden, errors.ErrForbidden, "Access denied", "You do not have permission to delete this user.")
		return
	}

	// Delete user using the service
	err = h.userService.DeleteUser(r.Context(), id)
	if err != nil {
		// Check if it's a "not found" error
		if err.(*errors.AppError).Code == errors.ErrResourceNotFound {
			log.Printf("User not found: %v", err)
			errors.WriteErrorResponse(w, http.StatusNotFound, errors.ErrResourceNotFound, "User not found", fmt.Sprintf("No user found with ID %d", id))
			return
		}

		log.Printf("Failed to delete user: %v", err)
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to delete user", "An internal server error occurred while deleting the user.")
		return
	}

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// RegisterUserRoutes registers all user-related routes with the provided router.
// It defines public routes (POST) and protected routes (GET, PUT, DELETE) that require authentication.
// Admin-only routes are also defined for retrieving all users.
//
// Example:
//
//	router := mux.NewRouter()
//	userHandler.RegisterUserRoutes(router)
func (h *UserHandler) RegisterUserRoutes(router *mux.Router) {
	// Public routes (no authentication required)
	router.HandleFunc("/users", h.CreateUserHandler).Methods("POST")

	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/users").Subrouter()
	protectedRouter.Use(middleware.Authenticate) // Apply authentication middleware

	// Protected routes (require authentication)
	protectedRouter.HandleFunc("/{id:[0-9]+}", h.GetUserHandler).Methods("GET")
	protectedRouter.Handle("/{id:[0-9]+}", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.UpdateUserHandler), models.User{})).Methods("PUT")
	protectedRouter.HandleFunc("/{id:[0-9]+}", h.DeleteUserHandler).Methods("DELETE")

	// Create a subrouter for admin-only routes
	adminRouter := protectedRouter.PathPrefix("").Subrouter()
	adminRouter.Use(middleware.CheckAdmin) // Apply admin check middleware

	// Admin-only routes (require authentication and admin privileges)
	adminRouter.HandleFunc("", h.GetAllUsersHandler).Methods("GET")
}
