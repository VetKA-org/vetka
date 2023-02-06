package repo

import (
	"context"
	"fmt"

	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/pkg/postgres"
	"github.com/jackc/pgx/v5"
	uuid "github.com/satori/go.uuid"
)

type UsersRepo struct {
	*postgres.Postgres
}

func NewUsersRepo(pg *postgres.Postgres) *UsersRepo {
	return &UsersRepo{pg}
}

func (r *UsersRepo) List(ctx context.Context) ([]entity.User, error) {
	rows, err := r.Pool.Query(ctx, "SELECT id, login FROM users")
	if err != nil {
		return nil, fmt.Errorf("UsersRepo - List - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var (
		id    uuid.UUID
		login string
	)

	rv := make([]entity.User, 0)
	_, err = pgx.ForEachRow(rows, []any{&id, &login}, func() error {
		rv = append(rv, entity.User{ID: id, Login: login})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("UsersRepo - List - r.Pool.ForEachRow: %w", err)
	}

	return rv, nil
}
