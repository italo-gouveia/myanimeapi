// pkg/auth/auth.go
// This package defines functions to hash and compare passwords using bcrypt.
// It is used by the user model to hash and compare passwords.
package auth

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the given password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return "", err
	}
	return string(hash), nil
}

// CheckPasswordHash compares the hashed password with the plain text password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Printf("Password mismatch: %v", err)
		return false
	}
	return true
}
