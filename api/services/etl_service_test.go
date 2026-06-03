package services

import (
	"context"
	"errors"
	"myanimeapi/api/adapters/cache"
	"myanimeapi/internal/logger"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestLogger() *logger.Logger { return logger.New() }

// spyCache records every DeleteByPattern call so tests can assert on
// which patterns were invalidated, without needing a real Redis instance.
type spyCache struct {
	cache.NoopCache
	deleted []string
	failOn  string // if set, DeleteByPattern returns an error for this pattern
}

func (s *spyCache) DeleteByPattern(_ context.Context, pattern string) error {
	s.deleted = append(s.deleted, pattern)
	if s.failOn == pattern {
		return errors.New("redis: connection refused")
	}
	return nil
}

// TestInvalidateCache_DeletesAllPatterns verifies that invalidateCache issues
// DeleteByPattern for every expected key namespace after a sync.
func TestInvalidateCache_DeletesAllPatterns(t *testing.T) {
	spy := &spyCache{}
	svc := &ETLService{
		logger: newTestLogger(),
		cache:  spy,
	}

	svc.invalidateCache(context.Background())

	expected := []string{
		cache.PatternAllAnimeDetails,
		cache.PatternAllAnimes,
		cache.PatternAllGenres,
		cache.PatternAllTags,
	}
	assert.Equal(t, expected, spy.deleted,
		"invalidateCache must delete all four cache namespaces in order")
}

// TestInvalidateCache_ToleratesPartialFailure verifies that a Redis error on
// one pattern does not prevent the remaining patterns from being invalidated.
func TestInvalidateCache_ToleratesPartialFailure(t *testing.T) {
	spy := &spyCache{failOn: cache.PatternAllAnimes}
	svc := &ETLService{
		logger: newTestLogger(),
		cache:  spy,
	}

	// Must not panic or return an error — invalidation is best-effort.
	require.NotPanics(t, func() {
		svc.invalidateCache(context.Background())
	})

	// All four patterns were still attempted despite the mid-run error.
	assert.Len(t, spy.deleted, 4,
		"all patterns must be attempted even when one DeleteByPattern fails")
}

// TestInvalidateCache_RespectsContextCancellation verifies that passing an
// already-cancelled context does not cause a panic (the logger must handle it).
func TestInvalidateCache_RespectsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	// Let the context expire before the call.
	time.Sleep(1 * time.Millisecond)

	spy := &spyCache{}
	svc := &ETLService{
		logger: newTestLogger(),
		cache:  spy,
	}

	require.NotPanics(t, func() {
		svc.invalidateCache(ctx)
	})
}
