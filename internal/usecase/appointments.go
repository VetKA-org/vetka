package usecase

import (
	"context"
	"fmt"
	"time"

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

func (uc *AppointmentsUseCase) List(
	ctx context.Context,
	patientID *uuid.UUID,
) ([]entity.Appointment, error) {
	appointments, err := uc.appointmentsRepo.List(ctx, patientID)
	if err != nil {
		return nil, fmt.Errorf("AppointmentsUseCase - List - uc.appointmentsRepo.List: %w", err)
	}

	return appointments, nil
}

func (uc *AppointmentsUseCase) Create(
	ctx context.Context,
	patientID uuid.UUID,
	assigneeID uuid.UUID,
	scheduledFor time.Time,
	reason string,
	details *string,
) error {
	tx, err := uc.appointmentsRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("AppointmentsUseCase - Create - uc.appointmentsRepo.BeginTx: %w", err)
	}

	defer tx.Release()

	if err := uc.appointmentsRepo.Create(
		ctx,
		tx,
		patientID,
		assigneeID,
		scheduledFor,
		reason,
		details,
	); err != nil {
		return fmt.Errorf("AppointmentsUseCase - Create - uc.appointmentsRepo.Create: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("AppointmentsUseCase - Create - tx.Commit: %w", err)
	}

	return nil
}

func (uc *AppointmentsUseCase) Update(ctx context.Context, id uuid.UUID) error {
	return nil
}
