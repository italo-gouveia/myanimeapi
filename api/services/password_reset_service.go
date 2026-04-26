package services

import (
	"context"
	"net/http"
	"time"

	"os"
	"strconv"

	"myanimeapi/api/auth"
	"myanimeapi/api/models"
	"myanimeapi/api/repositories"
	"myanimeapi/internal/db"
	apperrors "myanimeapi/internal/errors"
	"myanimeapi/internal/logger"

	"github.com/google/uuid"
)

// PasswordResetServiceInterface defines the interface for password reset operations
type PasswordResetServiceInterface interface {
	// RequestPasswordReset initiates a password reset request for a user
	RequestPasswordReset(ctx context.Context, email string) error
	// ResetPassword resets a user's password using a valid token
	ResetPassword(ctx context.Context, token, newPassword string) error
}

// PasswordResetService handles password reset operations
type PasswordResetService struct {
	userRepo    repositories.UserRepository
	emailSvc    *EmailService
	db          db.DBInterface
	tokenExpiry time.Duration
	logger      *logger.Logger
}

// resetTokenExpiry returns token expiry from RESET_TOKEN_EXPIRY_MINUTES env var, defaulting to 60m.
func resetTokenExpiry() time.Duration {
	if m := os.Getenv("RESET_TOKEN_EXPIRY_MINUTES"); m != "" {
		if mins, err := strconv.Atoi(m); err == nil && mins > 0 {
			return time.Duration(mins) * time.Minute
		}
	}
	return time.Hour
}

// NewPasswordResetService creates a new instance of PasswordResetService
func NewPasswordResetService(userRepo repositories.UserRepository, emailSvc *EmailService, db db.DBInterface) *PasswordResetService {
	return &PasswordResetService{
		userRepo:    userRepo,
		emailSvc:    emailSvc,
		db:          db,
		tokenExpiry: resetTokenExpiry(),
		logger:      logger.New(),
	}
}

// RequestPasswordReset initiates a password reset request
func (s *PasswordResetService) RequestPasswordReset(ctx context.Context, email string) error {
	s.logger.WithField("email", email).Info("Initiating password reset request")

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// Don't reveal if the email exists or not for security reasons
		s.logger.WithField("email", email).Info("Password reset requested for non-existent email")
		return nil
	}

	// Generate a unique token
	token := uuid.New().String()
	expiresAt := time.Now().Add(s.tokenExpiry)

	// Store the token in the database
	resetToken := &models.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	result := s.db.Create(ctx, resetToken)
	if result.Error != nil {
		s.logger.WithFields(map[string]interface{}{
			"email": email,
			"error": result.Error.Error(),
		}).Error("Failed to create reset token")
		return apperrors.NewError(apperrors.ErrInternalServer, "Failed to create reset token", result.Error.Error(), http.StatusInternalServerError, map[string]interface{}{
			"email": email,
		}, result.Error)
	}

	// Send password reset email
	if err := s.emailSvc.SendPasswordResetEmail(email, token); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"email": email,
			"error": err.Error(),
		}).Error("Failed to send reset email")
		return apperrors.NewError(apperrors.ErrInternalServer, "Failed to send reset email", err.Error(), http.StatusInternalServerError, map[string]interface{}{
			"email": email,
		}, err)
	}

	s.logger.WithField("email", email).Info("Password reset email sent successfully")
	return nil
}

// ResetPassword resets a user's password using a valid token
func (s *PasswordResetService) ResetPassword(ctx context.Context, token, newPassword string) error {
	s.logger.WithField("token", token).Info("Processing password reset request")

	// Find the token in the database
	var resetToken models.PasswordResetToken
	result := s.db.Where(ctx, "token = ? AND expires_at > ?", token, time.Now()).First(&resetToken)
	if result.Error != nil {
		s.logger.WithFields(map[string]interface{}{
			"token": token,
			"error": result.Error.Error(),
		}).Error("Invalid or expired reset token")
		return apperrors.NewError(apperrors.ErrInvalidInput, "Invalid or expired token", "The provided token is invalid or has expired", http.StatusBadRequest, map[string]interface{}{
			"token": token,
		}, result.Error)
	}

	// Hash the new password using Argon2 (consistent with user registration)
	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"token": token,
			"error": err.Error(),
		}).Error("Failed to hash new password")
		return apperrors.NewError(apperrors.ErrInternalServer, "Failed to hash password", err.Error(), http.StatusInternalServerError, map[string]interface{}{
			"token": token,
		}, err)
	}

	// Get the user
	user, err := s.userRepo.GetByID(ctx, resetToken.UserID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"token": token,
			"error": err.Error(),
		}).Error("Failed to find user for reset token")
		return apperrors.NewError(apperrors.ErrResourceNotFound, "User not found", err.Error(), http.StatusNotFound, map[string]interface{}{
			"token":   token,
			"user_id": resetToken.UserID,
		}, err)
	}

	// Type assert the user to *models.User
	userModel, ok := user.(*models.User)
	if !ok {
		s.logger.WithFields(map[string]interface{}{
			"token":   token,
			"user_id": resetToken.UserID,
		}).Error("Invalid user type")
		return apperrors.NewError(apperrors.ErrInternalServer, "Invalid user type", "Type assertion failed for user model", http.StatusInternalServerError, map[string]interface{}{
			"token":   token,
			"user_id": resetToken.UserID,
		}, nil)
	}

	// Update the user's password
	userModel.Password = hashedPassword
	if err := s.userRepo.Update(ctx, userModel); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"token":   token,
			"user_id": resetToken.UserID,
			"error":   err.Error(),
		}).Error("Failed to update user password")
		return apperrors.NewError(apperrors.ErrInternalServer, "Failed to update password", err.Error(), http.StatusInternalServerError, map[string]interface{}{
			"token":   token,
			"user_id": resetToken.UserID,
		}, err)
	}

	// Delete the used token
	result = s.db.Delete(ctx, &resetToken)
	if result.Error != nil {
		s.logger.WithFields(map[string]interface{}{
			"token": token,
			"error": result.Error.Error(),
		}).Error("Failed to delete reset token")
		return apperrors.NewError(apperrors.ErrInternalServer, "Failed to delete reset token", result.Error.Error(), http.StatusInternalServerError, map[string]interface{}{
			"token": token,
		}, result.Error)
	}

	s.logger.WithFields(map[string]interface{}{
		"token":   token,
		"user_id": resetToken.UserID,
	}).Info("Password reset completed successfully")
	return nil
}
