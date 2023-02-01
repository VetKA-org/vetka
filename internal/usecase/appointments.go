package usecase

import (
	"context"

	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/internal/repo"
	uuid "github.com/satori/go.uuid"
)

type AppointmentsUseCase struct {
	appointmentsRepo repo.Appointments
}

func NewAppointmentsUseCase(appointments repo.Appointments) *AppointmentsUseCase {
	return &AppointmentsUseCase{appointments}
}

func (uc *AppointmentsUseCase) List(ctx context.Context) ([]entity.Appointment, error) {
	return nil, nil
}

func (uc *AppointmentsUseCase) Create(ctx context.Context) error {
	return nil
}

func (uc *AppointmentsUseCase) Update(ctx context.Context, id uuid.UUID) error {
	return nil
}
