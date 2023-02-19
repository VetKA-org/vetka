package schema

import "github.com/VetKA-org/vetka/pkg/entity"

type ListRolesResponse struct {
	Data []entity.Role `json:"data"`
}
