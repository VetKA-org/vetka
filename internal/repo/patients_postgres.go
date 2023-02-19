package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/VetKA-org/vetka/pkg/entity"
	"github.com/VetKA-org/vetka/pkg/postgres"
	uuid "github.com/satori/go.uuid"
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
         species_id,
         gender,
         breed,
         aggressive,
         vaccinated_at,
         sterilized_at
     FROM
         patients`,
	); err != nil {
		return nil, fmt.Errorf("PatientsRepo - List - r.Select: %w", err)
	}

	return rv, nil
}

func (r *PatientsRepo) Register(
	ctx context.Context,
	tx postgres.Transaction,
	name string,
	speciesID uuid.UUID,
	gender entity.Gender,
	breed string,
	birth time.Time,
	aggressive bool,
	vaccinatedAt *time.Time,
	sterilizedAt *time.Time,
) error {
	if _, err := tx.Tx.Exec(
		ctx,
		`INSERT INTO patients
         (name, species_id, gender, breed, birth, aggressive, vaccinated_at, sterilized_at)
     VALUES
         ($1, $2, $3, $4, $5, $6, $7, $8)
    `,
		name,
		speciesID,
		gender,
		breed,
		birth,
		aggressive,
		vaccinatedAt,
		sterilizedAt,
	); err != nil {
		if postgres.IsEntityExists(err) {
			return entity.ErrPatientExists
		}

		return fmt.Errorf("PatientsRepo - Register - tx.Tx.Exec: %w", err)
	}

	return nil
}
