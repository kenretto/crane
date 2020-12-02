package sessions

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/kenretto/sessions"
	"github.com/kenretto/sessions/memstore"
	store "github.com/kenretto/sessions/redis"
	"github.com/spf13/viper"
	"sync"
	"time"
)

type Sessions struct {
	Driver           string           `mapstructure:"driver"`
	RedisStoreConfig RedisStoreConfig `mapstructure:"redis"`
	Key              string           `mapstructure:"key"`
	Name             string           `mapstructure:"name"`
	Domain           string           `mapstructure:"domain"`
	MaxAge           string           `mapstructure:"max_age"`
	HttpOnly         bool             `mapstructure:"http_only"`

	store sessions.Store
	conn  redis.Cmdable `mapstructure:"redis"`
	mu    sync.RWMutex
}

func (s *Sessions) OnChange(viper *viper.Viper) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = viper.Unmarshal(&s)
	switch s.Driver {
	case "redis":
		s.conn = s.RedisStoreConfig.NewRedis()
		s.store = store.NewStore(s.conn, []byte(s.Key))
	default:
		s.store = memstore.NewStore([]byte(s.Key))
	}
	duration, err := time.ParseDuration(s.MaxAge)
	if err != nil {
		panic(err)
	}
	s.store.Options(sessions.Options{MaxAge: int(duration / time.Second), Path: "/", Domain: s.Domain, HttpOnly: s.HttpOnly})
}

// Inject 启动session服务, 在自定义的路由代码中调用, 传入 *gin.Engine 对象
func (s *Sessions) Inject(engine *gin.Engine) gin.IRoutes {
	return engine.Use(sessions.Sessions(s.Name, s.store))
}

// Get 获取指定session
func Get(c *gin.Context, key string) string {
	sess := sessions.Default(c)
	val := sess.Get(key)
	if val != nil {
		return val.(string)
	}
	return ""
}

// Set 设置session
func Set(c *gin.Context, key, val string) {
	sess := sessions.Default(c)
	sess.Set(key, val)
	_ = sess.Save()
}

// Del 删除指定session
func Del(c *gin.Context, key string) {
	sess := sessions.Default(c)
	sess.Delete(key)
	_ = sess.Save()
}
