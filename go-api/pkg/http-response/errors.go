package httpresponse

import (
	"errors"
)

var ErrNotYours = errors.New("not yours")
var ErrSmthWrong = errors.New("что-то пошло не так")
var ErrWrongId = errors.New("неправильный id")

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
