package usecase

import (
	"context"
	"fmt"

	"github.com/VetKA-org/vetka/internal/repo"
	"github.com/VetKA-org/vetka/pkg/entity"
	uuid "github.com/satori/go.uuid"
)

type UsersUseCase struct {
	usersRepo repo.Users
	rolesRepo repo.Roles
}

func NewUsersUseCase(users repo.Users, roles repo.Roles) *UsersUseCase {
	return &UsersUseCase{users, roles}
}

func (uc *UsersUseCase) List(ctx context.Context) ([]entity.User, error) {
	users, err := uc.usersRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("UsersUseCase - List - uc.usersRepo.List: %w", err)
	}

	return users, nil
}

func (uc *UsersUseCase) Register(
	ctx context.Context,
	login, password string,
	roles []uuid.UUID,
) error {
	tx, err := uc.usersRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("UsersUseCase - Register - uc.usersRepo.BeginTx: %w", err)
	}

	defer tx.Release()

	userID, err := uc.usersRepo.Register(ctx, tx, login, password)
	if err != nil {
		return fmt.Errorf("UsersUseCase - Register - uc.usersRepo.Register: %w", err)
	}

	if len(roles) > 0 {
		if err := uc.rolesRepo.Assign(ctx, tx, userID, roles); err != nil {
			return fmt.Errorf("UsersUseCase - Register - uc.rolesRepo.Assign: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("UsersUseCase - Register - tx.Commit: %w", err)
	}

	return nil
}
