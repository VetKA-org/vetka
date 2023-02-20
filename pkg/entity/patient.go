package entity

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	ErrPatientExists   = errors.New("patient already exists")
	ErrPatientNotFound = errors.New("patient not found")
	ErrUnknownSpecies  = errors.New("species not found")
)

type Patient struct {
	ID           uuid.UUID  `json:"id" db:"patient_id"`
	Name         string     `json:"name"`
	SpeciesID    uuid.UUID  `json:"species_id"`
	Gender       Gender     `json:"gender"`
	Breed        string     `json:"breed"`
	Birth        time.Time  `json:"birth"`
	Aggressive   bool       `json:"aggressive"`
	VaccinatedAt *time.Time `json:"vaccinated_at"`
	SterilizedAt *time.Time `json:"sterilized_at"`
}
