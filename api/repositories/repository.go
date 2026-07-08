// Package repositories defines the data access layer interfaces for the MyAnimeAPI application.
// It provides a set of interfaces that define the contract for database operations,
// ensuring a clean separation between the business logic and data access layers.
//
// The package includes:
//   - Base Repository interface with common CRUD operations
//   - Specialized repository interfaces for each entity
//   - Methods for specific business requirements
//
// Each repository interface extends the base Repository interface and adds
// entity-specific methods as needed.
package repositories

import (
	"context"
	"myanimeapi/api/models"
)

// Repository defines the common repository methods for all entities.
// It provides a standard set of CRUD operations that can be implemented
// for any entity in the system.
//
// The interface includes:
//   - GetByID: Retrieve a single entity by its ID
//   - GetAll: Retrieve all entities with pagination support
//   - Create: Create a new entity
//   - Update: Update an existing entity
//   - Delete: Delete an entity by its ID
type Repository interface {
	// GetByID retrieves an entity by its ID.
	// Returns the entity if found, or an error if not found or if an error occurs.
	GetByID(ctx context.Context, id uint) (interface{}, error)

	// GetAll retrieves all entities with optional pagination.
	// Returns a slice of entities, the total count, and any error that occurred.
	// The page and limit parameters control pagination:
	//   - page: The page number (1-based)
	//   - limit: The number of items per page
	GetAll(ctx context.Context, page, limit int) ([]interface{}, int64, error)

	// Create creates a new entity.
	// Returns an error if the creation fails or if the entity already exists.
	Create(ctx context.Context, entity interface{}) error

	// Update updates an existing entity.
	// Returns an error if the update fails or if the entity doesn't exist.
	Update(ctx context.Context, entity interface{}) error

	// Delete deletes an entity by its ID.
	// Returns an error if the deletion fails or if the entity doesn't exist.
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

	// GetWithFilters retrieves animes with optional filtering, sorting, and pagination.
	// Filters are applied via subqueries to avoid JOIN-induced duplicate rows.
	GetWithFilters(ctx context.Context, filter models.AnimeFilter, page, limit int) ([]interface{}, int64, error)
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

// CharacterRepository defines methods for character-related database operations.
type CharacterRepository interface {
	Repository

	// GetByAnimeID retrieves all characters for a given anime (paginated).
	GetByAnimeID(ctx context.Context, animeID uint, page, limit int) ([]models.Character, int64, error)

	// GetByName retrieves characters whose name contains the query string.
	GetByName(ctx context.Context, name string, page, limit int) ([]models.Character, int64, error)

	// AddToAnime associates a set of characters with an anime.
	AddToAnime(ctx context.Context, animeID uint, characterIDs []uint) error

	// RemoveFromAnime dissociates a character from an anime.
	RemoveFromAnime(ctx context.Context, animeID, characterID uint) error
}

// WatchlistRepository defines methods for watchlist database operations
type WatchlistRepository interface {
	Repository

	// GetByUserID retrieves watchlist entries by user ID
	GetByUserID(ctx context.Context, userID uint) ([]models.WatchlistEntry, error)

	// GetByUserAndAnime retrieves a watchlist entry by user ID and anime ID
	GetByUserAndAnime(ctx context.Context, userID, animeID uint) (*models.WatchlistEntry, error)

	// UpsertByUserAndAnime creates or updates a watchlist entry for a user/anime pair
	UpsertByUserAndAnime(ctx context.Context, userID, animeID uint, status models.WatchlistStatus) (*models.WatchlistEntry, error)

	// DeleteByUserAndAnime deletes a watchlist entry by user ID and anime ID
	DeleteByUserAndAnime(ctx context.Context, userID, animeID uint) error
}
