package services

import (
	"context"
	"net/http"
	"time"

	"myanimeapi/api/models"
	"myanimeapi/api/repositories"
	"myanimeapi/internal/db"
	apperrors "myanimeapi/internal/errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// PasswordResetService handles password reset operations
type PasswordResetService struct {
	userRepo    repositories.UserRepository
	emailSvc    *EmailService
	db          db.DBInterface
	tokenExpiry time.Duration
}

// NewPasswordResetService creates a new instance of PasswordResetService
func NewPasswordResetService(userRepo repositories.UserRepository, emailSvc *EmailService, db db.DBInterface) *PasswordResetService {
	return &PasswordResetService{
		userRepo:    userRepo,
		emailSvc:    emailSvc,
		db:          db,
		tokenExpiry: time.Hour, // Token expires in 1 hour
	}
}

// RequestPasswordReset initiates a password reset request
func (s *PasswordResetService) RequestPasswordReset(ctx context.Context, email string) error {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// Don't reveal if the email exists or not for security reasons
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
		return apperrors.NewError(apperrors.ErrInternalServer, "Failed to create reset token", result.Error.Error(), http.StatusInternalServerError)
	}

	// Send password reset email
	if err := s.emailSvc.SendPasswordResetEmail(email, token); err != nil {
		return apperrors.NewError(apperrors.ErrInternalServer, "Failed to send reset email", err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// ResetPassword resets a user's password using a valid token
func (s *PasswordResetService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Find the token in the database
	var resetToken models.PasswordResetToken
	result := s.db.Where(ctx, "token = ? AND expires_at > ?", token, time.Now()).First(&resetToken)
	if result.Error != nil {
		return apperrors.NewError(apperrors.ErrInvalidInput, "Invalid or expired token", "The provided token is invalid or has expired", http.StatusBadRequest)
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperrors.NewError(apperrors.ErrInternalServer, "Failed to hash password", err.Error(), http.StatusInternalServerError)
	}

	// Get the user
	user, err := s.userRepo.GetByID(ctx, resetToken.UserID)
	if err != nil {
		return apperrors.NewError(apperrors.ErrResourceNotFound, "User not found", err.Error(), http.StatusNotFound)
	}

	// Type assert the user to *models.User
	userModel, ok := user.(*models.User)
	if !ok {
		return apperrors.NewError(apperrors.ErrInternalServer, "Invalid user type", "Type assertion failed for user model", http.StatusInternalServerError)
	}

	// Update the user's password
	userModel.Password = string(hashedPassword)
	if err := s.userRepo.Update(ctx, userModel); err != nil {
		return apperrors.NewError(apperrors.ErrInternalServer, "Failed to update password", err.Error(), http.StatusInternalServerError)
	}

	// Delete the used token
	result = s.db.Delete(ctx, &resetToken)
	if result.Error != nil {
		return apperrors.NewError(apperrors.ErrInternalServer, "Failed to delete reset token", result.Error.Error(), http.StatusInternalServerError)
	}

	return nil
}
