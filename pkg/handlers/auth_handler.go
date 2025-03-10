// internal/handlers/auth_handler.go
// This package defines the handlers for the authentication routes.
// It is used by the server to handle authentication-related requests.
// It provides handlers for user registration and authentication.
// It uses the auth package to hash and compare passwords.
// It uses the middleware package to generate JWT tokens.
// It uses the models package to interact with the database.
// It uses the gorilla/mux package to handle HTTP requests.
// It uses the http package to write HTTP responses.
// It uses the encoding/json package to encode and decode JSON data.
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"myanimeapi/internal/db"
	"myanimeapi/pkg/auth"
	"myanimeapi/pkg/middleware"
	"myanimeapi/pkg/models"

	"github.com/gorilla/mux"
)

// AuthHandler defines the handlers for authentication-related routes
type AuthHandler struct {
	DB db.DBInterface
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(db db.DBInterface) *AuthHandler {
	return &AuthHandler{DB: db}
}

// RegisterUserHandler registers a new user
// @Summary Register a new user
// @Description Register a new user with the provided data.
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User registration data"
// @Success 201 {object} models.User
// @Failure 400 {string} string "Invalid input or missing required fields"
// @Failure 409 {string} string "User with this username or email already exists"
// @Failure 500 {string} string "Failed to hash password or create user"
// @Router /auth/register [post]
func (h *AuthHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
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

	// Hash the user's password
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	payload.Password = hashedPassword

	// Create the user
	result = h.DB.Create(r.Context(), payload)
	if result.Error != nil {
		log.Printf("Error creating user: %v", result.Error)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	log.Printf("User %s registered successfully", payload.Username)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// AuthenticateHandler authenticates a user and returns a JWT token
// @Summary Authenticate a user
// @Description Authenticate a user and return a JWT token.
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.UserCredentials true "User credentials"
// @Success 200 {object} map[string]string "Returns a JWT token"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "User not found or invalid credentials"
// @Failure 500 {string} string "Failed to generate token"
// @Router /auth/authenticate [post]
func (h *AuthHandler) AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the validated and sanitized payload from the context
	payload, ok := r.Context().Value(middleware.ValidatedPayloadKey).(*models.UserCredentials)
	if !ok {
		http.Error(w, "Invalid payload", http.StatusInternalServerError)
		return
	}

	var user models.User
	result := h.DB.Where(r.Context(), "username = ?", payload.Username).First(&user)
	if result.Error != nil {
		log.Printf("User %s not found", payload.Username)
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Migrate the hash if necessary
	newHash, err := auth.MigrateHash(payload.Password, user.Password)
	if err != nil {
		log.Printf("Invalid credentials for user %s", payload.Username)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Update the database with the new hash if it was migrated
	if newHash != user.Password {
		user.Password = newHash
		result = h.DB.Save(r.Context(), &user)
		if result.Error != nil {
			log.Printf("Failed to update user hash: %v", result.Error)
			http.Error(w, "Failed to update user hash", http.StatusInternalServerError)
			return
		}
	}

	/*	if !auth.CheckPasswordHash(payload.Password, user.Password) {
		log.Printf("Invalid credentials for user %s", payload.Username)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}*/

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID, user.IsAdmin)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	log.Printf("User %s authenticated successfully", payload.Username)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// RegisterAuthRoutes registers all authentication-related routes
func (h *AuthHandler) RegisterAuthRoutes(router *mux.Router) {
	// Public routes (no authentication required)
	router.Handle("/auth/authenticate", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.AuthenticateHandler), models.UserCredentials{})).Methods("POST")
	router.Handle("/auth/register", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.RegisterUserHandler), models.User{})).Methods("POST")
}
