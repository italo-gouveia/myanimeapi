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
}

type clientInfo struct {
	count    int       // Number of requests made by the client
	lastSeen time.Time // Last time the client made a request
}

// NewRateLimiter creates a new RateLimiter instance
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*clientInfo),
	}
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

		log.Printf("IP: %s, Path: %s, Count: %d, LastSeen: %v", ip, r.URL.Path, rl.clients[ip].count, rl.clients[ip].lastSeen)

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
		if time.Since(rl.clients[ip].lastSeen) > window {
			rl.clients[ip].count = 0
		}

		// Check if the request count exceeds the limit
		if rl.clients[ip].count >= limit {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "Rate limit exceeded. Please try again later."}`))
			return
		}

		// Increment the request count and update the last seen time
		rl.clients[ip].count++
		rl.clients[ip].lastSeen = time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
