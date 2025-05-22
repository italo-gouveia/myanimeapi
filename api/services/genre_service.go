package services

import (
	"context"
	"fmt"
	"myanimeapi/api/models"
	"myanimeapi/api/repositories"
	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"
	"net/http"
)

// GenreServiceInterface defines the interface for genre operations
type GenreServiceInterface interface {
	GetGenreByID(ctx context.Context, id uint) (*models.Genre, error)
	GetAllGenres(ctx context.Context, page, limit int) ([]models.Genre, int64, error)
	CreateGenre(ctx context.Context, genre *models.Genre) error
	UpdateGenre(ctx context.Context, genre *models.Genre) error
	DeleteGenre(ctx context.Context, id uint) error
	GetGenresByIDs(ctx context.Context, ids []uint) ([]models.Genre, error)
	SearchGenres(ctx context.Context, query string, page, limit int) ([]models.Genre, int64, error)
	BulkCreateGenres(ctx context.Context, genres []models.Genre) error
	BulkDeleteGenres(ctx context.Context, ids []uint) error
}

// GenreService handles business logic for genre operations
type GenreService struct {
	genreRepo repositories.GenreRepository
	logger    *logger.Logger
}

// NewGenreService creates a new GenreService instance
func NewGenreService(genreRepo repositories.GenreRepository) *GenreService {
	return &GenreService{
		genreRepo: genreRepo,
		logger:    logger.New(),
	}
}

// GetGenreByID retrieves a genre by ID
func (s *GenreService) GetGenreByID(ctx context.Context, id uint) (*models.Genre, error) {
	s.logger.WithField("genre_id", id).Info("Retrieving genre")

	genreInterface, err := s.genreRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"genre_id": id,
			"error":    err.Error(),
		}).Error("Failed to retrieve genre")
		return nil, errors.NewError(errors.ErrResourceNotFound, "Genre not found",
			fmt.Sprintf("Genre with ID %d not found", id), http.StatusNotFound,
			map[string]interface{}{
				"genre_id": id,
			},
			err)
	}

	genre, ok := genreInterface.(*models.Genre)
	if !ok {
		s.logger.WithField("genre_id", id).Error("Invalid genre type")
		return nil, errors.NewError(errors.ErrInternalServer, "Invalid genre type", "Type assertion failed", http.StatusInternalServerError,
			map[string]interface{}{
				"genre_id": id,
			},
			nil)
	}

	s.logger.WithFields(map[string]interface{}{
		"genre_id":   id,
		"genre_name": genre.Name,
	}).Info("Successfully retrieved genre")
	return genre, nil
}

// GetAllGenres retrieves all genres with pagination
func (s *GenreService) GetAllGenres(ctx context.Context, page, limit int) ([]models.Genre, int64, error) {
	s.logger.WithFields(map[string]interface{}{
		"page":  page,
		"limit": limit,
	}).Info("Retrieving all genres")

	genresInterface, total, err := s.genreRepo.GetAll(ctx, page, limit)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"page":  page,
			"limit": limit,
			"error": err.Error(),
		}).Error("Failed to retrieve genres")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to retrieve genres",
			err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"page":  page,
				"limit": limit,
			},
			err)
	}

	genres := make([]models.Genre, len(genresInterface))
	for i, genreInterface := range genresInterface {
		genre, ok := genreInterface.(*models.Genre)
		if !ok {
			s.logger.WithFields(map[string]interface{}{
				"index": i,
			}).Error("Invalid genre type in slice")
			return nil, 0, errors.NewError(errors.ErrInternalServer, "Invalid genre type", fmt.Sprintf("Type assertion failed at index %d", i), http.StatusInternalServerError,
				map[string]interface{}{
					"index": i,
				},
				nil)
		}
		genres[i] = *genre
	}

	s.logger.WithFields(map[string]interface{}{
		"count": len(genres),
		"total": total,
	}).Info("Successfully retrieved genres")
	return genres, total, nil
}

// CreateGenre creates a new genre
func (s *GenreService) CreateGenre(ctx context.Context, genre *models.Genre) error {
	s.logger.WithField("genre_name", genre.Name).Info("Creating new genre")

	// Check if genre already exists by getting all genres and checking the name
	genres, _, err := s.GetAllGenres(ctx, 1, 1000) // Get all genres to check for duplicates
	if err != nil {
		return err
	}

	for _, existingGenre := range genres {
		if existingGenre.Name == genre.Name {
			s.logger.WithField("genre_name", genre.Name).Warning("Genre with same name already exists")
			return errors.NewError(errors.ErrConflict, "Name already exists",
				fmt.Sprintf("Genre with name '%s' already exists", genre.Name), http.StatusConflict,
				map[string]interface{}{
					"name": genre.Name,
				},
				nil)
		}
	}

	if err := s.genreRepo.Create(ctx, genre); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"genre_name": genre.Name,
			"error":      err.Error(),
		}).Error("Failed to create genre")
		return errors.NewError(errors.ErrInternalServer, "Failed to create genre",
			err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"name": genre.Name,
			},
			err)
	}

	s.logger.WithFields(map[string]interface{}{
		"genre_id":   genre.ID,
		"genre_name": genre.Name,
	}).Info("Successfully created genre")
	return nil
}

// UpdateGenre updates an existing genre
func (s *GenreService) UpdateGenre(ctx context.Context, genre *models.Genre) error {
	s.logger.WithFields(map[string]interface{}{
		"genre_id":   genre.ID,
		"genre_name": genre.Name,
	}).Info("Updating genre")

	// Get existing genre
	existingGenre, err := s.GetGenreByID(ctx, genre.ID)
	if err != nil {
		return err
	}

	// Check if name is being changed and if it already exists
	if genre.Name != existingGenre.Name {
		genres, _, err := s.GetAllGenres(ctx, 1, 1000) // Get all genres to check for duplicates
		if err != nil {
			return err
		}

		for _, g := range genres {
			if g.Name == genre.Name && g.ID != genre.ID {
				s.logger.WithField("genre_name", genre.Name).Warning("Genre with same name already exists")
				return errors.NewError(errors.ErrConflict, "Name already exists",
					fmt.Sprintf("Genre with name '%s' is already taken by another genre", genre.Name), http.StatusConflict,
					map[string]interface{}{
						"name": genre.Name,
					},
					nil)
			}
		}
	}

	if err := s.genreRepo.Update(ctx, genre); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"genre_id":   genre.ID,
			"genre_name": genre.Name,
			"error":      err.Error(),
		}).Error("Failed to update genre")
		return errors.NewError(errors.ErrInternalServer, "Failed to update genre",
			err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"genre_id": genre.ID,
				"name":     genre.Name,
			},
			err)
	}

	s.logger.WithFields(map[string]interface{}{
		"genre_id":   genre.ID,
		"genre_name": genre.Name,
	}).Info("Successfully updated genre")
	return nil
}

// DeleteGenre deletes a genre by ID
func (s *GenreService) DeleteGenre(ctx context.Context, id uint) error {
	s.logger.WithField("genre_id", id).Info("Deleting genre")

	// Check if genre exists
	existingGenre, err := s.GetGenreByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.genreRepo.Delete(ctx, id); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"genre_id": id,
			"error":    err.Error(),
		}).Error("Failed to delete genre")
		return errors.NewError(errors.ErrInternalServer, "Failed to delete genre",
			err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"genre_id": id,
			},
			err)
	}

	s.logger.WithFields(map[string]interface{}{
		"genre_id":   id,
		"genre_name": existingGenre.Name,
	}).Info("Successfully deleted genre")
	return nil
}

// GetGenresByIDs retrieves genres by their IDs
func (s *GenreService) GetGenresByIDs(ctx context.Context, ids []uint) ([]models.Genre, error) {
	s.logger.WithField("genre_ids", ids).Info("Retrieving genres by IDs")

	genres, err := s.genreRepo.GetByIDs(ctx, ids)
	if err != nil {
		s.logger.WithFields(map[string]interface{}{
			"genre_ids": ids,
			"error":     err.Error(),
		}).Error("Failed to retrieve genres by IDs")
		return nil, errors.NewError(errors.ErrInternalServer, "Failed to retrieve genres by IDs",
			err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"genre_ids": ids,
			},
			err)
	}

	s.logger.WithFields(map[string]interface{}{
		"genre_ids": ids,
		"count":     len(genres),
	}).Info("Successfully retrieved genres by IDs")
	return genres, nil
}

// SearchGenres searches for genres based on a query string
func (s *GenreService) SearchGenres(ctx context.Context, query string, page, limit int) ([]models.Genre, int64, error) {
	s.logger.WithFields(map[string]interface{}{
		"query": query,
		"page":  page,
		"limit": limit,
	}).Info("Searching genres")

	result := s.genreRepo.Search(ctx, query, page, limit)
	if result.Error != nil {
		s.logger.WithFields(map[string]interface{}{
			"query": query,
			"error": result.Error.Error(),
		}).Error("Failed to search genres")
		return nil, 0, errors.NewError(errors.ErrInternalServer, "Failed to search genres",
			result.Error.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"query": query,
				"page":  page,
				"limit": limit,
			},
			result.Error)
	}

	s.logger.WithFields(map[string]interface{}{
		"query": query,
		"count": len(result.Genres),
		"total": result.Total,
	}).Info("Successfully searched genres")
	return result.Genres, result.Total, nil
}

// BulkCreateGenres creates multiple genres at once
func (s *GenreService) BulkCreateGenres(ctx context.Context, genres []models.Genre) error {
	s.logger.WithField("count", len(genres)).Info("Bulk creating genres")

	if err := s.genreRepo.BulkCreate(ctx, genres); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"count": len(genres),
			"error": err.Error(),
		}).Error("Failed to bulk create genres")
		return errors.NewError(errors.ErrInternalServer, "Failed to bulk create genres",
			err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"count": len(genres),
			},
			err)
	}

	s.logger.WithField("count", len(genres)).Info("Successfully bulk created genres")
	return nil
}

// BulkDeleteGenres deletes multiple genres at once
func (s *GenreService) BulkDeleteGenres(ctx context.Context, ids []uint) error {
	s.logger.WithField("genre_ids", ids).Info("Bulk deleting genres")

	if err := s.genreRepo.BulkDelete(ctx, ids); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"genre_ids": ids,
			"error":     err.Error(),
		}).Error("Failed to bulk delete genres")
		return errors.NewError(errors.ErrInternalServer, "Failed to bulk delete genres",
			err.Error(), http.StatusInternalServerError,
			map[string]interface{}{
				"genre_ids": ids,
			},
			err)
	}

	s.logger.WithField("genre_ids", ids).Info("Successfully bulk deleted genres")
	return nil
}
