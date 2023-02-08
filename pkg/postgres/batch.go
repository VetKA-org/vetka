package postgres

import (
	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/jackc/pgx/v5"
)

type Batch struct {
	batch *pgx.Batch
}

func (b Batch) Queue(query string, arguments ...any) {
	b.batch.Queue(query, arguments...)
}

type BatchResults struct {
	results pgx.BatchResults
	log     *logger.Logger
}

func (b BatchResults) Exec() error {
	_, err := b.results.Exec()

	return err
}

func (b BatchResults) Close() {
	if err := b.results.Close(); err != nil {
		b.log.Error().Err(err).Msg("BatchResults - Close - r.batch.Close")
	}
}
