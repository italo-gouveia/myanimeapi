package mocks

import (
	"testing"
)

// MockManager centralizes all mocks used in tests
type MockManager struct {
	T *testing.T

	// Repository mocks - will be implemented when necessary
	// UserRepo     *MockUserRepository
	// AnimeRepo    *MockAnimeRepository
	// ReviewRepo   *MockReviewRepository
	// GenreRepo    *MockGenreRepository
	// TagRepo      *MockTagRepository
	// FavoriteRepo *MockFavoriteRepository

	// Service mocks - will be implemented when necessary
	// AuthService     *MockAuthService
	// EmailService    *MockEmailService
	// StorageService  *MockStorageService
	// PasswordService *MockPasswordResetService

	// Utility mocks - will be implemented when necessary
	// Logger *MockLogger
}

// NewMockManager creates a new mock manager with all mocks
func NewMockManager(t *testing.T) *MockManager {
	return &MockManager{
		T: t,
		// Mocks will be added as necessary
	}
}

// SetupDefaultExpectations configures common mock behaviors
func (m *MockManager) SetupDefaultExpectations() {
	// Setup common mock behaviors that most tests will need
	// Will be implemented when the mocks are available
}

// VerifyAllMocks ensures all expected mock calls were made
func (m *MockManager) VerifyAllMocks() {
	// This will be called automatically by gomock.Controller.Finish()
	// but we can add custom verification logic here if needed
}
