package repo

import "github.com/VetKA-org/vetka/pkg/postgres"

type PatientsRepo struct {
	*postgres.Postgres
}

func NewPatientsRepo(pg *postgres.Postgres) *PatientsRepo {
	return &PatientsRepo{pg}
}
