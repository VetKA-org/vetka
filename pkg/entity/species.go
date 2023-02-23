package entity

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

var ErrSpeciesNotFound = errors.New("species not found")

type Species struct {
	ID    uuid.UUID `json:"id" db:"species_id"`
	Title string    `json:"title"`
}
