package schema

import (
	"time"

	"github.com/VetKA-org/vetka/pkg/entity"
	uuid "github.com/satori/go.uuid"
)

type RegisterPatientRequest struct {
	Name       string        `json:"name" binding:"required,max=32"`
	SpeciesID  uuid.UUID     `json:"species_id" binding:"required"`
	Gender     entity.Gender `json:"gender" binding:"required,oneof=male female"`
	Breed      string        `json:"breed" binding:"required,max=64"`
	Birth      time.Time     `json:"birth" time_format:"2006-01-02" time_utc:"1" binding:"required"`
	Aggressive bool          `json:"aggressive"`
	// NB (alkurbatov): Ask lll to ignore struct tags, see:
	// https://github.com/walle/lll/issues/11
	VaccinatedAt *time.Time `json:"vaccinated_at" time_format:"2006-01-02" time_utc:"1" binding:"omitempty,gtefield=Birth"` //nolint:lll // tags
	SterilizedAt *time.Time `json:"sterilized_at" time_format:"2006-01-02" time_utc:"1" binding:"omitempty,gtefield=Birth"` //nolint:lll // tags
}

type RegisterPatientResponse struct {
	ID uuid.UUID `json:"id"`
}

type ListPatientsResponse struct {
	Data []entity.Patient `json:"data"`
}
