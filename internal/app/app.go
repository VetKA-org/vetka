// Package app configures and runs application.
package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VetKA-org/vetka/internal/config"
	v1 "github.com/VetKA-org/vetka/internal/controller/http/v1"
	"github.com/VetKA-org/vetka/internal/repo"
	"github.com/VetKA-org/vetka/internal/usecase"
	"github.com/VetKA-org/vetka/pkg/httpserver"
	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/VetKA-org/vetka/pkg/postgres"
	"github.com/VetKA-org/vetka/pkg/redis"
	"github.com/gin-gonic/gin"
)

const _defaultShutdownTimeout = 10 * time.Second

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	log := logger.New(cfg.LogLevel)
	log.Info().Msg(cfg.String())

	// Redis connection.
	rdb, err := redis.New(cfg.RedisURI)
	if err != nil {
		log.Fatal().Err(err).Msg("app - Run - redis.New")
	}

	// Postgres DB connection.
	pg, err := postgres.New(cfg.DatabaseURI, log)
	if err != nil {
		log.Fatal().Err(err).Msg("app - Run - database.New")
	}

	// Domain.
	repos := repo.New(log, pg, rdb)
	useCases := usecase.New(log, cfg, repos)

	// HTTP Server.
	handler := gin.New()
	v1.NewRouter(handler, log, cfg, useCases)
	httpServer := httpserver.New(handler, cfg.RunAddress)

	// Waiting signal.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	select {
	case s := <-interrupt:
		log.Info().Msg("app - Run - interrupt: signal " + s.String())
	case err = <-httpServer.Notify():
		log.Error().Err(err).Msg("app - Run - httpServer.Notify")
	}

	// Shutdown
	log.Info().Msg("Shutting down...")

	stopped := make(chan struct{})

	ctx, cancel := context.WithTimeout(context.Background(), _defaultShutdownTimeout)
	defer cancel()

	go shutdown(ctx, stopped, log, httpServer, pg, rdb)

	select {
	case <-stopped:
		log.Info().Msg("Server shutdown successful")

	case <-ctx.Done():
		log.Warn().Msgf("Exceeded %s shutdown timeout, exit forcibly", _defaultShutdownTimeout)
	}
}

func shutdown(
	ctx context.Context,
	notify chan<- struct{},
	log *logger.Logger,
	httpServer *httpserver.Server,
	pg *postgres.Postgres,
	rdb *redis.Redis,
) {
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("app - shutdown - httpServer.Shutdown")
	}

	if err := rdb.Close(); err != nil {
		log.Error().Err(err).Msg("app - shutdown - rdb.Close")
	}

	pg.Close()

	close(notify)
}
