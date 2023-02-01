// Package redis implements connection to standalone Redis service.
package redis

import (
	"github.com/redis/go-redis/v9"
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
