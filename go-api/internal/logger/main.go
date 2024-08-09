package logger

import "log/slog"

var impl Interface

type Interface interface {
	Info(err error)
	Error(err error)
}

func init() {
	impl = &Slogger{
		slog: slog.Default(),
	}
}

func Info(err error) {
	impl.Info(err)
}

func Error(err error) {
	impl.Error(err)
}
