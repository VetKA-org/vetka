package entity

import uuid "github.com/satori/go.uuid"

type User struct {
	ID    uuid.UUID
	Login string
}
