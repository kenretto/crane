package redis

import "log"

// ILogger logger interface
type ILogger interface {
	Warn(args ...interface{})
	Error(args ...interface{})
}

// DefaultLogger default logger
type DefaultLogger struct {
}

// Error error...
func (d DefaultLogger) Error(args ...interface{}) {
	log.Println(args...)
}

// Warn warn...
func (d DefaultLogger) Warn(args ...interface{}) {
	log.Println(args...)
}

// NewDefaultLogger default logger
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{}
}
