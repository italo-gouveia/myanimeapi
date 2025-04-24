package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// ParseUint parses a string to uint
func ParseUint(s string) (uint, error) {
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}

// WriteErrorResponse writes an error response to the HTTP response writer
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		// If we can't encode the error response, we can't do much more
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}

// WriteJSONResponse writes a JSON response to the HTTP response writer
func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// If we can't encode the response, we can't do much more
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetPaginationParams extracts pagination parameters from the request
func GetPaginationParams(r *http.Request) (page, limit int) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, _ = strconv.Atoi(pageStr)
	limit, _ = strconv.Atoi(limitStr)

	// Default values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	return page, limit
}
