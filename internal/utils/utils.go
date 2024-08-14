package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomString generates a random string of the specified length.
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
