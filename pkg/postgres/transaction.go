package postgres

import (
	"context"
	"errors"

	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transaction struct {
	Tx   pgx.Tx
	conn *pgxpool.Conn
	log  *logger.Logger
}

func (t Transaction) Commit(ctx context.Context) error {
	defer t.conn.Release()

	defer func() {
		if err := t.Tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			t.log.Error().Err(err).Msg("Transaction - Commit - t.Tx.Rollback")
		}
	}()

	return t.Tx.Commit(ctx)
}

func (t Transaction) SendBatch(ctx context.Context, batch Batch) BatchResults {
	resp := t.Tx.SendBatch(ctx, batch.batch)

	return BatchResults{resp, t.log}
}
