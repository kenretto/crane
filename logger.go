package crane

import (
	"github.com/kenretto/crane/server"
	"github.com/sirupsen/logrus"
)

// Logger encapsulating logrus
type Logger struct {
	*logrus.Entry
}

// WithFields encapsulating logrus
func (logger *Logger) WithFields(fields server.Fields) server.ILogger {
	return &Logger{logger.Logger.WithFields(logrus.Fields(fields))}
}
