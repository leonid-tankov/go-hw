package repository

type Logger interface {
	Debug(msg string, a ...any)
	Info(msg string, a ...any)
	Error(msg string, a ...any)
	Fatal(msg string, a ...any)
}
