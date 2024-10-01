package service

import (
	"errors"
)

var ErrdomainNotFound = errors.New("not found")

func NewArgumentsError(msg string) ArgumentsError {
	if msg == "" {
		msg = "переданы неверные аргументы"
	}
	return ArgumentsError{
		msg: msg,
	}
}

type ArgumentsError struct {
	msg string
}

func (e ArgumentsError) Error() string {
	return e.msg
}
