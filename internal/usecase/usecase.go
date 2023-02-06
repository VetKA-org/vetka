package usecase

import (
	"context"

	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/internal/repo"
	uuid "github.com/satori/go.uuid"
)

type Auth interface {
	Login(ctx context.Context, username, password string) error
}

type Users interface {
	List(ctx context.Context) ([]entity.User, error)
	Register(ctx context.Context, login, password string) error
}

type Patients interface {
	List(ctx context.Context) ([]entity.Patient, error)
	Register(ctx context.Context) error
	ListAppointments(ctx context.CancelFunc, id uuid.UUID) error
}

type Appointments interface {
	List(ctx context.Context) ([]entity.Appointment, error)
	Create(ctx context.Context) error
	Update(ctx context.Context, id uuid.UUID) error
}

type Queue interface {
	List(ctx context.Context) ([]uuid.UUID, error)
	Enqueue(ctx context.Context, id uuid.UUID) error
	MoveUp(ctx context.Context, id uuid.UUID) error
	MoveDown(ctx context.Context, id uuid.UUID) error
	Dequeue(ctx context.Context, id uuid.UUID) error
}

type UseCases struct {
	Appointments Appointments
	Auth         Auth
	Patients     Patients
	Users        Users
	Queue        Queue
}

func New(repos *repo.Repositories) *UseCases {
	return &UseCases{
		Appointments: NewAppointmentsUseCase(repos.Appointments),
		Auth:         NewAuthUseCase(repos.Users),
		Patients:     NewPatientsUseCase(repos.Patients, repos.Appointments),
		Users:        NewUsersUseCase(repos.Users),
		Queue:        NewQueueUseCase(repos.Queue),
	}
}
