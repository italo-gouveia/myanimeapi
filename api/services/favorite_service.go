package services

import (
	"context"
	"fmt"
	"myanimeapi/api/models"
	"myanimeapi/api/repositories"
	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"
	"net/http"
	"time"
)

// FavoriteServiceInterface defines the interface for favorite operations
type FavoriteServiceInterface interface {
	// AddFavorite adds an anime to a user's favorites
	AddFavorite(ctx context.Context, userID uint, animeID uint) (*models.Favorite, error)
	// RemoveFavorite removes an anime from a user's favorites
	RemoveFavorite(ctx context.Context, userID uint, animeID uint) error
	// GetFavorites retrieves all favorites for a user
	GetFavorites(ctx context.Context, userID uint) ([]models.Favorite, error)
}

// FavoriteService handles business logic for favorite operations
// It implements the FavoriteServiceInterface.
type FavoriteService struct {
	favoriteRepo repositories.FavoriteRepository
	userRepo     repositories.UserRepository
	animeRepo    repositories.AnimeRepository
	logger       *logger.Logger
}

// NewFavoriteService creates a new FavoriteService instance
func NewFavoriteService(
	favoriteRepo repositories.FavoriteRepository,
	userRepo repositories.UserRepository,
	animeRepo repositories.AnimeRepository,
) *FavoriteService {
	return &FavoriteService{
		favoriteRepo: favoriteRepo,
		userRepo:     userRepo,
		animeRepo:    animeRepo,
		logger:       logger.New(),
	}
}

// AddFavorite adds an anime to a user's favorites
func (s *FavoriteService) AddFavorite(ctx context.Context, userID uint, animeID uint) (*models.Favorite, error) {
	s.logger.WithFields(map[string]interface{}{
		"user_id":  userID,
		"anime_id": animeID,
	}).Info("Starting to add favorite")

	// Get user from database
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to get user from database")
		return nil, errors.NewError(errors.ErrUnauthorized, "User not found", err.Error(), http.StatusUnauthorized,
			map[string]interface{}{
				"user_id": userID,
			},
			err)
	}
	userModel, ok := user.(*models.User)
	if !ok {
		s.logger.WithField("user_id", userID).Error("Invalid user type")
		return nil, errors.NewError(errors.ErrInternalServer, "Invalid user type", "Type assertion failed for user model", http.StatusInternalServerError,
			map[string]interface{}{
				"user_id": userID,
			},
			nil)
	}
	s.logger.WithField("username", userModel.Username).Info("Found user in database")

	// Check if anime exists
	anime, err := s.animeRepo.GetByID(ctx, animeID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"anime_id": animeID,
			"error":    err.Error(),
		}).Error("Failed to find anime")
		return nil, errors.NewError(errors.ErrResourceNotFound, "Anime not found", err.Error(), http.StatusNotFound,
			map[string]interface{}{
				"anime_id": animeID,
			},
			err)
	}
	animeModel, ok := anime.(*models.Anime)
	if !ok {
		s.logger.WithField("anime_id", animeID).Error("Invalid anime type")
		return nil, errors.NewError(errors.ErrInternalServer, "Invalid anime type", "Type assertion failed for anime model", http.StatusInternalServerError,
			map[string]interface{}{
				"anime_id": animeID,
			},
			nil)
	}
	s.logger.WithField("title", animeModel.Title).Info("Found anime in database")

	// Check if favorite already exists
	_, err = s.favoriteRepo.GetByUserAndAnime(ctx, userID, animeID)
	if err == nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id":  userID,
			"anime_id": animeID,
		}).Warning("Anime already in favorites for user")
		return nil, errors.NewError(errors.ErrConflict, "Anime already in favorites", fmt.Sprintf("User %d already has anime %d in favorites", userID, animeID), http.StatusConflict,
			map[string]interface{}{
				"user_id":  userID,
				"anime_id": animeID,
			},
			nil)
	}
	s.logger.Info("Anime not already in favorites, proceeding to create")

	// Create new favorite
	favorite := &models.Favorite{
		UserID:    userID,
		AnimeID:   animeID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.favoriteRepo.Create(ctx, favorite); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id":  userID,
			"anime_id": animeID,
			"error":    err.Error(),
		}).Error("Failed to create favorite")
		return nil, errors.NewError(errors.ErrInternalServer, "Failed to add favorite", err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"user_id":  userID,
				"anime_id": animeID,
			},
			err)
	}
	s.logger.WithFields(map[string]interface{}{
		"user_id":  userID,
		"anime_id": animeID,
	}).Info("Successfully created favorite")

	// Get the created favorite with preloaded data
	result, err := s.favoriteRepo.GetByID(ctx, favorite.ID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"favorite_id": favorite.ID,
			"error":       err.Error(),
		}).Error("Failed to get created favorite")
		return nil, errors.NewError(errors.ErrInternalServer, "Failed to load favorite details", err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"favorite_id": favorite.ID,
				"user_id":     userID,
				"anime_id":    animeID,
			},
			err)
	}
	createdFavorite, ok := result.(*models.Favorite)
	if !ok {
		s.logger.WithField("favorite_id", favorite.ID).Error("Invalid favorite type")
		return nil, errors.NewError(errors.ErrInternalServer, "Invalid favorite type", "Type assertion failed for favorite model", http.StatusInternalServerError,
			map[string]interface{}{
				"favorite_id": favorite.ID,
				"user_id":     userID,
				"anime_id":    animeID,
			},
			nil)
	}
	s.logger.Info("Successfully retrieved created favorite with related data")

	return createdFavorite, nil
}

// RemoveFavorite removes an anime from a user's favorites
func (s *FavoriteService) RemoveFavorite(ctx context.Context, userID uint, animeID uint) error {
	s.logger.WithFields(map[string]interface{}{
		"user_id":  userID,
		"anime_id": animeID,
	}).Info("Starting to remove favorite")

	// Check if favorite exists
	_, err := s.favoriteRepo.GetByUserAndAnime(ctx, userID, animeID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id":  userID,
			"anime_id": animeID,
			"error":    err.Error(),
		}).Error("Favorite not found")
		return errors.NewError(errors.ErrResourceNotFound, "Favorite not found", err.Error(), http.StatusNotFound,
			map[string]interface{}{
				"user_id":  userID,
				"anime_id": animeID,
			},
			err)
	}

	// Delete favorite
	if err := s.favoriteRepo.DeleteByUserAndAnime(ctx, userID, animeID); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id":  userID,
			"anime_id": animeID,
			"error":    err.Error(),
		}).Error("Failed to delete favorite")
		return errors.NewError(errors.ErrInternalServer, "Failed to remove favorite", err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"user_id":  userID,
				"anime_id": animeID,
			},
			err)
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":  userID,
		"anime_id": animeID,
	}).Info("Successfully removed favorite")
	return nil
}

// GetFavorites retrieves all favorites for a user
func (s *FavoriteService) GetFavorites(ctx context.Context, userID uint) ([]models.Favorite, error) {
	s.logger.WithField("user_id", userID).Info("Retrieving favorites for user")

	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to find user")
		return nil, errors.NewError(errors.ErrResourceNotFound, "User not found", err.Error(), http.StatusNotFound,
			map[string]interface{}{
				"user_id": userID,
			},
			err)
	}

	// Get favorites
	favorites, err := s.favoriteRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to retrieve favorites")
		return nil, errors.NewError(errors.ErrInternalServer, "Failed to retrieve favorites", err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"user_id": userID,
			},
			err)
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id": userID,
		"count":   len(favorites),
	}).Info("Successfully retrieved favorites for user")
	return favorites, nil
}
