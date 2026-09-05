package main

import (
	"moviesapi/internal/shared"

	"github.com/gin-gonic/gin"
)

func (app *Application) routes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	// Prometheus RED Metrics Middleware
	router.Use(shared.PrometheusGinMiddleware("movie-service"))

	router.GET("/healthcheck", app.healthcheck)
	router.GET("/metrics", gin.WrapH(shared.MetricsHandler()))

	return router
}
