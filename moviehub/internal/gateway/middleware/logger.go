package middleware

import (
	"moviesapi/internal/shared"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// responseWriter is a custom wrapper around http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logger returns a middleware that logs incoming HTTP requests using zap
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // default status code
		}

		// Process request
		next.ServeHTTP(rw, r)

		elapsed := time.Since(start)

		// Log request details
		shared.Log.Info("incoming request",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.String("client_ip", r.RemoteAddr),
			zap.Int("status", rw.statusCode),
			zap.Duration("elapsed", elapsed),
		)
	})
}
