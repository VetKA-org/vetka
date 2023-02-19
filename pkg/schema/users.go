package schema

import uuid "github.com/satori/go.uuid"

type RegisterUserRequest struct {
	Login    string      `json:"login" binding:"required,max=128"`
	Password string      `json:"password" binding:"required"`
	Roles    []uuid.UUID `json:"roles" binding:"required"`
}
