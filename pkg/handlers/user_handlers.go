// internal/handlers/user_handlers.go
package handlers

import (
	"encoding/json"
	"myanimeapi/pkg/middleware"
	"myanimeapi/pkg/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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
func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	isAdmin, ok := r.Context().Value("is_admin").(bool)
	if !ok || !isAdmin {
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
		if err != nil {
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}
	}

	// Parse limit number
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit number", http.StatusBadRequest)
			return
		}
	}

	// Calculate offset
	offset := (page - 1) * limit

	var users []models.User
	if err := database.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	// Set response header and encode the result
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
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
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var user models.User
	if err := database.First(&user, id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
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
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := database.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
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
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var user models.User
	if err := database.First(&user, id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := database.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
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
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := database.Delete(&models.User{}, id).Error; err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func RegisterUserRoutes(router *mux.Router) {
	// Public routes (no authentication required)
	router.HandleFunc("/users", CreateUserHandler).Methods("POST")

	// Create a subrouter for protected routes
	protectedRouter := router.PathPrefix("/users").Subrouter()
	protectedRouter.Use(middleware.Authenticate) // Apply authentication middleware

	// Protected routes (require authentication)
	protectedRouter.HandleFunc("/{id:[0-9]+}", GetUserHandler).Methods("GET")
	protectedRouter.HandleFunc("/{id:[0-9]+}", UpdateUserHandler).Methods("PUT")
	protectedRouter.HandleFunc("/{id:[0-9]+}", DeleteUserHandler).Methods("DELETE")

	// Create a subrouter for admin-only routes
	adminRouter := protectedRouter.PathPrefix("").Subrouter()
	adminRouter.Use(middleware.CheckAdmin) // Apply admin check middleware

	// Admin-only routes (require authentication and admin privileges)
	adminRouter.HandleFunc("", GetAllUsersHandler).Methods("GET")
}
