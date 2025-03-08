// internal/handlers/user_handlers.go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"myanimeapi/internal/db"
	"myanimeapi/pkg/middleware"
	"myanimeapi/pkg/models"

	"github.com/gorilla/mux"
)

// UserHandler defines the handlers for user-related routes
type UserHandler struct {
	DB db.DBInterface
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(db db.DBInterface) *UserHandler {
	return &UserHandler{DB: db}
}

// GetAllUsersHandler retrieves paginated user entries (admin access required)
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

	// Extract pagination parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Default values
	page := 1
	limit := 10
	var err error

	// Parse page number
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			log.Printf("Invalid page number: %v", err)
			http.Error(w, "Invalid page number. Must be a positive integer.", http.StatusBadRequest)
			return
		}
	}

	// Parse limit number
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 100 {
			log.Printf("Invalid limit number: %v", err)
			http.Error(w, "Invalid limit number. Must be a positive integer between 1 and 100.", http.StatusBadRequest)
			return
		}
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

// GetUserHandler retrieves a user by ID
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

	// Convert the ID from string to uint
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
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

// CreateUserHandler creates a new user
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
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Invalid input: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	result := h.DB.Create(r.Context(), &user)
	if result.Error != nil {
		log.Printf("Failed to create user: %v", result.Error)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	log.Printf("User %s created successfully", user.Username)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// UpdateUserHandler updates an existing user
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

	// Convert the ID from string to uint
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var user models.User
	result := h.DB.First(r.Context(), &user, id)
	if result.Error != nil {
		log.Printf("User not found: %v", result.Error)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Invalid input: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

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

// DeleteUserHandler deletes a user
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

	// Convert the ID from string to uint
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
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

// RegisterUserRoutes registers all user-related routes
func (h *UserHandler) RegisterUserRoutes(router *mux.Router) {
	// Public routes (no authentication required)
	router.HandleFunc("/users", h.CreateUserHandler).Methods("POST")

	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/users").Subrouter()
	protectedRouter.Use(middleware.Authenticate) // Apply authentication middleware

	// Protected routes (require authentication)
	protectedRouter.HandleFunc("/{id:[0-9]+}", h.GetUserHandler).Methods("GET")
	protectedRouter.HandleFunc("/{id:[0-9]+}", h.UpdateUserHandler).Methods("PUT")
	protectedRouter.HandleFunc("/{id:[0-9]+}", h.DeleteUserHandler).Methods("DELETE")

	// Create a subrouter for admin-only routes
	adminRouter := protectedRouter.PathPrefix("").Subrouter()
	adminRouter.Use(middleware.CheckAdmin) // Apply admin check middleware

	// Admin-only routes (require authentication and admin privileges)
	adminRouter.HandleFunc("", h.GetAllUsersHandler).Methods("GET")
}
