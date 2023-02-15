package repo

import (
	"context"
	"fmt"

	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/pkg/redis"
	uuid "github.com/satori/go.uuid"
)

const _queueKey = "patients_queue"

type QueueRepo struct {
	*redis.Redis
}

func NewQueueRepo(rdb *redis.Redis) *QueueRepo {
	return &QueueRepo{rdb}
}

func (r *QueueRepo) List(ctx context.Context) ([]uuid.UUID, error) {
	rawData, err := r.Redis.Client.LRange(ctx, _queueKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("QueueRepo - List - r.Client.LRange: %w", err)
	}

	data := make([]uuid.UUID, len(rawData))

	for i, value := range rawData {
		converted, err := uuid.FromString(value)
		if err != nil {
			return nil, fmt.Errorf("QueueRepo - List - uuid.FromString: %w", err)
		}

		data[i] = converted
	}

	return data, nil
}

func (r *QueueRepo) Enqueue(ctx context.Context, id uuid.UUID) error {
	pos, err := r.Redis.LFirstMatch(ctx, _queueKey, id)
	if err != nil {
		return fmt.Errorf("QueueRepo - Enqueue - r.Redis.LFirstMatch: %w", err)
	}

	if pos != -1 {
		return entity.ErrPatientExists
	}

	if _, err := r.Redis.Client.RPush(ctx, _queueKey, id.String()).Result(); err != nil {
		return fmt.Errorf("QueueRepo - Enqueue - r.Redis.Client.RPush: %w", err)
	}

	return nil
}

func (r *QueueRepo) Dequeue(ctx context.Context, id uuid.UUID) error {
	count, err := r.Redis.Client.LRem(ctx, _queueKey, 1, id.String()).Result()
	if err != nil {
		return fmt.Errorf("QueueRepo - Dequeue - r.Redis.Client.LRem: %w", err)
	}

	if count == 0 {
		return entity.ErrPatientNotFound
	}

	return nil
}
