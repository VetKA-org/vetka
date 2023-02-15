package repo

import (
	"context"
	"fmt"

	"github.com/VetKA-org/vetka/pkg/redis"
	uuid "github.com/satori/go.uuid"
)

const _patientsQueueKey = "patients_queue"

type QueueRepo struct {
	*redis.Redis
}

func NewQueueRepo(rdb *redis.Redis) *QueueRepo {
	return &QueueRepo{rdb}
}

func (r *QueueRepo) List(ctx context.Context) ([]uuid.UUID, error) {
	rawData, err := r.Client.ZRange(ctx, _patientsQueueKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("QueueRepo - List - r.Client.ZRange: %w", err)
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
