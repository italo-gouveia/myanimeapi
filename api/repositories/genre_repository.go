package repositories

import (
	"context"
	"myanimeapi/api/models"
)

// GenreRepository defines the interface for genre operations
type GenreRepository interface {
	CreateGenre(ctx context.Context, genre *models.Genre) error
	GetGenreByID(ctx context.Context, id uint) (*models.Genre, error)
	GetAllGenres(ctx context.Context) ([]models.Genre, error)
	UpdateGenre(ctx context.Context, genre *models.Genre) error
	DeleteGenre(ctx context.Context, id uint) error
	GetGenresByIDs(ctx context.Context, ids []uint) ([]models.Genre, error)
}
