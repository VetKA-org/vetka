package entity

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

const (
	Administrator = "administrator"
	Doctor        = "doctor"
	HeadDoctor    = "head-doctor"
)

var ErrRoleNotFound = errors.New("role not found")

type Role struct {
	ID   uuid.UUID `json:"id" db:"role_id"`
	Name string    `json:"name"`
}
