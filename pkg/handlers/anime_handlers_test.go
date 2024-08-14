package handlers

import (
	"encoding/json"
	"myanimeapi/pkg/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type Database interface {
	First(out interface{}, where ...interface{}) *gorm.DB
}

type Handler struct {
	DB Database
}

func NewHandler(db Database) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract the anime ID from the URL
	id := r.URL.Path[len("/anime/"):]

	// Create an instance of the Anime model
	var anime models.Anime

	// Query the database for the anime with the given ID
	result := h.DB.First(&anime, id)
	if result.Error != nil {
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the anime data to JSON and write it to the response
	if err := json.NewEncoder(w).Encode(anime); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func TestGetAnimeHandler(t *testing.T) {
	// Create a new instance of the mock database
	mockDB := new(MockDB)

	// Create an instance of the handler with the mock database
	handler := NewHandler(mockDB)

	// Define the anime data to return from the mock
	animeID := "1"
	expectedAnime := &models.Anime{
		Model: gorm.Model{
			ID: 1, // Ensure the ID matches what you expect
		},
		Title:       "Naruto",
		Description: "A story about a ninja.",
		Rating:      8.5,
	}

	// Set up the mock expectations
	mockDB.On("First", mock.AnythingOfType("*models.Anime"), animeID).Run(func(args mock.Arguments) {
		out := args.Get(0).(*models.Anime)
		*out = *expectedAnime
	}).Return(&gorm.DB{})

	// Create a new HTTP request for testing
	req, err := http.NewRequest("GET", "/anime/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Create the handler function and pass in the ResponseRecorder and Request
	handler.ServeHTTP(rr, req)

	// Assert that the HTTP status code is 200 OK
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert that the response body contains the expected anime data
	expectedResponse := `{"id":1,"title":"Naruto","description":"A story about a ninja.","rating":8.5}`
	assert.JSONEq(t, expectedResponse, rr.Body.String())

	// Assert that the mock expectations were met
	mockDB.AssertExpectations(t)
}
