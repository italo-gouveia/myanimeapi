package repositories

import (
	"context"
	"myanimeapi/api/models"

	"gorm.io/gorm"
)

// GenreRepositoryImpl implements the GenreRepository interface
type GenreRepositoryImpl struct {
	db *gorm.DB
}

// NewGenreRepository creates a new instance of GenreRepositoryImpl
func NewGenreRepository(db *gorm.DB) GenreRepository {
	return &GenreRepositoryImpl{db: db}
}

// CreateGenre creates a new genre
func (r *GenreRepositoryImpl) CreateGenre(ctx context.Context, genre *models.Genre) error {
	return r.db.WithContext(ctx).Create(genre).Error
}

// GetGenreByID retrieves a genre by its ID
func (r *GenreRepositoryImpl) GetGenreByID(ctx context.Context, id uint) (*models.Genre, error) {
	var genre models.Genre
	if err := r.db.WithContext(ctx).First(&genre, id).Error; err != nil {
		return nil, err
	}
	return &genre, nil
}

// GetAllGenres retrieves all genres
func (r *GenreRepositoryImpl) GetAllGenres(ctx context.Context) ([]models.Genre, error) {
	var genres []models.Genre
	if err := r.db.WithContext(ctx).Find(&genres).Error; err != nil {
		return nil, err
	}
	return genres, nil
}

// UpdateGenre updates an existing genre
func (r *GenreRepositoryImpl) UpdateGenre(ctx context.Context, genre *models.Genre) error {
	return r.db.WithContext(ctx).Save(genre).Error
}

// DeleteGenre deletes a genre by its ID
func (r *GenreRepositoryImpl) DeleteGenre(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Genre{}, id).Error
}

// GetGenresByIDs retrieves genres by their IDs
func (r *GenreRepositoryImpl) GetGenresByIDs(ctx context.Context, ids []uint) ([]models.Genre, error) {
	var genres []models.Genre
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&genres).Error; err != nil {
		return nil, err
	}
	return genres, nil
}
