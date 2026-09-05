package movie

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Plrepo struct {
	c *mongo.Collection
}

func PLrepoHandler(c *mongo.Collection) *Plrepo {
	return &Plrepo{c: c}
}

// CreatePlaylist creates a new playlist
func (pcoll *Plrepo) CreatePlaylist(ctx context.Context, playlist *Playlist) error {
	playlist.CreatedAt = time.Now()
	playlist.UpdatedAt = time.Now()
	if playlist.MovieIDs == nil {
		playlist.MovieIDs = []primitive.ObjectID{}
	}

	result, err := pcoll.c.InsertOne(ctx, playlist)
	if err != nil {
		return err
	}

	playlist.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetPlaylistByID retrieves a playlist by its ID
func (pcoll *Plrepo) GetPlaylistByID(ctx context.Context, id string) (*Playlist, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid playlist id: %v", err)
	}

	var playlist Playlist
	err = pcoll.c.FindOne(ctx, bson.M{"_id": objID}).Decode(&playlist)
	if err != nil {
		return nil, err
	}

	return &playlist, nil
}

// GetPlaylistsByUser retrieves all playlists for a given user
func (pcoll *Plrepo) GetPlaylistsByUser(ctx context.Context, userID string) ([]Playlist, error) {
	cursor, err := pcoll.c.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var playlists []Playlist
	if err = cursor.All(ctx, &playlists); err != nil {
		return nil, err
	}

	return playlists, nil
}

// AddMovieToPlaylist adds a movie to the playlist
func (pcoll *Plrepo) AddMovieToPlaylist(ctx context.Context, playlistID string, movieID string) error {
	pObjID, err := primitive.ObjectIDFromHex(playlistID)
	if err != nil {
		return fmt.Errorf("invalid playlist id: %v", err)
	}

	mObjID, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		return fmt.Errorf("invalid movie id: %v", err)
	}

	update := bson.M{
		"$addToSet": bson.M{"movie_ids": mObjID},
		"$set":      bson.M{"updated_at": time.Now()},
	}

	result, err := pcoll.c.UpdateByID(ctx, pObjID, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// RemoveMovieFromPlaylist removes a movie from the playlist
func (pcoll *Plrepo) RemoveMovieFromPlaylist(ctx context.Context, playlistID string, movieID string) error {
	pObjID, err := primitive.ObjectIDFromHex(playlistID)
	if err != nil {
		return fmt.Errorf("invalid playlist id: %v", err)
	}

	mObjID, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		return fmt.Errorf("invalid movie id: %v", err)
	}

	update := bson.M{
		"$pull": bson.M{"movie_ids": mObjID},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := pcoll.c.UpdateByID(ctx, pObjID, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// DeletePlaylist deletes a playlist
func (pcoll *Plrepo) DeletePlaylist(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid playlist id: %v", err)
	}

	result, err := pcoll.c.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
