package repo

import (
	"context"
	"fmt"
	"strconv"

	"github.com/VetKA-org/vetka/pkg/entity"
	"github.com/VetKA-org/vetka/pkg/postgres"
)

type SpeciesRepo struct {
	*postgres.Postgres
}

func NewSpeciesRepo(pg *postgres.Postgres) *SpeciesRepo {
	return &SpeciesRepo{pg}
}

func (r *SpeciesRepo) List(ctx context.Context, title string) ([]entity.Species, error) {
	var values []interface{}

	query := "SELECT * FROM animal_species"

	if title != "" {
		values = append(values, title+"%")
		query += " WHERE title ILIKE $" + strconv.Itoa(len(values))
	}

	rv := make([]entity.Species, 0)
	if err := r.Select(ctx, &rv, query, values...); err != nil {
		return nil, fmt.Errorf("SpeciesRepo - List - r.Select: %w", err)
	}

	return rv, nil
}
