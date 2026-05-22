// api/middleware/metrics_middleware.go
// Package middleware provides HTTP middleware utilities for handling requests.
// This file exposes Prometheus metrics for every HTTP request handled by the API:
//
//   - http_requests_total          — counter  {method, path, status}
//   - http_request_duration_seconds — histogram {method, path}
//   - http_requests_in_flight       — gauge
package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// httpRequestsTotal counts every completed request.
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "myanimeapi",
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests completed.",
		},
		[]string{"method", "path", "status"},
	)

	// httpRequestDuration records latency per request.
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "myanimeapi",
			Name:      "http_request_duration_seconds",
			Help:      "Duration of HTTP requests in seconds.",
			Buckets:   prometheus.DefBuckets, // 5ms … 10s
		},
		[]string{"method", "path"},
	)

	// httpRequestsInFlight tracks concurrent requests.
	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "myanimeapi",
			Name:      "http_requests_in_flight",
			Help:      "Current number of HTTP requests being served.",
		},
	)
)

// metricsResponseWriter wraps http.ResponseWriter to capture the status code.
type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (mrw *metricsResponseWriter) WriteHeader(code int) {
	mrw.statusCode = code
	mrw.ResponseWriter.WriteHeader(code)
}

// MetricsMiddleware instruments every HTTP request with Prometheus metrics.
// It records request count (labelled by method, path and status code),
// request duration (labelled by method and path) and the number of
// requests currently in flight.
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		mrw := &metricsResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		next.ServeHTTP(mrw, r)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(mrw.statusCode)
		path := sanitizePath(r.URL.Path)

		httpRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
	})
}

// sanitizePath replaces high-cardinality path segments (numeric IDs) with ":id"
// so that Prometheus label sets stay bounded and easy to query.
//
// Example:  /v1/animes/42/genres  →  /v1/animes/:id/genres
func sanitizePath(path string) string {
	var result []byte
	i := 0
	for i < len(path) {
		if path[i] == '/' {
			result = append(result, '/')
			i++
			// Collect the next segment.
			j := i
			isNum := true
			for j < len(path) && path[j] != '/' {
				if path[j] < '0' || path[j] > '9' {
					isNum = false
				}
				j++
			}
			if isNum && j > i {
				result = append(result, []byte(":id")...)
			} else {
				result = append(result, []byte(path[i:j])...)
			}
			i = j
		} else {
			result = append(result, path[i])
			i++
		}
	}
	return string(result)
}
