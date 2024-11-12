package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	log zerolog.Logger
}

func NewLogger(level zerolog.Level) *Logger {
	return &Logger{
		log: zerolog.New(os.Stdout).
			With().
			Timestamp().
			Logger().
			Level(level),
	}
}

func (l *Logger) GetLogger() *zerolog.Logger {
	return &l.log
}
