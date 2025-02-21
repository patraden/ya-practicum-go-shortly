package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// Logger is a wrapper around zerolog.Logger to provide structured logging.
type Logger struct {
	log zerolog.Logger
}

// NewLogger initializes and returns a new Logger instance with the specified log level.
func NewLogger(level zerolog.Level) *Logger {
	return &Logger{
		log: zerolog.New(os.Stdout).
			With().
			Timestamp().
			Logger().
			Level(level),
	}
}

// GetLogger returns the zerolog.Logger instance for custom log messages.
func (l *Logger) GetLogger() *zerolog.Logger {
	return &l.log
}
