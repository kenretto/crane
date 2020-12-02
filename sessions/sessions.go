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

// Sessions session related function processing
type Sessions struct {
	Driver           string           `mapstructure:"driver"`
	RedisStoreConfig RedisStoreConfig `mapstructure:"redis"`
	Key              string           `mapstructure:"key"`
	Name             string           `mapstructure:"name"`
	Domain           string           `mapstructure:"domain"`
	MaxAge           string           `mapstructure:"max_age"`
	HTTPOnly         bool             `mapstructure:"http_only"`

	store sessions.Store
	conn  redis.Cmdable `mapstructure:"redis"`
	mu    sync.RWMutex
}

// OnChange reinitialize when configuration file changes
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
	s.store.Options(sessions.Options{MaxAge: int(duration / time.Second), Path: "/", Domain: s.Domain, HttpOnly: s.HTTPOnly})
}

// Inject start the session service, call it in the custom routing code, and import it. *gin.Engine  object
func (s *Sessions) Inject(engine *gin.Engine) gin.IRoutes {
	return engine.Use(sessions.Sessions(s.Name, s.store))
}

// Get gets the specified session
func Get(c *gin.Context, key string) string {
	sess := sessions.Default(c)
	val := sess.Get(key)
	if val != nil {
		return val.(string)
	}
	return ""
}

// Set set session
func Set(c *gin.Context, key, val string) {
	sess := sessions.Default(c)
	sess.Set(key, val)
	_ = sess.Save()
}

// Del delete the specified session
func Del(c *gin.Context, key string) {
	sess := sessions.Default(c)
	sess.Delete(key)
	_ = sess.Save()
}
