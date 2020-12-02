package server

import (
	"fmt"
	"log"
	"strings"
)

type Fields map[string]interface{}

type ILogger interface {
	Println(args ...interface{})
	Error(args ...interface{})
	Info(args ...interface{})
	Fatalln(args ...interface{})
	WithFields(fields Fields) ILogger
}

type DefaultLogger struct {
	args Fields
}

func NewDefaultLogger(args Fields) *DefaultLogger {
	return &DefaultLogger{args: args}
}

func (d *DefaultLogger) Println(args ...interface{}) {
	var b strings.Builder
	for key, value := range d.args {
		b.WriteString(fmt.Sprintf("%s:%+v    ", key, value))
	}

	var msg = []interface{}{b.String()}
	msg = append(msg, args...)
	log.Println(msg...)
}

func (d *DefaultLogger) Error(args ...interface{}) {
	var b strings.Builder
	for key, value := range d.args {
		b.WriteString(fmt.Sprintf("%s:%+v    ", key, value))
	}

	var msg = []interface{}{b.String()}
	msg = append(msg, args...)
	log.Println(msg...)
}

func (d *DefaultLogger) Info(args ...interface{}) {
	var b strings.Builder
	for key, value := range d.args {
		b.WriteString(fmt.Sprintf("%s:%+v    ", key, value))
	}

	var msg = []interface{}{b.String()}
	msg = append(msg, args...)
	log.Println(msg...)
}

func (d *DefaultLogger) Fatalln(args ...interface{}) {
	var b strings.Builder
	for key, value := range d.args {
		b.WriteString(fmt.Sprintf("%s:%+v    ", key, value))
	}

	var msg = []interface{}{b.String()}
	msg = append(msg, args...)
	log.Fatalln(msg...)
}

func (d DefaultLogger) WithFields(fields Fields) ILogger {
	var logger = NewDefaultLogger(make(Fields))
	for key, value := range fields {
		logger.args[key] = value
	}
	return logger
}
