package usecase

import (
	"context"
	"fmt"

	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/internal/repo"
)

type UsersUseCase struct {
	usersRepo repo.Users
}

func NewUsersUseCase(users repo.Users) *UsersUseCase {
	return &UsersUseCase{users}
}

func (uc *UsersUseCase) List(ctx context.Context) ([]entity.User, error) {
	users, err := uc.usersRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("UsersUseCase - List - uc.usersRepo.List: %w", err)
	}

	return users, nil
}

func (uc *UsersUseCase) Register(ctx context.Context, username, password string) error {
	return nil
}
