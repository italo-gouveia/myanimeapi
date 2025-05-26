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

// TagRepository defines the interface for tag operations
type TagRepository interface {
	Repository

	// GetByIDs retrieves tags by their IDs
	GetByIDs(ctx context.Context, ids []uint) ([]models.Tag, error)
}

// TagRepositoryImpl implements the TagRepository interface
type TagRepositoryImpl struct {
	db     db.DBInterface
	logger *logger.Logger
}

// NewTagRepository creates a new TagRepositoryImpl instance
func NewTagRepository(db db.DBInterface) TagRepository {
	return &TagRepositoryImpl{
		db:     db,
		logger: logger.New(),
	}
}

// GetByID retrieves a tag by its ID
func (r *TagRepositoryImpl) GetByID(ctx context.Context, id uint) (interface{}, error) {
	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Retrieving tag by ID")

	var tag models.Tag
	if err := r.db.WithContext(ctx).First(&tag, id).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to retrieve tag")
		return nil, errors.NewError(errors.ErrResourceNotFound, "Tag not found", fmt.Sprintf("Tag with ID %d not found", id), http.StatusNotFound, map[string]interface{}{"id": id}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Successfully retrieved tag")
	return &tag, nil
}

// GetAll retrieves all tags with optional pagination
func (r *TagRepositoryImpl) GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error) {
	r.logger.WithFields(map[string]interface{}{
		"page":  page,
		"limit": limit,
	}).Info("Retrieving all tags")

	var tags []models.Tag
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Tag{}).Count(&total).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to count tags")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to count tags", err.Error(), http.StatusInternalServerError, nil, err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve tags with pagination
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&tags).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to retrieve tags")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve tags", err.Error(), http.StatusInternalServerError, nil, err)
	}

	// Convert to interface slice
	result := make([]interface{}, len(tags))
	for i, tag := range tags {
		result[i] = &tag
	}

	r.logger.WithFields(map[string]interface{}{
		"count": len(tags),
		"total": total,
	}).Info("Successfully retrieved tags")
	return result, total, nil
}

// Create creates a new tag
func (r *TagRepositoryImpl) Create(ctx context.Context, entity interface{}) error {
	r.logger.Info("Creating new tag")

	tag, ok := entity.(*models.Tag)
	if !ok {
		r.logger.Error("Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.Tag", http.StatusBadRequest, nil, nil)
	}

	if err := r.db.WithContext(ctx).Create(tag).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to create tag")
		return errors.NewError(errors.ErrInternalServer, "Failed to create tag", err.Error(), http.StatusInternalServerError, nil, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": tag.ID,
	}).Info("Successfully created tag")
	return nil
}

// Update updates an existing tag
func (r *TagRepositoryImpl) Update(ctx context.Context, entity interface{}) error {
	r.logger.Info("Updating tag")

	tag, ok := entity.(*models.Tag)
	if !ok {
		r.logger.Error("Invalid entity type")
		return errors.NewError(errors.ErrInvalidInput, "Invalid entity type", "Expected *models.Tag", http.StatusBadRequest, nil, nil)
	}

	if err := r.db.WithContext(ctx).Save(tag).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to update tag")
		return errors.NewError(errors.ErrInternalServer, "Failed to update tag", err.Error(), http.StatusInternalServerError, nil, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": tag.ID,
	}).Info("Successfully updated tag")
	return nil
}

// Delete deletes a tag by its ID
func (r *TagRepositoryImpl) Delete(ctx context.Context, id uint) error {
	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Deleting tag")

	if err := r.db.WithContext(ctx).Delete(&models.Tag{}, id).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to delete tag")
		return errors.NewError(errors.ErrInternalServer, "Failed to delete tag", err.Error(), http.StatusInternalServerError, map[string]interface{}{"id": id}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Successfully deleted tag")
	return nil
}

// GetByIDs retrieves tags by their IDs
func (r *TagRepositoryImpl) GetByIDs(ctx context.Context, ids []uint) ([]models.Tag, error) {
	r.logger.WithFields(map[string]interface{}{
		"ids": ids,
	}).Info("Retrieving tags by IDs")

	var tags []models.Tag
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&tags).Error; err != nil {
		r.logger.WithFields(map[string]interface{}{
			"ids":   ids,
			"error": err.Error(),
		}).Error("Failed to retrieve tags")
		return nil, errors.NewError(errors.ErrInternalServer, "Failed to retrieve tags", err.Error(), http.StatusInternalServerError, map[string]interface{}{"ids": ids}, err)
	}

	r.logger.WithFields(map[string]interface{}{
		"count": len(tags),
	}).Info("Successfully retrieved tags")
	return tags, nil
}
