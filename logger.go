package crane

import (
	"github.com/kenretto/crane/server"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func (logger *Logger) WithFields(fields server.Fields) server.ILogger  {
	return &Logger{logger.Logger.WithFields(logrus.Fields(fields)).Logger}
}