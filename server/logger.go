package server

import (
	"fmt"
	"log"
	"strings"
)

// Fields other info
type Fields map[string]interface{}

// ILogger methods to be implemented for loggers injected into HTTPServer
type ILogger interface {
	Println(args ...interface{})
	Error(args ...interface{})
	Info(args ...interface{})
	Fatalln(args ...interface{})
	WithFields(fields Fields) ILogger
}

// DefaultLogger default logger
type DefaultLogger struct {
	args Fields
}

// NewDefaultLogger default logger
func NewDefaultLogger(args Fields) *DefaultLogger {
	return &DefaultLogger{args: args}
}

// Println println...
func (d *DefaultLogger) Println(args ...interface{}) {
	var b strings.Builder
	for key, value := range d.args {
		b.WriteString(fmt.Sprintf("%s:%+v    ", key, value))
	}

	var msg = []interface{}{b.String()}
	msg = append(msg, args...)
	log.Println(msg...)
}

// Error error...
func (d *DefaultLogger) Error(args ...interface{}) {
	var b strings.Builder
	for key, value := range d.args {
		b.WriteString(fmt.Sprintf("%s:%+v    ", key, value))
	}

	var msg = []interface{}{b.String()}
	msg = append(msg, args...)
	log.Println(msg...)
}

// Info info...
func (d *DefaultLogger) Info(args ...interface{}) {
	var b strings.Builder
	for key, value := range d.args {
		b.WriteString(fmt.Sprintf("%s:%+v    ", key, value))
	}

	var msg = []interface{}{b.String()}
	msg = append(msg, args...)
	log.Println(msg...)
}

// Fatalln fatalln
func (d *DefaultLogger) Fatalln(args ...interface{}) {
	var b strings.Builder
	for key, value := range d.args {
		b.WriteString(fmt.Sprintf("%s:%+v    ", key, value))
	}

	var msg = []interface{}{b.String()}
	msg = append(msg, args...)
	log.Fatalln(msg...)
}

// WithFields Set other information to be recorded in the log, and a new object will be returned
func (d DefaultLogger) WithFields(fields Fields) ILogger {
	var logger = NewDefaultLogger(make(Fields))
	for key, value := range fields {
		logger.args[key] = value
	}
	return logger
}
