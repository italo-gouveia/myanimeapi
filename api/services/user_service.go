// Package services implements the business logic layer of the MyAnimeAPI application.
// It provides service implementations that handle the core business operations,
// working with repositories for data access and implementing business rules.
//
// The package includes:
//   - User management (registration, authentication, profile updates)
//   - Anime management (CRUD operations, favorites)
//   - Review management
//   - Genre and tag management
//
// Each service is responsible for:
//   - Implementing business rules and validations
//   - Coordinating between different repositories
//   - Handling errors and providing meaningful error messages
//   - Logging operations for debugging and monitoring
package services

import (
	"context"
	"fmt"
	"log"
	"myanimeapi/api/auth"
	"myanimeapi/api/models"
	"myanimeapi/api/repositories"
	"myanimeapi/internal/errors"
	"net/http"
	"time"
)

// UserService handles business logic for user operations.
// It provides methods for user management including:
//   - User registration and authentication
//   - Profile management (view, update, delete)
//   - Password management
//   - Account deactivation
//
// The service ensures:
//   - Data validation and business rule enforcement
//   - Proper error handling and logging
//   - Secure password handling
//   - Unique username and email constraints
type UserService struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new UserService instance.
// It takes a UserRepository implementation as a dependency,
// following the dependency injection pattern.
func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	log.Printf("UserService.GetUserByID: Retrieving user with ID %d", id)

	userInterface, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		log.Printf("UserService.GetUserByID: Failed to retrieve user: %v", err)
		return nil, errors.NewError(errors.ErrInternalServer, "Failed to retrieve user", err.Error(), http.StatusInternalServerError)
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		log.Printf("UserService.GetUserByID: Invalid user type returned from repository")
		return nil, errors.NewError(errors.ErrInternalServer, "Invalid user type returned from repository", "Type assertion failed for user model", http.StatusInternalServerError)
	}

	log.Printf("UserService.GetUserByID: Successfully retrieved user with ID %d", id)
	return user, nil
}

// GetAllUsers retrieves all users with pagination
func (s *UserService) GetAllUsers(ctx context.Context, page, limit int) ([]models.User, int64, error) {
	log.Printf("UserService.GetAllUsers: Retrieving all users with page %d and limit %d", page, limit)

	usersInterface, total, err := s.userRepo.GetAll(ctx, page, limit)
	if err != nil {
		log.Printf("UserService.GetAllUsers: Failed to retrieve users: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve users", err.Error(), http.StatusInternalServerError)
	}

	// Convert interface slice to User slice
	users := make([]models.User, len(usersInterface))
	for i, userInterface := range usersInterface {
		user, ok := userInterface.(models.User)
		if !ok {
			log.Printf("UserService.GetAllUsers: Invalid user type in slice at index %d", i)
			return nil, 0, errors.NewError(errors.ErrInternalServer, "Invalid user type in repository response", fmt.Sprintf("Type assertion failed at index %d", i), http.StatusInternalServerError)
		}
		users[i] = user
	}

	log.Printf("UserService.GetAllUsers: Successfully retrieved %d users", len(users))
	return users, total, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	log.Printf("UserService.CreateUser: Creating new user with username %s", user.Username)

	// Check if username already exists
	existingUser, err := s.userRepo.GetByUsername(ctx, user.Username)
	if err == nil && existingUser != nil {
		log.Printf("UserService.CreateUser: Username %s already exists", user.Username)
		return errors.NewError(errors.ErrConflict, "Username already exists", fmt.Sprintf("Username '%s' is already taken", user.Username), http.StatusConflict)
	}

	// Check if email already exists
	existingUser, err = s.userRepo.GetByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		log.Printf("UserService.CreateUser: Email %s already exists", user.Email)
		return errors.NewError(errors.ErrConflict, "Email already exists", fmt.Sprintf("Email '%s' is already registered", user.Email), http.StatusConflict)
	}

	// Set timestamps
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Create user
	if err := s.userRepo.Create(ctx, user); err != nil {
		log.Printf("UserService.CreateUser: Failed to create user: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to create user", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("UserService.CreateUser: Successfully created user with ID %d", user.ID)
	return nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	log.Printf("UserService.UpdateUser: Updating user with ID %d", user.ID)

	// Check if user exists
	existingUser, err := s.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		log.Printf("UserService.UpdateUser: Failed to retrieve user: %v", err)
		return errors.NewError(errors.ErrResourceNotFound, "User not found", err.Error(), http.StatusNotFound)
	}

	// Check if username is being changed and if it already exists
	if user.Username != existingUser.(*models.User).Username {
		usernameUser, err := s.userRepo.GetByUsername(ctx, user.Username)
		if err == nil && usernameUser != nil && usernameUser.ID != user.ID {
			log.Printf("UserService.UpdateUser: Username %s already exists", user.Username)
			return errors.NewError(errors.ErrConflict, "Username already exists", fmt.Sprintf("Username '%s' is already taken by another user", user.Username), http.StatusConflict)
		}
	}

	// Check if email is being changed and if it already exists
	if user.Email != existingUser.(*models.User).Email {
		emailUser, err := s.userRepo.GetByEmail(ctx, user.Email)
		if err == nil && emailUser != nil && emailUser.ID != user.ID {
			log.Printf("UserService.UpdateUser: Email %s already exists", user.Email)
			return errors.NewError(errors.ErrConflict, "Email already exists", fmt.Sprintf("Email '%s' is already registered by another user", user.Email), http.StatusConflict)
		}
	}

	// Set updated timestamp
	user.UpdatedAt = time.Now()

	// Update user
	if err := s.userRepo.Update(ctx, user); err != nil {
		log.Printf("UserService.UpdateUser: Failed to update user: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to update user", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("UserService.UpdateUser: Successfully updated user with ID %d", user.ID)
	return nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	log.Printf("UserService.DeleteUser: Deleting user with ID %d", id)

	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		log.Printf("UserService.DeleteUser: Failed to retrieve user: %v", err)
		return errors.NewError(errors.ErrResourceNotFound, "User not found", err.Error(), http.StatusNotFound)
	}

	// Delete user
	if err := s.userRepo.Delete(ctx, id); err != nil {
		log.Printf("UserService.DeleteUser: Failed to delete user: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to delete user", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("UserService.DeleteUser: Successfully deleted user with ID %d", id)
	return nil
}

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	log.Printf("UserService.GetUserByUsername: Retrieving user with username %s", username)

	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		log.Printf("UserService.GetUserByUsername: Failed to retrieve user: %v", err)
		return nil, errors.NewError(errors.ErrResourceNotFound, "User not found", fmt.Sprintf("No user found with username '%s'", username), http.StatusNotFound)
	}

	log.Printf("UserService.GetUserByUsername: Successfully retrieved user with username %s", username)
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	log.Printf("UserService.GetUserByEmail: Retrieving user with email %s", email)

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		log.Printf("UserService.GetUserByEmail: Failed to retrieve user: %v", err)
		return nil, errors.NewError(errors.ErrResourceNotFound, "User not found", fmt.Sprintf("No user found with email '%s'", email), http.StatusNotFound)
	}

	log.Printf("UserService.GetUserByEmail: Successfully retrieved user with email %s", email)
	return user, nil
}

// Update updates a user in the database.
func (s *UserService) Update(ctx context.Context, user *models.User) error {
	return s.userRepo.Update(ctx, user)
}

// GetByID retrieves a user by ID.
func (s *UserService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	userInterface, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user, ok := userInterface.(*models.User)
	if !ok {
		return nil, fmt.Errorf("invalid user type returned from repository")
	}
	return user, nil
}

// GetByEmail retrieves a user by email.
func (s *UserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// GetByUsername retrieves a user by username.
func (s *UserService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.userRepo.GetByUsername(ctx, username)
}

// ValidateUser validates a user's credentials.
func (s *UserService) ValidateUser(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.NewError(errors.ErrUnauthorized, "Account is deactivated", "This account has been deactivated", http.StatusUnauthorized)
	}

	valid, err := auth.CheckPasswordHash(password, user.Password)
	if err != nil || !valid {
		return nil, errors.NewError(errors.ErrUnauthorized, "Invalid password", "The provided password is incorrect", http.StatusUnauthorized)
	}

	return user, nil
}
