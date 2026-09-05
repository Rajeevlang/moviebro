package movie

import "github.com/gin-gonic/gin"

// RegisterRoutes registers the movie routes
func RegisterRoutes(router *gin.Engine, handler *Handler) {
	movies := router.Group("/api/v1/movies")
	{
		movies.GET("/", handler.GetAllMovies)
		movies.POST("/", handler.CreateMovie)
		movies.GET("/:id", handler.GetMovie)
		movies.PATCH("/:id", handler.UpdateMovie)
		movies.DELETE("/:id", handler.DeleteMovie)
		movies.POST("/:id/reviews", handler.AddReview)
	}

	Playlist := router.Group("/api/v1/list")
	{
		Playlist.POST("/", handler.CreatePlaylist)
		Playlist.GET("/:id", handler.GetPlaylist)
		Playlist.GET("/user/:user_id", handler.GetUserPlaylists)
		Playlist.DELETE("/:id", handler.DeletePlaylist)
		Playlist.POST("/:id/movies/:movie_id", handler.AddMovieToPlaylist)
		Playlist.DELETE("/:id/movies/:movie_id", handler.RemoveMovieFromPlaylist)
	}
}
