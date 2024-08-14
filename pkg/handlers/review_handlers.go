package handlers

import (
	"encoding/json"
	"myanimeapi/pkg/middleware"
	"myanimeapi/pkg/models"
	"net/http"

	"github.com/gorilla/mux"
)

func GetReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var review models.Review
	if err := database.Preload("User").Preload("Anime").First(&review, id).Error; err != nil {
		http.Error(w, "Review not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

func CreateReviewHandler(w http.ResponseWriter, r *http.Request) {
	var review models.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	// Check if the user exists
	var user models.User
	if err := database.First(&user, review.UserID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if the anime exists
	var anime models.Anime
	if err := database.First(&anime, review.AnimeID).Error; err != nil {
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	// Create the review
	if err := database.Create(&review).Error; err != nil {
		http.Error(w, "Failed to create review", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}

func UpdateReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var review models.Review
	if err := database.Preload("User").Preload("Anime").First(&review, id).Error; err != nil {
		http.Error(w, "Review not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Ensure the user and anime still exist
	var user models.User
	if err := database.First(&user, review.UserID).Error; err != nil {
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	var anime models.Anime
	if err := database.First(&anime, review.AnimeID).Error; err != nil {
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	if err := database.Save(&review).Error; err != nil {
		http.Error(w, "Failed to update review", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

func DeleteReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := database.Delete(&models.Review{}, id); err != nil {
		http.Error(w, "Failed to delete review", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func RegisterReviewRoutes(router *mux.Router) {
	router.Use(middleware.Authenticate)
	router.HandleFunc("/reviews", CreateReviewHandler).Methods("POST")
	router.HandleFunc("/reviews/{id:[0-9]+}", GetReviewHandler).Methods("GET")
	router.HandleFunc("/reviews/{id:[0-9]+}", UpdateReviewHandler).Methods("PUT")
	router.HandleFunc("/reviews/{id:[0-9]+}", DeleteReviewHandler).Methods("DELETE")
}
