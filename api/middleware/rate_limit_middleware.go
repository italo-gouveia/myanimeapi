// api/middleware/rate_limit_middleware.go
// Package middleware provides HTTP middleware utilities for handling requests.
// This file defines a rate-limiting middleware that restricts the number of requests a client can make within a specified time window.
package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"
)

// RateLimiter is a struct that holds rate-limiting data for clients.
// It tracks the number of requests made by each client within a specified time window.
type RateLimiter struct {
	mu      sync.Mutex             // Mutex to ensure thread-safe access to the clients map
	clients map[string]*clientInfo // Map of client IPs to their request information
	clock   func() time.Time       // Custom clock function for testing purposes
}

// clientInfo holds information about a client's request activity.
type clientInfo struct {
	count    int       // Number of requests made by the client within the current time window
	lastSeen time.Time // Timestamp of the client's last request
}

// NewRateLimiter creates and returns a new RateLimiter instance.
// It initializes the clients map and sets the default clock to the current time.
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*clientInfo),
		clock:   time.Now, // Default to real time
	}
}

// SetClock sets a custom clock function for testing purposes.
// This allows the rate limiter to use a simulated time source instead of the system clock.
func (rl *RateLimiter) SetClock(clock func() time.Time) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.clock = clock
}

// RateLimitMiddleware is an HTTP middleware that enforces rate limits based on the client's IP address.
// It limits the number of requests a client can make within a specified time window.
// The rate limit is stricter for authentication-related endpoints (e.g., /auth/authenticate, /auth/register).
//
// Example usage:
//
//	rl := NewRateLimiter()
//	http.Handle("/path", rl.RateLimitMiddleware(myHandler))
//
// This will enforce rate limits for requests to "/path".
func (rl *RateLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.Get()
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
			log.WithFields(map[string]interface{}{
				"ip":         ip,
				"path":       r.URL.Path,
				"count":      rl.clients[ip].count,
				"limit":      limit,
				"window":     window.String(),
				"last_seen":  rl.clients[ip].lastSeen,
				"user_agent": r.UserAgent(),
				"request_id": r.Context().Value(RequestIDContextKey),
			}).Warning("Rate limit exceeded")

			errors.WriteErrorResponse(w, http.StatusTooManyRequests, errors.ErrTooManyRequests, "Rate limit exceeded", "Please try again later.", nil)
			return
		}

		// Increment the request count and update the last seen time
		rl.clients[ip].count++
		rl.clients[ip].lastSeen = rl.clock()

		// Log the current state for debugging
		log.WithFields(map[string]interface{}{
			"ip":         ip,
			"path":       r.URL.Path,
			"count":      rl.clients[ip].count,
			"limit":      limit,
			"window":     window.String(),
			"last_seen":  rl.clients[ip].lastSeen,
			"user_agent": r.UserAgent(),
			"request_id": r.Context().Value(RequestIDContextKey),
		}).Debug("Rate limit check passed")

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
