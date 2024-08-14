// internal/handlers/user_handlers.go
package handlers

import (
	"encoding/json"
	"myanimeapi/pkg/middleware"
	"myanimeapi/pkg/models"
	"net/http"

	"github.com/gorilla/mux"
)

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
	router.Use(middleware.Authenticate)
	router.HandleFunc("/users", CreateUserHandler).Methods("POST")
	router.HandleFunc("/users/{id:[0-9]+}", GetUserHandler).Methods("GET")
	router.HandleFunc("/users/{id:[0-9]+}", UpdateUserHandler).Methods("PUT")
	router.HandleFunc("/users/{id:[0-9]+}", DeleteUserHandler).Methods("DELETE")
}
