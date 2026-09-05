package movie

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// `json:""  bson:""  binding:""`

// Person represents an individual involved in a movie (Actor, Director, etc.)
type Person struct {
	Name string `json:"name" bson:"name"`
	Role string `json:"role,omitempty" bson:"role,omitempty"` // e.g., "Lead Actor", "Supporting Actor"
}

// Rating aggregates the movie's review scores
type Rating struct {
	Average float64 `json:"average" bson:"average"`
	Count   int     `json:"count" bson:"count"`
}

// Review represents a single user review
type Review struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty" `
	UserID    string             `json:"user_id" bson:"user_id"`
	Comment   string             `json:"comment" bson:"comment"`
	Score     float64            `json:"score" bson:"score"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

// Movie represents a movie in our catalog
type Movie struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	ReleaseYear int                `json:"release_year" bson:"release_year"`
	Genres      []string           `json:"genres" bson:"genres"`                       // e.g., ["Action", "Sci-Fi"]
	Directors   []Person           `json:"directors" bson:"directors"`                 // Embedded directors
	Actors      []Person           `json:"actors" bson:"actors"`                       // Embedded actors
	Poster      string             `json:"poster" bson:"poster"`                       // URL to movie poster
	Rating      Rating             `json:"rating" bson:"rating"`                       // Cached aggregate rating
	Reviews     []Review           `json:"reviews,omitempty" bson:"reviews,omitempty"` // Embedded reviews (can be separated later)
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type MovieSearchQuery struct {
	Page      int64   `form:"page,default=1"`
	Limit     int64   `form:"limit,default=10"`
	Year      int     `form:"year"`
	Actor     string  `form:"actor"`
	Director  string  `form:"director"`
	MinRating float64 `form:"min_rating"` // e.g. ?min_rating=8.5
}

type CreateMovieInput struct {
	Title       string   `json:"title" binding:"required,min=2,max=100"`
	Description string   `json:"description" binding:"required,max=500"`
	ReleaseYear int      `json:"release_year" binding:"required,gte=1888,lte=2100"`
	Genres      []string `json:"genres" binding:"required,min=1"` // Must have at least 1 genre
	Poster      string   `json:"poster" binding:"required,url"`   // Must be a valid URL
}

type Playlist struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	UserID      string               `json:"user_id" bson:"user_id"` // String to hold Postgres UUID
	Name        string               `json:"name" bson:"name"`
	Description string               `json:"description" bson:"description"`
	IsPublic    bool                 `json:"is_public" bson:"is_public"`
	IsDefault   bool                 `json:"is_default" bson:"is_default"`
	MovieIDs    []primitive.ObjectID `json:"movie_ids" bson:"movie_ids"` // Array of Mongo ObjectIDs
	CreatedAt   time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at" bson:"updated_at"`
}
type CreatePlaylist struct {
	Userid       string  `json:"user_id" binding:"required"`
	Visiblity    bool    `json:"visiblity"`
	Playlistname string  `json:"playlist_name" binding:"required"`
	Description  *string `json:"description"`
}
