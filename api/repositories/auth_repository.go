package repositories

import (
	"context"

	"myanimeapi/api/models"
	"myanimeapi/internal/db"
)

// AuthRepository defines the interface for authentication-related database operations.
type AuthRepository interface {
	// GetUserByUsername retrieves a user by their username.
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	// GetUserByEmail retrieves a user by their email.
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	// CreateUser creates a new user in the database.
	CreateUser(ctx context.Context, user *models.User) error
	// UpdateUser updates an existing user in the database.
	UpdateUser(ctx context.Context, user *models.User) error
}

// authRepository implements the AuthRepository interface.
type authRepository struct {
	db db.DBInterface
}

// NewAuthRepository creates a new instance of AuthRepository.
func NewAuthRepository(db db.DBInterface) AuthRepository {
	return &authRepository{db: db}
}

// GetUserByUsername retrieves a user by their username.
func (r *authRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email.
func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// CreateUser creates a new user in the database.
func (r *authRepository) CreateUser(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Create(user)
	return result.Error
}

// UpdateUser updates an existing user in the database.
func (r *authRepository) UpdateUser(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Save(user)
	return result.Error
}
