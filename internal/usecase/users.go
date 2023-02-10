package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/internal/repo"
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
		return fmt.Errorf("UsersUseCase - Register - uc.usersRepo.Pool.BeginTx: %w", err)
	}

	userID, err := uc.usersRepo.Register(ctx, tx, login, password)
	if err != nil {
		if errors.Is(err, entity.ErrUserExists) {
			return err
		}

		return fmt.Errorf("UsersUseCase - Register - uc.usersRepo.Register: %w", err)
	}

	if len(roles) > 0 {
		if err := uc.rolesRepo.Assign(ctx, tx, userID, roles); err != nil {
			return fmt.Errorf("UsersUseCase - Register - uc.rolesRepo.Assign: %w", err)
		}
	}

	return tx.Commit(ctx)
}
