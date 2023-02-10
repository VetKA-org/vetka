// Package postgres implements postgres connection.
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

type Postgres struct {
	Pool *pgxpool.Pool
	log  *logger.Logger
}

func New(url string, log *logger.Logger) (*Postgres, error) {
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("Postgres - New - pgxpool.ParseConfig: %w", err)
	}

	pg := &Postgres{log: log}
	attempts := _defaultConnAttempts

	for attempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Info().Msgf("Postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultConnTimeout)

		attempts--
	}

	if err != nil {
		return nil, fmt.Errorf("Postgres - New - attempts == 0: %w", err)
	}

	return pg, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

func (p *Postgres) BeginTx(ctx context.Context) (Transaction, error) {
	conn, err := p.Pool.Acquire(ctx)
	if err != nil {
		return Transaction{}, fmt.Errorf("Postgres - BeginTx - r.Pool.Acquire: %w", err)
	}

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return Transaction{}, fmt.Errorf("Postgres - BeginTx - conn.BeginTx: %w", err)
	}

	return Transaction{tx, conn, p.log}, nil
}

func (p *Postgres) NewBatch() Batch {
	pb := new(pgx.Batch)

	return Batch{pb}
}

func (p *Postgres) Select(
	ctx context.Context,
	dst interface{},
	query string,
	args ...interface{},
) error {
	return pgxscan.Select(ctx, p.Pool, dst, query, args...)
}
