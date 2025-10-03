package scuter

type Logger interface {
	Printf(format string, v ...any)
}
