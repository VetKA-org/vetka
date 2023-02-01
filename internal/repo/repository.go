package repo

import (
	"context"

	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/pkg/postgres"
	"github.com/VetKA-org/vetka/pkg/redis"
)

type Appointments interface{}

type Patients interface{}

type Users interface {
	List(ctx context.Context) ([]entity.User, error)
}

type Queue interface{}

type Repositories struct {
	Appointments Appointments
	Patients     Patients
	Users        Users
	Queue        Queue
}

func New(pg *postgres.Postgres, rdb *redis.Redis) *Repositories {
	return &Repositories{
		Appointments: NewAppointmentsRepo(pg),
		Patients:     NewPatientsRepo(pg),
		Users:        NewUsersRepo(pg),
		Queue:        NewQueueRepo(rdb),
	}
}
