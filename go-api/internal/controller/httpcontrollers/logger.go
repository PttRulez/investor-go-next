package httpcontrollers

type Logger interface {
	Info(s string)
	Error(err error)
}
