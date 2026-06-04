// api/middleware/https_middleware.go
// Package middleware provides HTTP middleware utilities for handling requests in the MyAnimeAPI application.
// It includes middleware functions for error handling, authentication, logging, and more.
//
// Middleware functions in this package are designed to wrap HTTP handlers and provide additional
// functionality such as panic recovery, request validation, and access control.
package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"
)

// validHostPattern matches valid hostname:port combinations to prevent open redirect attacks.
var validHostPattern = regexp.MustCompile(`^[a-zA-Z0-9.\-]+(:\d{1,5})?$`)

// safeRedirectHost returns a validated host from the request, falling back to "localhost"
// if the host contains characters that could enable an open redirect attack.
func safeRedirectHost(r *http.Request) string {
	host := r.Host
	if host == "" {
		host = r.URL.Host
	}
	if !validHostPattern.MatchString(host) {
		return "localhost"
	}
	return host
}

// HTTPSMiddleware ensures that all requests are handled securely over HTTPS.
// It redirects HTTP requests to HTTPS and sets security headers.
func HTTPSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.Get()
		logFields := map[string]interface{}{
			"method":      r.Method,
			"path":        r.URL.Path,
			"remote_addr": r.RemoteAddr,
			"user_agent":  r.UserAgent(),
		}

		// Check if the request is already HTTPS
		if r.TLS == nil {
			// Get the host from the request, validated against a strict pattern to
			// prevent open redirect attacks via a crafted Host header (gosecurity:S5146).
			host := safeRedirectHost(r)

			// Construct the HTTPS URL using RequestURI to safely include path and query.
			// r.URL.RequestURI() returns only the path and query string (no scheme or host).
			httpsURL := "https://" + host + r.URL.RequestURI()

			// Log the redirect
			log.WithFields(logFields).Info("Redirecting HTTP request to HTTPS")

			// Redirect to HTTPS
			http.Redirect(w, r, httpsURL, http.StatusPermanentRedirect)
			return
		}

		// Set security headers
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Check for secure cookies
		if strings.HasPrefix(r.URL.Path, "/auth") {
			// Ensure cookies are secure for auth endpoints
			for _, cookie := range r.Cookies() {
				if !cookie.Secure {
					// Log insecure cookie
					log.WithFields(logFields).WithField("cookie_name", cookie.Name).Warning("Insecure cookie detected on auth endpoint")

					// Return error response
					errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput,
						"Security Error",
						"Invalid cookie configuration",
						logFields)
					return
				}
			}
		}

		// Continue with the request
		next.ServeHTTP(w, r)
	})
}

// IsSecureRequest checks if the request is coming over HTTPS
func IsSecureRequest(r *http.Request) bool {
	return r.TLS != nil
}

// GetProtocol returns the protocol used in the request
func GetProtocol(r *http.Request) string {
	if IsSecureRequest(r) {
		return "https"
	}
	return "http"
}
