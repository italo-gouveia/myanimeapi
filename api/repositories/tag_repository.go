package repositories

import (
	"context"
	"myanimeapi/api/models"
)

// TagRepository defines the interface for tag operations
type TagRepository interface {
	CreateTag(ctx context.Context, tag *models.Tag) error
	GetTagByID(ctx context.Context, id uint) (*models.Tag, error)
	GetAllTags(ctx context.Context) ([]models.Tag, error)
	UpdateTag(ctx context.Context, tag *models.Tag) error
	DeleteTag(ctx context.Context, id uint) error
	GetTagsByIDs(ctx context.Context, ids []uint) ([]models.Tag, error)
}
