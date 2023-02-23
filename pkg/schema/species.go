package schema

import "github.com/VetKA-org/vetka/pkg/entity"

type ListSpeciesResponse struct {
	Data []entity.Species `json:"data"`
}
