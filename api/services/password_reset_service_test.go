package services

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResetTokenExpiry_Default(t *testing.T) {
	os.Unsetenv("RESET_TOKEN_EXPIRY_MINUTES")
	d := resetTokenExpiry()
	assert.Equal(t, time.Hour, d)
}

func TestResetTokenExpiry_FromEnv(t *testing.T) {
	t.Setenv("RESET_TOKEN_EXPIRY_MINUTES", "30")
	d := resetTokenExpiry()
	assert.Equal(t, 30*time.Minute, d)
}

func TestResetTokenExpiry_InvalidEnv(t *testing.T) {
	t.Setenv("RESET_TOKEN_EXPIRY_MINUTES", "notanumber")
	d := resetTokenExpiry()
	assert.Equal(t, time.Hour, d)
}

func TestResetTokenExpiry_ZeroEnv(t *testing.T) {
	t.Setenv("RESET_TOKEN_EXPIRY_MINUTES", "0")
	d := resetTokenExpiry()
	assert.Equal(t, time.Hour, d)
}

func TestNewPasswordResetService(t *testing.T) {
	userRepo := &mockUserRepo{}
	svc := NewPasswordResetService(userRepo, nil, nil)
	require.NotNil(t, svc)
	assert.Equal(t, time.Hour, svc.tokenExpiry)
}
