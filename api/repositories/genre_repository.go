package repositories

import (
	"myanimeapi/api/models"

	"gorm.io/gorm"
)

// GenreRepository defines the interface for genre operations
type GenreRepository interface {
	Create(genre *models.Genre) error
	GetAll() ([]models.Genre, error)
	GetByID(id uint) (*models.Genre, error)
	GetByName(name string) (*models.Genre, error)
	Update(genre *models.Genre) error
	Delete(id uint) error
}

// GenreRepositoryImpl implements the GenreRepository interface
type GenreRepositoryImpl struct {
	db *gorm.DB
}

// NewGenreRepository creates a new instance of GenreRepository
func NewGenreRepository(db *gorm.DB) GenreRepository {
	return &GenreRepositoryImpl{db: db}
}

// Create adds a new genre to the database
func (r *GenreRepositoryImpl) Create(genre *models.Genre) error {
	return r.db.Create(genre).Error
}

// GetAll retrieves all genres from the database
func (r *GenreRepositoryImpl) GetAll() ([]models.Genre, error) {
	var genres []models.Genre
	err := r.db.Find(&genres).Error
	return genres, err
}

// GetByID retrieves a genre by its ID
func (r *GenreRepositoryImpl) GetByID(id uint) (*models.Genre, error) {
	var genre models.Genre
	err := r.db.First(&genre, id).Error
	if err != nil {
		return nil, err
	}
	return &genre, nil
}

// GetByName retrieves a genre by its name
func (r *GenreRepositoryImpl) GetByName(name string) (*models.Genre, error) {
	var genre models.Genre
	err := r.db.Where("name = ?", name).First(&genre).Error
	if err != nil {
		return nil, err
	}
	return &genre, nil
}

// Update modifies an existing genre
func (r *GenreRepositoryImpl) Update(genre *models.Genre) error {
	return r.db.Save(genre).Error
}

// Delete removes a genre from the database
func (r *GenreRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Genre{}, id).Error
}
