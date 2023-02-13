package usecase

import (
	"context"
	"time"

	"github.com/VetKA-org/vetka/internal/config"
	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/internal/repo"
	uuid "github.com/satori/go.uuid"
)

type Appointments interface {
	List(ctx context.Context, patientID *uuid.UUID) ([]entity.Appointment, error)
	Create(
		ctx context.Context,
		patientID uuid.UUID,
		assigneeID uuid.UUID,
		scheduledFor time.Time,
		reason string,
		details *string,
	) error
	Update(ctx context.Context, id uuid.UUID) error
}

type Auth interface {
	Login(ctx context.Context, login, password string) (entity.JWTToken, error)
}

type Patients interface {
	List(ctx context.Context) ([]entity.Patient, error)
	Register(
		ctx context.Context,
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

type Queue interface {
	List(ctx context.Context) ([]uuid.UUID, error)
	Enqueue(ctx context.Context, id uuid.UUID) error
	MoveUp(ctx context.Context, id uuid.UUID) error
	MoveDown(ctx context.Context, id uuid.UUID) error
	Dequeue(ctx context.Context, id uuid.UUID) error
}

type Roles interface {
	List(ctx context.Context) ([]entity.Role, error)
}

type Species interface {
	List(ctx context.Context) ([]entity.Species, error)
}

type Users interface {
	List(ctx context.Context) ([]entity.User, error)
	Register(ctx context.Context, login, password string, roles []uuid.UUID) error
}

type UseCases struct {
	Appointments Appointments
	Auth         Auth
	Patients     Patients
	Queue        Queue
	Roles        Roles
	Species      Species
	Users        Users
}

func New(cfg *config.Config, repos *repo.Repositories) *UseCases {
	return &UseCases{
		Appointments: NewAppointmentsUseCase(repos.Appointments),
		Auth:         NewAuthUseCase(cfg.Secret, repos.Users, repos.Roles),
		Patients:     NewPatientsUseCase(repos.Patients),
		Queue:        NewQueueUseCase(repos.Queue),
		Roles:        NewRolesUseCase(repos.Roles),
		Species:      NewSpeciesUseCase(repos.Species),
		Users:        NewUsersUseCase(repos.Users, repos.Roles),
	}
}
