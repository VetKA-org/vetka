package entity

import uuid "github.com/satori/go.uuid"

type Species struct {
	ID    uuid.UUID `json:"id" db:"species_id"`
	Title string    `json:"title"`
}
