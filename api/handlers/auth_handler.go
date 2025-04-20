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
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterUserHandler handles user registration requests.
// @Summary Register a new user
// @Description Register a new user with the provided credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User registration data"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 409 {object} models.Response
// @Router /auth/register [post]
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
	json.NewEncoder(w).Encode(response)
}

// AuthenticateHandler handles user authentication requests.
// @Summary Authenticate a user
// @Description Authenticate a user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.UserCredentials true "User credentials"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Router /auth/login [post]
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
	json.NewEncoder(w).Encode(response)
}

// RegisterAuthRoutes registers all auth-related routes with a *mux.Router
func (h *AuthHandler) RegisterAuthRoutes(router *mux.Router) {
	router.HandleFunc("/auth/register", h.RegisterUserHandler).Methods("POST")
	router.HandleFunc("/auth/login", h.AuthenticateHandler).Methods("POST")
}
