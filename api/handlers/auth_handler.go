// pkg/handlers/auth_handler.go
// Package handlers provides HTTP handlers for authentication-related routes in the MyAnimeAPI application.
// It defines methods to handle user registration and authentication, including password hashing and JWT token generation.
// The package uses the Gorilla Mux router for routing, GORM for database interactions, and middleware for request validation and token generation.
//
// Example usage:
//
//	db := // initialize your database connection
//	authHandler := handlers.NewAuthHandler(db)
//	router := mux.NewRouter()
//	authHandler.RegisterAuthRoutes(router)
//
//	http.ListenAndServe(":8080", router)
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"myanimeapi/api/auth"
	"myanimeapi/api/middleware"
	"myanimeapi/api/models"
	"myanimeapi/internal/db"

	"github.com/gorilla/mux"
)

// AuthHandler defines the handlers for authentication-related routes.
// It contains a database interface for interacting with the database.
type AuthHandler struct {
	DB db.DBInterface
}

// NewAuthHandler creates a new instance of AuthHandler.
// It accepts a database interface and returns a pointer to an AuthHandler.
//
// Example:
//
//	db := // initialize your database connection
//	authHandler := NewAuthHandler(db)
func NewAuthHandler(db db.DBInterface) *AuthHandler {
	return &AuthHandler{DB: db}
}

// RegisterUserHandler registers a new user in the database.
// It validates the input payload, checks for existing users with the same username or email,
// hashes the password, and creates the user. If successful, it returns the created user as a JSON response.
//
// @Summary Register a new user
// @Description Register a new user with the provided data
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User registration data"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string "Invalid input or missing required fields"
// @Failure 409 {object} map[string]string "User with this username or email already exists"
// @Failure 500 {object} map[string]string "Failed to hash password or create user"
// @Router /v1/auth/register [post]
// @Example
//
//	{
//	  "username": "john_doe",
//	  "email": "john@example.com",
//	  "password": "password123"
//	}
//
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
// @Security []
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

// AuthenticateHandler authenticates a user and returns a JWT token.
// It validates the input payload, checks the user's credentials, and generates a JWT token if the credentials are valid.
// If the credentials are invalid or the token generation fails, it returns an error response.
//
// @Summary Authenticate a user
// @Description Authenticate a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.UserCredentials true "User credentials"
// @Success 200 {object} map[string]string "Returns a JWT token"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "User not found or invalid credentials"
// @Failure 500 {object} map[string]string "Failed to generate token"
// @Router /v1/auth/authenticate [post]
// @Example
//
//	{
//	  "username": "john_doe",
//	  "password": "password123"
//	}
//
// @ExampleResponse
//
//	{
//	  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
//	}
//
// @Security []
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

// RegisterAuthRoutes registers all authentication-related routes with the provided router.
// It defines public routes for user registration and authentication.
//
// Example:
//
//	router := mux.NewRouter()
//	authHandler.RegisterAuthRoutes(router)
func (h *AuthHandler) RegisterAuthRoutes(router *mux.Router) {
	// Public routes (no authentication required)
	router.Handle("/auth/authenticate", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.AuthenticateHandler), models.UserCredentials{})).Methods("POST")
	router.Handle("/auth/register", middleware.ValidateAndSanitizePayload(http.HandlerFunc(h.RegisterUserHandler), models.User{})).Methods("POST")
}
