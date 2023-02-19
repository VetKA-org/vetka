package usecase

import (
	"context"
	"fmt"

	"github.com/VetKA-org/vetka/internal/repo"
	"github.com/VetKA-org/vetka/pkg/entity"
)

type RolesUseCase struct {
	rolesRepo repo.Roles
}

func NewRolesUseCase(roles repo.Roles) *RolesUseCase {
	return &RolesUseCase{roles}
}

func (uc *RolesUseCase) List(ctx context.Context) ([]entity.Role, error) {
	roles, err := uc.rolesRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("RolesUseCase - List - uc.rolesRepo.List: %w", err)
	}

	return roles, nil
}
