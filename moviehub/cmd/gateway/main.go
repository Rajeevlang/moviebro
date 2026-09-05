package main

import (
	"moviesapi/internal/gateway"
	"moviesapi/internal/gateway/middleware"
	"moviesapi/internal/shared"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	shared.InitLogger()
	defer shared.Log.Sync()

	// Load config
	config, err := shared.LoadConfig(".")
	if err != nil {
		shared.Log.Fatal("cannot load config", zap.Error(err))
	}

	// Use the SAME secret the user-service signs JWTs with,
	// otherwise token validation in Auth() always fails.
	middleware.SetJWTSecret(config.SECERETKEY)

	// Create Movie Proxy
	movieProxy, err := gateway.NewMovieProxy(config.MovieServiceURL)
	if err != nil {
		shared.Log.Fatal("failed to create movie proxy", zap.Error(err))
	}

	// Create User Proxy
	userProxy, err := gateway.NewUserProxy(config.UserServiceURL)
	if err != nil {
		shared.Log.Fatal("failed to create user proxy", zap.Error(err))
	}

	// Setup middlewares
	// e.g. 100 requests per second, burst size 200
	rateLimit := middleware.RateLimit(100, 200)

	// Wrap the proxy handlers in middleware chain
	movieProxyHandler := middleware.Logger(rateLimit(middleware.Auth(movieProxy)))
	userProxyHandler := middleware.Logger(rateLimit(userProxy)) // No Auth needed to login/register

	// Setup routing mux
	mux := http.NewServeMux()

	// Healthcheck for load balancers, Docker healthchecks and CI/CD smoke tests
	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"api-gateway"}`))
	})

	// Prometheus Metrics Endpoint for scraping
	mux.Handle("/metrics", shared.MetricsHandler())

	// Forward to Movie proxy
	mux.Handle("/api/v1/movies/", movieProxyHandler)
	mux.Handle("/api/v1/movies", movieProxyHandler)
	mux.Handle("/api/v1/list/", movieProxyHandler)
	mux.Handle("/api/v1/list", movieProxyHandler)

	// Forward to User proxy
	mux.Handle("/api/v1/users/", userProxyHandler)
	mux.Handle("/api/v1/users", userProxyHandler)

	// Serve the frontend UI
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/", fs)

	handler := http.Handler(mux)

	// Enable CORS only for origins explicitly allow-listed via
	// CORS_ALLOWED_ORIGINS (e.g. the Cloudflare Pages domain).
	if config.CORSAllowedOrigins != "" {
		cors := middleware.CORS(strings.Split(config.CORSAllowedOrigins, ","))
		handler = cors(handler)
	}

	shared.Log.Info("Starting API Gateway",
		zap.String("port", config.GatewayPort),
		zap.String("cors_origins", config.CORSAllowedOrigins))

	serverPort := config.GatewayPort
	if serverPort == "" {
		serverPort = "8080" // fallback
	}

	monitoredMux := shared.PrometheusHTTPMiddleware("api-gateway", handler)

	if err := http.ListenAndServe(":"+serverPort, monitoredMux); err != nil {
		shared.Log.Fatal("API Gateway server failed", zap.Error(err))
	}
}
