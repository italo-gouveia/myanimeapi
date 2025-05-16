// api/middleware/https_middleware.go
// Package middleware provides HTTP middleware functions for the MyAnimeAPI application.
// This file contains middleware for enforcing HTTPS connections.
package middleware

import (
	"net/http"
	"strings"

	"myanimeapi/internal/errors"
	"myanimeapi/internal/logger"

	"github.com/sirupsen/logrus"
)

// allowedDomains is a list of trusted domains to which HTTPS redirection is allowed.
var allowedDomains = []string{"localhost", "myanimeapi.onrender.com", "myanimeapi.com"}

// isAllowedDomain checks if the provided host is in the list of allowed domains.
// It ensures that the host matches one of the trusted domains to prevent open redirect vulnerabilities.
func isAllowedDomain(host string) bool {
	for _, domain := range allowedDomains {
		if strings.EqualFold(host, domain) {
			return true
		}
	}
	return false
}

// HTTPSRedirectMiddleware is a middleware that redirects HTTP requests to HTTPS.
// It checks if the request is already using HTTPS by inspecting the "X-Forwarded-Proto" header.
// If the request is not using HTTPS, it redirects the client to the HTTPS version of the URL.
// The middleware ensures that redirection only occurs for trusted domains listed in `allowedDomains`.
//
// Example:
//
//	Usage in a router:
//	router := http.NewServeMux()
//	router.HandleFunc("/", myHandler)
//	http.ListenAndServe(":80", middleware.HTTPSRedirectMiddleware(router))
//
//	When a client accesses "http://example.com", they will be redirected to "https://example.com".
func HTTPSRedirectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.Get()

		// Check if the request is not HTTPS
		if r.Header.Get("X-Forwarded-Proto") != "https" {
			host := r.Host
			if isAllowedDomain(host) {
				// Hardcode the HTTPS URL for trusted domains
				redirectURL := "https://" + host
				log.WithFields(logrus.Fields{
					"host":         host,
					"redirect_url": redirectURL,
				}).Info("Redirecting HTTP request to HTTPS")
				http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
				return
			}
			// If the domain is not allowed, return a 403 Forbidden error
			log.WithField("host", host).Warn("Domain not allowed for HTTPS redirection")
			errors.WriteErrorResponse(w, http.StatusForbidden, errors.ErrForbidden, "Forbidden: Domain not allowed for HTTPS redirection", "The requested domain is not allowed for HTTPS redirection.")
			return
		}
		// If the request is already HTTPS, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
