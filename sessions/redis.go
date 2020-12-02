package sessions

import (
	"context"
	"crypto/tls"
	"github.com/go-redis/redis/v8"
	"net"
	"time"
)

// RedisStoreConfig redis config
type RedisStoreConfig struct {
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

func (rsc RedisStoreConfig) parseDuration(s string) time.Duration {
	if s == "" {
		return 0
	}
	duration, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}

	return duration
}

// NewRedis get a storage operation instance of redis according to the configuration information
func (rsc RedisStoreConfig) NewRedis() redis.Cmdable {
	var r redis.Cmdable
	switch rsc.RedisType {
	case "failover":
		r = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:         rsc.MasterName,
			SentinelAddrs:      rsc.Addresses,
			SentinelPassword:   rsc.SentinelPassword,
			RouteByLatency:     rsc.RouteByLatency,
			RouteRandomly:      rsc.RouteRandomly,
			SlaveOnly:          rsc.SlaveOnly,
			Dialer:             rsc.Dialer,
			OnConnect:          rsc.OnConnect,
			Username:           rsc.Username,
			Password:           rsc.Password,
			DB:                 rsc.DB,
			MaxRetries:         rsc.MaxRetries,
			MinRetryBackoff:    rsc.parseDuration(rsc.MinRetryBackoff),
			MaxRetryBackoff:    rsc.parseDuration(rsc.MaxRetryBackoff),
			DialTimeout:        rsc.parseDuration(rsc.DialTimeout),
			ReadTimeout:        rsc.parseDuration(rsc.ReadTimeout),
			WriteTimeout:       rsc.parseDuration(rsc.WriteTimeout),
			PoolSize:           rsc.PoolSize,
			MinIdleConns:       rsc.MinIdleConns,
			MaxConnAge:         rsc.parseDuration(rsc.MaxConnAge),
			PoolTimeout:        rsc.parseDuration(rsc.PoolTimeout),
			IdleTimeout:        rsc.parseDuration(rsc.IdleTimeout),
			IdleCheckFrequency: rsc.parseDuration(rsc.IdleCheckFrequency),
			TLSConfig:          rsc.TLSConfig,
		})
	case "cluster":
		r = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:              rsc.Addresses,
			NewClient:          rsc.NewClient,
			MaxRedirects:       rsc.MaxRetries,
			ReadOnly:           rsc.ReadOnly,
			RouteByLatency:     rsc.RouteByLatency,
			RouteRandomly:      rsc.RouteRandomly,
			ClusterSlots:       rsc.ClusterSlots,
			Dialer:             rsc.Dialer,
			OnConnect:          rsc.OnConnect,
			Username:           rsc.Username,
			Password:           rsc.Password,
			MaxRetries:         rsc.MaxRetries,
			MinRetryBackoff:    rsc.parseDuration(rsc.MinRetryBackoff),
			MaxRetryBackoff:    rsc.parseDuration(rsc.MaxRetryBackoff),
			DialTimeout:        rsc.parseDuration(rsc.DialTimeout),
			ReadTimeout:        rsc.parseDuration(rsc.ReadTimeout),
			WriteTimeout:       rsc.parseDuration(rsc.WriteTimeout),
			PoolSize:           rsc.PoolSize,
			MinIdleConns:       rsc.MinIdleConns,
			MaxConnAge:         rsc.parseDuration(rsc.MaxConnAge),
			PoolTimeout:        rsc.parseDuration(rsc.PoolTimeout),
			IdleTimeout:        rsc.parseDuration(rsc.IdleTimeout),
			IdleCheckFrequency: rsc.parseDuration(rsc.IdleCheckFrequency),
			TLSConfig:          rsc.TLSConfig,
		})
	case "universal":
		r = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:              rsc.Addresses,
			DB:                 rsc.DB,
			Dialer:             rsc.Dialer,
			OnConnect:          rsc.OnConnect,
			Username:           rsc.Username,
			Password:           rsc.Password,
			SentinelPassword:   rsc.SentinelPassword,
			MaxRetries:         rsc.MaxRetries,
			MinRetryBackoff:    rsc.parseDuration(rsc.MinRetryBackoff),
			MaxRetryBackoff:    rsc.parseDuration(rsc.MaxRetryBackoff),
			DialTimeout:        rsc.parseDuration(rsc.DialTimeout),
			ReadTimeout:        rsc.parseDuration(rsc.ReadTimeout),
			WriteTimeout:       rsc.parseDuration(rsc.WriteTimeout),
			PoolSize:           rsc.PoolSize,
			MinIdleConns:       rsc.MinIdleConns,
			MaxConnAge:         rsc.parseDuration(rsc.MaxConnAge),
			PoolTimeout:        rsc.parseDuration(rsc.PoolTimeout),
			IdleTimeout:        rsc.parseDuration(rsc.IdleTimeout),
			IdleCheckFrequency: rsc.parseDuration(rsc.IdleCheckFrequency),
			MaxRedirects:       rsc.MaxRedirects,
			ReadOnly:           rsc.ReadOnly,
			RouteByLatency:     rsc.RouteByLatency,
			RouteRandomly:      rsc.Randomly,
			MasterName:         rsc.MasterName,
			TLSConfig:          rsc.TLSConfig,
		})
	case "default":
		r = redis.NewClient(&redis.Options{
			Network:            rsc.Network,
			Addr:               rsc.Addr,
			Dialer:             rsc.Dialer,
			OnConnect:          rsc.OnConnect,
			Username:           rsc.Username,
			Password:           rsc.Password,
			DB:                 rsc.DB,
			MaxRetries:         rsc.MaxRetries,
			MinRetryBackoff:    rsc.parseDuration(rsc.MinRetryBackoff),
			MaxRetryBackoff:    rsc.parseDuration(rsc.MaxRetryBackoff),
			DialTimeout:        rsc.parseDuration(rsc.DialTimeout),
			ReadTimeout:        rsc.parseDuration(rsc.ReadTimeout),
			WriteTimeout:       rsc.parseDuration(rsc.WriteTimeout),
			PoolSize:           rsc.PoolSize,
			MinIdleConns:       rsc.MinIdleConns,
			MaxConnAge:         rsc.parseDuration(rsc.MaxConnAge),
			PoolTimeout:        rsc.parseDuration(rsc.PoolTimeout),
			IdleTimeout:        rsc.parseDuration(rsc.IdleTimeout),
			IdleCheckFrequency: rsc.parseDuration(rsc.IdleCheckFrequency),
			TLSConfig:          rsc.TLSConfig,
			Limiter:            rsc.Limiter,
		})
	}

	return r
}
