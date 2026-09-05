package middleware

import (
	"moviesapi/internal/shared"
	"net"
	"net/http"
	"sync"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// rateLimiter struct holds the limiters and a mutex for thread-safe access
type rateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	r        rate.Limit
	b        int
}

func newRateLimiter(r rate.Limit, b int) *rateLimiter {
	return &rateLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		b:        b,
	}
}

// getLimiter returns a rate.Limiter for a given IP address
func (rl *rateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.r, rl.b)
		rl.limiters[ip] = limiter
	}

	return limiter
}

// RateLimit is a middleware that rate limits incoming requests by IP
func RateLimit(requestsPerSecond rate.Limit, burst int) func(http.Handler) http.Handler {
	limiter := newRateLimiter(requestsPerSecond, burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract IP address from RemoteAddr
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			if !limiter.getLimiter(ip).Allow() {
				shared.Log.Warn("rate limit exceeded", zap.String("ip", ip))
				http.Error(w, "429 Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
