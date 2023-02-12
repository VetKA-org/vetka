package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/pkg/postgres"
	uuid "github.com/satori/go.uuid"
)

type AppointmentsRepo struct {
	*postgres.Postgres
}

func NewAppointmentsRepo(pg *postgres.Postgres) *AppointmentsRepo {
	return &AppointmentsRepo{pg}
}

func (r *AppointmentsRepo) List(ctx context.Context) ([]entity.Appointment, error) {
	rv := make([]entity.Appointment, 0)
	if err := r.Select(ctx, &rv, "SELECT * FROM appointments"); err != nil {
		return nil, fmt.Errorf("AppointmentsRepo - List - r.Select: %w", err)
	}

	return rv, nil
}

func (r *AppointmentsRepo) Create(
	ctx context.Context,
	tx postgres.Transaction,
	patientID uuid.UUID,
	assigneeID uuid.UUID,
	scheduledFor time.Time,
	reason string,
	details *string,
) error {
	if _, err := tx.Tx.Exec(
		ctx,
		`INSERT INTO appointments
         (patient_id, assignee_id, scheduled_for, status, reason, details)
     VALUES
         ($1, $2, $3, 'scheduled', $4, $5)
   `,
		patientID,
		assigneeID,
		scheduledFor,
		reason,
		details,
	); err != nil {
		return fmt.Errorf("AppointmentsRepo - Create - tx.Tx.Exec: %w", err)
	}

	return nil
}
