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
	db     db.DBInterface
	logger *logger.Logger
}

// NewGenreRepository creates a new GenreRepositoryImpl instance
func NewGenreRepository(db db.DBInterface) GenreRepository {
	return &GenreRepositoryImpl{
		db:     db,
		logger: logger.New(),
	}
}

// GetByID retrieves a genre by its ID
func (r *GenreRepositoryImpl) GetByID(ctx context.Context, id uint) (interface{}, error) {
	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Retrieving genre by ID")

	var genre models.Genre
	if err := r.db.WithContext(ctx).First(&genre, id).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to retrieve genre")
		return nil, errors.NewError(errors.ErrResourceNotFound, "Genre not found", fmt.Sprintf("Genre with ID %d not found", id), http.StatusNotFound, map[string]interface{}{"id": id}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Successfully retrieved genre")
	return &genre, nil
}

// GetAll retrieves all genres with optional pagination
func (r *GenreRepositoryImpl) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	r.logger.WithFields(map[string]interface{}{
		"page":  page,
		"limit": limit,
	}).Info("Retrieving all genres")

	var genres []models.Genre
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Genre{}).Count(&total).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to count genres")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count genres", err.Error(), http.StatusInternalServerError, nil, err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve genres with pagination
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&genres).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to retrieve genres")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve genres", err.Error(), http.StatusInternalServerError, nil, err)
	}

	// Convert to interface slice
	result := make([]interface{}, len(genres))
	for i, genre := range genres {
		result[i] = &genre
	}

	r.logger.WithFields(map[string]interface{}{
		"count": len(genres),
		"total": total,
	}).Info("Successfully retrieved genres")
	return result, total, nil
}

// Create creates a new genre
func (r *GenreRepositoryImpl) Create(ctx context.Context, entity interface{}) error {
	r.logger.Info("Creating new genre")

	genre, ok := entity.(*models.Genre)
	if !ok {
		r.logger.Error("Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.Genre", http.StatusBadRequest, nil, nil)
	}

	if err := r.db.WithContext(ctx).Create(genre).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to create genre")
		return errors.NewError(errors.ErrInternalServer, "Failed to create genre", err.Error(), http.StatusInternalServerError, nil, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": genre.ID,
	}).Info("Successfully created genre")
	return nil
}

// Update updates an existing genre
func (r *GenreRepositoryImpl) Update(ctx context.Context, entity interface{}) error {
	r.logger.Info("Updating genre")

	genre, ok := entity.(*models.Genre)
	if !ok {
		r.logger.Error("Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.Genre", http.StatusBadRequest, nil, nil)
	}

	if err := r.db.WithContext(ctx).Save(genre).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to update genre")
		return errors.NewError(errors.ErrInternalServer, "Failed to update genre", err.Error(), http.StatusInternalServerError, nil, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": genre.ID,
	}).Info("Successfully updated genre")
	return nil
}

// Delete deletes a genre by its ID
func (r *GenreRepositoryImpl) Delete(ctx context.Context, id uint) error {
	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Deleting genre")

	if err := r.db.WithContext(ctx).Delete(&models.Genre{}, id).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to delete genre")
		return errors.NewError(errors.ErrInternalServer, "Failed to delete genre", err.Error(), http.StatusInternalServerError, map[string]interface{}{"id": id}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Successfully deleted genre")
	return nil
}

// GetByIDs retrieves genres by their IDs
func (r *GenreRepositoryImpl) GetByIDs(ctx context.Context, ids []uint) ([]models.Genre, error) {
	r.logger.WithFields(map[string]interface{}{
		"ids": ids,
	}).Info("Retrieving genres by IDs")

	var genres []models.Genre
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&genres).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"ids":   ids,
			"error": err.Error(),
		}).Error("Failed to retrieve genres")
		return nil, errors.NewError(errors.ErrInternalServer, "Failed to retrieve genres", err.Error(), http.StatusInternalServerError, map[string]interface{}{"ids": ids}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"count": len(genres),
	}).Info("Successfully retrieved genres")
	return genres, nil
}

type SearchResult struct {
	Genres []models.Genre
	Total  int64
	Error  error
}

func (r *GenreRepositoryImpl) Search(ctx context.Context, query string, page, limit int) SearchResult {
	r.logger.WithFields(map[string]interface{}{
		"query": query,
		"page":  page,
		"limit": limit,
	}).Info("Searching genres")

	var genres []models.Genre
	var total int64

	offset := (page - 1) * limit

	// Search in both name and description
	result := r.db.WithContext(ctx).
		Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").
		Count(&total)

	if result.Error != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": result.Error.Error(),
		}).Error("Failed to count genres during search")
		return SearchResult{Error: result.Error}
	}

	result = r.db.WithContext(ctx).Offset(offset).Limit(limit).
		Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").
		Find(&genres)

	if result.Error != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": result.Error.Error(),
		}).Error("Failed to retrieve genres during search")
		return SearchResult{Error: result.Error}
	}

	r.logger.WithFields(map[string]interface{}{
		"count": len(genres),
		"total": total,
	}).Info("Successfully searched genres")
	return SearchResult{
		Genres: genres,
		Total:  total,
	}
}

// BulkCreate creates multiple genres at once
func (r *GenreRepositoryImpl) BulkCreate(ctx context.Context, genres []models.Genre) error {
	r.logger.WithFields(map[string]interface{}{
		"count": len(genres),
	}).Info("Bulk creating genres")

	if err := r.db.WithContext(ctx).Create(&genres).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to bulk create genres")
		return errors.NewError(errors.ErrInternalServer, "Failed to bulk create genres", err.Error(), http.StatusInternalServerError, nil, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"count": len(genres),
	}).Info("Successfully bulk created genres")
	return nil
}

// BulkDelete deletes multiple genres at once
func (r *GenreRepositoryImpl) BulkDelete(ctx context.Context, ids []uint) error {
	r.logger.WithFields(map[string]interface{}{
		"ids": ids,
	}).Info("Bulk deleting genres")

	if err := r.db.WithContext(ctx).Delete(&models.Genre{}, ids).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to bulk delete genres")
		return errors.NewError(errors.ErrInternalServer, "Failed to bulk delete genres", err.Error(), http.StatusInternalServerError, map[string]interface{}{"ids": ids}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"count": len(ids),
	}).Info("Successfully bulk deleted genres")
	return nil
}
