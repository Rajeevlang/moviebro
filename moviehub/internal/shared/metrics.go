package shared

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTPRequestsTotal counts total HTTP requests partitioned by service, method, path, and status code.
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed, partitioned by service, method, path, and status code.",
		},
		[]string{"service", "method", "path", "status"},
	)

	// HTTPRequestDurationSeconds tracks the latency of HTTP requests in seconds.
	HTTPRequestDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of HTTP request latency in seconds.",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"service", "method", "path", "status"},
	)

	// HTTPRequestsInFlight tracks the current number of active in-flight requests.
	HTTPRequestsInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being served concurrently.",
		},
		[]string{"service"},
	)

	// DBQueryDurationSeconds tracks database query latency (PostgreSQL / MongoDB / Redis).
	DBQueryDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Latency of database queries in seconds.",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
		},
		[]string{"service", "database", "operation"},
	)
)

// PrometheusGinMiddleware returns a Gin middleware that records RED metrics for Gin applications.
func PrometheusGinMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" || c.Request.URL.Path == "/healthcheck" {
			c.Next()
			return
		}

		start := time.Now()
		HTTPRequestsInFlight.WithLabelValues(serviceName).Inc()
		defer HTTPRequestsInFlight.WithLabelValues(serviceName).Dec()

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		duration := time.Since(start).Seconds()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		HTTPRequestsTotal.WithLabelValues(serviceName, c.Request.Method, path, status).Inc()
		HTTPRequestDurationSeconds.WithLabelValues(serviceName, c.Request.Method, path, status).Observe(duration)
	}
}

// PrometheusHTTPMiddleware wraps standard net/http handlers to record RED metrics.
func PrometheusHTTPMiddleware(serviceName string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" || r.URL.Path == "/healthcheck" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		HTTPRequestsInFlight.WithLabelValues(serviceName).Inc()
		defer HTTPRequestsInFlight.WithLabelValues(serviceName).Dec()

		rw := &statusResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(rw.status)

		HTTPRequestsTotal.WithLabelValues(serviceName, r.Method, r.URL.Path, status).Inc()
		HTTPRequestDurationSeconds.WithLabelValues(serviceName, r.Method, r.URL.Path, status).Observe(duration)
	})
}

// MetricsHandler returns the Prometheus HTTP metrics handler.
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *statusResponseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
