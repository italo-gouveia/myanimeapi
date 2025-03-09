// pkg/validation/validation.go
package validation

import (
	"errors"
	"strconv"
)

// ValidateID validates and converts a string ID to a uint.
func ValidateID(idStr string) (uint, error) {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, errors.New("invalid ID format")
	}
	return uint(id), nil
}

// ValidatePagination validates and extracts pagination parameters (page and limit).
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
