package services

import (
	"context"
	"fmt"
	"log"
	"myanimeapi/api/models"
	"myanimeapi/api/repositories"
	"myanimeapi/internal/errors"
	"net/http"
	"time"
)

// FavoriteService handles business logic for favorite operations
type FavoriteService struct {
	favoriteRepo repositories.FavoriteRepository
	userRepo     repositories.UserRepository
	animeRepo    repositories.AnimeRepository
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
	}
}

// AddFavorite adds an anime to a user's favorites
func (s *FavoriteService) AddFavorite(ctx context.Context, userID uint, animeID uint) (*models.Favorite, error) {
	log.Printf("AddFavorite: Starting to add favorite for user %d and anime %d", userID, animeID)

	// Get user from database
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		log.Printf("AddFavorite: Failed to get user from database: %v", err)
		return nil, errors.NewError(errors.ErrUnauthorized, "User not found", err.Error(), http.StatusUnauthorized)
	}
	userModel, ok := user.(*models.User)
	if !ok {
		log.Printf("AddFavorite: Invalid user type")
		return nil, errors.NewError(errors.ErrInternalServer, "Invalid user type", "Type assertion failed for user model", http.StatusInternalServerError)
	}
	log.Printf("AddFavorite: Found user in database: %s", userModel.Username)

	// Check if anime exists
	anime, err := s.animeRepo.GetByID(ctx, animeID)
	if err != nil {
		log.Printf("AddFavorite: Failed to find anime: %v", err)
		return nil, errors.NewError(errors.ErrResourceNotFound, "Anime not found", err.Error(), http.StatusNotFound)
	}
	animeModel, ok := anime.(*models.Anime)
	if !ok {
		log.Printf("AddFavorite: Invalid anime type")
		return nil, errors.NewError(errors.ErrInternalServer, "Invalid anime type", "Type assertion failed for anime model", http.StatusInternalServerError)
	}
	log.Printf("AddFavorite: Found anime in database: %s", animeModel.Title)

	// Check if favorite already exists
	_, err = s.favoriteRepo.GetByUserAndAnime(ctx, userID, animeID)
	if err == nil {
		log.Printf("AddFavorite: Anime already in favorites for user %d", userID)
		return nil, errors.NewError(errors.ErrConflict, "Anime already in favorites", fmt.Sprintf("User %d already has anime %d in favorites", userID, animeID), http.StatusConflict)
	}
	log.Printf("AddFavorite: Anime not already in favorites, proceeding to create")

	// Create new favorite
	favorite := &models.Favorite{
		UserID:    userID,
		AnimeID:   animeID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.favoriteRepo.Create(ctx, favorite); err != nil {
		log.Printf("AddFavorite: Failed to create favorite: %v", err)
		return nil, errors.NewError(errors.ErrInternalServer, "Failed to add favorite", err.Error(), http.StatusInternalServerError)
	}
	log.Printf("AddFavorite: Successfully created favorite for user %d and anime %d", userID, animeID)

	// Get the created favorite with preloaded data
	result, err := s.favoriteRepo.GetByID(ctx, favorite.ID)
	if err != nil {
		log.Printf("AddFavorite: Failed to get created favorite: %v", err)
		return nil, errors.NewError(errors.ErrInternalServer, "Failed to load favorite details", err.Error(), http.StatusInternalServerError)
	}
	createdFavorite, ok := result.(*models.Favorite)
	if !ok {
		log.Printf("AddFavorite: Invalid favorite type")
		return nil, errors.NewError(errors.ErrInternalServer, "Invalid favorite type", "Type assertion failed for favorite model", http.StatusInternalServerError)
	}
	log.Printf("AddFavorite: Successfully retrieved created favorite with related data")

	return createdFavorite, nil
}

// RemoveFavorite removes an anime from a user's favorites
func (s *FavoriteService) RemoveFavorite(ctx context.Context, userID uint, animeID uint) error {
	log.Printf("RemoveFavorite: Starting to remove favorite for user %d and anime %d", userID, animeID)

	// Check if favorite exists
	_, err := s.favoriteRepo.GetByUserAndAnime(ctx, userID, animeID)
	if err != nil {
		log.Printf("RemoveFavorite: Favorite not found: %v", err)
		return errors.NewError(errors.ErrResourceNotFound, "Favorite not found", err.Error(), http.StatusNotFound)
	}

	// Delete favorite
	if err := s.favoriteRepo.DeleteByUserAndAnime(ctx, userID, animeID); err != nil {
		log.Printf("RemoveFavorite: Failed to delete favorite: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to remove favorite", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("RemoveFavorite: Successfully removed favorite for user %d and anime %d", userID, animeID)
	return nil
}

// GetFavorites retrieves all favorites for a user
func (s *FavoriteService) GetFavorites(ctx context.Context, userID uint) ([]models.Favorite, error) {
	log.Printf("GetFavorites: Retrieving favorites for user %d", userID)

	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		log.Printf("GetFavorites: Failed to find user: %v", err)
		return nil, errors.NewError(errors.ErrResourceNotFound, "User not found", err.Error(), http.StatusNotFound)
	}

	// Get favorites
	favorites, err := s.favoriteRepo.GetByUserID(ctx, userID)
	if err != nil {
		log.Printf("GetFavorites: Failed to retrieve favorites: %v", err)
		return nil, errors.NewError(errors.ErrInternalServer, "Failed to retrieve favorites", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("GetFavorites: Successfully retrieved %d favorites for user %d", len(favorites), userID)
	return favorites, nil
}
