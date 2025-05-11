package services

import (
	"context"
	"fmt"
	"log"
	"myanimeapi/api/models"
	"myanimeapi/api/repositories"
	"myanimeapi/internal/errors"
	"net/http"
	"time"
)

// AnimeServiceInterface defines the interface for anime service operations
type AnimeServiceInterface interface {
	GetAnimeByID(ctx context.Context, id uint) (*models.Anime, error)
	GetAllAnimes(ctx context.Context, page, limit int) ([]*models.Anime, int64, error)
	CreateAnime(ctx context.Context, anime *models.Anime) error
	UpdateAnime(ctx context.Context, anime *models.Anime) error
	DeleteAnime(ctx context.Context, id uint) error
	GetAnimesByTitle(ctx context.Context, title string, page, limit int) ([]*models.Anime, int64, error)
	GetAnimesByGenre(ctx context.Context, genre string, page, limit int) ([]*models.Anime, int64, error)
	GetReviewsForAnime(ctx context.Context, animeID uint, page, limit int) ([]*models.Review, int64, error)
	AddGenresToAnime(ctx context.Context, animeID uint, genreIDs []uint) error
	RemoveGenresFromAnime(ctx context.Context, animeID uint, genreIDs []uint) error
	AddTagsToAnime(ctx context.Context, animeID uint, tagIDs []uint) error
	RemoveTagsFromAnime(ctx context.Context, animeID uint, tagIDs []uint) error
}

// AnimeService handles business logic for anime operations
type AnimeService struct {
	animeRepo  repositories.AnimeRepository
	genreRepo  repositories.GenreRepository
	tagRepo    repositories.TagRepository
	reviewRepo repositories.ReviewRepository
}

// NewAnimeService creates a new AnimeService instance
func NewAnimeService(
	animeRepo repositories.AnimeRepository,
	genreRepo repositories.GenreRepository,
	tagRepo repositories.TagRepository,
	reviewRepo repositories.ReviewRepository,
) *AnimeService {
	return &AnimeService{
		animeRepo:  animeRepo,
		genreRepo:  genreRepo,
		tagRepo:    tagRepo,
		reviewRepo: reviewRepo,
	}
}

// GetAnimeByID retrieves an anime by ID
func (s *AnimeService) GetAnimeByID(ctx context.Context, id uint) (*models.Anime, error) {
	log.Printf("AnimeService.GetAnimeByID: Retrieving anime with ID %d", id)

	animeInterface, err := s.animeRepo.GetByID(ctx, id)
	if err != nil {
		log.Printf("AnimeService.GetAnimeByID: Failed to retrieve anime: %v", err)
		return nil, err
	}

	anime, ok := animeInterface.(*models.Anime)
	if !ok {
		log.Printf("AnimeService.GetAnimeByID: Invalid anime type")
		return nil, errors.NewError(errors.ErrInternalServer, "Invalid anime type", "Type assertion failed", http.StatusInternalServerError)
	}

	return anime, nil
}

// GetAllAnimes retrieves all animes with pagination
func (s *AnimeService) GetAllAnimes(ctx context.Context, page, limit int) ([]*models.Anime, int64, error) {
	log.Printf("AnimeService.GetAllAnimes: Retrieving all animes with page %d and limit %d", page, limit)

	animesInterface, total, err := s.animeRepo.GetAll(ctx, page, limit)
	if err != nil {
		log.Printf("AnimeService.GetAllAnimes: Failed to retrieve animes: %v", err)
		return nil, 0, err
	}

	animes := make([]*models.Anime, len(animesInterface))
	for i, animeInterface := range animesInterface {
		anime, ok := animeInterface.(*models.Anime)
		if !ok {
			log.Printf("AnimeService.GetAllAnimes: Invalid anime type in slice at index %d", i)
			return nil, 0, errors.NewError(errors.ErrInternalServer, "Invalid anime type", fmt.Sprintf("Type assertion failed at index %d", i), http.StatusInternalServerError)
		}
		animes[i] = anime
	}

	return animes, total, nil
}

// CreateAnime creates a new anime
func (s *AnimeService) CreateAnime(ctx context.Context, anime *models.Anime) error {
	log.Printf("AnimeService.CreateAnime: Creating new anime with title %s", anime.Title)

	// Check if title already exists
	existingAnimes, _, err := s.animeRepo.GetByTitle(ctx, anime.Title, 1, 1)
	if err == nil && len(existingAnimes) > 0 {
		log.Printf("AnimeService.CreateAnime: Title %s already exists", anime.Title)
		return errors.NewError(errors.ErrConflict, "Title already exists", fmt.Sprintf("Anime with title '%s' already exists", anime.Title), http.StatusConflict)
	}

	// Set timestamps
	anime.CreatedAt = time.Now()
	anime.UpdatedAt = time.Now()

	// Create anime
	if err := s.animeRepo.Create(ctx, anime); err != nil {
		log.Printf("AnimeService.CreateAnime: Failed to create anime: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to create anime", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("AnimeService.CreateAnime: Successfully created anime with ID %d", anime.ID)
	return nil
}

// UpdateAnime updates an existing anime
func (s *AnimeService) UpdateAnime(ctx context.Context, anime *models.Anime) error {
	log.Printf("AnimeService.UpdateAnime: Updating anime with ID %d", anime.ID)

	// Check if anime exists
	existingAnime, err := s.animeRepo.GetByID(ctx, anime.ID)
	if err != nil {
		log.Printf("AnimeService.UpdateAnime: Failed to retrieve anime: %v", err)
		return errors.NewError(errors.ErrResourceNotFound, "Anime not found", err.Error(), http.StatusNotFound)
	}

	// Check if title is being changed and if it already exists
	if anime.Title != existingAnime.(*models.Anime).Title {
		titleAnimes, _, err := s.animeRepo.GetByTitle(ctx, anime.Title, 1, 1)
		if err == nil && len(titleAnimes) > 0 && titleAnimes[0].ID != anime.ID {
			log.Printf("AnimeService.UpdateAnime: Title %s already exists", anime.Title)
			return errors.NewError(errors.ErrConflict, "Title already exists", fmt.Sprintf("Anime with title '%s' is already taken by another anime", anime.Title), http.StatusConflict)
		}
	}

	// Set updated timestamp
	anime.UpdatedAt = time.Now()

	// Update anime
	if err := s.animeRepo.Update(ctx, anime); err != nil {
		log.Printf("AnimeService.UpdateAnime: Failed to update anime: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to update anime", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("AnimeService.UpdateAnime: Successfully updated anime with ID %d", anime.ID)
	return nil
}

// DeleteAnime deletes an anime by ID
func (s *AnimeService) DeleteAnime(ctx context.Context, id uint) error {
	log.Printf("AnimeService.DeleteAnime: Deleting anime with ID %d", id)

	// Check if anime exists
	_, err := s.animeRepo.GetByID(ctx, id)
	if err != nil {
		log.Printf("AnimeService.DeleteAnime: Failed to retrieve anime: %v", err)
		return errors.NewError(errors.ErrResourceNotFound, "Anime not found", err.Error(), http.StatusNotFound)
	}

	// Delete anime
	if err := s.animeRepo.Delete(ctx, id); err != nil {
		log.Printf("AnimeService.DeleteAnime: Failed to delete anime: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to delete anime", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("AnimeService.DeleteAnime: Successfully deleted anime with ID %d", id)
	return nil
}

// GetAnimesByTitle retrieves animes by title with pagination
func (s *AnimeService) GetAnimesByTitle(ctx context.Context, title string, page, limit int) ([]*models.Anime, int64, error) {
	log.Printf("AnimeService.GetAnimesByTitle: Retrieving animes with title %s", title)

	animes, total, err := s.animeRepo.GetByTitle(ctx, title, page, limit)
	if err != nil {
		log.Printf("AnimeService.GetAnimesByTitle: Failed to retrieve animes: %v", err)
		return nil, 0, err
	}

	// Convert to pointer slice
	result := make([]*models.Anime, len(animes))
	for i := range animes {
		result[i] = &animes[i]
	}

	return result, total, nil
}

// GetAnimesByGenre retrieves animes by genre with pagination
func (s *AnimeService) GetAnimesByGenre(ctx context.Context, genre string, page, limit int) ([]*models.Anime, int64, error) {
	log.Printf("AnimeService.GetAnimesByGenre: Retrieving animes with genre %s", genre)

	animes, total, err := s.animeRepo.GetByGenre(ctx, genre, page, limit)
	if err != nil {
		log.Printf("AnimeService.GetAnimesByGenre: Failed to retrieve animes: %v", err)
		return nil, 0, err
	}

	// Convert to pointer slice
	result := make([]*models.Anime, len(animes))
	for i := range animes {
		result[i] = &animes[i]
	}

	return result, total, nil
}

// GetReviewsForAnime retrieves paginated reviews for an anime
func (s *AnimeService) GetReviewsForAnime(ctx context.Context, animeID uint, page, limit int) ([]*models.Review, int64, error) {
	log.Printf("AnimeService.GetReviewsForAnime: Retrieving reviews for anime with ID %d", animeID)

	// Check if anime exists
	_, err := s.animeRepo.GetByID(ctx, animeID)
	if err != nil {
		log.Printf("AnimeService.GetReviewsForAnime: Failed to retrieve anime: %v", err)
		return nil, 0, err
	}

	reviews, total, err := s.reviewRepo.GetByAnimeID(ctx, animeID, page, limit)
	if err != nil {
		log.Printf("AnimeService.GetReviewsForAnime: Failed to retrieve reviews: %v", err)
		return nil, 0, err
	}

	// Convert to pointer slice
	result := make([]*models.Review, len(reviews))
	for i := range reviews {
		result[i] = &reviews[i]
	}

	return result, total, nil
}

// AddGenresToAnime adds genres to an anime
func (s *AnimeService) AddGenresToAnime(ctx context.Context, animeID uint, genreIDs []uint) error {
	log.Printf("AnimeService.AddGenresToAnime: Adding genres to anime with ID %d", animeID)

	// Check if anime exists
	animeInterface, err := s.animeRepo.GetByID(ctx, animeID)
	if err != nil {
		log.Printf("AnimeService.AddGenresToAnime: Failed to retrieve anime: %v", err)
		return err
	}

	anime, ok := animeInterface.(*models.Anime)
	if !ok {
		log.Printf("AnimeService.AddGenresToAnime: Invalid anime type")
		return errors.NewError(errors.ErrInternalServer, "Invalid anime type", "Type assertion failed", http.StatusInternalServerError)
	}

	// Get genres by IDs
	genres, err := s.genreRepo.GetByIDs(ctx, genreIDs)
	if err != nil {
		log.Printf("AnimeService.AddGenresToAnime: Failed to retrieve genres: %v", err)
		return err
	}

	// Add genres to anime
	anime.Genres = append(anime.Genres, genres...)
	if err := s.animeRepo.Update(ctx, anime); err != nil {
		log.Printf("AnimeService.AddGenresToAnime: Failed to update anime: %v", err)
		return err
	}

	return nil
}

// RemoveGenresFromAnime removes genres from an anime
func (s *AnimeService) RemoveGenresFromAnime(ctx context.Context, animeID uint, genreIDs []uint) error {
	log.Printf("AnimeService.RemoveGenresFromAnime: Removing genres from anime with ID %d", animeID)

	// Check if anime exists
	anime, err := s.animeRepo.GetByID(ctx, animeID)
	if err != nil {
		log.Printf("AnimeService.RemoveGenresFromAnime: Failed to retrieve anime: %v", err)
		return errors.NewError(errors.ErrResourceNotFound, "Anime not found", err.Error(), http.StatusNotFound)
	}

	animeObj, ok := anime.(*models.Anime)
	if !ok {
		log.Printf("AnimeService.RemoveGenresFromAnime: Invalid anime type")
		return errors.NewError(errors.ErrInternalServer, "Invalid anime type", "Type assertion failed", http.StatusInternalServerError)
	}

	// Create a map of genre IDs to remove
	genreMap := make(map[uint]bool)
	for _, id := range genreIDs {
		genreMap[id] = true
	}

	// Filter out genres to remove
	var filteredGenres []models.Genre
	for _, genre := range animeObj.Genres {
		if !genreMap[genre.ID] {
			filteredGenres = append(filteredGenres, genre)
		}
	}

	// Update anime with filtered genres
	animeObj.Genres = filteredGenres
	if err := s.animeRepo.Update(ctx, animeObj); err != nil {
		log.Printf("AnimeService.RemoveGenresFromAnime: Failed to update anime: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to update anime", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("AnimeService.RemoveGenresFromAnime: Successfully removed genres from anime with ID %d", animeID)
	return nil
}

// AddTagsToAnime adds tags to an anime
func (s *AnimeService) AddTagsToAnime(ctx context.Context, animeID uint, tagIDs []uint) error {
	log.Printf("AnimeService.AddTagsToAnime: Adding tags to anime with ID %d", animeID)

	// Check if anime exists
	animeInterface, err := s.animeRepo.GetByID(ctx, animeID)
	if err != nil {
		log.Printf("AnimeService.AddTagsToAnime: Failed to retrieve anime: %v", err)
		return err
	}

	anime, ok := animeInterface.(*models.Anime)
	if !ok {
		log.Printf("AnimeService.AddTagsToAnime: Invalid anime type")
		return errors.NewError(errors.ErrInternalServer, "Invalid anime type", "Type assertion failed", http.StatusInternalServerError)
	}

	// Get tags by IDs
	tags, err := s.tagRepo.GetByIDs(ctx, tagIDs)
	if err != nil {
		log.Printf("AnimeService.AddTagsToAnime: Failed to retrieve tags: %v", err)
		return err
	}

	// Add tags to anime
	anime.Tags = append(anime.Tags, tags...)
	if err := s.animeRepo.Update(ctx, anime); err != nil {
		log.Printf("AnimeService.AddTagsToAnime: Failed to update anime: %v", err)
		return err
	}

	return nil
}

// RemoveTagsFromAnime removes tags from an anime
func (s *AnimeService) RemoveTagsFromAnime(ctx context.Context, animeID uint, tagIDs []uint) error {
	log.Printf("AnimeService.RemoveTagsFromAnime: Removing tags from anime with ID %d", animeID)

	// Check if anime exists
	anime, err := s.animeRepo.GetByID(ctx, animeID)
	if err != nil {
		log.Printf("AnimeService.RemoveTagsFromAnime: Failed to retrieve anime: %v", err)
		return errors.NewError(errors.ErrResourceNotFound, "Anime not found", err.Error(), http.StatusNotFound)
	}

	animeObj, ok := anime.(*models.Anime)
	if !ok {
		log.Printf("AnimeService.RemoveTagsFromAnime: Invalid anime type")
		return errors.NewError(errors.ErrInternalServer, "Invalid anime type", "Type assertion failed", http.StatusInternalServerError)
	}

	// Create a map of tag IDs to remove
	tagMap := make(map[uint]bool)
	for _, id := range tagIDs {
		tagMap[id] = true
	}

	// Filter out tags to remove
	var filteredTags []models.Tag
	for _, tag := range animeObj.Tags {
		if !tagMap[tag.ID] {
			filteredTags = append(filteredTags, tag)
		}
	}

	// Update anime with filtered tags
	animeObj.Tags = filteredTags
	if err := s.animeRepo.Update(ctx, animeObj); err != nil {
		log.Printf("AnimeService.RemoveTagsFromAnime: Failed to update anime: %v", err)
		return errors.NewError(errors.ErrInternalServer, "Failed to update anime", err.Error(), http.StatusInternalServerError)
	}

	log.Printf("AnimeService.RemoveTagsFromAnime: Successfully removed tags from anime with ID %d", animeID)
	return nil
}
