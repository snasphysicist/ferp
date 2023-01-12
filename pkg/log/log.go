package log

import (
	"go.uber.org/zap"
)

// logger is the global logger instance
// TODO: avoid package level mutable state
var logger *zap.SugaredLogger

// Initialise initialises the logger
func Initialise() (func(), error) {
	l, err := zap.NewProduction()
	if err != nil {
		return func() {}, err
	}
	logger = l.Sugar()
	return func() { _ = l.Sync() }, nil
}

// L provides safe access to the global logger
func L() Logger {
	return logger
}

// Logger is returned from this package instead of the full logger
// to limit the operations client code is allowed to take on it
type Logger interface {
	Errorf(string, ...interface{})
	Infof(string, ...interface{})
	Debugf(string, ...interface{})
}
