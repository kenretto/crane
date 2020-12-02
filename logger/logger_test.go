package logger

import (
	"github.com/kenretto/crane/configurator"
	"testing"
)

func TestLogger_Instance(t *testing.T) {
	var logger = new(Logger)
	var c, err = configurator.NewConfigurator("testdata/logger.yaml")
	if err != nil {
		t.Error(err)
	}
	c.Add("logger", logger)
	logger.Instance().Info("this is info message")
	logger.Instance().Warn("this is warn message")
	logger.Instance().Error("this is error message")
	logger.Instance().Trace("this is trace message")
}
