// api/handlers/auth_handler.go
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
	"net/http"

	"myanimeapi/api/models"
	"myanimeapi/api/services"
	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// RequestIDContextKey is the key used to store the request ID in the context
const RequestIDContextKey = "request_id"

// AuthHandler handles authentication-related HTTP requests.
// It contains an auth service for handling authentication business logic.
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new instance of AuthHandler.
// It accepts an auth service and returns a pointer to an AuthHandler.
//
// Example:
//
//	authService := services.NewAuthService(userRepo)
//	authHandler := NewAuthHandler(authService)
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterUserHandler handles user registration requests.
// It validates the input payload and uses the service to register the user.
// If successful, it returns the created user as a JSON response.
//
// @Summary Register a new user
// @Description Register a new user with the provided credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User registration information"
// @Success 201 {object} models.User
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /auth/register [post]
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
func (h *AuthHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.WithFields(logrus.Fields{
			"path":       r.URL.Path,
			"method":     r.Method,
			"request_id": r.Context().Value(RequestIDContextKey),
			"error":      err,
		}).Warn("Failed to decode registration request body")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid input", "Details about the error")
		return
	}

	log.WithFields(logrus.Fields{
		"path":       r.URL.Path,
		"method":     r.Method,
		"request_id": r.Context().Value(RequestIDContextKey),
		"username":   user.Username,
		"email":      user.Email,
	}).Info("Processing user registration request")

	if err := h.authService.RegisterUser(r.Context(), &user); err != nil {
		log.WithFields(logrus.Fields{
			"path":       r.URL.Path,
			"method":     r.Method,
			"request_id": r.Context().Value(RequestIDContextKey),
			"error":      err,
			"username":   user.Username,
			"email":      user.Email,
		}).Error("Failed to register user")
		errors.WriteAppErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.WithFields(logrus.Fields{
			"path":       r.URL.Path,
			"method":     r.Method,
			"request_id": r.Context().Value(RequestIDContextKey),
			"error":      err,
			"username":   user.Username,
			"email":      user.Email,
		}).Error("Failed to encode registration response")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.")
		return
	}

	log.WithFields(logrus.Fields{
		"path":       r.URL.Path,
		"method":     r.Method,
		"request_id": r.Context().Value(RequestIDContextKey),
		"username":   user.Username,
		"email":      user.Email,
		"user_id":    user.ID,
	}).Info("User registered successfully")
}

// AuthenticateHandler handles user authentication requests.
// It validates the credentials and generates a JWT token upon successful authentication.
// If successful, it returns the token as a JSON response.
//
// @Summary Authenticate a user
// @Description Authenticate a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.UserCredentials true "User credentials"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /auth/login [post]
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
//	  "message": "Authentication successful",
//	  "data": {
//	    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
//	  }
//	}
func (h *AuthHandler) AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	var credentials models.UserCredentials

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		log.WithFields(logrus.Fields{
			"path":       r.URL.Path,
			"method":     r.Method,
			"request_id": r.Context().Value(RequestIDContextKey),
			"error":      err,
		}).Warn("Failed to decode authentication request body")
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid input", "Details about the error")
		return
	}

	log.WithFields(logrus.Fields{
		"path":       r.URL.Path,
		"method":     r.Method,
		"request_id": r.Context().Value(RequestIDContextKey),
		"username":   credentials.Username,
	}).Info("Processing authentication request")

	token, err := h.authService.AuthenticateUser(r.Context(), &credentials)

	if err != nil {
		log.WithFields(logrus.Fields{
			"path":       r.URL.Path,
			"method":     r.Method,
			"request_id": r.Context().Value(RequestIDContextKey),
			"error":      err,
			"username":   credentials.Username,
		}).Warn("Authentication failed")
		errors.WriteAppErrorResponse(w, err)
		return
	}

	response := models.AuthResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.WithFields(logrus.Fields{
			"path":       r.URL.Path,
			"method":     r.Method,
			"request_id": r.Context().Value(RequestIDContextKey),
			"error":      err,
			"username":   credentials.Username,
		}).Error("Failed to encode authentication response")
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.")
		return
	}

	log.WithFields(logrus.Fields{
		"path":       r.URL.Path,
		"method":     r.Method,
		"request_id": r.Context().Value(RequestIDContextKey),
		"username":   credentials.Username,
	}).Info("Authentication successful")
}

// RegisterAuthRoutes registers all auth-related routes with a *mux.Router.
// It sets up the routes for user registration and authentication.
//
// Routes registered:
// - POST /auth/register - Register a new user
// - POST /auth/login - Authenticate a user
func (h *AuthHandler) RegisterAuthRoutes(router *mux.Router) {
	router.HandleFunc("/auth/register", h.RegisterUserHandler).Methods("POST")
	router.HandleFunc("/auth/login", h.AuthenticateHandler).Methods("POST")
}
