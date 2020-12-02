package timer

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

// driver 全局保存驱动
var driver StoreDriver

// StoreDriver 定时器状态保存驱动
type StoreDriver interface {
	// Set 设置指定任务的某个状态对应的值
	Set(name JobName, key, val string)
	// Get 获取某个任务的某个状态
	Get(name JobName, key string) string
}

// RedisStore redis 存储
type RedisStore struct {
	redis redis.UniversalClient
}

// Set 设置指定任务的某个状态对应的值
func (r *RedisStore) Set(name JobName, key, val string) {
	r.redis.Set(context.Background(), fmt.Sprintf("cron_status_%s_%s", name, key), val, 0)
}

// Get 获取某个任务的某个状态
func (r *RedisStore) Get(name JobName, key string) string {
	return r.redis.Get(context.Background(), fmt.Sprintf("cron_status_%s_%s", name, key)).Val()
}

// UseRedisStore 将默认全局变量 driver 注册为 RedisStore
func UseRedisStore(r redis.UniversalClient) {
	driver = &RedisStore{redis: r}
}

// UseCustomStore 使用自定义的存储
func UseCustomStore(store StoreDriver) {
	driver = store
}
