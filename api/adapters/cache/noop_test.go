// api/adapters/cache/noop_test.go
package cache_test

import (
	"context"
	"testing"
	"time"

	"myanimeapi/api/adapters/cache"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoopCache_GetReturnsCacheMiss(t *testing.T) {
	c := cache.NewNoopCache()
	data, err := c.Get(context.Background(), "any-key")
	assert.Nil(t, data)
	assert.ErrorIs(t, err, cache.ErrCacheMiss)
}

func TestNoopCache_SetIsNoop(t *testing.T) {
	c := cache.NewNoopCache()
	err := c.Set(context.Background(), "key", []byte("value"), time.Minute)
	require.NoError(t, err)

	// After Set, Get still returns a cache miss (nothing is stored)
	data, err := c.Get(context.Background(), "key")
	assert.Nil(t, data)
	assert.ErrorIs(t, err, cache.ErrCacheMiss)
}

func TestNoopCache_DeleteIsNoop(t *testing.T) {
	c := cache.NewNoopCache()
	assert.NoError(t, c.Delete(context.Background(), "k1", "k2"))
	// zero keys — must also be fine
	assert.NoError(t, c.Delete(context.Background()))
}

func TestNoopCache_DeleteByPatternIsNoop(t *testing.T) {
	c := cache.NewNoopCache()
	assert.NoError(t, c.DeleteByPattern(context.Background(), "animes:*"))
}

func TestNoopCache_PingAlwaysSucceeds(t *testing.T) {
	c := cache.NewNoopCache()
	assert.NoError(t, c.Ping(context.Background()))
}

func TestNoopCache_ImplementsCacheInterface(t *testing.T) {
	// Compile-time assertion: *NoopCache must satisfy CacheInterface.
	var _ cache.CacheInterface = cache.NewNoopCache()
}
