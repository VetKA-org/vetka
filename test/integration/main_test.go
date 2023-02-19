package integration_test

import (
	"context"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/stretchr/testify/suite"
)

const (
	_defaultStartTimeout = 5 * time.Second
	_defaultStopTimeout  = 5 * time.Second
)

func runCmd(ctx context.Context, path string, args ...string) error {
	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

func TestMain(m *testing.M) {
	log := logger.New("debug")
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := runCmd(
			ctx,
			"../../cmd/app/app",
			"-s", os.Getenv("SECRET"),
			"-d", os.Getenv("DATABASE_URI"),
		); err != nil {
			log.Error().Err(err).Msg("TestMain - runCmd")
		}
	}()

	time.Sleep(_defaultStartTimeout)

	code := m.Run()

	cancel()

	time.Sleep(_defaultStopTimeout)
	os.Exit(code)
}

func TestAuth(t *testing.T) { //nolint:paralleltest //sync test
	suite.Run(t, new(AuthTestSuite))
}

func TestUsers(t *testing.T) { //nolint:paralleltest //sync test
	suite.Run(t, new(UsersTestSuite))
}
