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
         patient_id,
         name,
         a.title,
         gender,
         breed,
         aggressive,
         vaccinated_at,
         sterilized_at
     FROM
         patients
     INNER JOIN
         animal_species AS a
     USING(species_id)`,
	); err != nil {
		return nil, fmt.Errorf("PatientsRepo - List - r.Select: %w", err)
	}

	return rv, nil
}
