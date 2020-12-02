package captcha

import "log"

// ILogger Log interface required by this package
type ILogger interface {
	Error(args ...interface{})
}

// DefaultLogger The default logger implementation
type DefaultLogger struct {
}

func (d DefaultLogger) Error(args ...interface{}) {
	log.Println(args...)
}

// NewDefaultLogger default logger
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{}
}
