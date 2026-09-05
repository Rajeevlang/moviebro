package movie

import (
	"context"
	"encoding/json"
	"fmt"
	"moviesapi/internal/shared"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func StartRedisConsumer(ctx context.Context, plrepo *Plrepo) {
	// create consumer group if it doesn't exist

	err := shared.RedisDB.XGroupCreateMkStream(ctx, shared.UserEventsStream, "movie_service_group", "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		shared.Log.Error("failed to create consumer group", zap.Error(err))
	}

	for {
		select {
		case <-ctx.Done():
			shared.Log.Info("Stopping Redis consumer")
			return
		default:
			// Read from stream
			streams, err := shared.RedisDB.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    "movie_service_group",
				Consumer: "consumer1",
				Streams:  []string{shared.UserEventsStream, ">"},
				Count:    10,
				Block:    5 * time.Second,
				NoAck:    false,
			}).Result()

			if err != nil {
				if err != redis.Nil {
					shared.Log.Error("failed to read from redis stream", zap.Error(err))
					time.Sleep(2 * time.Second)
				}
				continue
			}

			for _, stream := range streams {
				for _, msg := range stream.Messages {
					processMessage(ctx, plrepo, msg)
				}
			}
		}
	}
}

func processMessage(ctx context.Context, plrepo *Plrepo, msg redis.XMessage) {
	eventType, ok := msg.Values["event_type"].(string)
	if !ok || eventType != "user_created" {
		// Acknowledge and skip if not user_created
		shared.RedisDB.XAck(ctx, shared.UserEventsStream, "movie_service_group", msg.ID)
		return
	}

	payloadStr, ok := msg.Values["payload"].(string)
	if !ok {
		shared.RedisDB.XAck(ctx, shared.UserEventsStream, "movie_service_group", msg.ID)
		return
	}

	var event shared.UserCreatedEvent
	if err := json.Unmarshal([]byte(payloadStr), &event); err != nil {
		shared.Log.Error("failed to unmarshal user_created event", zap.Error(err))
		shared.RedisDB.XAck(ctx, shared.UserEventsStream, "movie_service_group", msg.ID)
		return
	}

	// Create automatic playlist
	playlist := &Playlist{
		ID:          primitive.NewObjectID(),
		UserID:      event.UserID,
		Name:        "Watch Later",
		Description: fmt.Sprintf("Automatic Watch Later playlist for %s", event.Username),
		IsPublic:    false,
		MovieIDs:    []primitive.ObjectID{},
	}

	if err := plrepo.CreatePlaylist(ctx, playlist); err != nil {
		shared.Log.Error("failed to create automatic playlist", zap.Error(err))
		return // Do not ack so it can be retried
	}

	// Create automatic playlist
	playlist2 := &Playlist{
		ID:          primitive.NewObjectID(),
		UserID:      event.UserID,
		Name:        "favorites",
		Description: fmt.Sprintf("favourates playlist for %s", event.Username),
		IsPublic:    false,
		MovieIDs:    []primitive.ObjectID{},
	}

	if err := plrepo.CreatePlaylist(ctx, playlist2); err != nil {
		shared.Log.Error("failed to create automatic playlist", zap.Error(err))
		return // Do not ack so it can be retried
	}

	shared.Log.Info("Created automatic playlists watch later and favories ", zap.String("user_id", event.UserID))

	// Acknowledge the message
	shared.RedisDB.XAck(ctx, shared.UserEventsStream, "movie_service_group", msg.ID)
}
