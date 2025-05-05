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

// GenreServiceInterface defines the interface for genre operations
type GenreServiceInterface interface {
	GetGenreByID(ctx context.Context, id uint) (*models.Genre, error)
	GetAllGenres(ctx context.Context, page, limit int) ([]models.Genre, int64, error)
	CreateGenre(ctx context.Context, genre *models.Genre) error
	UpdateGenre(ctx context.Context, genre *models.Genre) error
	DeleteGenre(ctx context.Context, id uint) error
	GetGenresByIDs(ctx context.Context, ids []uint) ([]models.Genre, error)
}

// GenreService handles business logic for genre operations
type GenreService struct {
	genreRepo repositories.GenreRepository
}

// NewGenreService creates a new GenreService instance
func NewGenreService(genreRepo repositories.GenreRepository) *GenreService {
	return &GenreService{
		genreRepo: genreRepo,
	}
}

// GetGenreByID retrieves a genre by ID
func (s *GenreService) GetGenreByID(ctx context.Context, id uint) (*models.Genre, error) {
	log.Printf("GenreService.GetGenreByID: Retrieving genre with ID %d", id)

	genreInterface, err := s.genreRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	genre, ok := genreInterface.(*models.Genre)
	if !ok {
		log.Printf("GenreService.GetGenreByID: Invalid genre type in repository response")
		return nil, errors.NewError(errors.ErrInternalServer, "Invalid genre type", "Type assertion failed", http.StatusInternalServerError)
	}

	return genre, nil
}

// GetAllGenres retrieves all genres with pagination
func (s *GenreService) GetAllGenres(ctx context.Context, page, limit int) ([]models.Genre, int64, error) {
	log.Printf("GenreService.GetAllGenres: Retrieving all genres with page %d and limit %d", page, limit)

	genresInterface, total, err := s.genreRepo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

	genres := make([]models.Genre, len(genresInterface))
	for i, genreInterface := range genresInterface {
		genre, ok := genreInterface.(*models.Genre)
		if !ok {
			log.Printf("GenreService.GetAllGenres: Invalid genre type in slice at index %d", i)
			return nil, 0, errors.NewError(errors.ErrInternalServer, "Invalid genre type", fmt.Sprintf("Type assertion failed at index %d", i), http.StatusInternalServerError)
		}
		genres[i] = *genre
	}

	return genres, total, nil
}

// CreateGenre creates a new genre
func (s *GenreService) CreateGenre(ctx context.Context, genre *models.Genre) error {
	log.Printf("GenreService.CreateGenre: Creating new genre with name %s", genre.Name)

	// Check if genre already exists by getting all genres and checking the name
	genres, _, err := s.GetAllGenres(ctx, 1, 1000) // Get all genres to check for duplicates
	if err != nil {
		return err
	}

	for _, existingGenre := range genres {
		if existingGenre.Name == genre.Name {
			log.Printf("GenreService.CreateGenre: Name %s already exists", genre.Name)
			return errors.NewError(errors.ErrConflict, "Name already exists", fmt.Sprintf("Genre with name '%s' already exists", genre.Name), http.StatusConflict)
		}
	}

	return s.genreRepo.Create(ctx, genre)
}

// UpdateGenre updates an existing genre
func (s *GenreService) UpdateGenre(ctx context.Context, genre *models.Genre) error {
	log.Printf("GenreService.UpdateGenre: Updating genre with ID %d", genre.ID)

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
				log.Printf("GenreService.UpdateGenre: Name %s already exists", genre.Name)
				return errors.NewError(errors.ErrConflict, "Name already exists", fmt.Sprintf("Genre with name '%s' is already taken by another genre", genre.Name), http.StatusConflict)
			}
		}
	}

	return s.genreRepo.Update(ctx, genre)
}

// DeleteGenre deletes a genre by ID
func (s *GenreService) DeleteGenre(ctx context.Context, id uint) error {
	log.Printf("GenreService.DeleteGenre: Deleting genre with ID %d", id)

	// Check if genre exists
	if _, err := s.GetGenreByID(ctx, id); err != nil {
		return err
	}

	return s.genreRepo.Delete(ctx, id)
}

// GetGenresByIDs retrieves genres by their IDs
func (s *GenreService) GetGenresByIDs(ctx context.Context, ids []uint) ([]models.Genre, error) {
	log.Printf("GenreService.GetGenresByIDs: Retrieving genres with IDs %v", ids)
	return s.genreRepo.GetByIDs(ctx, ids)
}
