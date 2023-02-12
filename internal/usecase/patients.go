package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/internal/repo"
	uuid "github.com/satori/go.uuid"
)

type PatientsUseCase struct {
	patientsRepo     repo.Patients
	appointmentsRepo repo.Appointments
}

func NewPatientsUseCase(patients repo.Patients, appointments repo.Appointments) *PatientsUseCase {
	return &PatientsUseCase{patients, appointments}
}

func (uc *PatientsUseCase) List(ctx context.Context) ([]entity.Patient, error) {
	patients, err := uc.patientsRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("PatientsUseCase - List - uc.patientsRepo.List: %w", err)
	}

	return patients, nil
}

func (uc *PatientsUseCase) Register(
	ctx context.Context,
	name string,
	speciesID uuid.UUID,
	gender entity.Gender,
	breed string,
	birth time.Time,
	aggressive bool,
	vaccinatedAt *time.Time,
	sterilizedAt *time.Time,
) error {
	tx, err := uc.patientsRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("PatientsUseCase - Register - uc.patientsRepo.BeginTx: %w", err)
	}

	defer tx.Release()

	if err := uc.patientsRepo.Register(
		ctx,
		tx,
		name,
		speciesID,
		gender,
		breed,
		birth,
		aggressive,
		vaccinatedAt,
		sterilizedAt,
	); err != nil {
		return fmt.Errorf("PatientsUseCase - Register - uc.patientsRepo.Register: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("PatientsUseCase - Register - tx.Commit: %w", err)
	}

	return nil
}

func (uc *PatientsUseCase) ListAppointments(ctx context.CancelFunc, id uuid.UUID) error {
	return nil
}
