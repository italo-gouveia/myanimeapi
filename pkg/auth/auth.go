// pkg/auth/auth.go
// This package defines functions to hash and compare passwords using Argon2.
// Argon2 is a modern and secure password hashing algorithm that is resistant to GPU-based attacks.
package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

// Define the parameters for Argon2
var (
	argon2Time    uint32 = 1         // Number of iterations
	argon2Memory  uint32 = 64 * 1024 // 64 MB of memory
	argon2Threads uint8  = 4         // Number of threads
	argon2KeyLen  uint32 = 32        // Length of the generated hash
)

// HashPassword hashes the given password using Argon2
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

// CheckPasswordHash compares the hashed password with the plain text password
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
	comparisonHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(hash)))

	// Compare the hashes in a constant-time manner
	if subtle.ConstantTimeCompare(hash, comparisonHash) == 1 {
		return true
	}

	log.Printf("Password mismatch")
	return false
}

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
