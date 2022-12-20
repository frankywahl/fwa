package fwa

// Logger is used by the worker to write logs
type Logger interface {
	Info(...interface{})
}

type noopLogger struct{}

func (l *noopLogger) Debugf(string, ...interface{}) {}
func (l *noopLogger) Infof(string, ...interface{})  {}
func (l *noopLogger) Errorf(string, ...interface{}) {}
func (l *noopLogger) Debug(...interface{})          {}
func (l *noopLogger) Info(...interface{})           {}
func (l *noopLogger) Error(...interface{})          {}
