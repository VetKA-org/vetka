package repo

import (
	"context"
	"fmt"

	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/pkg/postgres"
	"github.com/jackc/pgx/v5"
	uuid "github.com/satori/go.uuid"
)

type RolesRepo struct {
	*postgres.Postgres
}

func NewRolesRepo(pg *postgres.Postgres) *RolesRepo {
	return &RolesRepo{pg}
}

func (r *RolesRepo) List(ctx context.Context) ([]entity.Role, error) {
	rows, err := r.Pool.Query(ctx, "SELECT id, name FROM roles")
	if err != nil {
		return nil, fmt.Errorf("RolesRepo - List - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var (
		id   uuid.UUID
		name string
	)

	rv := make([]entity.Role, 0)
	_, err = pgx.ForEachRow(rows, []any{&id, &name}, func() error {
		rv = append(rv, entity.Role{ID: id, Name: name})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("RolesRepo - List - r.Pool.ForEachRow: %w", err)
	}

	return rv, nil
}

func (r *RolesRepo) Assign(
	ctx context.Context,
	tx postgres.Transaction,
	userID uuid.UUID,
	roles []uuid.UUID,
) error {
	batch := r.NewBatch()

	for _, roleID := range roles {
		batch.Queue(
			"INSERT INTO users_roles (user_id, role_id) VALUES ($1, $2)",
			userID,
			roleID,
		)
	}

	batchResp := tx.SendBatch(ctx, batch)
	defer batchResp.Close()

	for i := 0; i < len(roles); i++ {
		if err := batchResp.Exec(); err != nil {
			return fmt.Errorf("RolesRepo - Assign - batchResp.Exec: %w", err)
		}
	}

	return nil
}
