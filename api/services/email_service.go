package services

import (
	"fmt"
	"net/smtp"
	"os"
)

// EmailService handles sending emails
type EmailService struct {
	from     string
	password string
	smtpHost string
	smtpPort string
}

// NewEmailService creates a new email service
func NewEmailService() *EmailService {
	return &EmailService{
		from:     os.Getenv("EMAIL_FROM"),
		password: os.Getenv("EMAIL_PASSWORD"),
		smtpHost: os.Getenv("SMTP_HOST"),
		smtpPort: os.Getenv("SMTP_PORT"),
	}
}

// SendPasswordResetEmail sends a password reset email
func (s *EmailService) SendPasswordResetEmail(to, resetToken string) error {
	subject := "Password Reset Request"
	body := fmt.Sprintf(`
		Hello,
		
		You have requested to reset your password. Click the link below to reset your password:
		
		%s/reset-password?token=%s
		
		If you did not request this, please ignore this email.
		
		This link will expire in 1 hour.
		
		Best regards,
		MyAnimeAPI Team
	`, os.Getenv("FRONTEND_URL"), resetToken)

	message := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)

	auth := smtp.PlainAuth("", s.from, s.password, s.smtpHost)
	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)

	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(message))
}
