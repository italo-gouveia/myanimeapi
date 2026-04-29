package services

import (
	"context"
	"encoding/json"
	"fmt"
	"myanimeapi/api/adapters/cache"
	"myanimeapi/api/models"
	"myanimeapi/api/repositories"
	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"
	"net/http"
)

const (
	defaultPaginationLimit = 1000
)

// tagListCache is used to serialise paginated tag lists into/out of cache.
type tagListCache struct {
	Items []models.Tag `json:"items"`
	Total int64        `json:"total"`
}

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
	log     *logger.Logger
	cache   cache.CacheInterface
}

// NewTagService creates a new TagService instance
func NewTagService(tagRepo repositories.TagRepository, cacheImpl cache.CacheInterface) *TagService {
	return &TagService{
		tagRepo: tagRepo,
		log:     logger.New(),
		cache:   cacheImpl,
	}
}

// GetTagByID retrieves a tag by ID
func (s *TagService) GetTagByID(ctx context.Context, id uint) (*models.Tag, error) {
	s.log.WithField("id", id).Info("Retrieving tag by ID")

	// Cache-aside: check cache first
	cacheKey := cache.KeyTag(id)
	if data, err := s.cache.Get(ctx, cacheKey); err == nil {
		var tag models.Tag
		if jsonErr := json.Unmarshal(data, &tag); jsonErr == nil {
			s.log.WithField("id", id).Info("Tag cache hit")
			return &tag, nil
		}
	}

	tagInterface, err := s.tagRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag by ID: %w", err)
	}

	tag, ok := tagInterface.(*models.Tag)
	if !ok {
		s.log.WithField("id", id).Error("Invalid tag type in repository response")
		return nil, errors.NewError(errors.ErrInternalServer, "Invalid tag type", "Type assertion failed", http.StatusInternalServerError, map[string]interface{}{
			"id": id,
		}, nil)
	}

	// Populate cache
	if data, jsonErr := json.Marshal(tag); jsonErr == nil {
		_ = s.cache.Set(ctx, cacheKey, data, cache.TTLTag)
	}

	return tag, nil
}

// GetAllTags retrieves all tags with pagination
func (s *TagService) GetAllTags(ctx context.Context, page, limit int) ([]models.Tag, int64, error) {
	s.log.WithFields(map[string]interface{}{
		"page":  page,
		"limit": limit,
	}).Info("Retrieving all tags")

	// Cache-aside: check cache first
	cacheKey := cache.KeyTagList(page, limit)
	if data, err := s.cache.Get(ctx, cacheKey); err == nil {
		var cached tagListCache
		if jsonErr := json.Unmarshal(data, &cached); jsonErr == nil {
			s.log.WithFields(map[string]interface{}{"page": page, "limit": limit}).Info("Tag list cache hit")
			return cached.Items, cached.Total, nil
		}
	}

	tagsInterface, total, err := s.tagRepo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all tags: %w", err)
	}

	tags := make([]models.Tag, len(tagsInterface))
	for i, tagInterface := range tagsInterface {
		tag, ok := tagInterface.(*models.Tag)
		if !ok {
			s.log.WithField("index", i).Error("Invalid tag type in slice")
			return nil, 0, errors.NewError(errors.ErrInternalServer, "Invalid tag type", fmt.Sprintf("Type assertion failed at index %d", i), http.StatusInternalServerError, map[string]interface{}{
				"index": i,
			}, nil)
		}
		tags[i] = *tag
	}

	// Populate cache
	if data, jsonErr := json.Marshal(tagListCache{Items: tags, Total: total}); jsonErr == nil {
		_ = s.cache.Set(ctx, cacheKey, data, cache.TTLList)
	}

	return tags, total, nil
}

// CreateTag creates a new tag
func (s *TagService) CreateTag(ctx context.Context, tag *models.Tag) error {
	s.log.WithField("name", tag.Name).Info("Creating new tag")

	// Check if tag already exists by getting all tags and checking the name
	tags, _, err := s.GetAllTags(ctx, 1, defaultPaginationLimit)
	if err != nil {
		return fmt.Errorf("failed to check for existing tags: %w", err)
	}

	for _, existingTag := range tags {
		if existingTag.Name == tag.Name {
			s.log.WithField("name", tag.Name).Warning("Tag name already exists")
			return errors.NewError(errors.ErrConflict, "Name already exists", fmt.Sprintf("Tag with name '%s' already exists", tag.Name), http.StatusConflict, map[string]interface{}{
				"name": tag.Name,
			}, nil)
		}
	}

	if err := s.tagRepo.Create(ctx, tag); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	// Invalidate list caches
	_ = s.cache.DeleteByPattern(ctx, cache.PatternAllTags)

	return nil
}

// UpdateTag updates an existing tag
func (s *TagService) UpdateTag(ctx context.Context, tag *models.Tag) error {
	s.log.WithField("id", tag.ID).Info("Updating tag")

	// Get existing tag
	existingTag, err := s.GetTagByID(ctx, tag.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing tag: %w", err)
	}

	// Check if name is being changed and if it already exists
	if tag.Name != existingTag.Name {
		tags, _, err := s.GetAllTags(ctx, 1, defaultPaginationLimit)
		if err != nil {
			return fmt.Errorf("failed to check for existing tags: %w", err)
		}

		for _, t := range tags {
			if t.Name == tag.Name && t.ID != tag.ID {
				s.log.WithField("name", tag.Name).Warning("Tag name already exists")
				return errors.NewError(errors.ErrConflict, "Name already exists", fmt.Sprintf("Tag with name '%s' is already taken by another tag", tag.Name), http.StatusConflict, map[string]interface{}{
					"name": tag.Name,
					"id":   tag.ID,
				}, nil)
			}
		}
	}

	if err := s.tagRepo.Update(ctx, tag); err != nil {
		return fmt.Errorf("failed to update tag: %w", err)
	}

	// Invalidate single-item and list caches
	_ = s.cache.Delete(ctx, cache.KeyTag(tag.ID))
	_ = s.cache.DeleteByPattern(ctx, cache.PatternAllTags)

	return nil
}

// DeleteTag deletes a tag by ID
func (s *TagService) DeleteTag(ctx context.Context, id uint) error {
	s.log.WithField("id", id).Info("Deleting tag")

	// Check if tag exists
	if _, err := s.GetTagByID(ctx, id); err != nil {
		return fmt.Errorf("failed to verify tag existence: %w", err)
	}

	if err := s.tagRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	// Invalidate single-item and list caches
	_ = s.cache.Delete(ctx, cache.KeyTag(id))
	_ = s.cache.DeleteByPattern(ctx, cache.PatternAllTags)

	return nil
}

// GetTagsByIDs retrieves tags by their IDs
func (s *TagService) GetTagsByIDs(ctx context.Context, ids []uint) ([]models.Tag, error) {
	s.log.WithField("ids", ids).Info("Retrieving tags by IDs")

	tags, err := s.tagRepo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags by IDs: %w", err)
	}

	return tags, nil
}
