package repo

import (
	"context"
	"time"

	"github.com/VetKA-org/vetka/pkg/entity"
	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/VetKA-org/vetka/pkg/postgres"
	"github.com/VetKA-org/vetka/pkg/redis"
	uuid "github.com/satori/go.uuid"
)

type Appointments interface {
	BeginTx(ctx context.Context) (postgres.Transaction, error)

	List(ctx context.Context, patientID *uuid.UUID) ([]entity.Appointment, error)
	Create(
		ctx context.Context,
		tx postgres.Transaction,
		patientID uuid.UUID,
		assigneeID uuid.UUID,
		scheduledFor time.Time,
		reason string,
		details *string,
	) error
	Update(
		ctx context.Context,
		tx postgres.Transaction,
		id uuid.UUID,
		status entity.ApptStatus,
	) error
}

type Patients interface {
	BeginTx(ctx context.Context) (postgres.Transaction, error)

	List(ctx context.Context) ([]entity.Patient, error)
	Register(
		ctx context.Context,
		tx postgres.Transaction,
		name string,
		speciesID uuid.UUID,
		gender entity.Gender,
		breed string,
		birth time.Time,
		aggressive bool,
		vaccinatedAt *time.Time,
		sterilizedAt *time.Time,
	) error
}

type Users interface {
	BeginTx(ctx context.Context) (postgres.Transaction, error)

	List(ctx context.Context) ([]entity.User, error)
	Register(ctx context.Context, tx postgres.Transaction, login, password string) (uuid.UUID, error)
	Verify(ctx context.Context, login, password string) (entity.User, error)
}

type Roles interface {
	BeginTx(ctx context.Context) (postgres.Transaction, error)

	List(ctx context.Context, name string) ([]entity.Role, error)
	Assign(ctx context.Context, tx postgres.Transaction, userID uuid.UUID, roles []uuid.UUID) error
	Get(ctx context.Context, userID uuid.UUID) ([]entity.Role, error)
}

type Species interface {
	List(ctx context.Context) ([]entity.Species, error)
}

type Queue interface {
	List(ctx context.Context) ([]uuid.UUID, error)
	Enqueue(ctx context.Context, id uuid.UUID) error
	Dequeue(ctx context.Context, id uuid.UUID) error
	MoveUp(ctx context.Context, id uuid.UUID) error
	MoveDown(ctx context.Context, id uuid.UUID) error
}

type Repositories struct {
	Appointments Appointments
	Patients     Patients
	Users        Users
	Roles        Roles
	Species      Species
	Queue        Queue
}

func New(log *logger.Logger, pg *postgres.Postgres, rdb *redis.Redis) *Repositories {
	return &Repositories{
		Appointments: NewAppointmentsRepo(pg),
		Patients:     NewPatientsRepo(pg),
		Users:        NewUsersRepo(pg),
		Roles:        NewRolesRepo(pg),
		Species:      NewSpeciesRepo(pg),
		Queue:        NewQueueRepo(rdb),
	}
}
