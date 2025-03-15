// pkg/middleware/https_middleware.go
// pkg/middleware/https_middleware.go
// Package middleware provides HTTP middleware functions for the MyAnimeAPI application.
// This file contains middleware for enforcing HTTPS connections.
package middleware

import "net/http"

// RedirectToHTTPS is a middleware that redirects HTTP requests to HTTPS.
// It checks the "X-Forwarded-Proto" header to determine if the request is already using HTTPS.
// If the request is not using HTTPS, it redirects the client to the HTTPS version of the URL.
//
// Example:
//
//		Usage in a router:
//	router := http.NewServeMux()
//	router.HandleFunc("/", myHandler)
//	http.ListenAndServe(":80", middleware.RedirectToHTTPS(router))
//
//	 	When a client accesses "http://example.com", they will be redirected to "https://example.com".
func RedirectToHTTPS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Forwarded-Proto") != "https" {
			http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}
