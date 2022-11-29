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

// Errorf logs an error level message
func Errorf(t string, args ...interface{}) {
	logger.Errorf(t, args...)
}

// Infof logs an information level message
func Infof(t string, args ...interface{}) {
	logger.Infof(t, args...)
}

// Debugf logs a debug level message
func Debugf(t string, args ...interface{}) {
	logger.Debugf(t, args...)
}
