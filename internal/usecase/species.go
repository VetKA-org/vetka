package usecase

import (
	"context"
	"fmt"

	"github.com/VetKA-org/vetka/internal/repo"
	"github.com/VetKA-org/vetka/pkg/entity"
)

type SpeciesUseCase struct {
	speciesRepo repo.Species
}

func NewSpeciesUseCase(species repo.Species) *SpeciesUseCase {
	return &SpeciesUseCase{species}
}

func (uc *SpeciesUseCase) List(ctx context.Context) ([]entity.Species, error) {
	species, err := uc.speciesRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("SpeciesUseCase - List - uc.speciesRepo.List: %w", err)
	}

	return species, nil
}
