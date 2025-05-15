// api/utils/validation.go
// Package validation provides utility functions for validating and parsing input data.
// It includes functions for validating IDs and pagination parameters.
package utils

import (
	"errors"
	"math"
	"strconv"
)

// ValidateID validates and converts a string ID to a uint.
// It ensures the ID is a valid positive unsigned integer and fits within the range of uint.
//
// Example:
//
//	id, err := ValidateID("123")
//	if err != nil {
//	    log.Fatalf("Invalid ID: %v", err)
//	}
func ValidateID(idStr string) (uint, error) {
	id, err := strconv.ParseUint(idStr, 10, 64) // Parse as uint64 to detect overflow beyond typical int/uint32 ranges
	if err != nil {
		var numErr *strconv.NumError
		if errors.As(err, &numErr) && numErr.Err == strconv.ErrRange {
			return 0, errors.New("invalid ID: out of range")
		}
		return 0, errors.New("invalid ID format")
	}

	if id == 0 {
		return 0, errors.New("invalid ID: must be a positive number")
	}

	// Check if it fits in a platform-dependent uint (typically uint32 or uint64)
	// This check is more relevant if your model's ID field is specifically uint32
	// and you want to ensure it doesn't overflow that before casting.
	// If model.ID is uint (which can be uint64 on 64-bit systems), this specific check might be redundant
	// as ParseUint(..., 64) already ensures it fits uint64.
	// However, if `uint` is `uint32` on a target system, this is important.
	// For simplicity and wider compatibility, we ensure it fits in what `uint` can represent.
	if id > math.MaxUint {
		return 0, errors.New("invalid ID: exceeds maximum value for uint type on this system")
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
