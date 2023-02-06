package entity

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

var ErrUserExists = errors.New("user already exists")

type User struct {
	ID       uuid.UUID `json:"id"`
	Login    string    `json:"login"`
	Password string    `json:"-"`
}
