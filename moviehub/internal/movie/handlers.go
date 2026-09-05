package movie

import (
	"context"
	"errors"
	"moviesapi/internal/shared"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// Handler holds the repository
type Handler struct {
	repo   *MovieRepository
	plrepo *Plrepo
	rdb    *Redisdb
}

// NewHandler creates a new Handler
func NewHandler(repo *MovieRepository, p *Plrepo, rdb *Redisdb) *Handler {
	return &Handler{
		repo:   repo,
		plrepo: p,
		rdb:    rdb,
	}
}

// get requests

// Method	URI	Purpose
// GET	/v1/movies	List movies (paginated, filtered, sorted). Query params: ?page=2&per_page=20&genre=sci-fi&q=interstellar&sort=-average_rating
// GET	/v1/movies/{movieId}	Get a single movie with full detail (optionally include cast/crew via ?include=cast,crew)
// POST	/v1/movies	(Admin) Create a new movie
// PUT	/v1/movies/{movieId}	(Admin) Replace movie data
// PATCH	/v1/movies/{movieId}	(Admin) Partial update
// DELETE	/v1/movies/{movieId}	(Admin) Delete a movie

// CreateMovie handles the creation of a new movie
func (h *Handler) CreateMovie(c *gin.Context) {
	// 1. Create an input struct with validation tags
	var input struct {
		Title       string   `json:"title" binding:"required,min=2"`
		Description string   `json:"description" binding:"required"`
		ReleaseYear int      `json:"release_year" binding:"required,min=1888"`
		Genres      []string `json:"genres" binding:"required,min=1"`
		Poster      string   `json:"poster" binding:"required,url"`
	}

	// 2. Bind and Validate
	if err := c.ShouldBindJSON(&input); err != nil {
		var errs validator.ValidationErrors
		v := shared.New() // Using your custom validator struct!

		if errors.As(err, &errs) {
			for _, f := range errs {
				// f.Field() gets the struct field name, f.Tag() gets the validation rule that failed
				v.AddError(f.Field(), "failed validation on tag: "+f.Tag())
			}
		} else {
			v.AddError("request", "invalid JSON format")
		}

		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": v.Errors})
		return // CRITICAL: Must return so it doesn't proceed to database insertion!
	}

	// 3. Map input to our actual Database Model
	movie := Movie{
		ID:          primitive.NewObjectID(),
		Title:       input.Title,
		Description: input.Description,
		ReleaseYear: input.ReleaseYear,
		Genres:      input.Genres,
		Poster:      input.Poster,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 4. Save to Database
	if err := h.repo.CreateMovie(c.Request.Context(), &movie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "try again in some time"})
		return
	}

	c.JSON(http.StatusCreated, movie)
}

// GetMovie retrieves a movie by ID
func (h *Handler) GetMovie(c *gin.Context) {
	id := c.Param("id")

	// 1. Try to get from Cache
	if cachedMovie, err := h.rdb.Getmovie(c.Request.Context(), id); err == nil {
		c.JSON(http.StatusOK, cachedMovie)
		return
	}

	// 2. Fallback to Database
	movie, err := h.repo.GetMovieByID(c.Request.Context(), id)
	if err != nil {
		// Assuming the error is a not found error, could be more specific
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	// 3. Save to Cache asynchronously so we don't block the response
	go func() {
		if err := h.rdb.Addmovie(context.Background(), id, movie); err != nil {
			shared.Log.Error("Failed to cache movie", zap.Error(err))
		}
	}()

	c.JSON(http.StatusOK, movie)
}

// GetAllMovies handles retrieving a list of movies based on query parameters
func (h *Handler) GetAllMovies(c *gin.Context) {
	var query MovieSearchQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	movies, err := h.repo.GetAllMovies(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		return
	}

	// If no movies found, return an empty array instead of null
	if movies == nil {
		movies = []Movie{}
	}

	c.JSON(http.StatusOK, movies)
}

// AddReview handles the submission of a new review for a movie
func (h *Handler) AddReview(c *gin.Context) {
	movieID := c.Param("id")

	var input struct {
		UserID  string  `json:"user_id" binding:"required"`
		Comment string  `json:"comment"`
		Score   float64 `json:"score" binding:"required,min=1,max=10"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review data or score out of bounds"})
		return
	}

	newReview := Review{
		ID:        primitive.NewObjectID(),
		UserID:    input.UserID,
		Comment:   input.Comment,
		Score:     input.Score,
		CreatedAt: time.Now(),
	}

	if err := h.repo.AddReviewAndUpdateRating(c.Request.Context(), movieID, newReview); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add review"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Review added successfully", "review": newReview})
}

// UpdateMovie handles partial updates to a movie via PATCH
func (h *Handler) UpdateMovie(c *gin.Context) {
	id := c.Param("id")

	// We use pointers so we can differentiate between "empty string" and "not provided"
	var input struct {
		Title       *string `json:"title" binding:"omitempty,min=2"`
		Description *string `json:"description" binding:"omitempty"`
		ReleaseYear *int    `json:"release_year" binding:"omitempty,min=1888"`
		Poster      *string `json:"poster" binding:"omitempty,url"`
	}

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

	// Dynamically build the update map based on provided fields
	updateData := bson.M{}
	if input.Title != nil {
		updateData["title"] = *input.Title
	}
	if input.Description != nil {
		updateData["description"] = *input.Description
	}
	if input.ReleaseYear != nil {
		updateData["release_year"] = *input.ReleaseYear
	}
	if input.Poster != nil {
		updateData["poster"] = *input.Poster
	}

	// If nothing was provided to update
	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields provided for update"})
		return
	}

	if err := h.repo.UpdateMovie(c.Request.Context(), id, updateData); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
		return
	}

	// Invalidate Cache
	go h.rdb.Delmovie(context.Background(), id)

	c.JSON(http.StatusOK, gin.H{"message": "Movie updated successfully"})
}

// DeleteMovie handles the removal of a movie by ID
func (h *Handler) DeleteMovie(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.DeleteMovie(c.Request.Context(), id); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie"})
		return
	}

	// Invalidate Cache
	go h.rdb.Delmovie(context.Background(), id)

	// 204 No Content is the standard response for a successful deletion
	c.JSON(http.StatusNoContent, nil)
}
