package usecase

import (
	"context"
	"fmt"
	"sync"

	"github.com/VetKA-org/vetka/internal/repo"
	uuid "github.com/satori/go.uuid"
)

type QueueUseCase struct {
	queueRepo repo.Queue
	mu        sync.RWMutex
}

func NewQueueUseCase(queue repo.Queue) *QueueUseCase {
	return &QueueUseCase{queueRepo: queue}
}

func (uc *QueueUseCase) List(ctx context.Context) ([]uuid.UUID, error) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	patients, err := uc.queueRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("QueueUseCase - List - uc.QueueRepo.List: %w", err)
	}

	return patients, nil
}

func (uc *QueueUseCase) Enqueue(ctx context.Context, id uuid.UUID) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if err := uc.queueRepo.Enqueue(ctx, id); err != nil {
		return fmt.Errorf("QueueUseCase - Enqueue - uc.QueueRepo.Enqueue: %w", err)
	}

	return nil
}

func (uc *QueueUseCase) MoveUp(ctx context.Context, id uuid.UUID) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	return nil
}

func (uc *QueueUseCase) MoveDown(ctx context.Context, id uuid.UUID) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	return nil
}

func (uc *QueueUseCase) Dequeue(ctx context.Context, id uuid.UUID) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if err := uc.queueRepo.Dequeue(ctx, id); err != nil {
		return fmt.Errorf("QueueUseCase - Dequeue - uc.QueueRepo.Dequeue: %w", err)
	}

	return nil
}
