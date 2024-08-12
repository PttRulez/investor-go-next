package httpcontrollers

type Logger interface {
	Info(err error)
	Error(err error)
}
