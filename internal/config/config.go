package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
)

var errSecretNotSet = errors.New("secret key is required")

type Config struct {
	RunAddress  string           `env:"RUN_ADDRESS"`
	DatabaseURI entity.SecretURI `env:"DATABASE_URI"`
	RedisURI    entity.SecretURI `env:"REDIS_URI"`
	Secret      entity.Secret    `env:"SECRET"`
	LogLevel    string           `env:"LOG_LEVEL"`
}

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	runAddress := flag.StringP(
		"run-address",
		"a",
		"0.0.0.0:8080",
		"address:port of the service",
	)
	databaseURI := flag.StringP(
		"database-uri",
		"d",
		"postgres://postgres:postgres@127.0.0.1:5432/vetka?sslmode=disable",
		"full Postgres database connection URI",
	)
	redisURI := flag.StringP(
		"redis-uri",
		"r",
		"redis://:redis@127.0.0.1:6379/0",
		"full Redis connection URI",
	)
	secret := flag.StringP(
		"secret",
		"s",
		"",
		"secret key to sign JWT tokens",
	)
	logLevel := flag.StringP(
		"log-level",
		"l",
		"info",
		"log level of the service",
	)

	flag.Parse()

	cfg := &Config{
		RunAddress:  *runAddress,
		DatabaseURI: entity.SecretURI(*databaseURI),
		RedisURI:    entity.SecretURI(*redisURI),
		Secret:      entity.Secret(*secret),
		LogLevel:    *logLevel,
	}

	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	if cfg.Secret == "" {
		return nil, errSecretNotSet
	}

	return cfg, nil
}

func (c *Config) String() string {
	var sb strings.Builder

	sb.WriteString("Configuration:\n")
	sb.WriteString(fmt.Sprintf("\t\tListening address: %s\n", c.RunAddress))
	sb.WriteString(fmt.Sprintf("\t\tDatabase URI: %s\n", c.DatabaseURI))
	sb.WriteString(fmt.Sprintf("\t\tRedis URI: %s\n", c.RedisURI))
	sb.WriteString(fmt.Sprintf("\t\tSecret: %s\n", c.Secret))
	sb.WriteString(fmt.Sprintf("\t\tLog level: %s\n", c.LogLevel))

	return sb.String()
}
