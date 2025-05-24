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
	"myanimeapi/api/auth"
	"myanimeapi/api/models"
	"myanimeapi/api/repositories"
	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"
	"net/http"
	"time"
)

// UserServiceInterface defines the interface for user service operations
type UserServiceInterface interface {
	// GetAllUsers retrieves all users with pagination
	GetAllUsers(ctx context.Context, page, limit int) ([]models.User, int64, error)
	// CreateUser creates a new user
	CreateUser(ctx context.Context, user *models.User) error
	// DeleteUser deletes a user by ID
	DeleteUser(ctx context.Context, id uint) error
	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uint) (*models.User, error)
	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	// GetByUsername retrieves a user by username
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	// ValidateUser validates a user's credentials
	ValidateUser(ctx context.Context, username, password string) (*models.User, error)
	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, user *models.User) error
	// ChangePassword changes a user's password
	ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword string) error
	// DeactivateAccount deactivates a user's account
	DeactivateAccount(ctx context.Context, userID uint, password string) error
}

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
	logger   *logger.Logger
}

// NewUserService creates a new UserService instance.
// It takes a UserRepository implementation as a dependency,
// following the dependency injection pattern.
func NewUserService(userRepo repositories.UserRepository) UserServiceInterface {
	return &UserService{
		userRepo: userRepo,
		logger:   logger.New(),
	}
}

// GetAllUsers retrieves all users with pagination
func (s *UserService) GetAllUsers(ctx context.Context, page, limit int) ([]models.User, int64, error) {
	s.logger.WithFields(map[string]interface{}{
		"page":  page,
		"limit": limit,
	}).Info("Retrieving all users")

	usersInterface, total, err := s.userRepo.GetAll(ctx, page, limit)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"page":  page,
			"limit": limit,
			"error": err.Error(),
		}).Error("Failed to retrieve users")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve users", err.Error(), http.StatusInternalServerError, map[string]interface{}{
			"page":  page,
			"limit": limit,
		}, err)
	}

	users := make([]models.User, len(usersInterface))
	for i, userInterface := range usersInterface {
		user, ok := userInterface.(*models.User)
		if !ok {
			s.logger.WithFields(map[string]interface{}{
				"index": i,
			}).Error("Invalid user type in slice")
			return nil, 0, errors.NewError(errors.ErrInternalServer, "Invalid user type in repository response", fmt.Sprintf("Type assertion failed at index %d", i), http.StatusInternalServerError, map[string]interface{}{
				"index": i,
			}, nil)
		}
		users[i] = *user
	}

	s.logger.WithFields(map[string]interface{}{
		"count": len(users),
		"total": total,
	}).Info("Successfully retrieved users")
	return users, total, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	s.logger.WithFields(map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
	}).Info("Creating new user")

	existingUser, err := s.userRepo.GetByUsername(ctx, user.Username)
	if err == nil && existingUser != nil {
		s.logger.WithField("username", user.Username).Warning("Username already exists")
		return errors.NewError(errors.ErrConflict, "Username already exists", fmt.Sprintf("Username '%s' is already taken", user.Username), http.StatusConflict, map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
		}, nil)
	}

	existingUser, err = s.userRepo.GetByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		s.logger.WithField("email", user.Email).Warning("Email already exists")
		return errors.NewError(errors.ErrConflict, "Email already exists", fmt.Sprintf("Email '%s' is already registered", user.Email), http.StatusConflict, map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
		}, nil)
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"username": user.Username,
			"error":    err.Error(),
		}).Error("Failed to hash password")
		return errors.NewError(errors.ErrInternalServer, "Failed to hash password", err.Error(), http.StatusInternalServerError, map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
		}, err)
	}
	user.Password = hashedPassword

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"username": user.Username,
			"error":    err.Error(),
		}).Error("Failed to create user")
		return errors.NewError(errors.ErrInternalServer, "Failed to create user", err.Error(), http.StatusInternalServerError, map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
		}, err)
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
	}).Info("Successfully created user")
	return nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	s.logger.WithField("user_id", id).Info("Deleting user")

	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id": id,
			"error":   err.Error(),
		}).Error("Failed to retrieve user")
		return errors.NewError(errors.ErrResourceNotFound, "User not found", err.Error(), http.StatusNotFound, map[string]interface{}{
			"user_id": id,
		}, err)
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id": id,
			"error":   err.Error(),
		}).Error("Failed to delete user")
		return errors.NewError(errors.ErrInternalServer, "Failed to delete user", err.Error(), http.StatusInternalServerError, map[string]interface{}{
			"user_id": id,
		}, err)
	}

	s.logger.WithField("user_id", id).Info("Successfully deleted user")
	return nil
}

// GetByID retrieves a user by ID.
func (s *UserService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	s.logger.WithField("user_id", id).Info("Retrieving user by ID")

	userInterface, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id": id,
			"error":   err.Error(),
		}).Error("Failed to retrieve user")
		return nil, errors.NewError(errors.ErrResourceNotFound, "User not found", err.Error(), http.StatusNotFound, map[string]interface{}{
			"user_id": id,
		}, err)
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		s.logger.WithField("user_id", id).Error("Invalid user type returned from repository")
		return nil, errors.NewError(errors.ErrInternalServer, "Invalid user type returned from repository", "Type assertion failed", http.StatusInternalServerError, map[string]interface{}{
			"user_id": id,
		}, nil)
	}

	s.logger.WithField("user_id", id).Info("Successfully retrieved user")
	return user, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	s.logger.WithField("email", email).Info("Retrieving user by email")

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"email": email,
			"error": err.Error(),
		}).Error("Failed to retrieve user")
		return nil, errors.NewError(errors.ErrResourceNotFound, "User not found", fmt.Sprintf("No user found with email '%s'", email), http.StatusNotFound, map[string]interface{}{
			"email": email,
		}, err)
	}

	s.logger.WithField("email", email).Info("Successfully retrieved user")
	return user, nil
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	s.logger.WithField("username", username).Info("Retrieving user by username")

	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"username": username,
			"error":    err.Error(),
		}).Error("Failed to retrieve user")
		return nil, errors.NewError(errors.ErrResourceNotFound, "User not found", fmt.Sprintf("No user found with username '%s'", username), http.StatusNotFound, map[string]interface{}{
			"username": username,
		}, err)
	}

	s.logger.WithField("username", username).Info("Successfully retrieved user")
	return user, nil
}

// ValidateUser validates a user's credentials.
func (s *UserService) ValidateUser(ctx context.Context, username, password string) (*models.User, error) {
	s.logger.WithField("username", username).Info("Validating user credentials")

	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"username": username,
			"error":    err.Error(),
		}).Error("Failed to retrieve user")
		return nil, errors.NewError(errors.ErrUnauthorized, "Invalid credentials", "Invalid username or password", http.StatusUnauthorized, map[string]interface{}{
			"username": username,
		}, err)
	}

	valid, err := auth.CheckPasswordHash(password, user.Password)
	if err != nil || !valid {
		s.logger.WithFields(map[string]interface{}{
			"username": username,
			"error":    err.Error(),
		}).Error("Invalid password")
		return nil, errors.NewError(errors.ErrUnauthorized, "Invalid password", "The provided password is incorrect", http.StatusUnauthorized, map[string]interface{}{
			"username": username,
		}, err)
	}

	s.logger.WithFields(map[string]interface{}{
		"username": username,
		"user_id":  user.ID,
	}).Info("Successfully validated user credentials")
	return user, nil
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	s.logger.WithFields(map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
	}).Info("Updating user")

	existingUser, err := s.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id": user.ID,
			"error":   err.Error(),
		}).Error("Failed to retrieve user")
		return errors.NewError(errors.ErrResourceNotFound, "User not found", err.Error(), http.StatusNotFound, map[string]interface{}{
			"user_id": user.ID,
		}, err)
	}

	// Check if username is being changed and if it's already taken
	if user.Username != existingUser.(*models.User).Username {
		existingUser, err = s.userRepo.GetByUsername(ctx, user.Username)
		if err == nil && existingUser != nil {
			s.logger.WithField("username", user.Username).Warning("Username already exists")
			return errors.NewError(errors.ErrConflict, "Username already exists", fmt.Sprintf("Username '%s' is already taken", user.Username), http.StatusConflict, map[string]interface{}{
				"username": user.Username,
				"user_id":  user.ID,
			}, nil)
		}
	}

	// Check if email is being changed and if it's already taken
	if user.Email != existingUser.(*models.User).Email {
		existingUser, err = s.userRepo.GetByEmail(ctx, user.Email)
		if err == nil && existingUser != nil {
			s.logger.WithField("email", user.Email).Warning("Email already exists")
			return errors.NewError(errors.ErrConflict, "Email already exists", fmt.Sprintf("Email '%s' is already registered", user.Email), http.StatusConflict, map[string]interface{}{
				"email":   user.Email,
				"user_id": user.ID,
			}, nil)
		}
	}

	// Update timestamps
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id": user.ID,
			"error":   err.Error(),
		}).Error("Failed to update user")
		return errors.NewError(errors.ErrInternalServer, "Failed to update user", err.Error(), http.StatusInternalServerError, map[string]interface{}{
			"user_id": user.ID,
		}, err)
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
	}).Info("Successfully updated user")
	return nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword string) error {
	s.logger.WithField("user_id", userID).Info("Changing user password")

	userInterface, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to retrieve user")
		return errors.NewError(errors.ErrResourceNotFound, "User not found", err.Error(), http.StatusNotFound, map[string]interface{}{
			"user_id": userID,
		}, err)
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		s.logger.WithField("user_id", userID).Error("Invalid user type returned from repository")
		return errors.NewError(errors.ErrInternalServer, "Invalid user type", "Type assertion failed", http.StatusInternalServerError, map[string]interface{}{
			"user_id": userID,
		}, nil)
	}

	valid, err := auth.CheckPasswordHash(currentPassword, user.Password)
	if err != nil || !valid {
		s.logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Invalid current password")
		return errors.NewError(errors.ErrUnauthorized, "Invalid current password", "The provided current password is incorrect", http.StatusUnauthorized, map[string]interface{}{
			"user_id": userID,
		}, err)
	}

	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to hash new password")
		return errors.NewError(errors.ErrInternalServer, "Failed to hash new password", err.Error(), http.StatusInternalServerError, map[string]interface{}{
			"user_id": userID,
		}, err)
	}
	user.Password = hashedPassword

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to update user password")
		return errors.NewError(errors.ErrInternalServer, "Failed to update user password", err.Error(), http.StatusInternalServerError, map[string]interface{}{
			"user_id": userID,
		}, err)
	}

	s.logger.WithField("user_id", userID).Info("Successfully updated user password")
	return nil
}

// DeactivateAccount deactivates a user's account
func (s *UserService) DeactivateAccount(ctx context.Context, userID uint, password string) error {
	s.logger.WithField("user_id", userID).Info("Deactivating user account")

	userInterface, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to retrieve user")
		return errors.NewError(errors.ErrResourceNotFound, "User not found", err.Error(), http.StatusNotFound, map[string]interface{}{
			"user_id": userID,
		}, err)
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		s.logger.WithField("user_id", userID).Error("Invalid user type returned from repository")
		return errors.NewError(errors.ErrInternalServer, "Invalid user type", "Type assertion failed", http.StatusInternalServerError, map[string]interface{}{
			"user_id": userID,
		}, nil)
	}

	valid, err := auth.CheckPasswordHash(password, user.Password)
	if err != nil || !valid {
		s.logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Invalid password")
		return errors.NewError(errors.ErrUnauthorized, "Invalid password", "The provided password is incorrect", http.StatusUnauthorized, map[string]interface{}{
			"user_id": userID,
		}, err)
	}

	user.IsActive = false
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to deactivate user account")
		return errors.NewError(errors.ErrInternalServer, "Failed to deactivate user account", err.Error(), http.StatusInternalServerError, map[string]interface{}{
			"user_id": userID,
		}, err)
	}

	s.logger.WithField("user_id", userID).Info("Successfully deactivated user account")
	return nil
}
