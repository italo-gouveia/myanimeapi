// internal/utils/utils.go
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
)

// GenerateRandomString generates a random string of the specified length.
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
