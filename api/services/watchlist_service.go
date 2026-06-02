package services

import (
	"context"
	"myanimeapi/api/models"
	"myanimeapi/api/repositories"
	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"
	"net/http"
)

// WatchlistServiceInterface defines the interface for watchlist operations
type WatchlistServiceInterface interface {
	// UpsertWatchlistEntry adds or updates an anime in a user's watchlist
	UpsertWatchlistEntry(ctx context.Context, userID, animeID uint, status models.WatchlistStatus) (*models.WatchlistEntry, error)
	// RemoveWatchlistEntry removes an anime from a user's watchlist
	RemoveWatchlistEntry(ctx context.Context, userID, animeID uint) error
	// GetWatchlist retrieves all watchlist entries for a user
	GetWatchlist(ctx context.Context, userID uint) ([]models.WatchlistEntry, error)
	// GetWatchlistEntry retrieves a specific watchlist entry for a user/anime pair
	GetWatchlistEntry(ctx context.Context, userID, animeID uint) (*models.WatchlistEntry, error)
}

// WatchlistService handles business logic for watchlist operations
type WatchlistService struct {
	watchlistRepo repositories.WatchlistRepository
	logger        *logger.Logger
}

// NewWatchlistService creates a new WatchlistService instance
func NewWatchlistService(watchlistRepo repositories.WatchlistRepository) *WatchlistService {
	return &WatchlistService{
		watchlistRepo: watchlistRepo,
		logger:        logger.New(),
	}
}

// UpsertWatchlistEntry adds or updates an anime in a user's watchlist
func (s *WatchlistService) UpsertWatchlistEntry(ctx context.Context, userID, animeID uint, status models.WatchlistStatus) (*models.WatchlistEntry, error) {
	s.logger.WithFields(map[string]interface{}{
		"user_id":  userID,
		"anime_id": animeID,
		"status":   status,
	}).Info("Upserting watchlist entry")

	entry, err := s.watchlistRepo.UpsertByUserAndAnime(ctx, userID, animeID, status)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id":  userID,
			"anime_id": animeID,
			"error":    err.Error(),
		}).Error("Failed to upsert watchlist entry")
		return nil, err
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":  userID,
		"anime_id": animeID,
		"status":   status,
	}).Info("Successfully upserted watchlist entry")
	return entry, nil
}

// RemoveWatchlistEntry removes an anime from a user's watchlist
func (s *WatchlistService) RemoveWatchlistEntry(ctx context.Context, userID, animeID uint) error {
	s.logger.WithFields(map[string]interface{}{
		"user_id":  userID,
		"anime_id": animeID,
	}).Info("Removing watchlist entry")

	// Check if entry exists first
	_, err := s.watchlistRepo.GetByUserAndAnime(ctx, userID, animeID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id":  userID,
			"anime_id": animeID,
			"error":    err.Error(),
		}).Error("Watchlist entry not found")
		return errors.NewError(errors.ErrResourceNotFound, "Watchlist entry not found", err.Error(), http.StatusNotFound,
			map[string]interface{}{"user_id": userID, "anime_id": animeID}, err)
	}

	if err := s.watchlistRepo.DeleteByUserAndAnime(ctx, userID, animeID); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id":  userID,
			"anime_id": animeID,
			"error":    err.Error(),
		}).Error("Failed to delete watchlist entry")
		return err
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":  userID,
		"anime_id": animeID,
	}).Info("Successfully removed watchlist entry")
	return nil
}

// GetWatchlist retrieves all watchlist entries for a user
func (s *WatchlistService) GetWatchlist(ctx context.Context, userID uint) ([]models.WatchlistEntry, error) {
	s.logger.WithField("user_id", userID).Info("Retrieving watchlist for user")

	entries, err := s.watchlistRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to retrieve watchlist")
		return nil, err
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id": userID,
		"count":   len(entries),
	}).Info("Successfully retrieved watchlist for user")
	return entries, nil
}

// GetWatchlistEntry retrieves a specific watchlist entry for a user/anime pair
func (s *WatchlistService) GetWatchlistEntry(ctx context.Context, userID, animeID uint) (*models.WatchlistEntry, error) {
	s.logger.WithFields(map[string]interface{}{
		"user_id":  userID,
		"anime_id": animeID,
	}).Info("Retrieving watchlist entry")

	entry, err := s.watchlistRepo.GetByUserAndAnime(ctx, userID, animeID)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id":  userID,
			"anime_id": animeID,
			"error":    err.Error(),
		}).Error("Failed to retrieve watchlist entry")
		return nil, err
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":  userID,
		"anime_id": animeID,
	}).Info("Successfully retrieved watchlist entry")
	return entry, nil
}
