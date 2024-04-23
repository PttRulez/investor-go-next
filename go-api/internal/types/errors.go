package types

import "errors"

var ErrNotYours = errors.New("not yours")
var ErrSmthWrong = errors.New("что-то пошло не так")
var ErrWrongId = errors.New("неправильный id")

type ErrSendToClient struct {
	msg string
}

func (e ErrSendToClient) Error() string {
	return e.msg
}

func NewErrSendToClient(msg string) ErrSendToClient {
	return ErrSendToClient{msg: msg}
}
