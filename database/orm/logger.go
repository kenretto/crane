package orm

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

var logLevel = map[string]logger.LogLevel{
	"silent": logger.Silent,
	"error":  logger.Error,
	"warn":   logger.Warn,
	"info":   logger.Info,
}

type iLogger struct {
	logger        *logrus.Entry
	level         logger.LogLevel
	SlowThreshold time.Duration
}

// LogMode 设置日志级别
func (log *iLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *log
	newLogger.level = level
	newLogger.SlowThreshold = 100 * time.Millisecond
	return &newLogger
}

func (log *iLogger) Info(_ context.Context, msg string, data ...interface{}) {
	if log.level >= logger.Info {
		log.logger.Info(msg, data)
	}
}

func (log *iLogger) Warn(_ context.Context, msg string, data ...interface{}) {
	if log.level >= logger.Warn {
		log.logger.Warn(msg, data)
	}
}

func (log *iLogger) Error(_ context.Context, msg string, data ...interface{}) {
	if log.level >= logger.Error {
		log.logger.Error(msg, data)
	}
}

func (log *iLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if log.level > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && log.level >= logger.Error:
			sql, rows := fc()
			log.logger.WithFields(logrus.Fields{
				"exec_file":     utils.FileWithLineNum(),
				"rows_affected": rows,
				"error":         err,
				"take_time":     float64(elapsed.Nanoseconds()) / 1e6,
				"raw_sql":       sql,
			}).Error()
		case elapsed > log.SlowThreshold && log.SlowThreshold != 0 && log.level >= logger.Warn:
			sql, rows := fc()
			log.logger.WithFields(logrus.Fields{
				"exec_file":     utils.FileWithLineNum(),
				"rows_affected": rows,
				"error":         err,
				"take_time":     float64(elapsed.Nanoseconds()) / 1e6,
				"raw_sql":       sql,
			}).Warn()
		case log.level >= logger.Info:
			sql, rows := fc()
			log.logger.WithFields(logrus.Fields{
				"exec_file":     utils.FileWithLineNum(),
				"rows_affected": rows,
				"error":         err,
				"take_time":     float64(elapsed.Nanoseconds()) / 1e6,
				"raw_sql":       sql,
			}).Trace()
		}
	}
}
