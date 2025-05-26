package repositories

import (
	"context"
	"fmt"
	"myanimeapi/api/models"
	"myanimeapi/internal/db"
	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"
	"net/http"
)

// UserRepositoryImpl implements the UserRepository interface
type UserRepositoryImpl struct {
	db     db.DBInterface
	logger *logger.Logger
}

// NewUserRepository creates a new UserRepositoryImpl instance
func NewUserRepository(db db.DBInterface) UserRepository {
	return &UserRepositoryImpl{
		db:     db,
		logger: logger.New(),
	}
}

// GetByID retrieves a user by its ID
func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uint) (interface{}, error) {
	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Retrieving user by ID")

	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to retrieve user")
		return nil, errors.NewError(errors.ErrResourceNotFound, "User not found", fmt.Sprintf("User with ID %d not found", id), http.StatusNotFound, map[string]interface{}{"id": id}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Successfully retrieved user")
	return &user, nil
}

// GetAll retrieves all users with optional pagination
func (r *UserRepositoryImpl) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	r.logger.WithFields(map[string]interface{}{
		"page":  page,
		"limit": limit,
	}).Info("Retrieving all users")

	var users []models.User
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to count users")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count users", err.Error(), http.StatusInternalServerError, nil, err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve users with pagination
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to retrieve users")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve users", err.Error(), http.StatusInternalServerError, nil, err)
	}

	// Convert to interface slice
	result := make([]interface{}, len(users))
	for i, user := range users {
		result[i] = &user
	}

	r.logger.WithFields(map[string]interface{}{
		"count": len(users),
		"total": total,
	}).Info("Successfully retrieved users")
	return result, total, nil
}

// Create creates a new user
func (r *UserRepositoryImpl) Create(ctx context.Context, entity interface{}) error {
	r.logger.Info("Creating new user")

	user, ok := entity.(*models.User)
	if !ok {
		r.logger.Error("Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.User", http.StatusBadRequest, nil, nil)
	}

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to create user")
		return errors.NewError(errors.ErrInternalServer, "Failed to create user", err.Error(), http.StatusInternalServerError, map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
		}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": user.ID,
	}).Info("Successfully created user")
	return nil
}

// Update updates an existing user
func (r *UserRepositoryImpl) Update(ctx context.Context, entity interface{}) error {
	r.logger.Info("Updating user")

	user, ok := entity.(*models.User)
	if !ok {
		r.logger.Error("Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.User", http.StatusBadRequest, nil, nil)
	}

	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to update user")
		return errors.NewError(errors.ErrInternalServer, "Failed to update user", err.Error(), http.StatusInternalServerError, map[string]interface{}{
			"id": user.ID,
		}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": user.ID,
	}).Info("Successfully updated user")
	return nil
}

// Delete deletes a user by its ID
func (r *UserRepositoryImpl) Delete(ctx context.Context, id uint) error {
	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Deleting user")

	if err := r.db.WithContext(ctx).Delete(&models.User{}, id).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to delete user")
		return errors.NewError(errors.ErrInternalServer, "Failed to delete user", err.Error(), http.StatusInternalServerError, map[string]interface{}{"id": id}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Successfully deleted user")
	return nil
}

// GetByUsername retrieves a user by username
func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	r.logger.WithFields(map[string]interface{}{
		"username": username,
	}).Info("Retrieving user by username")

	var user models.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"username": username,
			"error":    err.Error(),
		}).Error("Failed to retrieve user")
		return nil, errors.NewError(errors.ErrResourceNotFound, "User not found", fmt.Sprintf("User with username '%s' not found", username), http.StatusNotFound, map[string]interface{}{"username": username}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"username": username,
	}).Info("Successfully retrieved user")
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	r.logger.WithFields(map[string]interface{}{
		"email": email,
	}).Info("Retrieving user by email")

	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"email": email,
			"error": err.Error(),
		}).Error("Failed to retrieve user")
		return nil, errors.NewError(errors.ErrResourceNotFound, "User not found", fmt.Sprintf("User with email '%s' not found", email), http.StatusNotFound, map[string]interface{}{"email": email}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"email": email,
	}).Info("Successfully retrieved user")
	return &user, nil
}
