// api/adapters/cache/noop.go
// NoopCache is a silent no-op implementation of CacheInterface used when
// Redis is not configured or not available. Every Get returns ErrCacheMiss
// so services fall through to the database as if caching were not present.
package cache

import (
	"context"
	"time"
)

// NoopCache silently discards all writes and returns ErrCacheMiss on reads.
// It is the default when REDIS_URL is not set, enabling graceful degradation.
type NoopCache struct{}

// NewNoopCache returns a NoopCache instance.
func NewNoopCache() *NoopCache { return &NoopCache{} }

func (n *NoopCache) Get(_ context.Context, _ string) ([]byte, error) {
	return nil, ErrCacheMiss
}

func (n *NoopCache) Set(_ context.Context, _ string, _ []byte, _ time.Duration) error {
	return nil
}

func (n *NoopCache) Delete(_ context.Context, _ ...string) error { return nil }

func (n *NoopCache) DeleteByPattern(_ context.Context, _ string) error { return nil }

func (n *NoopCache) Ping(_ context.Context) error { return nil }
