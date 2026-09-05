package main

import (
	"net/http"

	"moviesapi/internal/shared"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler) {
	// Prometheus RED Metrics Middleware
	r.Use(shared.PrometheusGinMiddleware("user-service"))

	api := r.Group("/api/v1/users")
	{
		api.POST("/register", h.Register)
		api.POST("/login", h.Login)
		api.POST("/google", h.LoginGoogle)
	}

	// Healthcheck for Docker healthchecks and CI/CD smoke tests
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "user-service"})
	})

	// Prometheus Metrics Endpoint
	r.GET("/metrics", gin.WrapH(shared.MetricsHandler()))

	r.GET("/googlelogin", h.Gloginbefore)
	r.GET("/auth/callback", h.Callback)
}
