package repo

import (
	"context"
	"fmt"

	"github.com/VetKA-org/vetka/pkg/entity"
	"github.com/VetKA-org/vetka/pkg/postgres"
	uuid "github.com/satori/go.uuid"
)

type UsersRepo struct {
	*postgres.Postgres
}

func NewUsersRepo(pg *postgres.Postgres) *UsersRepo {
	return &UsersRepo{pg}
}

func (r *UsersRepo) List(ctx context.Context) ([]entity.User, error) {
	rv := make([]entity.User, 0)
	if err := r.Select(ctx, &rv, "SELECT user_id, login FROM users"); err != nil {
		return nil, fmt.Errorf("UsersRepo - List - r.Select: %w", err)
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
		"INSERT INTO users (login, password) VALUES ($1, crypt($2, gen_salt('bf', 8))) RETURNING user_id",
		login,
		password,
	).Scan(&id)
	if err != nil {
		if postgres.IsEntityExists(err) {
			return uuid.UUID{}, entity.ErrUserExists
		}

		return uuid.UUID{}, fmt.Errorf("UsersRepo - Register - tx.Tx.QueryRow.Scan: %w", err)
	}

	return id, nil
}

func (r *UsersRepo) Verify(ctx context.Context, login, password string) (entity.User, error) {
	var user entity.User

	err := r.Pool.
		QueryRow(
			ctx,
			"SELECT user_id, login FROM users WHERE login=$1 AND password = crypt($2, password)",
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
