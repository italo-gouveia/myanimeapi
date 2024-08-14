package main

/*import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAnimeE2E(t *testing.T) {
	// Start the server
	go main() // Ensure your main function starts the server

	// Create a new anime
	newAnime := map[string]interface{}{
		"title":       "One Piece",
		"description": "Pirate anime",
		"rating":      5.0,
	}
	jsonData, _ := json.Marshal(newAnime)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/anime", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check the status code
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Check the response body
	var responseBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.Equal(t, "One Piece", responseBody["title"])
}
*/
