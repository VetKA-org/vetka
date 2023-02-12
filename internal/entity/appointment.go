package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Appointment struct {
	ID           uuid.UUID `json:"id" db:"appointment_id"`
	PatientID    uuid.UUID `json:"patient_id"`
	AssigneeID   uuid.UUID `json:"assignee_id"`
	ScheduledFor time.Time `json:"scheduled_for"`
	Status       string    `json:"status"`
	Reason       string    `json:"reason"`
	Details      *string   `json:"details"`
}
