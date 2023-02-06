package usecase

import (
	"context"
	"errors"
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

func (uc *UsersUseCase) Register(ctx context.Context, login, password string) error {
	if err := uc.usersRepo.Register(ctx, login, password); err != nil {
		if errors.Is(err, entity.ErrUserExists) {
			return err
		}

		return fmt.Errorf("UsersUseCase - Register - uc.usersRepo.Register: %w", err)
	}

	return nil
}
