//go:build migrate

// Apply database migrations.
package app

import (
	"errors"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/golang-migrate/migrate/v4"
	// Migrate tools.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func runMigrations(url string) {
	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://migrations", url)
		if err == nil {
			break
		}

		log.Info().Msgf("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)

		attempts--
	}

	if err != nil {
		log.Fatal().Err(err).Msg("Migrate: postgres connect error")
	}

	err = m.Up()
	defer m.Close()

	if err != nil {
		log.Fatal().Err(err).Msg("Migrate: up error")
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Info().Msg("Migrate: no change")
		return
	}

	log.Info().Msg("Migrate: applying migrations: success")
}
