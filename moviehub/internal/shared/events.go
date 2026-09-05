package shared

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var RedisDB *redis.Client

const UserEventsStream = "user:events"

type UserCreatedEvent struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

func InitRedis(uri string) error {
	opts, err := redis.ParseURL(uri)
	if err != nil {
		// fallback to standard address
		opts = &redis.Options{
			Addr: uri,
		}
	}

	RedisDB = redis.NewClient(opts)
	if err := RedisDB.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	Log.Info("Connected to Redis", zap.String("uri", uri))
	return nil
}

func PublishUserCreatedEvent(ctx context.Context, userID, username string) error {
	event := UserCreatedEvent{
		UserID:   userID,
		Username: username,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = RedisDB.XAdd(ctx, &redis.XAddArgs{
		Stream: UserEventsStream,
		Values: map[string]interface{}{
			"event_type": "user_created",
			"payload":    string(data),
		},
	}).Err()

	if err != nil {
		return fmt.Errorf("failed to publish to redis stream: %w", err)
	}

	Log.Info("Published user_created event", zap.String("user_id", userID))
	return nil
}
