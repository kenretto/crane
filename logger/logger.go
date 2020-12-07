package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/medivh-jay/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"sync"
	"time"
)

// Logger logger
type Logger struct {
	Config struct {
		MaxAge       string `mapstructure:"max_age"`
		RotationTime string `mapstructure:"rotation_time"`
		Level        string `mapstructure:"level"`
		Path         string `mapstructure:"path"`
		ReportCaller bool   `mapstructure:"report_caller"`
	}

	logger     *logrus.Logger
	FilenameFn func(path, level string) string
	rw         sync.RWMutex
}

func (l *Logger) Node() string {
	return "logger"
}

// LogLevel logger level
func (l *Logger) LogLevel() logrus.Level {
	return map[string]logrus.Level{
		"panic": logrus.PanicLevel,
		"fatal": logrus.FatalLevel,
		"error": logrus.ErrorLevel,
		"warn":  logrus.WarnLevel,
		"info":  logrus.InfoLevel,
		"debug": logrus.DebugLevel,
		"trace": logrus.TraceLevel,
	}[l.Config.Level]
}

func (l *Logger) parseDuration(s string, defaultValue time.Duration) time.Duration {
	duration, err := time.ParseDuration(s)
	if err != nil {
		return defaultValue
	}
	return duration
}

func (l *Logger) defaultFilenameFn(path, level string) string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s/%s.%%Y%%m%%d.%s.log", path, level, hostname)
}

func (l *Logger) newWriter(level string) *rotatelogs.RotateLogs {
	var (
		writer *rotatelogs.RotateLogs
		err    error
	)
	if l.FilenameFn == nil {
		l.FilenameFn = l.defaultFilenameFn
	}

	writer, err = rotatelogs.New(
		l.FilenameFn(l.Config.Path, level),
		rotatelogs.WithMaxAge(l.parseDuration(l.Config.MaxAge, time.Duration(30*86400)*time.Second)),
		rotatelogs.WithRotationTime(l.parseDuration(l.Config.RotationTime, time.Hour)),
	)

	if err != nil {
		panic(err)
	}

	return writer
}

func (l *Logger) newInstance() {
	l.logger = logrus.New()
	l.logger.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	l.logger.SetNoLock()
	l.logger.Level = l.LogLevel()

	l.logger.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  l.newWriter("info"),
			logrus.ErrorLevel: l.newWriter("error"),
			logrus.DebugLevel: l.newWriter("debug"),
			logrus.WarnLevel:  l.newWriter("warn"),
			logrus.TraceLevel: l.newWriter("trace"),
		}, &logrus.JSONFormatter{}))

	l.logger.SetOutput(nilWriter{})
	l.logger.SetReportCaller(l.Config.ReportCaller)
}

// OnChange When the configuration file changes, the logger will be reinitialized
func (l *Logger) OnChange(viper *viper.Viper) {
	l.rw.Lock()
	defer l.rw.Unlock()
	_ = viper.Unmarshal(&l.Config)
	l.newInstance()
}

// Instance get logger instance
func (l *Logger) Instance() *logrus.Logger {
	l.rw.RLock()
	defer l.rw.RUnlock()
	return l.logger
}

// disable terminal output after adding hook for file writing
type nilWriter struct {
}

func (nilWriter) Write(_ []byte) (n int, err error) {
	return 0, nil
}
