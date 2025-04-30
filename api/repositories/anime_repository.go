package repositories

import (
	"context"
	"log"
	"myanimeapi/api/models"
	"myanimeapi/internal/db"
	"myanimeapi/internal/errors"
	"net/http"
)

// AnimeRepositoryImpl implements the AnimeRepository interface
type AnimeRepositoryImpl struct {
	db db.DBInterface
}

// NewAnimeRepository creates a new AnimeRepositoryImpl instance
func NewAnimeRepository(db db.DBInterface) AnimeRepository {
	return &AnimeRepositoryImpl{db: db}
}

// GetByID retrieves an anime by its ID
func (r *AnimeRepositoryImpl) GetByID(ctx context.Context, id uint) (interface{}, error) {
	log.Printf("AnimeRepository.GetByID: Retrieving anime with ID %d", id)

	var anime models.Anime
	if err := r.db.WithContext(ctx).Preload("Reviews").Preload("Favorites").First(&anime, id).Error; err != nil {
		log.Printf("AnimeRepository.GetByID: Failed to retrieve anime: %v", err)
		return nil, errors.NewError(errors.ErrResourceNotFound, "Anime not found", err.Error(), http.StatusNotFound)
	}

	log.Printf("AnimeRepository.GetByID: Successfully retrieved anime with ID %d", id)
	return &anime, nil
}

// GetAll retrieves all animes with optional pagination
func (r *AnimeRepositoryImpl) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	log.Printf("AnimeRepository.GetAll: Retrieving all animes with page %d and limit %d", page, limit)

	var animes []models.Anime
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Anime{}).Count(&total).Error; err != nil {
		log.Printf("AnimeRepository.GetAll: Failed to count animes: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count animes", err.Error(), http.StatusInternalServerError)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve animes with pagination
	if err := r.db.WithContext(ctx).Preload("Reviews").Preload("Favorites").Offset(offset).Limit(limit).Find(&animes).Error; err != nil {
		log.Printf("AnimeRepository.GetAll: Failed to retrieve animes: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve animes", err.Error(), http.StatusInternalServerError)
	}

	// Convert to interface slice
	result := make([]interface{}, len(animes))
	for i, anime := range animes {
		result[i] = anime
	}

	log.Printf("AnimeRepository.GetAll: Successfully retrieved %d animes", len(animes))
	return result, total, nil
}

// Create creates a new anime
func (r *AnimeRepositoryImpl) Create(ctx context.Context, entity interface{}) error {
	log.Printf("AnimeRepository.Create: Creating new anime")

	anime, ok := entity.(*models.Anime)
	if !ok {
		log.Printf("AnimeRepository.Create: Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.Anime", http.StatusBadRequest)
	}

	if err := r.db.WithContext(ctx).Create(anime).Error; err != nil {
		log.Printf("AnimeRepository.Create: Failed to create anime: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to create anime", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("AnimeRepository.Create: Successfully created anime with ID %d", anime.ID)
	return nil
}

// Update updates an existing anime
func (r *AnimeRepositoryImpl) Update(ctx context.Context, entity interface{}) error {
	log.Printf("AnimeRepository.Update: Updating anime")

	anime, ok := entity.(*models.Anime)
	if !ok {
		log.Printf("AnimeRepository.Update: Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.Anime", http.StatusBadRequest)
	}

	if err := r.db.WithContext(ctx).Save(anime).Error; err != nil {
		log.Printf("AnimeRepository.Update: Failed to update anime: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to update anime", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("AnimeRepository.Update: Successfully updated anime with ID %d", anime.ID)
	return nil
}

// Delete deletes an anime by its ID
func (r *AnimeRepositoryImpl) Delete(ctx context.Context, id uint) error {
	log.Printf("AnimeRepository.Delete: Deleting anime with ID %d", id)

	if err := r.db.WithContext(ctx).Delete(&models.Anime{}, id).Error; err != nil {
		log.Printf("AnimeRepository.Delete: Failed to delete anime: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to delete anime", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("AnimeRepository.Delete: Successfully deleted anime with ID %d", id)
	return nil
}

// GetByTitle retrieves animes by title
func (r *AnimeRepositoryImpl) GetByTitle(ctx context.Context, title string, page, limit int) ([]models.Anime, int64, error) {
	log.Printf("AnimeRepository.GetByTitle: Retrieving animes with title %s, page %d, limit %d", title, page, limit)

	var animes []models.Anime
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Anime{}).Where("title LIKE ?", "%"+title+"%").Count(&total).Error; err != nil {
		log.Printf("AnimeRepository.GetByTitle: Failed to count animes: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count animes", err.Error(), http.StatusInternalServerError)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve animes with pagination
	if err := r.db.WithContext(ctx).Preload("Reviews").Preload("Favorites").Where("title LIKE ?", "%"+title+"%").Offset(offset).Limit(limit).Find(&animes).Error; err != nil {
		log.Printf("AnimeRepository.GetByTitle: Failed to retrieve animes: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve animes", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("AnimeRepository.GetByTitle: Successfully retrieved %d animes", len(animes))
	return animes, total, nil
}

// GetByGenre retrieves animes by genre
func (r *AnimeRepositoryImpl) GetByGenre(ctx context.Context, genre string, page, limit int) ([]models.Anime, int64, error) {
	log.Printf("AnimeRepository.GetByGenre: Retrieving animes with genre %s, page %d, limit %d", genre, page, limit)

	var animes []models.Anime
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Anime{}).Where("genre = ?", genre).Count(&total).Error; err != nil {
		log.Printf("AnimeRepository.GetByGenre: Failed to count animes: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count animes", err.Error(), http.StatusInternalServerError)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve animes with pagination
	if err := r.db.WithContext(ctx).Preload("Reviews").Preload("Favorites").Where("genre = ?", genre).Offset(offset).Limit(limit).Find(&animes).Error; err != nil {
		log.Printf("AnimeRepository.GetByGenre: Failed to retrieve animes: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve animes", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("AnimeRepository.GetByGenre: Successfully retrieved %d animes", len(animes))
	return animes, total, nil
}

// GetReviewsForAnime retrieves reviews for a specific anime
func (r *AnimeRepositoryImpl) GetReviewsForAnime(ctx context.Context, animeID uint, page, limit int) ([]models.Review, int64, error) {
	log.Printf("AnimeRepository.GetReviewsForAnime: Retrieving reviews for anime with ID %d", animeID)

	var reviews []models.Review
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&models.Review{}).Where("anime_id = ?", animeID).Count(&total).Error; err != nil {
		log.Printf("AnimeRepository.GetReviewsForAnime: Failed to count reviews: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count reviews", err.Error(), http.StatusInternalServerError)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated reviews
	if err := r.db.WithContext(ctx).Where("anime_id = ?", animeID).Offset(offset).Limit(limit).Find(&reviews).Error; err != nil {
		log.Printf("AnimeRepository.GetReviewsForAnime: Failed to retrieve reviews: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve reviews", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("AnimeRepository.GetReviewsForAnime: Successfully retrieved %d reviews", len(reviews))
	return reviews, total, nil
}
