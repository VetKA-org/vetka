package repo

import (
	"context"

	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/VetKA-org/vetka/pkg/postgres"
	"github.com/VetKA-org/vetka/pkg/redis"
	uuid "github.com/satori/go.uuid"
)

type Appointments interface {
	BeginTx(ctx context.Context) (postgres.Transaction, error)
}

type Patients interface {
	BeginTx(ctx context.Context) (postgres.Transaction, error)

	List(ctx context.Context) ([]entity.Patient, error)
}

type Users interface {
	BeginTx(ctx context.Context) (postgres.Transaction, error)

	List(ctx context.Context) ([]entity.User, error)
	Register(ctx context.Context, tx postgres.Transaction, login, password string) (uuid.UUID, error)
	Verify(ctx context.Context, login, password string) (entity.User, error)
}

type Roles interface {
	BeginTx(ctx context.Context) (postgres.Transaction, error)

	List(ctx context.Context) ([]entity.Role, error)
	Assign(ctx context.Context, tx postgres.Transaction, userID uuid.UUID, roles []uuid.UUID) error
	Get(ctx context.Context, userID uuid.UUID) ([]entity.Role, error)
}

type Queue interface{}

type Repositories struct {
	Appointments Appointments
	Patients     Patients
	Users        Users
	Roles        Roles
	Queue        Queue
}

func New(log *logger.Logger, pg *postgres.Postgres, rdb *redis.Redis) *Repositories {
	return &Repositories{
		Appointments: NewAppointmentsRepo(pg),
		Patients:     NewPatientsRepo(pg),
		Users:        NewUsersRepo(pg),
		Roles:        NewRolesRepo(pg),
		Queue:        NewQueueRepo(rdb),
	}
}
