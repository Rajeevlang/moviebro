package movie

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redisdb struct {
	rdb *redis.Client
}

func NewRedisdb(client *redis.Client) *Redisdb {
	return &Redisdb{rdb: client}
}

func (r *Redisdb) Addmovie(ctx context.Context, id string, movie *Movie) error {
	data, err := json.Marshal(movie)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("movie:%s", id)
	return r.rdb.Set(ctx, key, data, 24*time.Hour).Err()
}

func (r *Redisdb) Getmovie(ctx context.Context, id string) (*Movie, error) {
	key := fmt.Sprintf("movie:%s", id)
	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var movie Movie
	if err := json.Unmarshal([]byte(val), &movie); err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *Redisdb) Delmovie(ctx context.Context, id string) error {
	key := fmt.Sprintf("movie:%s", id)
	return r.rdb.Del(ctx, key).Err()
}
