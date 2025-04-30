package repositories

import (
	"context"
	"myanimeapi/api/models"
)

// Repository defines the common repository methods for all entities
type Repository interface {
	// GetByID retrieves an entity by its ID
	GetByID(ctx context.Context, id uint) (interface{}, error)

	// GetAll retrieves all entities with optional pagination
	GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error)

	// Create creates a new entity
	Create(ctx context.Context, entity interface{}) error

	// Update updates an existing entity
	Update(ctx context.Context, entity interface{}) error

	// Delete deletes an entity by its ID
	Delete(ctx context.Context, id uint) error
}

// UserRepository defines the methods for user-related database operations
type UserRepository interface {
	Repository

	// GetByUsername retrieves a user by username
	GetByUsername(ctx context.Context, username string) (*models.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

// AnimeRepository defines the methods for anime-related database operations
type AnimeRepository interface {
	Repository

	// GetByTitle retrieves animes by title (partial match)
	GetByTitle(ctx context.Context, title string, page, limit int) ([]models.Anime, int64, error)

	// GetByGenre retrieves animes by genre
	GetByGenre(ctx context.Context, genre string, page, limit int) ([]models.Anime, int64, error)

	// GetReviewsForAnime retrieves reviews for a specific anime
	GetReviewsForAnime(ctx context.Context, animeID uint, page, limit int) ([]models.Review, int64, error)
}

// ReviewRepository defines the methods for review-related database operations
type ReviewRepository interface {
	Repository

	// GetByUserID retrieves reviews by user ID
	GetByUserID(ctx context.Context, userID uint, page, limit int) ([]models.Review, int64, error)

	// GetByAnimeID retrieves reviews by anime ID
	GetByAnimeID(ctx context.Context, animeID uint, page, limit int) ([]models.Review, int64, error)

	// GetByUserAndAnime retrieves a review by user ID and anime ID
	GetByUserAndAnime(ctx context.Context, userID, animeID uint) (*models.Review, error)
}

// FavoriteRepository defines the methods for favorite-related database operations
type FavoriteRepository interface {
	Repository

	// GetByUserID retrieves favorites by user ID
	GetByUserID(ctx context.Context, userID uint) ([]models.Favorite, error)

	// GetByUserAndAnime retrieves a favorite by user ID and anime ID
	GetByUserAndAnime(ctx context.Context, userID, animeID uint) (*models.Favorite, error)

	// DeleteByUserAndAnime deletes a favorite by user ID and anime ID
	DeleteByUserAndAnime(ctx context.Context, userID, animeID uint) error
}
