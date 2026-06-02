package services

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"myanimeapi/api/auth"
	"myanimeapi/api/middleware"
	"myanimeapi/api/models"
	"myanimeapi/api/repositories"
	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"
)

// AuthServiceInterface defines the interface for authentication service operations
type AuthServiceInterface interface {
	// RegisterUser handles user registration
	RegisterUser(ctx context.Context, user *models.User) error
	// AuthenticateUser handles user authentication and returns a JWT token
	AuthenticateUser(ctx context.Context, credentials *models.UserCredentials) (string, error)
}

// AuthService handles authentication-related business logic.
// It implements the AuthServiceInterface.
type AuthService struct {
	authRepo repositories.AuthRepository
	logger   *logger.Logger
}

// NewAuthService creates a new instance of AuthService.
// It returns an AuthServiceInterface implementation.
func NewAuthService(authRepo repositories.AuthRepository) AuthServiceInterface {
	return &AuthService{
		authRepo: authRepo,
		logger:   logger.New(),
	}
}

// RegisterUser handles user registration.
// It checks for existing users with the same username or email,
// hashes the password, and creates the user.
func (s *AuthService) RegisterUser(ctx context.Context, user *models.User) error {
	s.logger.WithFields(map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
	}).Info("Starting user registration")

	// Check if username exists
	existingUser, err := s.authRepo.GetUserByUsername(ctx, user.Username)
	if err == nil && existingUser != nil {
		s.logger.WithField("username", user.Username).Warning("Username already exists")
		return errors.NewError(errors.ErrConflict, "Username already exists",
			fmt.Sprintf("Username '%s' is already taken. Please choose a different username.", user.Username),
			http.StatusConflict,
			map[string]interface{}{
				"username": user.Username,
				"email":    user.Email,
			},
			nil)
	}

	// Check if email exists
	existingUser, err = s.authRepo.GetUserByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		s.logger.WithField("email", user.Email).Warning("Email already exists")
		return errors.NewError(errors.ErrConflict, "Email already exists",
			fmt.Sprintf("Email '%s' is already registered. Please use a different email address.", user.Email),
			http.StatusConflict,
			map[string]interface{}{
				"username": user.Username,
				"email":    user.Email,
			},
			nil)
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"username": user.Username,
			"error":    err.Error(),
		}).Error("Failed to hash password")
		return errors.NewError(errors.ErrInternalServer, "Failed to hash password",
			"An error occurred while securing your password. Please try again.",
			http.StatusInternalServerError,
			map[string]interface{}{
				"username": user.Username,
				"email":    user.Email,
			},
			err)
	}
	user.Password = hashedPassword

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Create the user
	if err := s.authRepo.CreateUser(ctx, user); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"username": user.Username,
			"error":    err.Error(),
		}).Error("Failed to create user")
		return errors.NewError(errors.ErrInternalServer, "Failed to create user",
			"An error occurred while creating your account. Please try again.",
			http.StatusInternalServerError,
			map[string]interface{}{
				"username": user.Username,
				"email":    user.Email,
			},
			err)
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
	}).Info("Successfully registered user")
	return nil
}

// AuthenticateUser handles user authentication.
// It verifies the user's credentials and returns a JWT token if valid.
func (s *AuthService) AuthenticateUser(ctx context.Context, credentials *models.UserCredentials) (string, error) {
	s.logger.WithField("username", credentials.Username).Info("Starting user authentication")

	// Get user by username
	user, err := s.authRepo.GetUserByUsername(ctx, credentials.Username)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"username": credentials.Username,
			"error":    err.Error(),
		}).Error("Failed to retrieve user")
		return "", errors.NewError(errors.ErrUnauthorized, "Invalid credentials",
			"The provided username or password is incorrect. Please check your credentials and try again.",
			http.StatusUnauthorized,
			map[string]interface{}{
				"username": credentials.Username,
			},
			err)
	}

	// Verify password and migrate hash if necessary
	newHash, err := auth.MigrateHash(credentials.Password, user.Password)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"username": credentials.Username,
			"error":    err.Error(),
		}).Error("Invalid password")
		return "", errors.NewError(errors.ErrUnauthorized, "Invalid credentials",
			"The provided username or password is incorrect. Please check your credentials and try again.",
			http.StatusUnauthorized,
			map[string]interface{}{
				"username": credentials.Username,
			},
			err)
	}

	// Update password hash if it was migrated
	if newHash != user.Password {
		user.Password = newHash
		user.UpdatedAt = time.Now()
		if err := s.authRepo.UpdateUser(ctx, user); err != nil {
			s.logger.WithFields(map[string]interface{}{
				"username": credentials.Username,
				"error":    err.Error(),
			}).Error("Failed to update user hash")
			return "", errors.NewError(errors.ErrInternalServer, "Failed to update user hash",
				"An error occurred while updating your account. Please try again.",
				http.StatusInternalServerError,
				map[string]interface{}{
					"username": credentials.Username,
				},
				err)
		}
	}

	// Convert user ID to string for token generation
	userIDStr := strconv.FormatUint(uint64(user.ID), 10)

	// Generate JWT token
	token, err := middleware.GenerateToken(userIDStr, user.IsAdmin, user.Role)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"username": credentials.Username,
			"error":    err.Error(),
		}).Error("Failed to generate token")
		return "", errors.NewError(errors.ErrInternalServer, "Failed to generate token",
			"An error occurred while generating your authentication token. Please try again.",
			http.StatusInternalServerError,
			map[string]interface{}{
				"username": credentials.Username,
			},
			err)
	}

	s.logger.WithFields(map[string]interface{}{
		"username": credentials.Username,
		"user_id":  user.ID,
	}).Info("Successfully authenticated user")
	return token, nil
}
