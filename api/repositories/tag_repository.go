package repositories

import (
	"context"
	"log"
	"myanimeapi/api/models"
	"myanimeapi/internal/db"
	"myanimeapi/internal/errors"
	"net/http"
)

// TagRepository defines the interface for tag operations
type TagRepository interface {
	Repository

	// GetByIDs retrieves tags by their IDs
	GetByIDs(ctx context.Context, ids []uint) ([]models.Tag, error)
}

// TagRepositoryImpl implements the TagRepository interface
type TagRepositoryImpl struct {
	db db.DBInterface
}

// NewTagRepository creates a new TagRepositoryImpl instance
func NewTagRepository(db db.DBInterface) TagRepository {
	return &TagRepositoryImpl{db: db}
}

// GetByID retrieves a tag by its ID
func (r *TagRepositoryImpl) GetByID(ctx context.Context, id uint) (interface{}, error) {
	log.Printf("TagRepository.GetByID: Retrieving tag with ID %d", id)

	var tag models.Tag
	if err := r.db.WithContext(ctx).First(&tag, id).Error; err != nil {
		log.Printf("TagRepository.GetByID: Failed to retrieve tag: %v", err)
		return nil, errors.NewError(errors.ErrResourceNotFound, "Tag not found", err.Error(), http.StatusNotFound)
	}

	log.Printf("TagRepository.GetByID: Successfully retrieved tag with ID %d", id)
	return &tag, nil
}

// GetAll retrieves all tags with optional pagination
func (r *TagRepositoryImpl) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	log.Printf("TagRepository.GetAll: Retrieving all tags with page %d and limit %d", page, limit)

	var tags []models.Tag
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Tag{}).Count(&total).Error; err != nil {
		log.Printf("TagRepository.GetAll: Failed to count tags: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count tags", err.Error(), http.StatusInternalServerError)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve tags with pagination
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&tags).Error; err != nil {
		log.Printf("TagRepository.GetAll: Failed to retrieve tags: %v", err)
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve tags", err.Error(), http.StatusInternalServerError)
	}

	// Convert to interface slice
	result := make([]interface{}, len(tags))
	for i, tag := range tags {
		result[i] = &tag
	}

	log.Printf("TagRepository.GetAll: Successfully retrieved %d tags", len(tags))
	return result, total, nil
}

// Create creates a new tag
func (r *TagRepositoryImpl) Create(ctx context.Context, entity interface{}) error {
	log.Printf("TagRepository.Create: Creating new tag")

	tag, ok := entity.(*models.Tag)
	if !ok {
		log.Printf("TagRepository.Create: Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.Tag", http.StatusBadRequest)
	}

	if err := r.db.WithContext(ctx).Create(tag).Error; err != nil {
		log.Printf("TagRepository.Create: Failed to create tag: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to create tag", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("TagRepository.Create: Successfully created tag with ID %d", tag.ID)
	return nil
}

// Update updates an existing tag
func (r *TagRepositoryImpl) Update(ctx context.Context, entity interface{}) error {
	log.Printf("TagRepository.Update: Updating tag")

	tag, ok := entity.(*models.Tag)
	if !ok {
		log.Printf("TagRepository.Update: Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.Tag", http.StatusBadRequest)
	}

	if err := r.db.WithContext(ctx).Save(tag).Error; err != nil {
		log.Printf("TagRepository.Update: Failed to update tag: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to update tag", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("TagRepository.Update: Successfully updated tag with ID %d", tag.ID)
	return nil
}

// Delete deletes a tag by its ID
func (r *TagRepositoryImpl) Delete(ctx context.Context, id uint) error {
	log.Printf("TagRepository.Delete: Deleting tag with ID %d", id)

	if err := r.db.WithContext(ctx).Delete(&models.Tag{}, id).Error; err != nil {
		log.Printf("TagRepository.Delete: Failed to delete tag: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to delete tag", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("TagRepository.Delete: Successfully deleted tag with ID %d", id)
	return nil
}

// GetByIDs retrieves tags by their IDs
func (r *TagRepositoryImpl) GetByIDs(ctx context.Context, ids []uint) ([]models.Tag, error) {
	log.Printf("TagRepository.GetByIDs: Retrieving tags with IDs %v", ids)

	var tags []models.Tag
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&tags).Error; err != nil {
		log.Printf("TagRepository.GetByIDs: Failed to retrieve tags: %v", err)
		return nil, errors.NewError(errors.ErrInternalServer, "Failed to retrieve tags", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("TagRepository.GetByIDs: Successfully retrieved %d tags", len(tags))
	return tags, nil
}
