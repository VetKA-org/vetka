package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Appointment struct {
	ID      uuid.UUID
	At      time.Time
	Patient uuid.UUID
}
