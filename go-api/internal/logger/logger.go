package logger

import (
	"log/slog"
)

type Logger struct {
	Slog *slog.Logger
}

func NewLogger(log *slog.Logger) *Logger {
	if log == nil {
		log = slog.Default()
	}
	return &Logger{
		Slog: log,
	}
}

func (l *Logger) Debug(s string) {
	l.Slog.Debug(s)
}

func (l *Logger) Error(err error) {
	l.Slog.Error(err.Error())
}

func (l *Logger) Info(s string) {
	l.Slog.Info(s)
}
