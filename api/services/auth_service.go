package services

import (
	"context"
	"fmt"
	"net/http"

	"myanimeapi/api/auth"
	"myanimeapi/api/middleware"
	"myanimeapi/api/models"
	"myanimeapi/api/repositories"
	"myanimeapi/internal/errors"
)

// AuthService handles authentication-related business logic.
type AuthService struct {
	authRepo repositories.AuthRepository
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(authRepo repositories.AuthRepository) *AuthService {
	return &AuthService{
		authRepo: authRepo,
	}
}

// RegisterUser handles user registration.
// It checks for existing users with the same username or email,
// hashes the password, and creates the user.
func (s *AuthService) RegisterUser(ctx context.Context, user *models.User) error {
	// Check if username exists
	existingUser, err := s.authRepo.GetUserByUsername(ctx, user.Username)
	if err == nil && existingUser != nil {
		return errors.NewError(errors.ErrConflict, "Username already exists",
			fmt.Sprintf("Username '%s' is already taken", user.Username), http.StatusConflict)
	}

	// Check if email exists
	existingUser, err = s.authRepo.GetUserByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return errors.NewError(errors.ErrConflict, "Email already exists",
			fmt.Sprintf("Email '%s' is already registered", user.Email), http.StatusConflict)
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	// Create the user
	if err := s.authRepo.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// AuthenticateUser handles user authentication.
// It verifies the user's credentials and returns a JWT token if valid.
func (s *AuthService) AuthenticateUser(ctx context.Context, credentials *models.UserCredentials) (string, error) {
	// Get user by username
	user, err := s.authRepo.GetUserByUsername(ctx, credentials.Username)
	if err != nil {
		return "", errors.NewError(errors.ErrUnauthorized, "Invalid credentials",
			"The provided username or password is incorrect", http.StatusUnauthorized)
	}

	// Verify password and migrate hash if necessary
	newHash, err := auth.MigrateHash(credentials.Password, user.Password)
	if err != nil {
		return "", errors.NewError(errors.ErrUnauthorized, "Invalid credentials",
			"The provided username or password is incorrect", http.StatusUnauthorized)
	}

	// Update password hash if it was migrated
	if newHash != user.Password {
		user.Password = newHash
		if err := s.authRepo.UpdateUser(ctx, user); err != nil {
			return "", fmt.Errorf("failed to update user hash: %w", err)
		}
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID, user.IsAdmin)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
