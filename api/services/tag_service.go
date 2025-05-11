package services

import (
	"context"
	"fmt"
	"log"
	"myanimeapi/api/models"
	"myanimeapi/api/repositories"
	"myanimeapi/internal/errors"
	"net/http"
)

// TagServiceInterface defines the interface for tag operations
type TagServiceInterface interface {
	GetTagByID(ctx context.Context, id uint) (*models.Tag, error)
	GetAllTags(ctx context.Context, page, limit int) ([]models.Tag, int64, error)
	CreateTag(ctx context.Context, tag *models.Tag) error
	UpdateTag(ctx context.Context, tag *models.Tag) error
	DeleteTag(ctx context.Context, id uint) error
	GetTagsByIDs(ctx context.Context, ids []uint) ([]models.Tag, error)
}

// TagService handles business logic for tag operations
type TagService struct {
	tagRepo repositories.TagRepository
}

// NewTagService creates a new TagService instance
func NewTagService(tagRepo repositories.TagRepository) *TagService {
	return &TagService{
		tagRepo: tagRepo,
	}
}

// GetTagByID retrieves a tag by ID
func (s *TagService) GetTagByID(ctx context.Context, id uint) (*models.Tag, error) {
	log.Printf("TagService.GetTagByID: Retrieving tag with ID %d", id)

	tagInterface, err := s.tagRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	tag, ok := tagInterface.(*models.Tag)
	if !ok {
		log.Printf("TagService.GetTagByID: Invalid tag type in repository response")
		return nil, errors.NewError(errors.ErrInternalServer, "Invalid tag type", "Type assertion failed", http.StatusInternalServerError)
	}

	return tag, nil
}

// GetAllTags retrieves all tags with pagination
func (s *TagService) GetAllTags(ctx context.Context, page, limit int) ([]models.Tag, int64, error) {
	log.Printf("TagService.GetAllTags: Retrieving all tags with page %d and limit %d", page, limit)

	tagsInterface, total, err := s.tagRepo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

	tags := make([]models.Tag, len(tagsInterface))
	for i, tagInterface := range tagsInterface {
		tag, ok := tagInterface.(*models.Tag)
		if !ok {
			log.Printf("TagService.GetAllTags: Invalid tag type in slice at index %d", i)
			return nil, 0, errors.NewError(errors.ErrInternalServer, "Invalid tag type", fmt.Sprintf("Type assertion failed at index %d", i), http.StatusInternalServerError)
		}
		tags[i] = *tag
	}

	return tags, total, nil
}

// CreateTag creates a new tag
func (s *TagService) CreateTag(ctx context.Context, tag *models.Tag) error {
	log.Printf("TagService.CreateTag: Creating new tag with name %s", tag.Name)

	// Check if tag already exists by getting all tags and checking the name
	tags, _, err := s.GetAllTags(ctx, 1, 1000) // Get all tags to check for duplicates
	if err != nil {
		return err
	}

	for _, existingTag := range tags {
		if existingTag.Name == tag.Name {
			log.Printf("TagService.CreateTag: Name %s already exists", tag.Name)
			return errors.NewError(errors.ErrConflict, "Name already exists", fmt.Sprintf("Tag with name '%s' already exists", tag.Name), http.StatusConflict)
		}
	}

	return s.tagRepo.Create(ctx, tag)
}

// UpdateTag updates an existing tag
func (s *TagService) UpdateTag(ctx context.Context, tag *models.Tag) error {
	log.Printf("TagService.UpdateTag: Updating tag with ID %d", tag.ID)

	// Get existing tag
	existingTag, err := s.GetTagByID(ctx, tag.ID)
	if err != nil {
		return err
	}

	// Check if name is being changed and if it already exists
	if tag.Name != existingTag.Name {
		tags, _, err := s.GetAllTags(ctx, 1, 1000) // Get all tags to check for duplicates
		if err != nil {
			return err
		}

		for _, t := range tags {
			if t.Name == tag.Name && t.ID != tag.ID {
				log.Printf("TagService.UpdateTag: Name %s already exists", tag.Name)
				return errors.NewError(errors.ErrConflict, "Name already exists", fmt.Sprintf("Tag with name '%s' is already taken by another tag", tag.Name), http.StatusConflict)
			}
		}
	}

	return s.tagRepo.Update(ctx, tag)
}

// DeleteTag deletes a tag by ID
func (s *TagService) DeleteTag(ctx context.Context, id uint) error {
	log.Printf("TagService.DeleteTag: Deleting tag with ID %d", id)

	// Check if tag exists
	if _, err := s.GetTagByID(ctx, id); err != nil {
		return err
	}

	return s.tagRepo.Delete(ctx, id)
}

// GetTagsByIDs retrieves tags by their IDs
func (s *TagService) GetTagsByIDs(ctx context.Context, ids []uint) ([]models.Tag, error) {
	log.Printf("TagService.GetTagsByIDs: Retrieving tags with IDs %v", ids)
	return s.tagRepo.GetByIDs(ctx, ids)
}
