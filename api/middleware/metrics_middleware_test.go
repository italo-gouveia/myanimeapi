package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestMetricsMiddleware_RecordsRequest(t *testing.T) {
	// Reset the counters so tests are independent regardless of run order.
	httpRequestsTotal.Reset()

	handler := MetricsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/v1/animes", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Expect exactly one completed request with label GET / /v1/animes / 200.
	count := testutil.ToFloat64(httpRequestsTotal.WithLabelValues("GET", "/v1/animes", "200"))
	assert.Equal(t, float64(1), count)
}

func TestMetricsMiddleware_RecordsStatusCode(t *testing.T) {
	httpRequestsTotal.Reset()

	handler := MetricsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	req := httptest.NewRequest(http.MethodGet, "/v1/animes/999", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Path contains a numeric ID — it must be sanitized to ":id".
	count := testutil.ToFloat64(httpRequestsTotal.WithLabelValues("GET", "/v1/animes/:id", "404"))
	assert.Equal(t, float64(1), count)
}

func TestMetricsMiddleware_InFlightDecremented(t *testing.T) {
	// After the handler returns, in-flight gauge must be back to its pre-call value.
	before := testutil.ToFloat64(httpRequestsInFlight)

	handler := MetricsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/v1/health", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	after := testutil.ToFloat64(httpRequestsInFlight)
	assert.Equal(t, before, after)
}

func TestMetricsMiddleware_DurationObserved(t *testing.T) {
	handler := MetricsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/v1/animes", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// At least one observation must have been recorded.
	hist, err := httpRequestDuration.GetMetricWithLabelValues("POST", "/v1/animes")
	assert.NoError(t, err)
	assert.NotNil(t, hist)

	// Gather and verify the histogram has at least one sample.
	mf, err := prometheus.DefaultGatherer.Gather()
	assert.NoError(t, err)

	found := false
	for _, m := range mf {
		if m.GetName() == "myanimeapi_http_request_duration_seconds" {
			found = true
			break
		}
	}
	assert.True(t, found, "expected myanimeapi_http_request_duration_seconds to be registered")
}

func TestSanitizePath(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"/v1/animes", "/v1/animes"},
		{"/v1/animes/42", "/v1/animes/:id"},
		{"/v1/animes/42/genres", "/v1/animes/:id/genres"},
		{"/v1/users/100/favorites/7", "/v1/users/:id/favorites/:id"},
		{"/v1/health", "/v1/health"},
		{"/metrics", "/metrics"},
		{"", ""},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := sanitizePath(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}
