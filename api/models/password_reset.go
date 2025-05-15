package models

import "time"

// PasswordResetRequest represents a request to reset a password
type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// PasswordResetResponse represents a response to a password reset request
type PasswordResetResponse struct {
	Message string `json:"message"`
}

// ResetPasswordRequest represents a request to set a new password
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// PasswordResetToken represents a password reset token in the database
type PasswordResetToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	Token     string    `gorm:"not null,unique"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
