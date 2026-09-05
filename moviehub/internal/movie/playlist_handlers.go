package movie

import (
	"errors"
	"moviesapi/internal/shared"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreatePlaylist handles the creation of a new playlist
func (h *Handler) CreatePlaylist(c *gin.Context) {
	var input CreatePlaylist

	if err := c.ShouldBindJSON(&input); err != nil {
		var errs validator.ValidationErrors
		v := shared.New()

		if errors.As(err, &errs) {
			for _, f := range errs {
				v.AddError(f.Field(), "failed validation on tag: "+f.Tag())
			}
		} else {
			v.AddError("request", "invalid JSON format")
		}

		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": v.Errors})
		return
	}

	description := ""
	if input.Description != nil {
		description = *input.Description
	}

	playlist := Playlist{
		ID:          primitive.NewObjectID(),
		UserID:      input.Userid,
		Name:        input.Playlistname,
		Description: description,
		IsPublic:    input.Visiblity,
	}

	if err := h.plrepo.CreatePlaylist(c.Request.Context(), &playlist); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create playlist"})
		return
	}

	c.JSON(http.StatusCreated, playlist)
}

// GetPlaylist handles retrieving a playlist by ID
func (h *Handler) GetPlaylist(c *gin.Context) {
	id := c.Param("id")

	playlist, err := h.plrepo.GetPlaylistByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve playlist"})
		return
	}

	c.JSON(http.StatusOK, playlist)
}

// GetUserPlaylists handles retrieving all playlists for a user
func (h *Handler) GetUserPlaylists(c *gin.Context) {
	userID := c.Param("user_id")

	playlists, err := h.plrepo.GetPlaylistsByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve playlists"})
		return
	}

	if playlists == nil {
		playlists = []Playlist{}
	}

	c.JSON(http.StatusOK, playlists)
}

// AddMovieToPlaylist handles adding a movie to a playlist
func (h *Handler) AddMovieToPlaylist(c *gin.Context) {
	playlistID := c.Param("id")
	movieID := c.Param("movie_id")

	if err := h.plrepo.AddMovieToPlaylist(c.Request.Context(), playlistID, movieID); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie to playlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie added to playlist successfully"})
}

// RemoveMovieFromPlaylist handles removing a movie from a playlist
func (h *Handler) RemoveMovieFromPlaylist(c *gin.Context) {
	playlistID := c.Param("id")
	movieID := c.Param("movie_id")

	if err := h.plrepo.RemoveMovieFromPlaylist(c.Request.Context(), playlistID, movieID); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove movie from playlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie removed from playlist successfully"})
}

// DeletePlaylist handles deleting a playlist
func (h *Handler) DeletePlaylist(c *gin.Context) {
	playlistID := c.Param("id")

	if err := h.plrepo.DeletePlaylist(c.Request.Context(), playlistID); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete playlist"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
