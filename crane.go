// Package crane ...
package crane

import (
	"github.com/gin-gonic/gin"
	"github.com/kenretto/crane/captcha"
	"github.com/kenretto/crane/configurator"
	"github.com/kenretto/crane/database/orm"
	"github.com/kenretto/crane/logger"
	"github.com/kenretto/crane/password"
	"github.com/kenretto/crane/redis"
	"github.com/kenretto/crane/server"
	"github.com/kenretto/crane/sessions"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Crane summarize sub-package configuration
type Crane struct {
	captcha      *captcha.Loader
	Configurator *configurator.Configurator
	orm          *orm.Loader
	logger       *logger.Logger
	Password     *password.Password
	Redis        *redis.Redis
	server       *server.HTTPServer
	sessions     *sessions.Sessions
}

// NewCrane 子目录的每个库其实都是可以单独使用的, 如果想要整个依赖, 建议使用这个方法来初始化
func NewCrane(config string) (crane *Crane, err error) {
	crane = new(Crane)
	crane.Configurator, err = configurator.NewConfigurator(config)
	if err != nil {
		return
	}

	crane.logger = new(logger.Logger)
	crane.Configurator.Add("logger", crane.logger)

	crane.captcha = captcha.NewCaptcha()
	crane.Configurator.Add("captcha", crane.captcha)

	crane.orm = orm.NewORM(logrus.NewEntry(crane.logger.Instance()))
	crane.Configurator.Add("database", crane.orm)

	crane.Password = new(password.Password)
	crane.Configurator.Add("password", crane.Password)

	crane.Redis = new(redis.Redis)
	crane.Configurator.Add("redis", crane.Redis)

	crane.sessions = new(sessions.Sessions)
	crane.Configurator.Add("sessions", crane.sessions)

	crane.server = server.NewHTTPServer(&Logger{crane.logger.Instance()})

	return
}

// Captcha 获取验证码操作
func (crane *Crane) Captcha() *captcha.Captcha {
	return crane.captcha.Instance().WithLogger(crane.logger.Instance())
}

// ORM 获取 gorm 的操作
func (crane *Crane) ORM(db ...string) *gorm.DB {
	return crane.orm.DB(db...)
}

// Logger 获取 logger
func (crane *Crane) Logger() *logrus.Logger {
	return crane.logger.Instance()
}

// Handler set handler
func (crane *Crane) Handler(handler func(router *gin.Engine)) {
	crane.server.Handler(handler)
}

// Run start service
func (crane *Crane) Run() {
	crane.Configurator.Add("server", crane.server)
	crane.server.Run()
}

// Sessions get sessions
func (crane *Crane) Sessions() *sessions.Sessions {
	return crane.sessions
}
