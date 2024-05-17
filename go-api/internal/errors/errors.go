package errors

import (
	"errors"
)

var ErrNotFound = errors.New("у нас такого нет")
var ErrNotYours = errors.New("не ваше")
var ErrSmthWrong = errors.New("что-то пошло не так")

type ErrSendToClient struct {
	msg    string
	Status int
}

func (e ErrSendToClient) Error() string {
	return e.msg
}

func NewErrSendToClient(msg string, status int) ErrSendToClient {
	return ErrSendToClient{msg: msg, Status: status}
}
