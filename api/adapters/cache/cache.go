// api/adapters/cache/cache.go
// Package cache defines the CacheInterface output port and shared key-building
// helpers used by all services that participate in caching.
//
// Two implementations are provided:
//   - RedisCache  — production, backed by Redis via go-redis/v9
//   - NoopCache   — silent no-op used when REDIS_URL is not configured or
//     when tests run without a real Redis instance
//
// Usage in a service:
//
//	type AnimeService struct {
//	    repo  repositories.AnimeRepository
//	    cache cache.CacheInterface
//	}
//
//	func (s *AnimeService) GetAnimeByID(ctx context.Context, id uint) (*models.Anime, error) {
//	    key  := cache.KeyAnime(id)
//	    data, err := s.cache.Get(ctx, key)
//	    if err == nil {
//	        // cache hit — unmarshal and return
//	    }
//	    // cache miss — fetch from DB, then Set
//	}
package cache

import (
	"context"
	"fmt"
	"time"
)

// TTL constants for each data category.
const (
	TTLAnime  = 10 * time.Minute
	TTLList   = 5 * time.Minute
	TTLGenre  = 30 * time.Minute
	TTLTag    = 30 * time.Minute
)

// ErrCacheMiss is returned by Get when the key does not exist in the cache.
var ErrCacheMiss = fmt.Errorf("cache miss")

// CacheInterface is the output port for caching operations.
// Any adapter (Redis, Memcached, in-memory) must implement this interface.
type CacheInterface interface {
	// Get retrieves a cached value by key.
	// Returns ErrCacheMiss when the key does not exist.
	Get(ctx context.Context, key string) ([]byte, error)

	// Set stores a value with the given TTL.
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Delete removes one or more keys.
	Delete(ctx context.Context, keys ...string) error

	// DeleteByPattern removes all keys matching a glob-style pattern.
	// Example: "anime:*" removes all anime cache entries.
	DeleteByPattern(ctx context.Context, pattern string) error

	// Ping checks connectivity to the cache backend.
	Ping(ctx context.Context) error
}

// ---- Key builders -------------------------------------------------------
// Centralise key construction so services never hard-code strings.

// KeyAnime returns the cache key for a single anime.
func KeyAnime(id uint) string { return fmt.Sprintf("anime:%d", id) }

// KeyAnimeList returns the cache key for a paginated anime list.
func KeyAnimeList(page, limit int) string {
	return fmt.Sprintf("animes:p%d:l%d", page, limit)
}

// KeyAnimeSearch returns the cache key for a title-search result page.
func KeyAnimeSearch(title string, page, limit int) string {
	return fmt.Sprintf("animes:search:%s:p%d:l%d", title, page, limit)
}

// KeyAnimeByGenre returns the cache key for a genre-filtered anime result page.
func KeyAnimeByGenre(genre string, page, limit int) string {
	return fmt.Sprintf("animes:genre:%s:p%d:l%d", genre, page, limit)
}

// KeyGenre returns the cache key for a single genre.
func KeyGenre(id uint) string { return fmt.Sprintf("genre:%d", id) }

// KeyGenreList returns the cache key for a paginated genre list.
func KeyGenreList(page, limit int) string {
	return fmt.Sprintf("genres:p%d:l%d", page, limit)
}

// KeyTag returns the cache key for a single tag.
func KeyTag(id uint) string { return fmt.Sprintf("tag:%d", id) }

// KeyTagList returns the cache key for a paginated tag list.
func KeyTagList(page, limit int) string {
	return fmt.Sprintf("tags:p%d:l%d", page, limit)
}

// Pattern constants used for bulk invalidation.
const (
	PatternAllAnimeDetails = "anime:*"  // individual detail keys  (anime:<id>)
	PatternAllAnimes       = "animes:*" // paginated list keys      (animes:p…)
	PatternAllGenres       = "genres:*"
	PatternAllTags         = "tags:*"
)
