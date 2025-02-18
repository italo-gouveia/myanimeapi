// internal/handlers/auth_handler.go
package handlers

import (
	"encoding/json"
	"net/http"

	"myanimeapi/pkg/auth"
	"myanimeapi/pkg/middleware"
	"myanimeapi/pkg/models"

	"github.com/gorilla/mux"
)

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if username and password are provided
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Check if a user with the same username already exists
	var existingUser models.User
	if err := database.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		http.Error(w, "User with this username already exists", http.StatusConflict)
		return
	}

	// Check if a user with the same email already exists (only if email is provided)
	if user.Email != "" {
		if err := database.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
			http.Error(w, "User with this email already exists", http.StatusConflict)
			return
		}
	}

	// Hash the user's password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	// Create the user
	if err := database.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// AuthenticateHandler authenticates a user and returns a JWT token
func AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := database.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if !auth.CheckPasswordHash(loginRequest.Password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID, user.IsAdmin)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// RegisterAnimeRoutes registers all anime-related routes
func RegisterAuthRoutes(router *mux.Router) {
	// Public routes (no authentication required)
	router.HandleFunc("/auth/authenticate", AuthenticateHandler).Methods("POST")
	router.HandleFunc("/auth/register", RegisterUserHandler).Methods("POST")
}
