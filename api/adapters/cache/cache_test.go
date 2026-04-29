// api/adapters/cache/cache_test.go
package cache_test

import (
	"testing"
	"time"

	"myanimeapi/api/adapters/cache"

	"github.com/stretchr/testify/assert"
)

// ---- TTL constants ----------------------------------------------------------

func TestTTLConstants(t *testing.T) {
	assert.Equal(t, 10*time.Minute, cache.TTLAnime)
	assert.Equal(t, 5*time.Minute, cache.TTLList)
	assert.Equal(t, 30*time.Minute, cache.TTLGenre)
	assert.Equal(t, 30*time.Minute, cache.TTLTag)
}

// ---- Key builders -----------------------------------------------------------

func TestKeyAnime(t *testing.T) {
	assert.Equal(t, "anime:1", cache.KeyAnime(1))
	assert.Equal(t, "anime:99", cache.KeyAnime(99))
}

func TestKeyAnimeList(t *testing.T) {
	assert.Equal(t, "animes:p1:l20", cache.KeyAnimeList(1, 20))
	assert.Equal(t, "animes:p3:l10", cache.KeyAnimeList(3, 10))
}

func TestKeyAnimeSearch(t *testing.T) {
	assert.Equal(t, "animes:search:naruto:p1:l10", cache.KeyAnimeSearch("naruto", 1, 10))
}

func TestKeyGenre(t *testing.T) {
	assert.Equal(t, "genre:5", cache.KeyGenre(5))
}

func TestKeyGenreList(t *testing.T) {
	assert.Equal(t, "genres:p2:l50", cache.KeyGenreList(2, 50))
}

func TestKeyTag(t *testing.T) {
	assert.Equal(t, "tag:7", cache.KeyTag(7))
}

func TestKeyTagList(t *testing.T) {
	assert.Equal(t, "tags:p1:l100", cache.KeyTagList(1, 100))
}

// ---- Pattern constants -------------------------------------------------------

func TestPatternConstants(t *testing.T) {
	assert.Equal(t, "animes:*", cache.PatternAllAnimes)
	assert.Equal(t, "genres:*", cache.PatternAllGenres)
	assert.Equal(t, "tags:*", cache.PatternAllTags)
}

// ---- ErrCacheMiss -----------------------------------------------------------

func TestErrCacheMiss(t *testing.T) {
	assert.EqualError(t, cache.ErrCacheMiss, "cache miss")
}
