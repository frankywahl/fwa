package fwa

// Logger is used by the worker to write logs
type Logger interface {
	Infof(string, ...interface{})
}

type noopLogger struct{}

func (l *noopLogger) Infof(string, ...interface{}) {}
