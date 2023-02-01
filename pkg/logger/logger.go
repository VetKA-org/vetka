// Package logger configures logging facility.
package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

func New(level string) *Logger {
	var lv zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		lv = zerolog.ErrorLevel
	case "warn":
		lv = zerolog.WarnLevel
	case "info":
		lv = zerolog.InfoLevel
	case "debug":
		lv = zerolog.DebugLevel
	default:
		lv = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(lv)

	output := zerolog.ConsoleWriter{Out: os.Stdout}
	output.TimeFormat = time.RFC822

	skipFrameCount := 3
	l := zerolog.New(output).
		With().
		Timestamp().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
		Logger()

	return &Logger{l}
}
