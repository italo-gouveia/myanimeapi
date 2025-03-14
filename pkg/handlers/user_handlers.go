// pkg/handlers/user_handlers.go
// Package handlers provides HTTP handlers for user-related routes in the MyAnimeAPI application.
// It defines methods to handle requests for retrieving, creating, updating, and deleting users.
// The package uses the Gorilla Mux router for routing, GORM for database interactions, and middleware for request validation and authentication.
//
// Example usage:
//
//	db := // initialize your database connection
//	userHandler := handlers.NewUserHandler(db)
//	router := mux.NewRouter()
//	userHandler.RegisterUserRoutes(router)
//
//	http.ListenAndServe(":8080", router)
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"myanimeapi/internal/db"
	"myanimeapi/pkg/middleware"
	"myanimeapi/pkg/models"
	"myanimeapi/pkg/validation"

	"github.com/gorilla/mux"
)

// UserHandler defines the handlers for user-related routes.
// It contains a database interface for interacting with the database.
type UserHandler struct {
	DB db.DBInterface
}

// NewUserHandler creates a new instance of UserHandler.
// It accepts a database interface and returns a pointer to a UserHandler.
//
// Example:
//
//	db := // initialize your database connection
//	userHandler := NewUserHandler(db)
func NewUserHandler(db db.DBInterface) *UserHandler {
	return &UserHandler{DB: db}
}

// GetAllUsersHandler retrieves a paginated list of all users.
// This endpoint is restricted to admin users only.
//
// @Summary Get all users (admin only)
// @Description Retrieve a paginated list of all users. Requires admin privileges.
// @Tags users
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10)"
// @Success 200 {array} models.User
// @Failure 400 {string} string "Invalid pagination parameters"
// @Failure 403 {string} string "Access denied"
// @Failure 500 {string} string "Failed to retrieve users"
// @Security ApiKeyAuth
// @Router /users [get]
func (h *UserHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	isAdmin, ok := r.Context().Value(middleware.IsAdminContextKey).(bool)
	if !ok || !isAdmin {
		log.Println("Access denied: user is not an admin")
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Validate pagination
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	page, limit, err := validation.ValidatePagination(pageStr, limitStr, 1, 10)
	if err != nil {
		log.Printf("Invalid pagination parameters: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Calculate offset
	offset := (page - 1) * limit

	var users []models.User
	result := h.DB.Offset(offset).Limit(limit).Find(&users)
	if result.Error != nil {
		log.Printf("Failed to retrieve users: %v", result.Error)
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	// Set response header and encode the result
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
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
// @Failure 404 {string} string "User not found"
// @Security ApiKeyAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Validate ID
	id, err := validation.ValidateID(idStr)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	result := h.DB.First(r.Context(), &user, id)
	if result.Error != nil {
		log.Printf("User not found: %v", result.Error)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// CreateUserHandler creates a new user.
// This endpoint is publicly accessible and does not require authentication.
//
// @Summary Create a new user
// @Description Create a new user with the provided data.
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 201 {object} models.User
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Failed to create user"
// @Router /users [post]
func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the validated and sanitized payload from the context
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.User)
	if !ok {
		http.Error(w, "Invalid payload", http.StatusInternalServerError)
		return
	}

	// Check if a user with the same username already exists
	var existingUser models.User
	result := h.DB.Where(r.Context(), "username = ?", payload.Username).First(&existingUser)
	if result.Error == nil {
		log.Printf("User with username %s already exists", payload.Username)
		http.Error(w, "User with this username already exists", http.StatusConflict)
		return
	}

	// Check if a user with the same email already exists
	result = h.DB.Where(r.Context(), "email = ?", payload.Email).First(&existingUser)
	if result.Error == nil {
		log.Printf("User with email %s already exists", payload.Email)
		http.Error(w, "User with this email already exists", http.StatusConflict)
		return
	}

	// Create the user
	result = h.DB.Create(r.Context(), payload)
	if result.Error != nil {
		log.Printf("Failed to create user: %v", result.Error)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	log.Printf("User %s created successfully", payload.Username)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// UpdateUserHandler updates an existing user.
// This endpoint requires authentication.
//
// @Summary Update a user
// @Description Update an existing user with the provided data. Requires authentication.
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.User true "Updated user data"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Failed to update user"
// @Security ApiKeyAuth
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Validate ID
	id, err := validation.ValidateID(idStr)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the validated and sanitized payload from the context
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.User)
	if !ok {
		http.Error(w, "Invalid payload", http.StatusInternalServerError)
		return
	}

	// Fetch the existing user
	var user models.User
	result := h.DB.First(r.Context(), &user, id)
	if result.Error != nil {
		log.Printf("User not found: %v", result.Error)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Update the user fields
	user.Username = payload.Username
	user.Email = payload.Email
	user.Password = payload.Password
	user.IsAdmin = payload.IsAdmin

	// Save the updated user
	result = h.DB.Save(r.Context(), &user)
	if result.Error != nil {
		log.Printf("Failed to update user: %v", result.Error)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	log.Printf("User %s updated successfully", user.Username)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// DeleteUserHandler deletes a user by their ID.
// This endpoint requires authentication.
//
// @Summary Delete a user
// @Description Delete a user by their ID. Requires authentication.
// @Tags users
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Failed to delete user"
// @Security ApiKeyAuth
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Validate ID
	id, err := validation.ValidateID(idStr)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := h.DB.Delete(r.Context(), &models.User{}, id)
	if result.Error != nil {
		log.Printf("Failed to delete user: %v", result.Error)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	log.Printf("User with ID %d deleted successfully", id)
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
