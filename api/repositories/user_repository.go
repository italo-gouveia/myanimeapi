package repositories

import (
	"context"
	"log"
	"myanimeapi/api/models"
	"myanimeapi/internal/db"
	"myanimeapi/internal/errors"
	"net/http"
)

// UserRepositoryImpl implements the UserRepository interface
type UserRepositoryImpl struct {
	db db.DBInterface
}

// NewUserRepository creates a new UserRepositoryImpl instance
func NewUserRepository(db db.DBInterface) UserRepository {
	return &UserRepositoryImpl{db: db}
}

// GetByID retrieves a user by its ID
func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uint) (interface{}, error) {
	log.Printf("UserRepository.GetByID: Retrieving user with ID %d", id)

	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		log.Printf("UserRepository.GetByID: Failed to retrieve user: %v", err)
		return nil, errors.NewError(errors.ErrResourceNotFound, "User not found", err.Error(), http.StatusNotFound)
	}

	log.Printf("UserRepository.GetByID: Successfully retrieved user with ID %d", id)
	return &user, nil
}

// GetAll retrieves all users with optional pagination
func (r *UserRepositoryImpl) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	log.Printf("UserRepository.GetAll: Retrieving all users with page %d and limit %d", page, limit)

	var users []models.User
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		log.Printf("UserRepository.GetAll: Failed to count users: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count users", err.Error(), http.StatusInternalServerError)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve users with pagination
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		log.Printf("UserRepository.GetAll: Failed to retrieve users: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve users", err.Error(), http.StatusInternalServerError)
	}

	// Convert to interface slice
	result := make([]interface{}, len(users))
	for i, user := range users {
		result[i] = &user
	}

	log.Printf("UserRepository.GetAll: Successfully retrieved %d users", len(users))
	return result, total, nil
}

// Create creates a new user
func (r *UserRepositoryImpl) Create(ctx context.Context, entity interface{}) error {
	log.Printf("UserRepository.Create: Creating new user")

	user, ok := entity.(*models.User)
	if !ok {
		log.Printf("UserRepository.Create: Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.User", http.StatusBadRequest)
	}

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		log.Printf("UserRepository.Create: Failed to create user: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to create user", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("UserRepository.Create: Successfully created user with ID %d", user.ID)
	return nil
}

// Update updates an existing user
func (r *UserRepositoryImpl) Update(ctx context.Context, entity interface{}) error {
	log.Printf("UserRepository.Update: Updating user")

	user, ok := entity.(*models.User)
	if !ok {
		log.Printf("UserRepository.Update: Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.User", http.StatusBadRequest)
	}

	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		log.Printf("UserRepository.Update: Failed to update user: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to update user", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("UserRepository.Update: Successfully updated user with ID %d", user.ID)
	return nil
}

// Delete deletes a user by its ID
func (r *UserRepositoryImpl) Delete(ctx context.Context, id uint) error {
	log.Printf("UserRepository.Delete: Deleting user with ID %d", id)

	if err := r.db.WithContext(ctx).Delete(&models.User{}, id).Error; err != nil {
		log.Printf("UserRepository.Delete: Failed to delete user: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to delete user", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("UserRepository.Delete: Successfully deleted user with ID %d", id)
	return nil
}

// GetByUsername retrieves a user by username
func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	log.Printf("UserRepository.GetByUsername: Retrieving user with username %s", username)

	var user models.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		log.Printf("UserRepository.GetByUsername: Failed to retrieve user: %v", err)
		return nil, errors.NewError(errors.ErrResourceNotFound, "User not found", err.Error(), http.StatusNotFound)
	}

	log.Printf("UserRepository.GetByUsername: Successfully retrieved user with username %s", username)
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	log.Printf("UserRepository.GetByEmail: Retrieving user with email %s", email)

	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		log.Printf("UserRepository.GetByEmail: Failed to retrieve user: %v", err)
		return nil, errors.NewError(errors.ErrResourceNotFound, "User not found", err.Error(), http.StatusNotFound)
	}

	log.Printf("UserRepository.GetByEmail: Successfully retrieved user with email %s", email)
	return &user, nil
}
