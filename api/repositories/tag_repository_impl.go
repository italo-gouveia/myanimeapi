package repositories

import (
	"context"
	"myanimeapi/api/models"

	"gorm.io/gorm"
)

// TagRepositoryImpl implements the TagRepository interface
type TagRepositoryImpl struct {
	db *gorm.DB
}

// NewTagRepository creates a new instance of TagRepositoryImpl
func NewTagRepository(db *gorm.DB) TagRepository {
	return &TagRepositoryImpl{db: db}
}

// CreateTag creates a new tag
func (r *TagRepositoryImpl) CreateTag(ctx context.Context, tag *models.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

// GetTagByID retrieves a tag by its ID
func (r *TagRepositoryImpl) GetTagByID(ctx context.Context, id uint) (*models.Tag, error) {
	var tag models.Tag
	if err := r.db.WithContext(ctx).First(&tag, id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetAllTags retrieves all tags
func (r *TagRepositoryImpl) GetAllTags(ctx context.Context) ([]models.Tag, error) {
	var tags []models.Tag
	if err := r.db.WithContext(ctx).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// UpdateTag updates an existing tag
func (r *TagRepositoryImpl) UpdateTag(ctx context.Context, tag *models.Tag) error {
	return r.db.WithContext(ctx).Save(tag).Error
}

// DeleteTag deletes a tag by its ID
func (r *TagRepositoryImpl) DeleteTag(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Tag{}, id).Error
}

// GetTagsByIDs retrieves tags by their IDs
func (r *TagRepositoryImpl) GetTagsByIDs(ctx context.Context, ids []uint) ([]models.Tag, error) {
	var tags []models.Tag
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}
