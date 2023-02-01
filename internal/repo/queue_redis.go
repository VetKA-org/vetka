package repo

import "github.com/VetKA-org/vetka/pkg/redis"

type QueueRepo struct {
	*redis.Redis
}

func NewQueueRepo(rdb *redis.Redis) *QueueRepo {
	return &QueueRepo{rdb}
}
