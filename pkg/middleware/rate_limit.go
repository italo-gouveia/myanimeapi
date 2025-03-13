// pkg/middleware/rate_limit.go
package middleware

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// RateLimiter is a struct to hold rate-limiting data
type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*clientInfo
	clock   func() time.Time // Custom clock for testing
}

type clientInfo struct {
	count    int       // Number of requests made by the client
	lastSeen time.Time // Last time the client made a request
}

// NewRateLimiter creates a new RateLimiter instance
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*clientInfo),
		clock:   time.Now, // Default to real time
	}
}

// SetClock sets a custom clock for testing
func (rl *RateLimiter) SetClock(clock func() time.Time) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.clock = clock
}

// RateLimitMiddleware is a middleware that limits the number of requests per IP
func (rl *RateLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.mu.Lock()
		defer rl.mu.Unlock()

		// Get the client's IP address
		ip := r.RemoteAddr

		// Initialize client info if it doesn't exist
		if _, exists := rl.clients[ip]; !exists {
			rl.clients[ip] = &clientInfo{}
		}

		// Determine the rate limit based on the request path
		var limit int
		var window time.Duration

		switch {
		case strings.HasPrefix(r.URL.Path, "/auth/authenticate") || strings.HasPrefix(r.URL.Path, "/auth/register"):
			// Stricter rate limit for authentication endpoints
			limit = 5
			window = time.Minute
		default:
			// Default rate limit for other endpoints
			limit = 50
			window = time.Minute
		}

		// Reset the count if the time window has passed
		if rl.clock().Sub(rl.clients[ip].lastSeen) > window {
			rl.clients[ip].count = 0
		}

		// Check if the request count exceeds the limit
		if rl.clients[ip].count >= limit {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			_, err := w.Write([]byte(`{"error": "Rate limit exceeded. Please try again later."}`))
			if err != nil {
				log.Printf("Failed to write response: %v", err)
			}
			return
		}

		// Increment the request count and update the last seen time
		rl.clients[ip].count++
		rl.clients[ip].lastSeen = rl.clock()

		// Log the current state for debugging
		log.Printf("IP: %s, Path: %s, Count: %d, LastSeen: %v", ip, r.URL.Path, rl.clients[ip].count, rl.clients[ip].lastSeen)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
