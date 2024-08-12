package logger

import "log/slog"

type Logger interface {
	Info(err error)
	Error(err error)
}

type Log struct {
	slog *slog.Logger
}

func NewLogger() Logger {
	return &Log{
		slog: slog.Default(),
	}
}

func (l *Log) Info(err error) {
	l.slog.Info(err.Error())
}

func (l *Log) Error(err error) {
	l.slog.Error(err.Error())
}
