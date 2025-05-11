package repositories

import (
	"context"
	"log"
	"myanimeapi/api/models"
	"myanimeapi/internal/db"
	"myanimeapi/internal/errors"
	"net/http"
)

// GenreRepository defines the interface for genre operations
type GenreRepository interface {
	Repository

	// GetByIDs retrieves genres by their IDs
	GetByIDs(ctx context.Context, ids []uint) ([]models.Genre, error)

	// Search retrieves genres based on a query
	Search(ctx context.Context, query string, page, limit int) SearchResult

	// BulkCreate creates multiple genres at once
	BulkCreate(ctx context.Context, genres []models.Genre) error

	// BulkDelete deletes multiple genres at once
	BulkDelete(ctx context.Context, ids []uint) error
}

// GenreRepositoryImpl implements the GenreRepository interface
type GenreRepositoryImpl struct {
	db db.DBInterface
}

// NewGenreRepository creates a new GenreRepositoryImpl instance
func NewGenreRepository(db db.DBInterface) GenreRepository {
	return &GenreRepositoryImpl{db: db}
}

// GetByID retrieves a genre by its ID
func (r *GenreRepositoryImpl) GetByID(ctx context.Context, id uint) (interface{}, error) {
	log.Printf("GenreRepository.GetByID: Retrieving genre with ID %d", id)

	var genre models.Genre
	if err := r.db.WithContext(ctx).First(&genre, id).Error; err != nil {
		log.Printf("GenreRepository.GetByID: Failed to retrieve genre: %v", err)
		return nil, errors.NewError(errors.ErrResourceNotFound, "Genre not found", err.Error(), http.StatusNotFound)
	}

	log.Printf("GenreRepository.GetByID: Successfully retrieved genre with ID %d", id)
	return &genre, nil
}

// GetAll retrieves all genres with optional pagination
func (r *GenreRepositoryImpl) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	log.Printf("GenreRepository.GetAll: Retrieving all genres with page %d and limit %d", page, limit)

	var genres []models.Genre
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Genre{}).Count(&total).Error; err != nil {
		log.Printf("GenreRepository.GetAll: Failed to count genres: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count genres", err.Error(), http.StatusInternalServerError)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve genres with pagination
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&genres).Error; err != nil {
		log.Printf("GenreRepository.GetAll: Failed to retrieve genres: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve genres", err.Error(), http.StatusInternalServerError)
	}

	// Convert to interface slice
	result := make([]interface{}, len(genres))
	for i, genre := range genres {
		result[i] = &genre
	}

	log.Printf("GenreRepository.GetAll: Successfully retrieved %d genres", len(genres))
	return result, total, nil
}

// Create creates a new genre
func (r *GenreRepositoryImpl) Create(ctx context.Context, entity interface{}) error {
	log.Printf("GenreRepository.Create: Creating new genre")

	genre, ok := entity.(*models.Genre)
	if !ok {
		log.Printf("GenreRepository.Create: Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.Genre", http.StatusBadRequest)
	}

	if err := r.db.WithContext(ctx).Create(genre).Error; err != nil {
		log.Printf("GenreRepository.Create: Failed to create genre: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to create genre", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("GenreRepository.Create: Successfully created genre with ID %d", genre.ID)
	return nil
}

// Update updates an existing genre
func (r *GenreRepositoryImpl) Update(ctx context.Context, entity interface{}) error {
	log.Printf("GenreRepository.Update: Updating genre")

	genre, ok := entity.(*models.Genre)
	if !ok {
		log.Printf("GenreRepository.Update: Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.Genre", http.StatusBadRequest)
	}

	if err := r.db.WithContext(ctx).Save(genre).Error; err != nil {
		log.Printf("GenreRepository.Update: Failed to update genre: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to update genre", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("GenreRepository.Update: Successfully updated genre with ID %d", genre.ID)
	return nil
}

// Delete deletes a genre by its ID
func (r *GenreRepositoryImpl) Delete(ctx context.Context, id uint) error {
	log.Printf("GenreRepository.Delete: Deleting genre with ID %d", id)

	if err := r.db.WithContext(ctx).Delete(&models.Genre{}, id).Error; err != nil {
		log.Printf("GenreRepository.Delete: Failed to delete genre: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to delete genre", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("GenreRepository.Delete: Successfully deleted genre with ID %d", id)
	return nil
}

// GetByIDs retrieves genres by their IDs
func (r *GenreRepositoryImpl) GetByIDs(ctx context.Context, ids []uint) ([]models.Genre, error) {
	log.Printf("GenreRepository.GetByIDs: Retrieving genres with IDs %v", ids)

	var genres []models.Genre
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&genres).Error; err != nil {
		log.Printf("GenreRepository.GetByIDs: Failed to retrieve genres: %v", err)
		return nil, errors.NewError(errors.ErrInternalServer, "Failed to retrieve genres", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("GenreRepository.GetByIDs: Successfully retrieved %d genres", len(genres))
	return genres, nil
}

type SearchResult struct {
	Genres []models.Genre
	Total  int64
	Error  error
}

func (r *GenreRepositoryImpl) Search(ctx context.Context, query string, page, limit int) SearchResult {
	var genres []models.Genre
	var total int64

	offset := (page - 1) * limit

	// Search in both name and description
	result := r.db.WithContext(ctx).
		Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").
		Count(&total)

	if result.Error != nil {
		return SearchResult{Error: result.Error}
	}

	result = r.db.WithContext(ctx).Offset(offset).Limit(limit).
		Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").
		Find(&genres)

	if result.Error != nil {
		return SearchResult{Error: result.Error}
	}

	return SearchResult{
		Genres: genres,
		Total:  total,
	}
}

// BulkCreate creates multiple genres at once
func (r *GenreRepositoryImpl) BulkCreate(ctx context.Context, genres []models.Genre) error {
	return r.db.WithContext(ctx).Create(&genres).Error
}

// BulkDelete deletes multiple genres at once
func (r *GenreRepositoryImpl) BulkDelete(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).Delete(&models.Genre{}, ids).Error
}
