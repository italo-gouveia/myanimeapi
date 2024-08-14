package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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

	// Hash the user's password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

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

	// Generate JWT token (replace with actual JWT generation)
	token := "MYANIMEAPI" + user.Username

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// GetAllUsersHandler retrieves paginated user entries (admin access required)
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

// RegisterAnimeRoutes registers all anime-related routes
func RegisterAuthRoutes(router *mux.Router) {
	router.HandleFunc("/authenticate", AuthenticateHandler).Methods("POST")
	router.HandleFunc("/register", RegisterUserHandler).Methods("POST")
	router.Use(middleware.CheckAdmin)
	router.HandleFunc("/users", GetAllUsersHandler).Methods("GET")
}
