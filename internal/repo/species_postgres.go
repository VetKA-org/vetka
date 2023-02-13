package repo

import (
	"context"
	"fmt"

	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/pkg/postgres"
)

type SpeciesRepo struct {
	*postgres.Postgres
}

func NewSpeciesRepo(pg *postgres.Postgres) *SpeciesRepo {
	return &SpeciesRepo{pg}
}

func (r *SpeciesRepo) List(ctx context.Context) ([]entity.Species, error) {
	rv := make([]entity.Species, 0)
	if err := r.Select(ctx, &rv, "SELECT * FROM animal_species"); err != nil {
		return nil, fmt.Errorf("SpeciesRepo - List - r.Select: %w", err)
	}

	return rv, nil
}
