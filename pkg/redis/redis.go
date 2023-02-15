// Package redis implements connection to standalone Redis service.
package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	uuid "github.com/satori/go.uuid"
)

type Redis struct {
	Client *redis.Client
}

func New(url string) (*Redis, error) {
	cfg, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	rdb := new(Redis)
	rdb.Client = redis.NewClient(cfg)

	return rdb, nil
}

func (r *Redis) Close() error {
	if r.Client == nil {
		return nil
	}

	return r.Client.Close()
}

func (r *Redis) LFirstMatch(ctx context.Context, key string, id uuid.UUID) (int64, error) {
	pos, err := r.Client.LPos(ctx, key, id.String(), redis.LPosArgs{Rank: 1}).Result()
	if err == nil {
		return pos, nil
	}

	if IsEntityNotFound(err) {
		return -1, nil
	}

	return -1, fmt.Errorf("Redis - ListFirst - r.Client.LPos: %w", err)
}
