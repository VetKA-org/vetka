package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/VetKA-org/vetka/pkg/postgres"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	uuid "github.com/satori/go.uuid"
)

type UsersRepo struct {
	log *logger.Logger
	*postgres.Postgres
}

func NewUsersRepo(log *logger.Logger, pg *postgres.Postgres) *UsersRepo {
	return &UsersRepo{log, pg}
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

func (r *UsersRepo) Register(ctx context.Context, login, password string) error {
	conn, err := r.Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("UsersRepo - Register - r.Pool.Acquire: %w", err)
	}

	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		return fmt.Errorf("UsersRepo - Register - conn.BeginTx: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			r.log.Error().Err(err).Msg("UsersRepo - Register - conn.Rollback")
		}
	}()

	if _, err := tx.Exec(
		ctx,
		"INSERT INTO users (login, password) VALUES ($1, crypt($2, gen_salt('bf', 8)))",
		login,
		password,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return entity.ErrUserExists
		}

		return fmt.Errorf("UsersRepo - Register - tx.Exec: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("UsersRepo - Register - tx.Commit: %w", err)
	}

	return nil
}
