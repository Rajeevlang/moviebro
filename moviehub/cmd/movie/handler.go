package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Application) healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": app.config.Environment,
		"message": "you can use to get movies"})
}
