package entity

import uuid "github.com/satori/go.uuid"

type Role struct {
	ID   uuid.UUID `json:"id" db:"role_id"`
	Name string    `json:"name"`
}
