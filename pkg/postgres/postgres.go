// Package postgres implements postgres connection.
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func New(url string, log *logger.Logger) (*Postgres, error) {
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("database - Postgres.New - pgxpool.ParseConfig: %w", err)
	}

	pg := new(Postgres)
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
		return nil, fmt.Errorf("database - NewPostgres - attempts == 0: %w", err)
	}

	return pg, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
