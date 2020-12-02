package orm

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"sync"
	"time"
)

type Node struct {
	LogLevel string `mapstructure:"log_level"`
	DSN      string `mapstructure:"dsn"`
	MaxIdle  int    `mapstructure:"max_idle"`
	MaxOpen  int    `mapstructure:"max_open"`
	Replicas struct {
		Connections []string `mapstructure:"connections"` // 所有复制集连接
		MaxIdle     int      `mapstructure:"max_idle"`
		MaxOpen     int      `mapstructure:"max_open"`
	} `mapstructure:"replicas"`
}

type Loader struct {
	nodes map[string]Node

	logger *logrus.Entry
	choice string
	rw     sync.RWMutex
	conns  map[string]*gorm.DB
}

func NewORM(logger *logrus.Entry) *Loader {
	loader := new(Loader)
	loader.logger = logger
	return loader
}

func (loader *Loader) OnChange(viper *viper.Viper) {
	loader.rw.Lock()
	defer loader.rw.Unlock()
	err := viper.Unmarshal(&loader.nodes)
	if err != nil {
		if loader.conns == nil {
			panic(err)
		}

	}
	loader.newInstance()
}

func (loader *Loader) newInstance() {
	conns := make(map[string]*gorm.DB)
	var err error
	for s, config := range loader.nodes {
		conns[s], err = gorm.Open(mysql.Open(config.DSN), &gorm.Config{
			Logger: &iLogger{
				logger: loader.logger,
				level:  logLevel[config.LogLevel],
			},
			PrepareStmt: true,
		})

		if err != nil {
			return
		}

		db, err := conns[s].DB()
		if err != nil {
			return
		}

		db.SetMaxIdleConns(config.MaxIdle)
		db.SetMaxOpenConns(config.MaxOpen)
		db.SetConnMaxLifetime(time.Hour)
		db.SetConnMaxIdleTime(time.Hour)

		var dialectors []gorm.Dialector
		for _, conn := range config.Replicas.Connections {
			dialectors = append(dialectors, mysql.Open(conn))
		}

		err = conns[s].Use(dbresolver.Register(dbresolver.Config{
			Replicas: dialectors,
		}).SetMaxIdleConns(config.Replicas.MaxIdle).SetMaxOpenConns(config.Replicas.MaxOpen).SetConnMaxLifetime(time.Hour).SetConnMaxIdleTime(time.Hour))
		if err != nil {
			return
		}
	}
	loader.conns = conns
}

func (loader *Loader) DB(db ...string) *gorm.DB {
	loader.rw.RLock()
	defer loader.rw.RUnlock()
	if len(db) == 1 {
		loader.choice = db[0]
	} else {
		loader.choice = "master"
	}
	return loader.conns[loader.choice]
}
