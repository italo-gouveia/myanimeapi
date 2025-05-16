// internal/utils/utils.go
// Package utils provides utility functions for common tasks.
// This file contains a function to generate random strings.
package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"myanimeapi/internal/logger"
)

// GenerateRandomString generates a cryptographically secure random string of the specified length.
// The string is encoded in hexadecimal format, so the output length will be twice the input length.
//
// Parameters:
//   - length: The desired length of the random string in bytes. Must be greater than 0.
//
// Returns:
//   - A random string encoded in hexadecimal format.
//   - An error if the length is invalid or if the random bytes cannot be generated.
//
// Example usage:
//
//	randomString, err := GenerateRandomString(16)
//	if err != nil {
//	    log.Fatalf("Failed to generate random string: %v", err)
//	}
//	fmt.Println(randomString)
func GenerateRandomString(length int) (string, error) {
	log := logger.Get()

	if length <= 0 {
		return "", fmt.Errorf("length must be greater than 0")
	}
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.WithError(err).Error("Error generating random bytes")
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
