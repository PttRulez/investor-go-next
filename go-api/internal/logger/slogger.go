package logger

import "log/slog"

type Slogger struct {
	slog *slog.Logger
}

func (l *Slogger) Info(err error) {
	l.slog.Info(err.Error())
}

func (l *Slogger) Error(err error) {
	l.slog.Error(err.Error())
}
