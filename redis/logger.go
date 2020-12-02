package redis

import "log"

type ILogger interface {
	Warn(args ...interface{})
	Error(args ...interface{})
}

type DefaultLogger struct {
}

func (d DefaultLogger) Error(args ...interface{}) {
	log.Println(args...)
}

func (d DefaultLogger) Warn(args ...interface{}) {
	log.Println(args...)
}

func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{}
}
