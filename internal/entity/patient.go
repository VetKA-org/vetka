package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Patient struct {
	ID           uuid.UUID  `json:"id" db:"patient_id"`
	Name         string     `json:"name"`
	Species      string     `json:"species" db:"title"`
	Gender       Gender     `json:"gender"`
	Breed        *string    `json:"breed"`
	Birth        time.Time  `json:"birth"`
	Aggressive   bool       `json:"aggressive"`
	VaccinatedAt *time.Time `json:"vaccinated_at"`
	SterilizedAt *time.Time `json:"sterilized_at"`
}
