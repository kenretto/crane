package timer

import (
	"github.com/sirupsen/logrus"
	"strings"
)

// ILogger 用于将 logrus 转为 log.Logger
type ILogger struct {
	*logrus.Entry
}

// Info 实现来自 cron 的logger接口
func (logger *ILogger) Info(_ string, _ ...interface{}) {
	//logger.Logger.WithField("kvs", keysAndValues).WithField("log_type", "cron").Info(msg)
}

// Error 实现来自 cron 的logger接口
func (logger *ILogger) Error(err error, msg string, keysAndValues ...interface{}) {
	var l = logger.Logger.WithField("error", err).WithField("log_type", "cron")
	for i := 0; i < len(keysAndValues); i += 2 {
		if i%2 == 0 {
			switch keysAndValues[i].(type) {
			case string:
				l = l.WithField(keysAndValues[i].(string), strings.Split(strings.ReplaceAll(keysAndValues[i+1].(string), "\t", ""), "\n"))
			}
		}
	}
	l.Error(msg)
}
