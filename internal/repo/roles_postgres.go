package repo

import (
	"context"
	"fmt"

	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/pkg/postgres"
	uuid "github.com/satori/go.uuid"
)

type RolesRepo struct {
	*postgres.Postgres
}

func NewRolesRepo(pg *postgres.Postgres) *RolesRepo {
	return &RolesRepo{pg}
}

func (r *RolesRepo) List(ctx context.Context) ([]entity.Role, error) {
	rv := make([]entity.Role, 0)
	if err := r.Select(ctx, &rv, "SELECT role_id, name FROM roles"); err != nil {
		return nil, fmt.Errorf("RolesRepo - List - r.Select: %w", err)
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
			if postgres.IsForeignKeyViolation(err, "users_roles_role_id_fkey") {
				return entity.ErrRoleNotFound
			}

			return fmt.Errorf("RolesRepo - Assign - batchResp.Exec: %w", err)
		}
	}

	return nil
}

func (r *RolesRepo) Get(ctx context.Context, userID uuid.UUID) ([]entity.Role, error) {
	rv := make([]entity.Role, 0)
	if err := r.Select(
		ctx,
		&rv,
		`SELECT
         role_id, name
     FROM (
         SELECT
             role_id
         FROM
             users_roles
         WHERE
             user_id = $1
     ) AS r
     INNER JOIN
         roles
     USING (role_id)`,
		userID,
	); err != nil {
		return nil, fmt.Errorf("RolesRepo - List - r.Select: %w", err)
	}

	return rv, nil
}
