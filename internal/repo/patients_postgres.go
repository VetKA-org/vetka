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
) (uuid.UUID, error) {
	var id uuid.UUID

	err := tx.Tx.QueryRow(
		ctx,
		`INSERT INTO patients
         (name, species_id, gender, breed, birth, aggressive, vaccinated_at, sterilized_at)
     VALUES
         ($1, $2, $3, $4, $5, $6, $7, $8)
     RETURNING patient_id
    `,
		name,
		speciesID,
		gender,
		breed,
		birth,
		aggressive,
		vaccinatedAt,
		sterilizedAt,
	).Scan(&id)
	if err != nil {
		if postgres.IsEntityExists(err) {
			return uuid.UUID{}, entity.ErrPatientExists
		}

		if postgres.IsForeignKeyViolation(err, "patients_species_id_fkey") {
			return uuid.UUID{}, entity.ErrSpeciesNotFound
		}

		return uuid.UUID{}, fmt.Errorf("PatientsRepo - Register - tx.Tx.QueryRow.Scan: %w", err)
	}

	return id, nil
}
