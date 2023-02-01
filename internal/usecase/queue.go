package usecase

import (
	"context"

	"github.com/VetKA-org/vetka/internal/repo"
	uuid "github.com/satori/go.uuid"
)

type QueueUseCase struct {
	queueRepo repo.Queue
}

func NewQueueUseCase(queue repo.Queue) *QueueUseCase {
	return &QueueUseCase{queue}
}

func (uc *QueueUseCase) List(ctx context.Context) ([]uuid.UUID, error) {
	return nil, nil
}

func (uc *QueueUseCase) Enqueue(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (uc *QueueUseCase) MoveUp(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (uc *QueueUseCase) MoveDown(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (uc *QueueUseCase) Dequeue(ctx context.Context, id uuid.UUID) error {
	return nil
}
