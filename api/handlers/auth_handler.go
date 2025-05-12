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

	"github.com/gorilla/mux"
)

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
// @Param user body models.User true "User registration data"
// @Success 201 {object} models.Response
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 409 {object} errors.ErrorResponse "Username or email already exists"
// @Failure 500 {object} errors.ErrorResponse "Failed to register user"
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
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid input", "Details about the error")
		return
	}

	if err := h.authService.RegisterUser(r.Context(), &user); err != nil {
		errors.WriteAppErrorResponse(w, err)
		return
	}

	response := models.Response{
		Status:  "success",
		Message: "User registered successfully",
		Data:    user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.")
		return
	}
}

// AuthenticateHandler handles user authentication requests.
// It validates the credentials and generates a JWT token upon successful authentication.
// If successful, it returns the token as a JSON response.
//
// @Summary Authenticate a user
// @Description Authenticate a user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.UserCredentials true "User credentials"
// @Success 200 {object} models.Response
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 401 {object} errors.ErrorResponse "Invalid credentials"
// @Failure 500 {object} errors.ErrorResponse "Failed to authenticate user"
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
	var credentials models.UserCredentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid input", "Details about the error")
		return
	}

	token, err := h.authService.AuthenticateUser(r.Context(), &credentials)
	if err != nil {
		errors.WriteAppErrorResponse(w, err)
		return
	}

	response := models.Response{
		Status:  "success",
		Message: "Authentication successful",
		Data: map[string]string{
			"token": token,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Failed to encode response", "An internal server error occurred while encoding the response.")
		return
	}
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
