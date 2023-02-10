package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Species string

const (
	Amphibian    = Species("amphibian")
	Bird         = Species("bird")
	Cat          = Species("cat")
	Dog          = Species("dog")
	ExoticAnimal = Species("exotic")
	Reptile      = Species("reptile")
	Rodent       = Species("rodent")
)

type Patient struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Species      Species    `json:"species"`
	Gender       Gender     `json:"gender"`
	Breed        *string    `json:"breed"`
	Birth        time.Time  `json:"birth"`
	Aggressive   bool       `json:"aggressive"`
	VaccinatedAt *time.Time `json:"vaccinated_at"`
	SterilizedAt *time.Time `json:"sterilized_at"`
}
