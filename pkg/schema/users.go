package schema

import (
	"github.com/VetKA-org/vetka/pkg/entity"
	uuid "github.com/satori/go.uuid"
)

type RegisterUserRequest struct {
	Login    string      `json:"login" binding:"required,max=128"`
	Password string      `json:"password" binding:"required"`
	Roles    []uuid.UUID `json:"roles" binding:"required"`
}

type ListUsersResponse struct {
	Data []entity.User `json:"data"`
}
