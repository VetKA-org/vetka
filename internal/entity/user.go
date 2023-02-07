package entity

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid login or password")
	ErrUserExists         = errors.New("user already exists")
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Login    string    `json:"login"`
	Password string    `json:"-"`
}
