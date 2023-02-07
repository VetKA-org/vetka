package config

import (
	"errors"

	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
)

var errSecretNotSet = errors.New("secret key is required")

type Config struct {
	RunAddress  string `env:"RUN_ADDRESS"`
	DatabaseURI string `env:"DATABASE_URI"`
	RedisURI    string `env:"REDIS_URI"`
	Secret      string `env:"SECRET"`
	LogLevel    string `env:"LOG_LEVEL"`
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
		"full database connection URI",
	)
	redisURI := flag.StringP(
		"redis-uri",
		"r",
		"redis://127.0.0.1:6379",
		"full redis connection URI",
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
		DatabaseURI: *databaseURI,
		RedisURI:    *redisURI,
		Secret:      *secret,
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
