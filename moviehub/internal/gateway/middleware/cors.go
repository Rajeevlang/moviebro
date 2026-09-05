package middleware

import (
	"net/http"
	"strings"
)

// CORS returns a middleware that adds Cross-Origin Resource Sharing headers
// for the configured allowed origins (e.g. the Cloudflare Pages domain).
// An empty allow-list disables CORS entirely: no CORS headers are emitted,
// which is correct for same-origin deployments.
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		o = strings.TrimSpace(strings.ToLower(o))
		if o != "" && o != "null" {
			allowed[o] = true
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && allowed[strings.ToLower(origin)] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
				w.Header().Set("Access-Control-Max-Age", "86400")
			}

			// Short-circuit preflight requests before Auth/RateLimit run,
			// so browsers never see a 401 on an unauthenticated OPTIONS call.
			if r.Method == http.MethodOptions && origin != "" &&
				r.Header.Get("Access-Control-Request-Method") != "" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
