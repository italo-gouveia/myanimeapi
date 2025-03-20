// api/utils/validation.go
// Package validation provides utility functions for validating and parsing input data.
// It includes functions for validating IDs and pagination parameters.
package utils

import (
	"errors"
	"strconv"
)

// ValidateID validates and converts a string ID to a uint.
// It ensures the ID is a valid unsigned integer and returns an error if the format is invalid.
//
// Example:
//
//	id, err := ValidateID("123")
//	if err != nil {
//	    log.Fatalf("Invalid ID: %v", err)
//	}
func ValidateID(idStr string) (uint, error) {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, errors.New("invalid ID format")
	}
	return uint(id), nil
}

// ValidatePagination validates and extracts pagination parameters (page and limit).
// It ensures the page and limit values are valid integers within acceptable ranges.
// If the page or limit values are not provided, default values are used.
//
// Example:
//
//	page, limit, err := ValidatePagination("1", "10", 1, 10)
//	if err != nil {
//	    log.Fatalf("Invalid pagination parameters: %v", err)
//	}
func ValidatePagination(pageStr, limitStr string, defaultPage, defaultLimit int) (page, limit int, err error) {
	// Parse page
	page = defaultPage
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			return 0, 0, errors.New("invalid page number. Must be a positive integer")
		}
	}

	// Parse limit
	limit = defaultLimit
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 100 {
			return 0, 0, errors.New("invalid limit number. Must be a positive integer between 1 and 100")
		}
	}

	return page, limit, nil
}
