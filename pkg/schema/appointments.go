package schema

import (
	"time"

	"github.com/VetKA-org/vetka/pkg/entity"
	uuid "github.com/satori/go.uuid"
)

type CreateAppointmentRequest struct {
	PatientID    uuid.UUID `json:"patient_id" binding:"required"`
	AssigneeID   uuid.UUID `json:"assignee_id" binding:"required"`
	ScheduledFor time.Time `json:"scheduled_for" time_utc:"1" binding:"required"`
	Reason       string    `json:"reason" binding:"required,max=255"`
	Details      *string   `json:"details"`
}

type UpdateAppointmentRequest struct {
	Status entity.ApptStatus `json:"status" binding:"required,oneof=scheduled opened closed canceled"`
}

type ListAppointmentsResponse struct {
	Data []entity.Appointment `json:"data"`
}
