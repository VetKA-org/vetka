package repo

import "github.com/VetKA-org/vetka/pkg/postgres"

type AppointmentsRepo struct {
	*postgres.Postgres
}

func NewAppointmentsRepo(pg *postgres.Postgres) *AppointmentsRepo {
	return &AppointmentsRepo{pg}
}
