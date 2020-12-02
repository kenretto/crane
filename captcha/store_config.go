package captcha

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net"
	"time"
)

// RedisStore Use redis to store the verification code content
type RedisStore struct {
	logger ILogger
	r      redis.Cmdable
}

// Verify Verify that the verification code is correct
func (s *RedisStore) Verify(id, answer string, clear bool) bool {
	return s.Get(id, clear) == answer
}

// Set store driver method
func (s *RedisStore) Set(id string, value string) {
	err := s.r.Set(context.Background(), fmt.Sprintf("catpcha_id_%s", id), value, time.Minute*10).Err()
	if err != nil {
		s.logger.Error(err)
	}
}

// Get store driver method
func (s *RedisStore) Get(id string, clear bool) string {
	val, err := s.r.Get(context.Background(), fmt.Sprintf("catpcha_id_%s", id)).Result()
	if err != nil {
		//logger.Logger("pkg.captcha").WithField("log_type", "pkg.Captcha.captcha").Error(err)
		return ""
	}
	if clear {
		err := s.r.Del(context.Background(), fmt.Sprintf("catpcha_id_%s", id)).Err()
		if err != nil {
			s.logger.Error(err)
			return ""
		}
	}
	return val
}

// SetLogger set logger
func (s *RedisStore) SetLogger(logger ILogger) {
	s.logger = logger
}

// StoreConfig redis config
type StoreConfig struct {
	RedisType string `mapstructure:"redis_type"`

	Network  string `mapstructure:"network"`
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`

	Addresses          []string `mapstructure:"addresses"`
	Username           string   `mapstructure:"username"`
	SentinelPassword   string   `mapstructure:"sentinel_password"`
	MaxRetries         int      `mapstructure:"max_retries"`
	MinRetryBackoff    string   `mapstructure:"min_retry_backoff"`
	MaxRetryBackoff    string   `mapstructure:"max_retry_backoff"`
	DialTimeout        string   `mapstructure:"dial_timeout"`
	ReadTimeout        string   `mapstructure:"read_timeout"`
	WriteTimeout       string   `mapstructure:"write_timeout"`
	PoolSize           int      `mapstructure:"pool_size"`
	MinIdleConns       int      `mapstructure:"min_idle_conns"`
	MaxConnAge         string   `mapstructure:"max_conn_age"`
	PoolTimeout        string   `mapstructure:"pool_timeout"`
	IdleTimeout        string   `mapstructure:"idle_timeout"`
	IdleCheckFrequency string   `mapstructure:"idle_check_frequency"`
	MaxRedirects       int      `mapstructure:"max_redirects"`
	ReadOnly           bool     `mapstructure:"read_only"`
	RouteByLatency     bool     `mapstructure:"route_by_latency"`
	RouteRandomly      bool     `mapstructure:"route_randomly"`
	SlaveOnly          bool     `mapstructure:"slave_only"`
	Randomly           bool     `mapstructure:"randomly"`
	MasterName         string   `mapstructure:"master_name"`

	NewClient    func(opt *redis.Options) *redis.Client
	TLSConfig    *tls.Config
	Dialer       func(ctx context.Context, network, addr string) (net.Conn, error)
	OnConnect    func(ctx context.Context, cn *redis.Conn) error
	Limiter      redis.Limiter
	ClusterSlots func(ctx context.Context) ([]redis.ClusterSlot, error)
}

func (config StoreConfig) parseDuration(s string) time.Duration {
	if s == "" {
		return 0
	}
	duration, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}

	return duration
}

// NewStore get the specified redis store according to the configuration information
func (config StoreConfig) NewStore() *RedisStore {
	var store = new(RedisStore)
	switch config.RedisType {
	case "failover":
		store.r = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:         config.MasterName,
			SentinelAddrs:      config.Addresses,
			SentinelPassword:   config.SentinelPassword,
			RouteByLatency:     config.RouteByLatency,
			RouteRandomly:      config.RouteRandomly,
			SlaveOnly:          config.SlaveOnly,
			Dialer:             config.Dialer,
			OnConnect:          config.OnConnect,
			Username:           config.Username,
			Password:           config.Password,
			DB:                 config.DB,
			MaxRetries:         config.MaxRetries,
			MinRetryBackoff:    config.parseDuration(config.MinRetryBackoff),
			MaxRetryBackoff:    config.parseDuration(config.MaxRetryBackoff),
			DialTimeout:        config.parseDuration(config.DialTimeout),
			ReadTimeout:        config.parseDuration(config.ReadTimeout),
			WriteTimeout:       config.parseDuration(config.WriteTimeout),
			PoolSize:           config.PoolSize,
			MinIdleConns:       config.MinIdleConns,
			MaxConnAge:         config.parseDuration(config.MaxConnAge),
			PoolTimeout:        config.parseDuration(config.PoolTimeout),
			IdleTimeout:        config.parseDuration(config.IdleTimeout),
			IdleCheckFrequency: config.parseDuration(config.IdleCheckFrequency),
			TLSConfig:          config.TLSConfig,
		})
	case "cluster":
		store.r = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:              config.Addresses,
			NewClient:          config.NewClient,
			MaxRedirects:       config.MaxRetries,
			ReadOnly:           config.ReadOnly,
			RouteByLatency:     config.RouteByLatency,
			RouteRandomly:      config.RouteRandomly,
			ClusterSlots:       config.ClusterSlots,
			Dialer:             config.Dialer,
			OnConnect:          config.OnConnect,
			Username:           config.Username,
			Password:           config.Password,
			MaxRetries:         config.MaxRetries,
			MinRetryBackoff:    config.parseDuration(config.MinRetryBackoff),
			MaxRetryBackoff:    config.parseDuration(config.MaxRetryBackoff),
			DialTimeout:        config.parseDuration(config.DialTimeout),
			ReadTimeout:        config.parseDuration(config.ReadTimeout),
			WriteTimeout:       config.parseDuration(config.WriteTimeout),
			PoolSize:           config.PoolSize,
			MinIdleConns:       config.MinIdleConns,
			MaxConnAge:         config.parseDuration(config.MaxConnAge),
			PoolTimeout:        config.parseDuration(config.PoolTimeout),
			IdleTimeout:        config.parseDuration(config.IdleTimeout),
			IdleCheckFrequency: config.parseDuration(config.IdleCheckFrequency),
			TLSConfig:          config.TLSConfig,
		})
	case "universal":
		store.r = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:              config.Addresses,
			DB:                 config.DB,
			Dialer:             config.Dialer,
			OnConnect:          config.OnConnect,
			Username:           config.Username,
			Password:           config.Password,
			SentinelPassword:   config.SentinelPassword,
			MaxRetries:         config.MaxRetries,
			MinRetryBackoff:    config.parseDuration(config.MinRetryBackoff),
			MaxRetryBackoff:    config.parseDuration(config.MaxRetryBackoff),
			DialTimeout:        config.parseDuration(config.DialTimeout),
			ReadTimeout:        config.parseDuration(config.ReadTimeout),
			WriteTimeout:       config.parseDuration(config.WriteTimeout),
			PoolSize:           config.PoolSize,
			MinIdleConns:       config.MinIdleConns,
			MaxConnAge:         config.parseDuration(config.MaxConnAge),
			PoolTimeout:        config.parseDuration(config.PoolTimeout),
			IdleTimeout:        config.parseDuration(config.IdleTimeout),
			IdleCheckFrequency: config.parseDuration(config.IdleCheckFrequency),
			MaxRedirects:       config.MaxRedirects,
			ReadOnly:           config.ReadOnly,
			RouteByLatency:     config.RouteByLatency,
			RouteRandomly:      config.Randomly,
			MasterName:         config.MasterName,
			TLSConfig:          config.TLSConfig,
		})
	case "default":
		store.r = redis.NewClient(&redis.Options{
			Network:            config.Network,
			Addr:               config.Addr,
			Dialer:             config.Dialer,
			OnConnect:          config.OnConnect,
			Username:           config.Username,
			Password:           config.Password,
			DB:                 config.DB,
			MaxRetries:         config.MaxRetries,
			MinRetryBackoff:    config.parseDuration(config.MinRetryBackoff),
			MaxRetryBackoff:    config.parseDuration(config.MaxRetryBackoff),
			DialTimeout:        config.parseDuration(config.DialTimeout),
			ReadTimeout:        config.parseDuration(config.ReadTimeout),
			WriteTimeout:       config.parseDuration(config.WriteTimeout),
			PoolSize:           config.PoolSize,
			MinIdleConns:       config.MinIdleConns,
			MaxConnAge:         config.parseDuration(config.MaxConnAge),
			PoolTimeout:        config.parseDuration(config.PoolTimeout),
			IdleTimeout:        config.parseDuration(config.IdleTimeout),
			IdleCheckFrequency: config.parseDuration(config.IdleCheckFrequency),
			TLSConfig:          config.TLSConfig,
			Limiter:            config.Limiter,
		})
	}
	store.logger = NewDefaultLogger()
	return store
}
