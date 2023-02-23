// Package redis implements connection to standalone Redis service.
package redis

import (
	"context"
	"fmt"

	"github.com/VetKA-org/vetka/pkg/entity"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func New(url entity.SecretURI) (*Redis, error) {
	cfg, err := redis.ParseURL(string(url))
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

func (r *Redis) LFirstMatch(ctx context.Context, key, value string) (int64, error) {
	pos, err := r.Client.LPos(ctx, key, value, redis.LPosArgs{Rank: 1}).Result()
	if err == nil {
		return pos, nil
	}

	if IsEntityNotFound(err) {
		return -1, nil
	}

	return -1, fmt.Errorf("Redis - ListFirst - r.Client.LPos: %w", err)
}

func (r *Redis) LSwap(
	ctx context.Context,
	key string,
	lPos int64,
	lValue string,
	rPos int64,
) error {
	rValue, err := r.Client.LIndex(ctx, key, rPos).Result()
	if err != nil {
		return fmt.Errorf("Redis - LSwap - r.Client.LIndex: %w", err)
	}

	if _, err := r.Client.LSet(ctx, key, rPos, lValue).Result(); err != nil {
		return fmt.Errorf("Redis - LSwap - r.Client.LSet: %w", err)
	}

	if _, err := r.Client.LSet(ctx, key, lPos, rValue).Result(); err != nil {
		return fmt.Errorf("Redis - LSwap - r.Client.LSet: %w", err)
	}

	return nil
}
