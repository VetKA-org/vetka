package usecase

import (
	"context"
	"fmt"

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

func (uc *PatientsUseCase) Register(ctx context.Context) error {
	return nil
}

func (uc *PatientsUseCase) ListAppointments(ctx context.CancelFunc, id uuid.UUID) error {
	return nil
}
