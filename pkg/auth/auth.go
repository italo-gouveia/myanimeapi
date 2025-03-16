// pkg/auth/auth.go
// Package auth provides functions for securely hashing and comparing passwords using the Argon2 algorithm.
// Argon2 is a modern, memory-hard password hashing algorithm designed to resist GPU-based attacks.
// This package also includes a utility function to migrate legacy bcrypt hashes to Argon2.
//
// The package defines default parameters for Argon2, including memory usage, iterations, and parallelism.
// These parameters can be adjusted to meet specific security requirements.
//
// Example usage:
//
//	hashedPassword, err := auth.HashPassword("mysecurepassword")
//	if err != nil {
//	    log.Fatalf("Error hashing password: %v", err)
//	}
//
//	isValid := auth.CheckPasswordHash("mysecurepassword", hashedPassword)
//	if isValid {
//	    log.Println("Password is valid")
//	} else {
//	    log.Println("Password is invalid")
//	}
//
//	migratedHash, err := auth.MigrateHash("mysecurepassword", oldBcryptHash)
//	if err != nil {
//	    log.Fatalf("Error migrating hash: %v", err)
//	}
package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log"
	"math"
	"strings"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

// Default Argon2 parameters for password hashing.
var (
	argon2Time    uint32 = 1         // Number of iterations
	argon2Memory  uint32 = 64 * 1024 // 64 MB of memory
	argon2Threads uint8  = 4         // Number of threads
	argon2KeyLen  uint32 = 32        // Length of the generated hash
)

// HashPassword hashes the given password using the Argon2 algorithm.
// It generates a random salt, hashes the password with the salt, and returns the encoded hash string.
// The encoded hash includes the Argon2 parameters, salt, and hash for later verification.
//
// Example:
//
//	hashedPassword, err := HashPassword("mysecurepassword")
//	if err != nil {
//	    log.Fatalf("Error hashing password: %v", err)
//	}
func HashPassword(password string) (string, error) {
	// Generate a random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		log.Printf("Error generating salt: %v", err)
		return "", err
	}

	// Hash the password using Argon2
	hash := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)

	// Encode the salt and hash to a string
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return the encoded string in a standardized format
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, argon2Memory, argon2Time, argon2Threads, b64Salt, b64Hash)
	return encodedHash, nil
}

// CheckPasswordHash compares a plain-text password with an encoded Argon2 hash.
// It extracts the parameters, salt, and hash from the encoded string, rehashes the password,
// and performs a constant-time comparison to verify the password.
//
// Example:
//
//	isValid := CheckPasswordHash("mysecurepassword", encodedHash)
//	if isValid {
//	    log.Println("Password is valid")
//	} else {
//	    log.Println("Password is invalid")
//	}
func CheckPasswordHash(password, encodedHash string) bool {
	// Split the encoded hash into its components
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		log.Printf("Invalid hash format")
		return false
	}

	// Extract the parameters from the encoded hash
	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		log.Printf("Error parsing version: %v", err)
		return false
	}
	if version != argon2.Version {
		log.Printf("Incompatible Argon2 version")
		return false
	}

	var memory uint32
	var time uint32
	var threads uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		log.Printf("Error parsing parameters: %v", err)
		return false
	}

	// Decode the salt and hash
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		log.Printf("Error decoding salt: %v", err)
		return false
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		log.Printf("Error decoding hash: %v", err)
		return false
	}

	// Hash the provided password with the same parameters
	hashLen := len(hash)
	if hashLen < 0 || hashLen > math.MaxUint32 {
		fmt.Errorf("invalid hash length")
		return false
	}
	comparisonHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(hashLen))

	// Compare the hashes in a constant-time manner
	if subtle.ConstantTimeCompare(hash, comparisonHash) == 1 {
		return true
	}

	log.Printf("Password mismatch")
	return false
}

// MigrateHash migrates a legacy bcrypt hash to an Argon2 hash.
// It verifies the password against the existing bcrypt hash and rehashes it using Argon2.
// If the existing hash is already an Argon2 hash, it is returned as-is.
//
// Example:
//
//	migratedHash, err := MigrateHash("mysecurepassword", oldBcryptHash)
//	if err != nil {
//	    log.Fatalf("Error migrating hash: %v", err)
//	}
func MigrateHash(password, existingHash string) (string, error) {
	// Check if the existing hash is bcrypt
	if strings.HasPrefix(existingHash, "$2a$") || strings.HasPrefix(existingHash, "$2b$") || strings.HasPrefix(existingHash, "$2y$") {
		// Verify the password using bcrypt
		err := bcrypt.CompareHashAndPassword([]byte(existingHash), []byte(password))
		if err != nil {
			return "", err // Password is incorrect
		}

		// Rehash the password using Argon2
		newHash, err := HashPassword(password)
		if err != nil {
			return "", err
		}

		// Return the new Argon2 hash
		return newHash, nil
	}

	// If the hash is already Argon2, return it as-is
	return existingHash, nil
}
