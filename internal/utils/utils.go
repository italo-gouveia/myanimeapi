// internal/utils/utils.go
// Package utils provides utility functions for common tasks.
// This file contains a function to generate random strings.
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
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
	if length <= 0 {
		return "", fmt.Errorf("length must be greater than 0")
	}
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Printf("Error generating random bytes: %v", err)
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
