package logger

import (
	"log/slog"
)

type Log struct {
	slog *slog.Logger
}

func NewLogger() *Log {
	return &Log{
		slog: slog.Default(),
	}
}

func (l *Log) Debug(s string) {
	l.slog.Debug(s)
}

func (l *Log) Error(err error) {
	l.slog.Error(err.Error())
}

func (l *Log) Info(s string) {
	l.slog.Info(s)
}
