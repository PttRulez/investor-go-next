package http_controllers

import "log/slog"

var logger LoggerInterface

type LoggerImpl struct {
	slogger *slog.Logger
}

type LoggerInterface interface {
	Info(err error)
	Error(err error)
}

var slogger *slog.Logger

func init() {
	slogger = slog.Default()
}

func Info(err error) {
	slogger.Info(err.Error())
}

func Error(err error) {
	slogger.Error(err.Error())
}
