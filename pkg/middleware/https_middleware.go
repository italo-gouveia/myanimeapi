// pkg/middleware/https_middleware.go
// Package middleware provides HTTP middleware functions for the MyAnimeAPI application.
// This file contains middleware for enforcing HTTPS connections.
package middleware

// allowedDomains is a list of trusted domains to which HTTPS redirection is allowed.
//var allowedDomains = []string{"localhost", "myanimeapi.onrender.com", "myanimeapi.com"}

// isAllowedDomain checks if the provided host is in the list of allowed domains.
// It ensures that the host matches one of the trusted domains to prevent open redirect vulnerabilities.
/*func isAllowedDomain(host string) bool {
	for _, domain := range allowedDomains {
		if strings.EqualFold(host, domain) {
			return true
		}
	}
	return false
}*/

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
/*func HTTPSRedirectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request is not HTTPS
		if r.Header.Get("X-Forwarded-Proto") != "https" {
			host := r.Host
			if isAllowedDomain(host) {
				// Hardcode the HTTPS URL for trusted domains
				redirectURL := "https://" + host
				http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
				return
			}
			// If the domain is not allowed, return a 403 Forbidden error
			http.Error(w, "Forbidden: Domain not allowed for HTTPS redirection", http.StatusForbidden)
			return
		}
		// If the request is already HTTPS, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
*/
