package usecase

import (
	"context"

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
	return nil, nil
}

func (uc *UsersUseCase) Register(ctx context.Context, username, password string) error {
	return nil
}
