package entity

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

type ApptStatus string

const (
	ApptScheduled = ApptStatus("scheduled")
	ApptOpened    = ApptStatus("opened")
	ApptClosed    = ApptStatus("closed")
	ApptCanceled  = ApptStatus("canceled")
)

var (
	ErrAptExists    = errors.New("appointment already exists")
	ErrAptNotFound  = errors.New("appointment not found")
	ErrAptHasMaxPos = errors.New("appointment is already first in the queue")
	ErrAptHasMinPos = errors.New("appointment is already last in the queue")
)

type Appointment struct {
	ID           uuid.UUID  `json:"id" db:"appointment_id"`
	PatientID    uuid.UUID  `json:"patient_id"`
	AssigneeID   uuid.UUID  `json:"assignee_id"`
	ScheduledFor time.Time  `json:"scheduled_for"`
	Status       ApptStatus `json:"status"`
	Reason       string     `json:"reason"`
	Details      *string    `json:"details"`
}
