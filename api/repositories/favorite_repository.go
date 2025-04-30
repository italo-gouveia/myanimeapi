package repositories

import (
	"context"
	"log"
	"myanimeapi/api/models"
	"myanimeapi/internal/db"
	"myanimeapi/internal/errors"
	"net/http"
)

// FavoriteRepositoryImpl implements the FavoriteRepository interface
type FavoriteRepositoryImpl struct {
	db db.DBInterface
}

// NewFavoriteRepository creates a new FavoriteRepositoryImpl instance
func NewFavoriteRepository(db db.DBInterface) FavoriteRepository {
	return &FavoriteRepositoryImpl{db: db}
}

// GetByID retrieves a favorite by its ID
func (r *FavoriteRepositoryImpl) GetByID(ctx context.Context, id uint) (interface{}, error) {
	log.Printf("FavoriteRepository.GetByID: Retrieving favorite with ID %d", id)

	var favorite models.Favorite
	if err := r.db.WithContext(ctx).Preload("User").Preload("Anime").First(&favorite, id).Error; err != nil {
		log.Printf("FavoriteRepository.GetByID: Failed to retrieve favorite: %v", err)
		return nil, errors.NewError(errors.ErrResourceNotFound, "Favorite not found", err.Error(), http.StatusNotFound)
	}

	log.Printf("FavoriteRepository.GetByID: Successfully retrieved favorite with ID %d", id)
	return &favorite, nil
}

// GetAll retrieves all favorites with optional pagination
func (r *FavoriteRepositoryImpl) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	log.Printf("FavoriteRepository.GetAll: Retrieving all favorites with page %d and limit %d", page, limit)

	var favorites []models.Favorite
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Favorite{}).Count(&total).Error; err != nil {
		log.Printf("FavoriteRepository.GetAll: Failed to count favorites: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count favorites", err.Error(), http.StatusInternalServerError)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve favorites with pagination
	if err := r.db.WithContext(ctx).Preload("User").Preload("Anime").Offset(offset).Limit(limit).Find(&favorites).Error; err != nil {
		log.Printf("FavoriteRepository.GetAll: Failed to retrieve favorites: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve favorites", err.Error(), http.StatusInternalServerError)
	}

	// Convert to interface slice
	result := make([]interface{}, len(favorites))
	for i, favorite := range favorites {
		result[i] = &favorite
	}

	log.Printf("FavoriteRepository.GetAll: Successfully retrieved %d favorites", len(favorites))
	return result, total, nil
}

// Create creates a new favorite
func (r *FavoriteRepositoryImpl) Create(ctx context.Context, entity interface{}) error {
	log.Printf("FavoriteRepository.Create: Creating new favorite")

	favorite, ok := entity.(*models.Favorite)
	if !ok {
		log.Printf("FavoriteRepository.Create: Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.Favorite", http.StatusBadRequest)
	}

	if err := r.db.WithContext(ctx).Create(favorite).Error; err != nil {
		log.Printf("FavoriteRepository.Create: Failed to create favorite: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to create favorite", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("FavoriteRepository.Create: Successfully created favorite with ID %d", favorite.ID)
	return nil
}

// Update updates an existing favorite
func (r *FavoriteRepositoryImpl) Update(ctx context.Context, entity interface{}) error {
	log.Printf("FavoriteRepository.Update: Updating favorite")

	favorite, ok := entity.(*models.Favorite)
	if !ok {
		log.Printf("FavoriteRepository.Update: Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.Favorite", http.StatusBadRequest)
	}

	if err := r.db.WithContext(ctx).Save(favorite).Error; err != nil {
		log.Printf("FavoriteRepository.Update: Failed to update favorite: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to update favorite", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("FavoriteRepository.Update: Successfully updated favorite with ID %d", favorite.ID)
	return nil
}

// Delete deletes a favorite by its ID
func (r *FavoriteRepositoryImpl) Delete(ctx context.Context, id uint) error {
	log.Printf("FavoriteRepository.Delete: Deleting favorite with ID %d", id)

	if err := r.db.WithContext(ctx).Delete(&models.Favorite{}, id).Error; err != nil {
		log.Printf("FavoriteRepository.Delete: Failed to delete favorite: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to delete favorite", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("FavoriteRepository.Delete: Successfully deleted favorite with ID %d", id)
	return nil
}

// GetByUserID retrieves favorites by user ID
func (r *FavoriteRepositoryImpl) GetByUserID(ctx context.Context, userID uint) ([]models.Favorite, error) {
	log.Printf("FavoriteRepository.GetByUserID: Retrieving favorites for user with ID %d", userID)

	var favorites []models.Favorite
	if err := r.db.WithContext(ctx).Preload("Anime").Where("user_id = ?", userID).Find(&favorites).Error; err != nil {
		log.Printf("FavoriteRepository.GetByUserID: Failed to retrieve favorites: %v", err)
		return nil, errors.NewError(errors.ErrInternalServer, "Failed to retrieve favorites", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("FavoriteRepository.GetByUserID: Successfully retrieved %d favorites for user with ID %d", len(favorites), userID)
	return favorites, nil
}

// GetByAnimeID retrieves favorites by anime ID
func (r *FavoriteRepositoryImpl) GetByAnimeID(ctx context.Context, animeID uint, page, limit int) ([]models.Favorite, int64, error) {
	log.Printf("FavoriteRepository.GetByAnimeID: Retrieving favorites for anime with ID %d", animeID)

	var favorites []models.Favorite
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Favorite{}).Where("anime_id = ?", animeID).Count(&total).Error; err != nil {
		log.Printf("FavoriteRepository.GetByAnimeID: Failed to count favorites: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count favorites", err.Error(), http.StatusInternalServerError)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve favorites with pagination
	if err := r.db.WithContext(ctx).Preload("User").Where("anime_id = ?", animeID).Offset(offset).Limit(limit).Find(&favorites).Error; err != nil {
		log.Printf("FavoriteRepository.GetByAnimeID: Failed to retrieve favorites: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve favorites", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("FavoriteRepository.GetByAnimeID: Successfully retrieved %d favorites for anime with ID %d", len(favorites), animeID)
	return favorites, total, nil
}

// GetByUserAndAnime retrieves a favorite by user ID and anime ID
func (r *FavoriteRepositoryImpl) GetByUserAndAnime(ctx context.Context, userID, animeID uint) (*models.Favorite, error) {
	log.Printf("FavoriteRepository.GetByUserAndAnime: Retrieving favorite for user with ID %d and anime with ID %d", userID, animeID)

	var favorite models.Favorite
	if err := r.db.WithContext(ctx).Where("user_id = ? AND anime_id = ?", userID, animeID).First(&favorite).Error; err != nil {
		log.Printf("FavoriteRepository.GetByUserAndAnime: Failed to retrieve favorite: %v", err)
		return nil, errors.NewError(errors.ErrResourceNotFound, "Favorite not found", err.Error(), http.StatusNotFound)
	}

	log.Printf("FavoriteRepository.GetByUserAndAnime: Successfully retrieved favorite with ID %d", favorite.ID)
	return &favorite, nil
}

// DeleteByUserAndAnime deletes a favorite by user ID and anime ID
func (r *FavoriteRepositoryImpl) DeleteByUserAndAnime(ctx context.Context, userID, animeID uint) error {
	log.Printf("FavoriteRepository.DeleteByUserAndAnime: Deleting favorite for user with ID %d and anime with ID %d", userID, animeID)

	if err := r.db.WithContext(ctx).Where("user_id = ? AND anime_id = ?", userID, animeID).Delete(&models.Favorite{}).Error; err != nil {
		log.Printf("FavoriteRepository.DeleteByUserAndAnime: Failed to delete favorite: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to delete favorite", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("FavoriteRepository.DeleteByUserAndAnime: Successfully deleted favorite")
	return nil
}
