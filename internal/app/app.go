// Package app configures and runs application.
package app

import (
	"os"
	"os/signal"
	"syscall"

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

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	log := logger.New(cfg.LogLevel)

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
	defer pg.Close()

	// Domain.
	repos := repo.New(pg, rdb)
	useCases := usecase.New(repos)

	// HTTP Server.
	handler := gin.New()
	v1.NewRouter(handler, log, useCases)
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
		log.Info().Msg("app - Run - interrupt: " + s.String())
	case err = <-httpServer.Notify():
		log.Error().Err(err).Msg("app - Run - httpServer.Notify")
	}

	// Shutdown
	log.Info().Msg("Shutting down...")

	err = httpServer.Shutdown()
	if err != nil {
		log.Error().Err(err).Msg("app - Run - httpServer.Shutdown")
	}

	err = rdb.Close()
	if err != nil {
		log.Error().Err(err).Msg("app - Run - redis.Close")
	}
}
