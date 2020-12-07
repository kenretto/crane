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
	"github.com/kenretto/daemon"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log"
	"sync"
)

type ICrane interface {
	IntegrationLogger()
	IntegrationCaptcha()
	IntegrationORM()
	IntegrationPassword()
	IntegrationRedis()
	IntegrationSession()
	IntegrationHTTPServer()
	Integration(bind configurator.IConfig)
	Get(node string) configurator.IConfig
	WithConfigurator(config string) error
	Captcha() *captcha.Captcha
	ORM(db ...string) *gorm.DB
	Logger() *logrus.Logger
	Handler(handler func(router *gin.Engine))
	Server() *server.HTTPServer
	Run()
	Sessions() *sessions.Sessions
	Redis() *redis.Redis
	Password() *password.Password
}

// Crane summarize sub-package configuration
//  如果需要对本结构内部提供的功能重新定制, 比如对 redis 连接操作附加 OnConnect 方法, 可以自己新建对象包裹此结构,
//  然后实现 IntegrationRedis 即可, see example/main.redisMain
type Crane struct {
	captcha      *captcha.Loader
	Configurator *configurator.Configurator
	orm          *orm.Loader
	logger       *logger.Logger
	password     *password.Password
	redis        *redis.Redis
	server       *server.HTTPServer
	sessions     *sessions.Sessions

	PIDSavePath string
	ServiceName string

	container map[string]configurator.IConfig
	mu        sync.RWMutex
}

func (crane *Crane) Node() string {
	return "server"
}

func (crane *Crane) OnChange(viper *viper.Viper) {
	crane.PIDSavePath = viper.GetString("pid")
	crane.ServiceName = viper.GetString("name")
}

func (crane *Crane) PidSavePath() string {
	return crane.PIDSavePath
}

func (crane *Crane) Name() string {
	return crane.ServiceName
}

func (crane *Crane) Start() {
	crane.Configurator.Add(crane.server)
	crane.server.Listen()
}

func (crane *Crane) Stop() error {
	crane.server.Stop()
	return nil
}

func (crane *Crane) Restart() error {
	return crane.Stop()
}

// Integration integration custom
//  The function of configurator.IConfig is implemented and the configuration can be loaded here, Then you can choose to manage the life cycle of the incoming object,
//  You can also use Crane.Get to get the specified object,
//  All objects that have been configured by this method will be saved
func (crane *Crane) Integration(bind configurator.IConfig) {
	crane.mu.Lock()
	defer crane.mu.Unlock()
	if crane.container == nil {
		crane.container = make(map[string]configurator.IConfig)
	}
	crane.container[bind.Node()] = bind
	crane.Configurator.Add(crane.container[bind.Node()])
}

func (crane *Crane) Get(node string) configurator.IConfig {
	crane.mu.RLock()
	defer crane.mu.RUnlock()
	return crane.container[node]
}

// IntegrationLogger integration logger
func (crane *Crane) IntegrationLogger() {
	crane.logger = new(logger.Logger)
	crane.Configurator.Add(crane.logger)
}

// IntegrationCaptcha integration captcha
func (crane *Crane) IntegrationCaptcha() {
	crane.captcha = captcha.NewCaptcha()
	crane.Configurator.Add(crane.captcha)
}

// IntegrationORM integration gorm
func (crane *Crane) IntegrationORM() {
	crane.orm = orm.NewORM(logrus.NewEntry(crane.logger.Instance()))
	crane.Configurator.Add(crane.orm)
}

// IntegrationPassword integration password
func (crane *Crane) IntegrationPassword() {
	crane.password = new(password.Password)
	crane.Configurator.Add(crane.password)
}

// IntegrationRedis integration redis
func (crane *Crane) IntegrationRedis() {
	crane.redis = new(redis.Redis)
	crane.Configurator.Add(crane.redis)
}

// IntegrationSession integration session
func (crane *Crane) IntegrationSession() {
	crane.sessions = new(sessions.Sessions)
	crane.Configurator.Add(crane.sessions)
}

// IntegrationHTTPServer integration http server
func (crane *Crane) IntegrationHTTPServer() {
	crane.server = server.NewHTTPServer(&Logger{logrus.NewEntry(crane.logger.Instance())})
}

func (crane *Crane) WithConfigurator(config string) error {
	var err error
	crane.Configurator, err = configurator.NewConfigurator(config)
	crane.Configurator.Add(crane)
	return err
}

// NewCrane 子目录的每个库其实都是可以单独使用的, 如果想要整个依赖, 建议使用这个方法来初始化
func NewCrane(config string) (crane ICrane, err error) {
	crane = new(Crane)
	err = crane.WithConfigurator(config)
	if err != nil {
		return
	}

	crane.IntegrationLogger()
	crane.IntegrationCaptcha()
	crane.IntegrationRedis()
	crane.IntegrationORM()
	crane.IntegrationPassword()
	crane.IntegrationSession()
	crane.IntegrationHTTPServer()
	daemon.Register(daemon.NewProcess(crane.(*Crane)))
	return
}

func NewCraneWithCustom(config string, crane ICrane) error {
	err := crane.WithConfigurator(config)
	if err != nil {
		return err
	}

	crane.IntegrationLogger()
	crane.IntegrationCaptcha()
	crane.IntegrationRedis()
	crane.IntegrationORM()
	crane.IntegrationPassword()
	crane.IntegrationSession()
	crane.IntegrationHTTPServer()
	if c, ok := crane.(*Crane); ok {
		daemon.Register(daemon.NewProcess(c))
	}
	return nil
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

func (crane *Crane) Server() *server.HTTPServer {
	return crane.server
}

// Run start service
func (crane *Crane) Run() {
	err := daemon.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

// Sessions get sessions
func (crane *Crane) Sessions() *sessions.Sessions {
	return crane.sessions
}

// Redis get redis
func (crane *Crane) Redis() *redis.Redis {
	return crane.redis
}

// Password get password
func (crane *Crane) Password() *password.Password {
	return crane.password
}
