package repo

import (
	"context"
	"fmt"

	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/pkg/postgres"
)

type PatientsRepo struct {
	*postgres.Postgres
}

func NewPatientsRepo(pg *postgres.Postgres) *PatientsRepo {
	return &PatientsRepo{pg}
}

func (r *PatientsRepo) List(ctx context.Context) ([]entity.Patient, error) {
	rv := make([]entity.Patient, 0)
	if err := r.Select(
		ctx,
		&rv,
		`SELECT
         p.id,
         p.name,
         s.name as species,
         p.gender,
         p.breed,
         p.aggressive,
         p.vaccinated_at,
         p.sterilized_at
     FROM
         patients as p
     LEFT OUTER JOIN animal_species as s ON (p.species = s.id)`,
	); err != nil {
		return nil, fmt.Errorf("PatientsRepo - List - r.Select: %w", err)
	}

	return rv, nil
}
