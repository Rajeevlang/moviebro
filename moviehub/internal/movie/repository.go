package movie

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// MovieRepository handles data operations for the Movie model
type MovieRepository struct {
	collection *mongo.Collection
}

// NewMovieRepository creates a new MovieRepository
func NewMovieRepository(collection *mongo.Collection) *MovieRepository {
	return &MovieRepository{
		collection: collection,
	}
}

// CreateMovie inserts a new movie into the database
func (r *MovieRepository) CreateMovie(ctx context.Context, movie *Movie) error {
	result, err := r.collection.InsertOne(ctx, movie)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		movie.ID = oid
	}

	return nil
}

// GetMovieByID retrieves a movie by its ID
func (r *MovieRepository) GetMovieByID(ctx context.Context, id string) (*Movie, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var movie Movie
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&movie)
	if err != nil {
		return nil, err
	}

	return &movie, nil
}

// GetAllMovies retrieves all movies matching the search query with pagination
func (r *MovieRepository) GetAllMovies(ctx context.Context, query MovieSearchQuery) ([]Movie, error) {
	filter := bson.M{}

	if query.Year != 0 {
		filter["release_year"] = query.Year
	}
	if query.Actor != "" {
		filter["actors.name"] = bson.M{"$regex": query.Actor, "$options": "i"}
	}
	if query.Director != "" {
		filter["directors.name"] = bson.M{"$regex": query.Director, "$options": "i"}
	}
	if query.MinRating > 0 {
		filter["rating.average"] = bson.M{"$gte": query.MinRating}
	}

	skip := (query.Page - 1) * query.Limit
	opts := options.Find().SetSkip(skip).SetLimit(query.Limit)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movies []Movie
	if err = cursor.All(ctx, &movies); err != nil {
		return nil, err
	}

	return movies, nil
}

// AddReviewAndUpdateRating adds a new review to a movie and atomically updates the running average rating
func (r *MovieRepository) AddReviewAndUpdateRating(ctx context.Context, movieID string, newReview Review) error {
	oid, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		return err
	}

	// 1. Fetch the CURRENT rating data
	var movie struct {
		Rating Rating `bson:"rating"`
	}
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&movie)
	if err != nil {
		return err
	}

	// 2. Calculate the NEW running average
	oldCount := movie.Rating.Count
	oldAverage := movie.Rating.Average

	newCount := oldCount + 1
	newAverage := ((oldAverage * float64(oldCount)) + newReview.Score) / float64(newCount)

	// 3. Atomically update the database
	update := bson.M{
		"$push": bson.M{"reviews": newReview},
		"$set": bson.M{
			"rating.average": newAverage,
			"rating.count":   newCount,
		},
	}

	_, err = r.collection.UpdateByID(ctx, oid, update)
	return err
}

// UpdateMovie partially updates a movie document based on provided fields
func (r *MovieRepository) UpdateMovie(ctx context.Context, id string, updateData bson.M) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Always update the UpdatedAt timestamp
	updateData["updated_at"] = time.Now()

	update := bson.M{
		"$set": updateData,
	}

	result, err := r.collection.UpdateByID(ctx, oid, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// DeleteMovie removes a movie from the database by its ID
func (r *MovieRepository) DeleteMovie(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// EnsureIndexes automatically creates the necessary database indexes to ensure fast queries
func (r *MovieRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "release_year", Value: -1}}, // Descending index on year
		},
		{
			Keys: bson.D{{Key: "actors.name", Value: 1}}, // Ascending index on actor names
		},
		{
			Keys: bson.D{{Key: "rating.average", Value: -1}}, // Descending index on average rating
		},
	}

	// CreateMany will only create the indexes if they don't already exist
	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}
