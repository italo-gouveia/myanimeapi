package repositories

import (
	"myanimeapi/api/models"

	"gorm.io/gorm"
)

// TagRepository defines the interface for tag operations
type TagRepository interface {
	Create(tag *models.Tag) error
	GetAll() ([]models.Tag, error)
	GetByID(id uint) (*models.Tag, error)
	GetByName(name string) (*models.Tag, error)
	Update(tag *models.Tag) error
	Delete(id uint) error
}

// TagRepositoryImpl implements the TagRepository interface
type TagRepositoryImpl struct {
	db *gorm.DB
}

// NewTagRepository creates a new instance of TagRepository
func NewTagRepository(db *gorm.DB) TagRepository {
	return &TagRepositoryImpl{db: db}
}

// Create adds a new tag to the database
func (r *TagRepositoryImpl) Create(tag *models.Tag) error {
	return r.db.Create(tag).Error
}

// GetAll retrieves all tags from the database
func (r *TagRepositoryImpl) GetAll() ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Find(&tags).Error
	return tags, err
}

// GetByID retrieves a tag by its ID
func (r *TagRepositoryImpl) GetByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetByName retrieves a tag by its name
func (r *TagRepositoryImpl) GetByName(name string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// Update modifies an existing tag
func (r *TagRepositoryImpl) Update(tag *models.Tag) error {
	return r.db.Save(tag).Error
}

// Delete removes a tag from the database
func (r *TagRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Tag{}, id).Error
}
