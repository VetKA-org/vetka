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

func (r *UsersRepo) Register(
	ctx context.Context,
	tx postgres.Transaction,
	login, password string,
) (uuid.UUID, error) {
	var id uuid.UUID

	err := tx.Tx.QueryRow(
		ctx,
		"INSERT INTO users (login, password) VALUES ($1, crypt($2, gen_salt('bf', 8))) RETURNING id",
		login,
		password,
	).Scan(&id)
	if err != nil {
		if postgres.IsEntityExists(err) {
			return uuid.UUID{}, entity.ErrUserExists
		}

		return uuid.UUID{}, fmt.Errorf("UsersRepo - Register - tx.Tx.Exec: %w", err)
	}

	return id, nil
}

func (r *UsersRepo) Verify(ctx context.Context, login, password string) (entity.User, error) {
	var user entity.User

	err := r.Pool.
		QueryRow(
			ctx,
			"SELECT id, login FROM users WHERE login=$1 AND password = crypt($2, password)",
			login,
			password,
		).
		Scan(&user.ID, &user.Login)
	if err != nil {
		if postgres.IsEmptyResponse(err) {
			return entity.User{}, entity.ErrInvalidCredentials
		}

		return entity.User{}, fmt.Errorf("UsersRepo - Verify - r.Pool.QueryRow.Scan: %w", err)
	}

	return user, nil
}
