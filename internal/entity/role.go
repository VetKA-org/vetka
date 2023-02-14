package entity

import uuid "github.com/satori/go.uuid"

const (
	Administrator = "administrator"
	Doctor        = "doctor"
	HeadDoctor    = "head-doctor"
)

type Role struct {
	ID   uuid.UUID `json:"id" db:"role_id"`
	Name string    `json:"name"`
}
