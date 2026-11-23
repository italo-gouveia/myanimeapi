package services

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEmailService_MissingConfig(t *testing.T) {
	require.NoError(t, os.Unsetenv("EMAIL_FROM"))
	require.NoError(t, os.Unsetenv("SMTP_HOST"))
	require.NoError(t, os.Unsetenv("SMTP_PORT"))

	svc, err := NewEmailService()
	assert.Nil(t, svc)
	assert.Error(t, err)
}

func TestNewEmailService_Success(t *testing.T) {
	t.Setenv("EMAIL_FROM", "test@example.com")
	t.Setenv("SMTP_HOST", "smtp.example.com")
	t.Setenv("SMTP_PORT", "587")

	svc, err := NewEmailService()
	require.NoError(t, err)
	assert.NotNil(t, svc)
}
